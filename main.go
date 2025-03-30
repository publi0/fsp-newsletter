package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Article struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

func main() {

	req, err := http.Get("https://www1.folha.uol.com.br/fsp/")
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()

	// Parse HTML
	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Map to store article groups
	articleGroups := make(map[string][]Article)

	// Find all article channels
	doc.Find(".c-channel").Each(func(i int, channel *goquery.Selection) {
		// Extract group name
		groupName := channel.Find(".c-channel__title a").First().Text()

		// Collect articles
		var articles []Article
		channel.Find(".c-channel__headline a").Each(func(j int, article *goquery.Selection) {
			title := article.Text()
			link, _ := article.Attr("href")
			articles = append(articles, Article{
				Title: title,
				Link:  link,
			})
		})

		articleGroups[groupName] = articles
	})

	// Find highest resolution cover image
	coverImg, _ := doc.Find("#capa-nacional").Attr("data-full-img")

	// Print results
	fmt.Println("Capa do dia (highest resolution):", coverImg)
	fmt.Println("\nArticle Groups:")
	for group, articles := range articleGroups {
		fmt.Printf("\n%s:\n", group)
		for _, article := range articles {
			fmt.Printf("- %s\n  %s\n", article.Title, article.Link)
		}
	}
}
