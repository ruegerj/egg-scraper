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
	fmt.Println("Let the egg hunt begin :)")
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("🐣 Enter an URL to start the hunt for some juicy eggs:")
	targetUrl, _ := reader.ReadString('\n')
	targetUrl = strings.Split(targetUrl, "\n")[0]

	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs: []string{targetUrl},
		ParseFunc: func(g *geziyor.Geziyor, r *client.Response) {
			r.HTMLDoc.Find("ul.products-grid>li.item>h2.product-name").Each(func(i int, s *goquery.Selection) {
				productUrl, exists := s.Find("a").Attr("href")
				if !exists {
					return
				}

				g.Exports <- productUrl
			})

		},
		Exporters: []export.Exporter{&export.PrettyPrint{}},
	}).Start()
}
