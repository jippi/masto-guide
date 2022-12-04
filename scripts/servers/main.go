package main

import (
	"io/ioutil"
	"math"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	"github.com/goodsign/monday"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v3"
)

var (
	config = &Config{}
	logger *logrus.Entry

	// Track when we started the program
	startTime time.Time

	// Worker coordination
	readerCh = make(chan Server)
	wg       sync.WaitGroup

	// HTTP client
	httpClient = resty.New().
			SetTimeout(5*time.Second).
			SetHeader("User-Agent", "Masto-Guide")

	// The latest release of Mastodon according to GitHub releases
	mastodonVersion *semver.Constraints

	serverErrors  = make(map[string]error)
	categoryIndex = make(map[string]int)
)

func main() {
	copenhagen, err := time.LoadLocation("Europe/Copenhagen")
	if err != nil {
		panic(err)
	}

	startTime = time.Now().In(copenhagen)

	if val, ok := os.LookupEnv("LOG_LEVEL"); ok && val == "debug" {
		log.SetLevel(log.DebugLevel)
	}

	logger = log.WithField("subsystem", "main")

	loadConfigFile()
	initializeTemplateRenderer()
	writeTerraform()

	if val, ok := os.LookupEnv("TF_ONLY"); ok && val == "1" {
		return
	}

	getLatestReleaseOfMastodon()

	// Start workers
	for w := 1; w <= int(math.Min(5, float64(len(config.Servers)))); w++ {
		go func(w int) {
			workLogger := logger.WithField("worker_id", w)
			workLogger.Debug("Starting worker")

			for job := range readerCh {
				fetchServerInformation(job, workLogger)
			}
		}(w)
	}

	// enqueue servers that need to be fetched
	for _, server := range config.Servers {
		wg.Add(1)
		readerCh <- *server
	}

	// Prevent any further writing to the channel
	close(readerCh)

	logger.Info("Waiting for all server information to arrive")
	wg.Wait()
	logger.Info("Done reading server information")

	// Sorting the servers by name to keep it consistent and fair
	for _, cat := range config.Categories {
		sort.Slice(cat.Servers, func(a, b int) bool {
			serverA, serverB := cat.Servers[a], cat.Servers[b]

			// Server A ✅ covenant, but Server B does not - A should come first
			if serverA.HasCommittedToServerCovenant() && !serverB.HasCommittedToServerCovenant() {
				return true
			}

			// Server A has no covenant, but Server A does - B should come first
			if !serverA.HasCommittedToServerCovenant() && serverB.HasCommittedToServerCovenant() {
				return false
			}

			// Sort by name for the remainder
			return serverA.Domain < serverB.Domain
		})
	}

	var markdownFile *os.File

	// Open the servers.md MarkDown markdownFile for writing
	markdownFile, err = os.OpenFile("../../docs/dk/servers.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logger.Fatal(err)
	}

	// Struct used in the template for rendering
	payload := struct {
		Categories []*Category
		Servers    []*Server
		UpdateAt   string
		Errors     map[string]error
	}{
		Categories: config.Categories,
		Servers:    config.Servers,
		UpdateAt:   monday.Format(startTime, "Monday, _2 January kl 15:04", monday.LocaleDaDK),
		Errors:     serverErrors,
	}

	logger.Info("Rendering MarkDown file")
	// Render and write the markdown template
	if err := indexTemplate.Execute(markdownFile, payload); err != nil {
		logger.WithError(err).Fatal("Could not render markdown file")
	}
	logger.Info("Rendering completed successfully")
}

func writeTerraform() {
	// Open the terraform file for writing
	terraformFile, err := os.OpenFile("monitoring/sites.tf", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logger.Fatal(err)
	}

	// Struct used in the template for rendering
	payload := struct {
		Servers         []*Server
		ExcludedDomains []string
	}{
		Servers:         getServersWithMonitoring(config.Servers),
		ExcludedDomains: getDomainsWithoutMonitoring(config.Servers),
	}

	logger.Info("Rendering Terraform file")
	// Render and write the markdown template
	if err := terraformTemplate.Execute(terraformFile, payload); err != nil {
		logger.WithError(err).Fatal("Could not render terraform file")
	}

	logger.Info("Rendering completed successfully")
}

func fetchServerInformation(server Server, log *log.Entry) {
	defer wg.Done()

	var lastError error

	// Very simplistic retry policy
	for i := 0; i < 5; i++ {
		logger := log.WithField("subsystem", "worker").WithField("attempt", i).WithField("server", server.Domain)
		logger.Info("Fetching server information")

		retry := func(err error) {
			lastError = err

			logger.Error(err)
			// We will sleep 1, 2, 3, 4, 5 seconds
			time.Sleep(time.Duration(i) * time.Second)
		}

		serverResponse := &ServerResponse{}
		if _, err := httpClient.R().SetResult(serverResponse).Get("https://" + server.Domain + "/api/v2/instance"); err != nil {
			retry(err)
			continue
		}

		// Copy config over
		serverResponse.MastodonCovenant = server.Covenant
		serverResponse.WithoutMonitoring = server.WithoutMonitoring

		// Categorize the server based on it's settings
		category := serverResponse.Categorize(server)

		// Append the server to the Category's server list
		// ordering doesn't matter here, we'll sort them later.
		category.Servers = append(category.Servers, serverResponse)

		return
	}

	serverErrors[server.Domain] = lastError
}

func getLatestReleaseOfMastodon() {
	logger := log.WithField("subsystem", "get-mastodon-release")
	logger.Debug("Finding latest release of Mastodon")

	serverResponse := GithubReleaseResponse{}
	resp, err := httpClient.R().SetBasicAuth("Bearer", os.Getenv("GITHUB_TOKEN")).SetResult(&serverResponse).Get("https://api.github.com/repos/mastodon/mastodon/releases/latest")
	if err != nil {
		logger.
			WithError(err).
			Fatal("Could not read latest release of Mastodon from GitHub API")
	}

	if resp.IsError() {
		logger.
			WithField("response_body", resp.String()).
			WithField("response_code", resp.StatusCode()).
			Fatal("Did not get HTTP 200 OK from GitHub API")
	}

	mastodonVersion, err = semver.NewConstraint(">= " + serverResponse.TagName)
	if err != nil {
		logger.
			WithError(err).
			Fatal("Could not create version constraint")
	}

	logger.Debug("Found version %s", serverResponse.TagName)
}

func loadConfigFile() {
	logger := log.WithField("subsystem", "config")
	logger.Debug("Loading configuration file")

	categories, err := ioutil.ReadFile("config/categories.yml")
	if err != nil {
		logger.Fatal(err)
	}

	if err := yaml.Unmarshal(categories, &config.Categories); err != nil {
		logger.Fatal(err)
	}

	for idx, cat := range config.Categories {
		categoryIndex[cat.ID] = idx
	}

	servers, err := ioutil.ReadFile("config/servers.yml")
	if err != nil {
		logger.Fatal(err)
	}
	if err := yaml.Unmarshal(servers, &config.Servers); err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Configuration file successfully loaded")
}

func getServersWithMonitoring(in []*Server) []*Server {
	monitored := make([]*Server, 0)

	for _, server := range in {
		if server.WithoutMonitoring {
			continue
		}

		monitored = append(monitored, server)
	}

	return monitored
}

func getDomainsWithoutMonitoring(in []*Server) []string {
	monitored := make([]string, 0)

	for _, server := range in {
		if !server.WithoutMonitoring {
			continue
		}

		monitored = append(monitored, server.Domain)
	}

	return monitored
}

func boolPtr(in bool) *bool {
	return &in
}
