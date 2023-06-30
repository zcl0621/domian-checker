package Job

import (
	"fmt"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var b *rod.Browser

func OpenBrowser() *rod.Browser {
	if b == nil {
		if path, exists := launcher.LookPath(); exists {
			u := launcher.New().Headless(false).Delete("use-mock-keychain").
				Set("no-sandbox").
				Set("user-data-dir", "/tmp/chrome").
				Set("proxy", fmt.Sprintf("https://%s:%s@%s:%s", "brd-customer-hl_3cf009f7-zone-data_center", "wqt22u1s0uyg", "brd.superproxy.io", "22225")).
				Bin(path).
				MustLaunch()
			b = rod.New().ControlURL(u).MustConnect()
		} else {
			panic("not found chrome")
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
	page := browser.MustPage(url)
	page.MustWaitLoad()
	return page
}

func SendSearch(page *rod.Page, search string) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			if b != nil {
				b.Close()
				b = nil
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
			fmt.Println(err)
			if b != nil {
				b.Close()
				b = nil
			}
		}
	}()
	result := page.MustElementX("/html/body/div/div/pre")
	return result.MustText()
}
