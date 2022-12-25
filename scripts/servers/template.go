package main

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	indexTemplate     *template.Template
	terraformTemplate *template.Template

	//go:embed template/index.ctmpl
	indexTemplateText string

	//go:embed template/terraform.ctmpl
	terraformTemplateText string
)

var tmplFuncs = template.FuncMap{
	"trimSpace": func(in string) string {
		return strings.TrimSpace(in)
	},
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
	"NumberFormat": func(in int) string {
		return message.NewPrinter(language.Danish).Sprintf("%d\n", in)
	},
	"IsCurrent": func(in ServerResponse) *bool {
		ver, err := semver.NewVersion(in.Version)
		if err != nil {
			panic(fmt.Errorf("Could not create SemVer for %s: %w", in.Domain, err))
		}

		if mastodonVersion.Check(ver) {
			return boolPtr(true)
		}

		return boolPtr(false)
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
	"TerraformID": func(in string) string {
		return strings.ReplaceAll(in, ".", "_")
	},
	"DD_SplitIntoFourX": func(in int) int {
		return (in * 6) % 12
	},
	"DD_SplitIntoFourY": func(in int) int {
		return (in * 2)
	},
	"ExcludedDomainsQuery": func(list []string) string {
		if len(list) == 0 {
			return ""
		}

		return fmt.Sprintf(" AND site NOT IN (%s)", strings.Join(list, ","))
	},
	"ExcludedDomainsEventQuery": func(list []string) string {
		if len(list) == 0 {
			return ""
		}

		return fmt.Sprintf(" -site:(%s)", strings.Join(list, ","))
	},
}

func initializeTemplateRenderer() {
	var err error

	indexTemplate, err = template.New("").Funcs(tmplFuncs).Parse(indexTemplateText)
	if err != nil {
		panic(err)
	}

	terraformTemplate, err = template.New("").Funcs(tmplFuncs).Parse(terraformTemplateText)
	if err != nil {
		panic(err)
	}
}
