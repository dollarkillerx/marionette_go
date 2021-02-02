package tests

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestChromeDP(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false), // debug使用
		// 禁用GPU，不显示GUI
		chromedp.DisableGPU,
		// 取消沙盒模式
		chromedp.NoSandbox,
		// 隐身模式启动
		chromedp.Flag("incognito", true),
		// 忽略证书错误
		chromedp.Flag("ignore-certificate-errors", true),
		// 窗口最大化
		chromedp.Flag("start-maximized", true),
		// 不加载图片, 提升速度
		chromedp.Flag("disable-images", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		// 禁用扩展
		chromedp.Flag("disable-extensions", true),
		// 禁止加载所有插件
		chromedp.Flag("disable-plugins", true),
		// 禁用浏览器应用
		chromedp.Flag("disable-software-rasterizer", true),
		//chromedp.Flag("user-data-dir", "./.cache"),
		// 设置UA，防止有些页面识别headless模式
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}

	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	c, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	// create context
	ctxt, cancel := chromedp.NewContext(c)
	defer cancel()

	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(ctxt, make([]chromedp.Action, 0, 1)...)

	// 给每个页面的爬取设置超时时间
	timeoutCtx, cancel := context.WithTimeout(ctxt, 20*time.Second)
	defer cancel()

	// run task list
	var body string
	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate("http://192.168.88.11:5920/#/user/login"),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		log.Fatalf("Failed getting body of %s: %v", "http://192.168.88.11:5920/#/user/login", err)
	}

	ioutil.WriteFile("a.html", []byte(body), 00666)

	list := []string{"http://www.baidu.com", "http://www.360.com", "http://hao123.com"}

	var tasks []chromedp.Action
	var taskResp []*string
	for _, v := range list {
		var body string

		tasks = append(tasks, chromedp.Tasks{
			chromedp.Navigate(v),
			chromedp.OuterHTML("html", &body),
		})
		taskResp = append(taskResp, &body)
	}
	if err := chromedp.Run(timeoutCtx, tasks...); err != nil {
		log.Fatalln(err)
	}
	for _, v := range taskResp {
		fmt.Println(*v)
	}

	//for i := range list {
	//	idx := i
	//	go func() {
	//		ctxt, cancel := chromedp.NewContext(c)
	//		defer cancel()
	//
	//		timeoutCtx, cancel := context.WithTimeout(ctxt, 20*time.Second)
	//		defer cancel()
	//
	//		// run task list
	//		var body string
	//		if err := chromedp.Run(timeoutCtx,
	//			chromedp.Tasks{
	//				chromedp.Navigate(list[idx]),
	//				chromedp.OuterHTML("html", &body),
	//			},
	//		); err != nil {
	//			log.Fatalf("Failed getting body of %s: %v", list[idx], err)
	//		}
	//
	//		fmt.Println(list[idx])
	//		ioutil.WriteFile(fmt.Sprintf("%d.html", idx), []byte(body), 00666)
	//	}()
	//}

	for {
		select {}
	}
}
