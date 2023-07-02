package Job

import (
	"dns-check/config"
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
			fmt.Println(err)
			if b != nil {
				b.Close()
				b = nil
			}
		}
	}()
	page := browser.Timeout(time.Second * 5).MustPage(url).CancelTimeout()
	page.Timeout(time.Second * 5).MustWaitLoad().CancelTimeout()
	return page
}

func SendSearch(page *rod.Page, search string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			if b != nil {
				page.Close()
			}
		}
	}()
	input := page.Timeout(time.Second * 5).MustElementX("/html/body/div/div/div/form/input[1]").CancelTimeout()
	input.Timeout(time.Second * 5).MustInput(search).CancelTimeout()
	sumbit := page.Timeout(time.Second * 5).MustElementX("/html/body/div/div/div/form/input[2]").CancelTimeout()
	sumbit.Timeout(time.Second * 5).MustClick().CancelTimeout()
	page.Timeout(time.Second).MustWaitNavigation()
}

func GetResult(page *rod.Page) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			if b != nil {
				page.Close()
			}
		}
	}()
	result := page.Timeout(time.Second * 5).MustElementX("/html/body/div/div/pre").CancelTimeout()
	return result.Timeout(time.Second * 5).MustText()
}
