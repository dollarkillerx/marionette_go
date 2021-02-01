// Command remote is a chromedp example demonstrating how to connect to an
// existing Chrome DevTools instance using a remote WebSocket URL.
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

func main() {
	// create allocator context for use with creating a browser context later
	allocatorContext, cancel := chromedp.NewRemoteAllocator(context.Background(), "ws://192.168.88.11:9222/devtools/page/54F0BB9D5FA2DC9CB1A938509597F79B")
	defer cancel()

	options := []chromedp.ExecAllocatorOption{
		//chromedp.Flag("headless", false), // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, cancel := chromedp.NewExecAllocator(allocatorContext, options...)
	defer cancel()

	// create context
	ctxt, cancel := chromedp.NewContext(c)
	defer cancel()

	fmt.Println("hello world")
	url := "http://192.168.88.11:5920/#/user/login"
	fmt.Println(url)

	// run task list
	var body string
	if err := chromedp.Run(ctxt,
		chromedp.Navigate(url),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		log.Fatalf("Failed getting body of %s: %v", url, err)
	}

	log.Printf("Body of %s starts with: \n", url)
	ioutil.WriteFile("a.html", []byte(body), 00666)
}
