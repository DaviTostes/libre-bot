package scrap

import (
	"context"
	"log"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func runActions(actions []chromedp.Action) error {
	userDataDir := "chrome-profile"

	opts := append(
		chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),

		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),

		chromedp.UserAgent(
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36",
		),

		chromedp.UserDataDir(userDataDir),
		chromedp.Flag("profile-directory", "Default"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := context.WithTimeout(allocCtx, 60*time.Second)

	taskCtx, cancel := chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cancel()

	return chromedp.Run(taskCtx, actions...)
}

func GenerateAffiliateLink(url string) (string, error) {
	urlTextAreaSelector := "#url-0"
	generateButtonSelector := ".button_generate-links"
	linkTextSelector := "#textfield-copyLink-1"
	var link string

	actions := []chromedp.Action{
		chromedp.Navigate("https://www.mercadolivre.com.br/afiliados/linkbuilder#hub"),

		chromedp.WaitVisible(urlTextAreaSelector, chromedp.ByID),
		chromedp.ScrollIntoView(urlTextAreaSelector, chromedp.ByID),

		chromedp.Click(urlTextAreaSelector, chromedp.ByID),
	}

	for _, char := range url {
		actions = append(
			actions,
			chromedp.SendKeys(urlTextAreaSelector, string(char), chromedp.ByID),
		)
		time.Sleep(1 * time.Microsecond)
	}

	actions = append(actions,
		chromedp.Click(generateButtonSelector, chromedp.ByQuery),
		chromedp.Value(linkTextSelector, &link, chromedp.ByID),
	)

	if err := runActions(actions); err != nil {
		return "", err
	}

	return link, nil
}

type PolyCard struct {
	Url  string
	Text string
}

func GetPolyCards() ([]PolyCard, error) {
	tasks := chromedp.Tasks{
		chromedp.Navigate("https://www.mercadolivre.com.br/afiliados/hub"),
		chromedp.WaitVisible("body", chromedp.ByQuery),
		chromedp.WaitVisible("#recommendations_column", chromedp.ByQuery),
	}

	polyCards := []PolyCard{}
	polyCardSelector := ".poly-card"

	tasks = append(tasks, chromedp.Tasks{
		chromedp.ActionFunc(func(ctx context.Context) error {
			var cardNodes []*cdp.Node
			if err := chromedp.Nodes(polyCardSelector, &cardNodes, chromedp.ByQueryAll).Do(ctx); err != nil {
				return err
			}

			for _, node := range cardNodes {
				var href string
				var text string
				var ok bool

				err := chromedp.AttributeValue("a", "href", &href, &ok, chromedp.ByQuery, chromedp.FromNode(node)).
					Do(ctx)
				if err != nil {
					continue
				}

				err = chromedp.Text("a", &text, chromedp.ByQuery, chromedp.FromNode(node)).
					Do(ctx)
				if err != nil {
					continue
				}

				if ok {
					polyCards = append(polyCards, PolyCard{Url: href, Text: text})
				}
			}

			return nil
		})})

	if err := runActions(tasks); err != nil {
		return nil, err
	}

	return polyCards, nil
}
