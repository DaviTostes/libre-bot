package main

import (
	"fmt"
	"librebot/internal/scrap"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	// "librebot/internal/whatsapp"
	"log"
)

var (
	isLoading = false
	mu        = sync.Mutex{}
)

func startLoading() {
	isLoading = true

	go func() {
		count := 0
		for {
			if isLoading {
				fmt.Print("*")
				time.Sleep(200 * time.Millisecond)
			}

			count++
			if count == 10 {
				fmt.Print("\r")
				fmt.Print(strings.Repeat(" ", 80))
				fmt.Print("\r")

				count = 0
			}
		}
	}()
}

func stopLoading() {
	isLoading = false
	fmt.Print("\r")
	fmt.Print(strings.Repeat(" ", 80))
	fmt.Print("\r")
	fmt.Println("")
}

func main() {
	// if err := whatsapp.ConnectToWhatsApp(); err != nil {
	// 	log.Fatalln(err)
	// }

	postInterval := os.Getenv("POST_INTERVAL")
	postIntervalInt, err := strconv.Atoi(postInterval)
	if err != nil {
		log.Fatalln(err)
	}

	need := (60 / postIntervalInt) * 24
	cards := []scrap.PolyCard{}

	fmt.Println("Need:", need)

	startLoading()

	for {
		newCards, err := scrap.GetPolyCards()
		if err != nil {
			log.Fatalln(err)
		}

		cards = append(cards, newCards...)

		if len(cards) >= need {
			break
		}
	}

	stopLoading()

	fmt.Println("Cards scrapped:", len(cards))
	for i := range cards {
		// link, err := scrap.GenerateAffiliateLink(cards[i].Url)
		// if err != nil {
		// 	continue
		// }

		fmt.Println(cards[i].Text /* , "->", link */)
	}
}
