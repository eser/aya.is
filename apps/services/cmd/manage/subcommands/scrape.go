package subcommands

import (
	"context"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/spf13/cobra"
)

func CmdScrape() *cobra.Command {
	scrapeCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "scrape",
		Short: "Scrapes a website",
		Long:  "Scrapes a website",
		RunE: func(cmd *cobra.Command, args []string) error {
			return execScrape(cmd.Context())
		},
	}

	return scrapeCmd
}

func execScrape(_ context.Context) error {
	tweetID := "1909221461195207130"
	url := "https://twitter.com/i/web/status/" + tweetID

	browser := rod.New().ControlURL(launcher.New().Headless(true).MustLaunch()).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(url)
	page.MustWaitLoad()
	time.Sleep(4 * time.Second) //nolint:mnd

	_, err := page.Elements("article")
	if err != nil {
		panic(err)
	}

	// if len(articles) > 1 {
	// 	fmt.Printf("➡️ Bu tweet bir cevap (reply). Conversation içinde.\n")
	// } else if len(articles) == 1 {
	// 	fmt.Printf("✅ Bu tweet bağımsız bir tweet.\n")
	// } else {
	// 	fmt.Println("⚠️ Tweet bulunamadı veya DOM yapısı değişmiş olabilir.")
	// }

	return nil
}
