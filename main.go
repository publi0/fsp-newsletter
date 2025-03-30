package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/resendlabs/resend-go"
)

type Article struct {
	Title  string
	Link   string
	Author string
}

type NewsletterData struct {
	Date        string
	CoverImgURL string
	Groups      map[string][]Article
}

func main() {
	resendAPIKey := os.Getenv("API_KEY")
	senderEmail := "newsletter@alerts.publio.dev"
	recipientEmail := "felipe@publio.dev"

	if len(os.Args) > 1 {
		if os.Args[1] == "send" {
			fmt.Println("Sending newsletter...")
			data, err := scrapeNewsData()
			if err != nil {
				log.Fatal(err)
			}
			err = sendNewsletter(data, resendAPIKey, senderEmail, recipientEmail)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Newsletter sent successfully!")
			return
		}
	}

	fmt.Println("Starting newsletter scheduler. Will send daily at 07:00 AM.")
	for {
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, now.Location())
		if now.After(nextRun) {
			nextRun = nextRun.Add(24 * time.Hour)
		}

		duration := nextRun.Sub(now)
		fmt.Printf("Next newsletter will be sent in %v\n", duration.Round(time.Second))

		time.Sleep(duration)

		data, err := scrapeNewsData()
		if err != nil {
			log.Println("Error scraping news:", err)
			continue
		}

		err = sendNewsletter(data, resendAPIKey, senderEmail, recipientEmail)
		if err != nil {
			log.Println("Error sending newsletter:", err)
			continue
		}

		fmt.Println("Newsletter sent successfully at", time.Now().Format("2006-01-02 15:04:05"))
	}
}

func scrapeNewsData() (NewsletterData, error) {
	req, err := http.Get("https://www1.folha.uol.com.br/fsp/")
	if err != nil {
		return NewsletterData{}, err
	}
	defer req.Body.Close()

	doc, err := goquery.NewDocumentFromReader(req.Body)
	if err != nil {
		return NewsletterData{}, err
	}

	articleGroups := make(map[string][]Article)

	doc.Find(".c-channel").Each(func(i int, channel *goquery.Selection) {
		groupName := channel.Find(".c-channel__title a").First().Text()
		var articles []Article

		channel.Find(".c-channel__headline a").Each(func(j int, article *goquery.Selection) {
			title := article.Text()
			link, _ := article.Attr("href")
			var author string

			kicker := article.Parent().Prev().Filter(".c-kicker")
			if kicker.Length() > 0 {
				if authorLink := kicker.Find("a"); authorLink.Length() > 0 {
					author = authorLink.Text()
				} else {
					author = kicker.Text()
				}
			}

			articles = append(articles, Article{
				Title:  title,
				Link:   link,
				Author: author,
			})
		})

		articleGroups[groupName] = articles
	})

	coverImg, _ := doc.Find("#capa-nacional").Attr("data-full-img")

	return NewsletterData{
		Date:        time.Now().Format("Monday, 02 January 2006"),
		CoverImgURL: coverImg,
		Groups:      articleGroups,
	}, nil
}

func sendNewsletter(data NewsletterData, apiKey, sender, recipient string) error {
	fmt.Println("Building newsletter...")
	t, err := template.New("newsletter").ParseFiles("template.html")
	if err != nil {
		return err
	}

	var htmlBody bytes.Buffer
	err = t.Execute(&htmlBody, data)
	if err != nil {
		return err
	}

	client := resend.NewClient(apiKey)
	subject := "Folha de São Paulo - Notícias do Dia " + data.Date

	params := &resend.SendEmailRequest{
		From:    sender,
		To:      []string{recipient},
		Subject: subject,
		Html:    htmlBody.String(),
	}

	fmt.Println("Sending email via Resend API...")
	_, err = client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Println("Email sent successfully!")
	return nil
}
