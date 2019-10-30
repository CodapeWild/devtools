package main

import (
	"context"
	"log"
	"os"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := chromedp.NewPool(chromedp.PortRange(10000, 20000))
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer pool.Shutdown()

	chromepool, err := pool.Allocate(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer chromepool.Release()

	nodes := []*cdp.Node{}
	err = chromepool.Run(ctx, chromedp.Tasks{
		// chromedp.Navigate("https://www.106888.com/betcenter"),
		// chromedp.WaitVisible(".menuItem___1Gogq", chromedp.ByQueryAll),
		// chromedp.Click("[value='HF_GDD11']", chromedp.ByQuery),
		// chromedp.WaitVisible(".countDownContainer___3rzLL", chromedp.ByQuery),
		// chromedp.Nodes("span.CircularProgressbarText___7oL13", &nodes, chromedp.ByQueryAll),

		// chromedp.Navigate("https://www.106888.com/results?gameUniqueId=HF_CQSSC"),
		// chromedp.WaitVisible(".singleResult_lotteryBall___JI3YE", chromedp.ByQueryAll),
		// chromedp.Nodes(".singleResult_lotteryBall___JI3YE", &nodes, chromedp.ByQueryAll),

		chromedp.Navigate("https://81cp.tw/lotteryV3/lotDetail.do?lotCode=HBK3"),
		chromedp.WaitVisible("ul.flip li.flip-clock-active div.down div.inn", chromedp.ByQueryAll),
		chromedp.Nodes("ul.flip li.flip-clock-active div.down div.inn,#current_issue", &nodes, chromedp.ByQueryAll),

		// chromedp.Navigate("https://81cp.tw/lottery/trendChart/index.do?lotCode=JSSB3"),
		// chromedp.WaitVisible(`"div:eq(1), div:eq(2)", "div.cl-30.clean"`, chromedp.ByQueryAll),
		// chromedp.Nodes(`"div:eq(1), div:eq(2)", "div.cl-30.clean"`, &nodes, chromedp.ByQueryAll),
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, v := range nodes {
		log.Println(v.Children[0].NodeValue)
	}
}

// func main() {
// 	log.SetOutput(os.Stdout)
// 	log.SetFlags(log.LstdFlags | log.Lshortfile)

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
// 	defer cancel()

// 	chrome, err := chromedp.New(ctx, chromedp.WithRunnerOptions(runner.Path("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"), runner.Flag("headless", true), runner.Flag("no-sandbox", true)))
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	nodes := []*cdp.Node{}
// 	err = chrome.Run(ctx, chromedp.Tasks{
// 		chromedp.Navigate("https://81cp.tw/lotteryV3/lotDetail.do?lotCode=CQSSC"),
// 		chromedp.WaitVisible("ul.flip li.flip-clock-active div.down div.inn, #current_issue", chromedp.ByQueryAll),
// 		chromedp.Nodes("ul.flip li.flip-clock-active div.down div.inn, #current_issue", &nodes, chromedp.ByQueryAll),
// 		// chromedp.Navigate("https://81cp.tw/lottery/trendChart/index.do?lotCode=GD11X5"),
// 		// chromedp.WaitVisible("div.number,div.openCode", chromedp.ByQueryAll),
// 		// chromedp.Nodes("div.number,div.openCode", &nodes, chromedp.ByQueryAll),
// 	})
// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	// if err = chrome.Shutdown(ctx); err != nil {
// 	// 	log.Fatal(err.Error())
// 	// }
// 	// if err = chrome.Wait(); err != nil {
// 	// 	log.Fatal(err.Error())
// 	// }

// 	for _, v := range nodes {
// 		log.Println(v.Children[0].NodeValue)
// 	}
// }
