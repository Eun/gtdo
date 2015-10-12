package page

import (
	"html/template"
	"net/url"

	"github.com/shurcooL/htmlg"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Tabs renders the html for <nav> element with tab header links.
func Tabs(path string, rawQuery string) template.HTML {
	query, _ := url.ParseQuery(rawQuery)

	selectedTab := query.Get("tab")

	var ns []*html.Node

	for _, tab := range []struct {
		id   string
		name string
		beta bool
	}{
		{id: "summary", name: "Summary"},
		{id: "", name: "Source Code"},
		{id: "imports", name: "Imports"},
	} {
		a := &html.Node{Type: html.ElementNode, Data: atom.A.String()}
		if tab.id == selectedTab {
			a.Attr = []html.Attribute{{Key: atom.Class.String(), Val: "selected"}}
		} else {
			query := query
			if tab.id == "" {
				query.Del("tab")
			} else {
				query.Set("tab", tab.id)
			}
			u := url.URL{
				Path:     path,
				RawQuery: query.Encode(),
			}
			a.Attr = []html.Attribute{
				{Key: atom.Href.String(), Val: u.String()},
			}
		}
		a.AppendChild(htmlg.Text(tab.name))
		if tab.beta {
			span := &html.Node{
				Type: html.ElementNode, Data: atom.Span.String(),
				Attr: []html.Attribute{{Key: atom.Class.String(), Val: "beta"}},
			}
			span.AppendChild(htmlg.Text("beta"))
			a.AppendChild(span)
		}
		ns = append(ns, a)
	}

	return htmlg.Render(ns...)
}