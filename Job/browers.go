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
				Delete("use-mock-keychain").
				Set("proxy", fmt.Sprintf("https://%s:%s@%s:%s", "brd-customer-hl_3cf009f7-zone-data_center", "wqt22u1s0uyg", "brd.superproxy.io", "22225")).
				Bin("/Applications/Google Chrome.app/Contents/MacOS/Google Chrome").
				MustLaunch()
			b = rod.New().ControlURL(u).MustConnect()
		} else {
			if path, exists := launcher.LookPath(); exists {
				u := launcher.New().Headless(true).Delete("use-mock-keychain").
					Set("no-sandbox").
					Set("proxy", fmt.Sprintf("https://%s:%s@%s:%s", "brd-customer-hl_3cf009f7-zone-data_center", "wqt22u1s0uyg", "brd.superproxy.io", "22225")).
					Bin(path).
					MustLaunch()
				b = rod.New().ControlURL(u).MustConnect()
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
	go func(doCancel func()) {
		ticker := time.NewTicker(time.Second * 20)
		select {
		case <-ticker.C:
			doCancel()
		}
	}(cancel)
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
