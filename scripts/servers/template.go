package main

import (
	_ "embed"
	"strings"
	"text/template"

	"github.com/Masterminds/semver"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var (
	indexTemplate *template.Template

	//go:embed template/index.ctmpl
	indexTemplateText string
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
	"IsCurrent": func(in string) *bool {
		ver, err := semver.NewVersion(in)
		if err != nil {
			panic(err)
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
}

func initializeTemplateRenderer() {
	var err error

	indexTemplate, err = template.New("").Funcs(tmplFuncs).Parse(indexTemplateText)
	if err != nil {
		panic(err)
	}
}
