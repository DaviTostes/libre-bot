# LibreBot

LibreBot is a Go CLI tool for scraping product cards from Mercado Livre's affiliate hub and generating affiliate links.

## Features

- Scrapes poly cards from the affiliate hub.
- Generates affiliate links for Mercado Livre products using browser automation.

## Installation

1. Ensure Go 1.25+ is installed.
2. Clone the repo: `git clone https://github.com/davitostes/libre-bot`
3. Install dependencies: `go mod tidy`

## Usage

Run the prod task with Taskfile:

```bash
task prod
```

Or directly:

```bash
go build -o bin/libre-bot ./cmd/cli/main.go && ./bin/libre-bot
```

It fetches up to 24 product cards and prints affiliate links.

## Dependencies

- google-chrome;
- Chrome's user data directory on the root of the project;
