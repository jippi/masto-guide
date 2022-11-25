package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	// Outcome
	categories = []*Category{
		OpenCategory,
		OpenCategoryWithApproval,
		ClosedCategory,
		PrivateCategory,
	}

	// Worker coordination
	readerCh = make(chan Server)
	wg       sync.WaitGroup
	lock     sync.Mutex

	// HTTP client
	httpClient = http.Client{
		Timeout: time.Second * 2, // Timeout after 2 seconds
	}

	// Template management
	serverTemplate *template.Template
	//go:embed template/server.ctmpl
	serverTemplateText string

	indexTemplate *template.Template
	//go:embed template/index.ctmpl
	indexTemplateText string
)

func main() {
	initTemplate()

	// Start workers
	for w := 1; w <= 10; w++ {
		go func() {
			for j := range readerCh {
				worker(j)
			}
		}()
	}

	// enqueue jobs
	for _, server := range servers {
		wg.Add(1)
		readerCh <- server
	}
	close(readerCh)

	fmt.Println("Waiting ....")
	wg.Wait()
	fmt.Println("Done! ....")

	for _, cat := range categories {
		sort.SliceStable(cat.Servers, func(a, b int) bool {
			return cat.Servers[a].Domain < cat.Servers[b].Domain
		})
	}

	file, err := os.OpenFile("../../docs/dk/servers.md", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	payload := struct {
		Categories []*Category
		UpdateAt   string
	}{
		Categories: categories,
		UpdateAt:   time.Now().Format(time.RFC822),
	}
	if err := indexTemplate.Execute(file, payload); err != nil {
		log.Fatal(err)
	}
}

func worker(in Server) {
	defer wg.Done()

	// Retry a couple of times
	for i := 0; i < 5; i++ {
		fmt.Println("Working on", in.URL, "attempt number", i)

		retry := func(err error) {
			fmt.Println(in.URL, i, err)
			time.Sleep(time.Duration(i) * time.Second)
		}

		req, err := http.NewRequest(http.MethodGet, in.URL+"/api/v2/instance", nil)
		if err != nil {
			retry(err)
			continue
		}
		req.Header.Set("User-Agent", "MastoGuide")

		res, getErr := httpClient.Do(req)
		if getErr != nil {
			retry(err)
			continue
		}

		if res.Body != nil {
			defer res.Body.Close()
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			retry(err)
			continue
		}

		serverResponse := ServerResponse{}
		jsonErr := json.Unmarshal(body, &serverResponse)
		if jsonErr != nil {
			retry(err)
			continue
		}

		serverResponse.MastodonCovenant = in.MastodonCovenant

		category := serverResponse.Categorize(in)
		category.Servers = append(category.Servers, serverResponse)
		return
	}
}

var tmplFuncs = template.FuncMap{
	"prefixWith": func(in, prefix string) string {
		in = strings.ReplaceAll(in, "\r", "")

		lines := strings.Split(in, "\n")
		for i, line := range lines {
			lines[i] = prefix + line
		}

		return strings.Join(lines, "\n")
	},
	"NoNewlines": func(in string) string {
		return strings.ReplaceAll(strings.ReplaceAll(in, "\n", " "), "\r", " ")
	},
	"BoolIcon": func(in *bool) string {
		if in == nil {
			return "❓"
		}

		if *in {
			return "✅"
		}

		return "❌"
	},
}

func initTemplate() {
	var err error
	serverTemplate, err = template.New("").Funcs(tmplFuncs).Parse(serverTemplateText)
	if err != nil {
		panic(err)
	}

	indexTemplate, err = template.New("").Funcs(tmplFuncs).Parse(indexTemplateText)
	if err != nil {
		panic(err)
	}
}
