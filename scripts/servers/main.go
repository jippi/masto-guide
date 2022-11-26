package main

import (
	"io/ioutil"
	"math"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/Masterminds/semver"
	log "github.com/sirupsen/logrus"
	"gopkg.in/resty.v1"
	"gopkg.in/yaml.v3"
)

var (
	config = &Config{}

	// Categories we're managing
	categories []*Category

	// Worker coordination
	readerCh = make(chan Server)
	wg       sync.WaitGroup

	// HTTP client
	httpClient = resty.New().
			SetTimeout(2*time.Second).
			SetHeader("User-Agent", "Masto-Guide")

	// The latest release of Mastodon according to GitHub releases
	mastodonVersion *semver.Constraints
)

func main() {
	// log.SetLevel(log.DebugLevel)
	logger := log.WithField("subsystem", "main")

	loadConfigFile()
	initializeTemplateRenderer()
	getLatestReleaseOfMastodon()

	// The order we want categories to show up in the markdown output
	categories = []*Category{
		config.Categories["open"],
		config.Categories["review"],
		config.Categories["invite"],
		config.Categories["private"],
	}

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
		readerCh <- server
	}

	// Prevent any further writing to the channel
	close(readerCh)

	logger.Info("Waiting for all server information to arrive")
	wg.Wait()
	logger.Info("Done reading server information")

	// Sorting the servers by name to keep it consistent and fair
	for _, cat := range categories {
		sort.SliceStable(cat.Servers, func(a, b int) bool {
			return cat.Servers[a].Domain < cat.Servers[b].Domain
		})
	}

	// Open the servers.md MarkDown file for writing
	file, err := os.OpenFile("../../docs/dk/servers.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logger.Fatal(err)
	}

	// Struct used in the template for rendering
	payload := struct {
		Categories []*Category
		UpdateAt   string
	}{
		Categories: categories,
		UpdateAt:   time.Now().Format(time.RFC822),
	}

	logger.Info("Rendering MarkDown file")
	// Render and write the markdown template
	if err := indexTemplate.Execute(file, payload); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Rendering completed successfully")
}

func fetchServerInformation(server Server, log *log.Entry) {
	defer wg.Done()

	// Very simplistic retry policy
	for i := 0; i < 5; i++ {
		logger := log.WithField("subsystem", "worker").WithField("attempt", i).WithField("server", server.URL)
		if i > 0 {
			logger.Info("Fetching server information")
		}

		retry := func(err error) {
			logger.Error(err)
			// We will sleep 1, 2, 3, 4, 5 seconds
			time.Sleep(time.Duration(i) * time.Second)
		}

		serverResponse := ServerResponse{}
		_, err := httpClient.R().SetResult(&serverResponse).Get(server.URL + "/api/v2/instance")
		if err != nil {
			retry(err)
			continue
		}

		// Copy the covenant setting over
		serverResponse.MastodonCovenant = server.Covenant

		// Categorize the server based on it's settings
		category := serverResponse.Categorize(server)

		// Append the server to the Category's server list
		// ordering doesn't matter here, we'll sort them later.
		category.Servers = append(category.Servers, serverResponse)

		return
	}
}

func getLatestReleaseOfMastodon() {
	logger := log.WithField("subsystem", "get-mastodon-release")
	logger.Debug("Finding latest release of Mastodon")

	serverResponse := GithubReleaseResponse{}
	_, err := httpClient.R().SetResult(&serverResponse).Get("https://api.github.com/repos/mastodon/mastodon/releases/latest")
	if err != nil {
		logger.Fatal(err)
	}

	mastodonVersion, err = semver.NewConstraint(">= " + serverResponse.TagName)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Found version %s", serverResponse.TagName)
}

func loadConfigFile() {
	logger := log.WithField("subsystem", "config")
	logger.Debug("Loading configuration file")

	serversContent, err := ioutil.ReadFile("../../server-config.yml")
	if err != nil {
		logger.Fatal(err)
	}

	if err := yaml.Unmarshal(serversContent, config); err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Configuration file successfully loaded")
}

func boolPtr(in bool) *bool {
	return &in
}
