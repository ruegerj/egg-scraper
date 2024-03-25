package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("🐣 Enter an URL to start the hunt for some juicy eggs:")
	targetUrl, _ := reader.ReadString('\n')
	targetUrl = strings.Split(targetUrl, "\n")[0]

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{targetUrl},
		ParseFunc: scrapeProductUrls,
		Exporters: []export.Exporter{&export.PrettyPrint{}},
	}).Start()
}

func scrapeProductUrls(g *geziyor.Geziyor, r *client.Response) {
	// find all product links on page
	r.HTMLDoc.Find("ul.products-grid>li.item>h2.product-name").Each(func(i int, s *goquery.Selection) {
		if productUrl, exists := s.Find("a").Attr("href"); exists {
			g.Get(productUrl, scrapeEgg)
		}
	})

	// find next page link and scrape
	if nextPage, found := r.HTMLDoc.Find("ol>li.next>a.next").First().Attr("href"); found {
		g.Get(nextPage, scrapeProductUrls)
	}
}

func scrapeEgg(g *geziyor.Geziyor, r *client.Response) {
	if r.HTMLDoc.Find("p.open-contest").First().Index() == -1 {
		return
	}

	g.Exports <- r.Request.URL.String()
}
