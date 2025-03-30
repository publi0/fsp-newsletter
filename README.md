# Folha de São Paulo Newsletter

A simple automated service that scrapes the "Folha de São Paulo" newspaper website and sends a formatted newsletter via email.

## Overview

This application fetches headlines and articles from the Folha de São Paulo website (https://www1.folha.uol.com.br/fsp/) and sends them as a formatted HTML newsletter to a configured email address. It can be run on-demand or scheduled to automatically send daily newsletters at 7:00 AM.

## Features

- Scrapes articles, headlines, authors, and links from Folha de São Paulo
- Groups articles by section/category
- Includes the newspaper cover image when available
- Formats content in a clean, readable HTML email
- Can be triggered manually or run as a scheduled service

## Requirements

- Go 1.24+
- Resend API key (for sending emails)
- Internet connection to access the Folha de São Paulo website

## Installation

1. Clone the repository
2. Install dependencies:

```
go mod tidy
```

## Configuration

Set the following environment variable:

- `API_KEY`: Your Resend API key

You can also modify the sender and recipient email addresses in the `main.go` file.

## Usage

### Run as a scheduled service

```
go run main.go
```

This will start the newsletter scheduler. The application will run continuously and send a newsletter every day at 7:00 AM.

### Send a newsletter immediately

```
go run main.go send
```

This will immediately scrape the news and send the newsletter without waiting for the scheduled time.

## How It Works

1. The application scrapes the Folha de São Paulo FSP section
2. It organizes articles by their respective categories
3. It formats the content using an HTML template
4. It sends the formatted newsletter via the Resend email API

## Dependencies

- [goquery](https://github.com/PuerkitoBio/goquery) - For HTML parsing and data extraction
- [resend-go](https://github.com/resendlabs/resend-go) - For sending emails via the Resend API

## Disclaimer

This project is for educational purposes only. The content belongs to Folha de São Paulo, and this tool is not affiliated with or endorsed by Folha de São Paulo. Please respect their terms of service and copyright when using this tool.
