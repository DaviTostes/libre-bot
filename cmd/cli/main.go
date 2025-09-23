package main

import (
	"bufio"
	"fmt"
	"librebot/internal/scrap"
	"librebot/internal/whatsapp"
	"log"
	"os"
	"strings"
)

type ActionFunc func(string)

var actionMap = map[string]ActionFunc{}

func populateOptions() {
	actionMap["connect"] = func(s string) {
		if s == "whatsapp" {
			if err := whatsapp.ConnectToWhatsApp(); err != nil {
				log.Fatalln(err)
			}
		} else {
			fmt.Println("options for connect: whatsapp")
		}
	}

	actionMap["get"] = func(s string) {
		if s == "cards" {
			cards, err := scrap.GetPolyCards()
			if err != nil {
				log.Fatalln(err)
			}

			for _, card := range cards {
				link, err := scrap.GenerateAffiliateLink(card.Url)
				if err != nil {
					log.Fatalln(err)
				}

				fmt.Println(card.Text, "->", link)
			}
		} else {
			fmt.Println("options for get: cards")
		}
	}

	actionMap["create"] = func(s string) {
		if s == "link" {
			fmt.Print("url: ")
			scanner := bufio.NewScanner(os.Stdin)
			if !scanner.Scan() {
				return
			}

			url := scanner.Text()
			if url == "" {
				fmt.Println("invalid url")
				return
			}

			link, err := scrap.GenerateAffiliateLink(url)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println("link:", link)
		} else {
			fmt.Println("options for create: link")
		}
	}
}

func main() {
	populateOptions()

	fmt.Println("welcome to LibreBot")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("#> ")

		if !scanner.Scan() {
			break
		}

		command := scanner.Text()
		if command == "help" {
			file, _ := os.ReadFile("help.txt")
			fmt.Println(string(file))

			continue
		}

		splitedCmd := strings.Split(command, " ")
		if len(splitedCmd) < 2 {
			fmt.Println("usage: <action> <option>")
			continue
		}

		option := actionMap[splitedCmd[0]]
		if option == nil {
			fmt.Println("invalid command")
			continue
		}

		option(splitedCmd[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("error reading from stdin:", err.Error())
	}
}
