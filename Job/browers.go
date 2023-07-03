package Job

import (
	"dns-check/config"
	"dns-check/logger"
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"time"
)

var b *rod.Browser

func OpenBrowser() *rod.Browser {
	if b == nil {
		if config.RunMode == "debug" {
			u := launcher.New().
				Headless(false).
				Proxy(fmt.Sprintf("https://%s:%s", "brd.superproxy.io", "22225")).
				Delete("use-mock-keychain").
				Bin("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome").
				MustLaunch()
			b = rod.New().ControlURL(u).MustConnect()
			b.MustIgnoreCertErrors(true)
			go b.MustHandleAuth(config.Conf.ProxyUserName, config.Conf.ProxyPassword)()
		} else {
			if path, exists := launcher.LookPath(); exists {
				u := launcher.New().Headless(true).Delete("use-mock-keychain").
					Set("no-sandbox").
					Delete("use-mock-keychain").
					Proxy(fmt.Sprintf("https://%s:%s", "brd.superproxy.io", "22225")).
					Bin(path).
					MustLaunch()
				b = rod.New().ControlURL(u).MustConnect()
				b.MustIgnoreCertErrors(true)
				go b.MustHandleAuth(config.Conf.ProxyUserName, config.Conf.ProxyPassword)()
			} else {
				panic("not found chrome")
			}
		}

	}
	return b
}

func OpenPage(browser *rod.Browser, url string) *rod.Page {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("OpenPage", logger.ERROR, nil, err.(error).Error())
			if b != nil {
				b.Close()
				b = nil
			}
		}
	}()
	page, cancel := browser.MustPage(url).WithCancel()
	go func(doCancel func(), p *rod.Page) {
		ticker := time.NewTicker(time.Second * 30)
		select {
		case <-ticker.C:
			if page != nil {
				doCancel()
				page.Close()
			}
		}
	}(cancel, page)
	page.MustWaitLoad()
	return page
}

func SendSearch(page *rod.Page, search string) {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("SendSearch", logger.ERROR, nil, err.(error).Error())
			if b != nil {
				page.Close()
			}
		}
	}()
	input := page.MustElementX("/html/body/div/div/div/form/input[1]")
	input.MustInput(search)
	sumbit := page.MustElementX("/html/body/div/div/div/form/input[2]")
	sumbit.MustClick()
	page.MustWaitNavigation()
}

func GetResult(page *rod.Page) string {
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("GetResult", logger.ERROR, nil, err.(error).Error())
			if b != nil {
				page.Close()
			}
		}
	}()
	result := page.MustElementX("/html/body/div/div/pre")
	return result.MustText()
}
