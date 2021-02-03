package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/dollarkillerx/marionette_go/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Server struct {
	chromeContext context.Context
}

func RunServers() {
	server := Server{}
	server.initChrome()
	server.run()
}

func (s *Server) run() {
	app := fiber.New()

	app.Get("/ssr", s.ssr)
	app.Post("/avaricious", s.avaricious)

	fmt.Println("Marionette Go Run ...")
	if err := app.Listen(config.CONF.ListenAddr); err != nil {
		log.Fatalln(err)
	}
}

func (s *Server) ssr(ctx *fiber.Ctx) error {
	url := ctx.OriginalURL()
	index := strings.Index(url, "ssr?q=")
	if index == -1 {
		ctx.Status(400)
		return errors.New("400")
	}

	url = url[index+6:]
	ctxt, cancelFunc := chromedp.NewContext(s.chromeContext)
	defer cancelFunc()
	timeoutCtx, cancel := context.WithTimeout(ctxt, 20*time.Second)
	defer cancel()
	var body string
	cookis := make([]*fiber.Cookie, 0)
	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(url),
		ShowCookies(&cookis),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		ctx.Status(400)
		return fmt.Errorf("Failed getting body of %s: %v \n", url, err)
	}

	for _, v := range cookis {
		ctx.Cookie(v)
	}
	ctx.Set("Content-Type", "text/html;charset=utf-8")
	ctx.WriteString(body)
	return nil
}

type cookie struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	HttpOnly bool   `json:"http_only"`
	Secure   bool   `json:"secure"`
}

type avaricious struct {
	Url     string   `json:"url"`
	Cookie  []cookie `json:"cookie"`
	Timeout int      `json:"timeout"` // sec
}

func (s *Server) avaricious(ctx *fiber.Ctx) error {
	ava := avaricious{}
	if err := ctx.BodyParser(&ava); err != nil {
		ctx.Status(400)
		return err
	}

	ctxt, cancelFunc := chromedp.NewContext(s.chromeContext)
	defer cancelFunc()
	timeoutCtx, cancel := context.WithTimeout(ctxt, time.Duration(ava.Timeout)*time.Second)
	defer cancel()
	var body string
	cookis := make([]*fiber.Cookie, 0)
	if err := chromedp.Run(timeoutCtx,
		chromedp.Navigate(ava.Url),
		SetCookie(ava.Cookie),
		ShowCookies(&cookis),
		chromedp.OuterHTML("html", &body),
	); err != nil {
		ctx.Status(400)
		return fmt.Errorf("Failed getting body of %s: %v \n", ava.Url, err)
	}

	for _, v := range cookis {
		ctx.Cookie(v)
	}
	ctx.Set("Content-Type", "text/html;charset=utf-8")
	ctx.JSON(map[string]interface{}{
		"cookies": cookis,
		"html":    body,
	})
	return nil
}

func SetCookie(cookies []cookie) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if len(cookies) == 0 {
			return nil
		}

		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		for _, c := range cookies {
			err := network.SetCookie(c.Name, c.Value).
				WithExpires(&expr).
				WithDomain(c.Domain).
				WithPath(c.Path).
				WithHTTPOnly(c.HttpOnly).
				WithSecure(c.Secure).
				Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func ShowCookies(cookies *[]*fiber.Cookie) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		cks, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return err
		}
		for _, cookie := range cks {
			*cookies = append(*cookies, &fiber.Cookie{
				Name:     cookie.Name,
				Value:    cookie.Value,
				Path:     cookie.Path,
				Domain:   cookie.Domain,
				Secure:   cookie.Secure,
				HTTPOnly: cookie.HTTPOnly,
				SameSite: cookie.SameSite.String(),
			})
		}
		return nil
	})
}

func (s *Server) initChrome() {
	options := []chromedp.ExecAllocatorOption{
		//chromedp.Flag("headless", false), // debug使用
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

	c, _ := chromedp.NewExecAllocator(context.Background(), options...)

	// create context
	ctxt, _ := chromedp.NewContext(c)

	// 执行一个空task, 用提前创建Chrome实例
	chromedp.Run(ctxt, make([]chromedp.Action, 0, 1)...)
	s.chromeContext = c
}
