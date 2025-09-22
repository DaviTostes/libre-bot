package main

import (
	"fmt"
	"librebot/internal/scrap"
	"log"
)

func main() {
	cards, err := scrap.GetPolyCards()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total:", len(cards))

	for _, card := range cards {
		link, err := scrap.GenerateAffiliateLink(card.Url)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(card.Text, "->", link)
	}
}
