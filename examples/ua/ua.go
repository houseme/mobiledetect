package main

import (
	"fmt"
	"strings"

	"github.com/houseme/mobiledetect/ua"
)

func main() {
	// The "New" function will create a new UserAgent object and it will parse
	// the given string. If you need to parse more strings, you can re-use
	// this object and call: ua.Parse("another string")
	// ua := ua.New("Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1")
	ua := ua.New("Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 mobile/14F89 Safari/602.1")
	fmt.Printf("UserAgent: %s\n", ua.UA())
	fmt.Printf("%v\n", ua.Mobile())           // => true
	fmt.Printf("%v\n", ua.Bot())              // => false
	fmt.Printf("Mozilla: %v\n", ua.Mozilla()) // => "5.0"
	fmt.Printf("Device: %v\n", ua.Device())   // => "Nexus One"

	fmt.Printf("Platform: %v\n", ua.Platform()) // => "Linux"
	fmt.Printf("OS: %v\n", ua.OS())             // => "Android 2.3.7"

	fmt.Printf("Name: %v\n", ua.Name())       // => "Android 2.3.7"
	fmt.Printf("Version: %v\n", ua.Version()) // => "Android 2.3.7"

	fmt.Printf("ShortOS: %v\n", ua.ShortOS())     // => "Android 2.3.7"
	fmt.Printf("OSVersion: %v\n", ua.OSVersion()) // => "Android 2.3.7"

	name, version := ua.Engine()
	fmt.Printf("Engine name: %v\n", name)       // => "AppleWebKit"
	fmt.Printf("Engine version: %v\n", version) // => "533.1"

	name, version = ua.Browser()
	fmt.Printf("Browser name: %v\n", name)       // => "Android"
	fmt.Printf("Browser version: %v\n", version) // => "4.0"

	fmt.Printf("user-agent string: %v\n", ua.Beautify()) // => "AppleWebKit 533.1"

	fmt.Printf("========================\n")

	// Let's see an example with a bot.

	ua.Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +https://www.google.com/bot.html)")

	fmt.Printf("%v\n", ua.Bot()) // => true

	name, version = ua.Browser()
	fmt.Printf("Browser name: %v\n", name)       // => Googlebot
	fmt.Printf("Browser version: %v\n", version) // => 2.1
	uaEx()
}

func uaEx() {
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/603.3.8 (KHTML, like Gecko) Version/10.1.2 Safari/603.3.8",
		"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) FxiOS/8.1.1b4948 Mobile/14F89 Safari/603.2.4",
		"Mozilla/5.0 (iPad; CPU OS 10_3_2 like Mac OS X) AppleWebKit/603.2.4 (KHTML, like Gecko) Version/10.0 Mobile/14F89 Safari/602.1",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.125 Mobile Safari/537.36",
		"Mozilla/5.0 (Android 4.3; Mobile; rv:54.0) Gecko/54.0 Firefox/54.0",
		"Mozilla/5.0 (Linux; Android 4.3; GT-I9300 Build/JSS15J) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/55.0.2883.91 Mobile Safari/537.36 OPR/42.9.2246.119956",
		"Opera/9.80 (Android; Opera Mini/28.0.2254/66.318; U; en) Presto/2.12.423 Version/12.16",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36 Google (+https://developers.google.com/+/web/snippet/)",
	}

	for _, s := range userAgents {
		ua := ua.New(s)
		fmt.Println()
		fmt.Println(ua.UA())
		fmt.Println(strings.Repeat("=", len(ua.UA())))
		fmt.Println("Name:", ua.Name(), "v: ", ua.Version())
		fmt.Println("OS:", ua.OS(), "v: ", ua.OSVersion())
		fmt.Println("Device:", ua.Device())
		if ua.Mobile() {
			fmt.Println("(Mobile)")
		}
		if ua.Tablet() {
			fmt.Println("(Tablet)")
		}
		if ua.Desktop() {
			fmt.Println("(Desktop)")
		}
		if ua.Bot() {
			fmt.Println("(Bot)")
		}
		if ua.URL() != "" {
			fmt.Println(ua.URL())
		}

		fmt.Printf("Mobile: %v\n", ua.Mobile())   // => true
		fmt.Printf("Bot: %v\n", ua.Bot())         // => false
		fmt.Printf("Mozilla: %v\n", ua.Mozilla()) // => "5.0"
		fmt.Printf("Device: %v\n", ua.Device())   // => "Nexus One"

		fmt.Printf("Platform: %v\n", ua.Platform()) // => "Linux"
		fmt.Printf("OS: %v\n", ua.OS())             // => "Android 2.3.7"

		name, version := ua.Engine()
		fmt.Printf("Engine.Name: %v\n", name)       // => "AppleWebKit"
		fmt.Printf("Engine.Version: %v\n", version) // => "533.1"

		name, version = ua.Browser()
		fmt.Printf("Browser.Name: %v\n", name)               // => "Android"
		fmt.Printf("Browser.version: %v\n", version)         // => "4.0"
		fmt.Printf("user-agent string: %v\n", ua.Beautify()) // => "AppleWebKit 533.1"
		fmt.Printf("========================\n")
		// Let's see an example with a bot.

		ua.Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +https://www.google.com/bot.html)")

		fmt.Printf("Bot: %v\n", ua.Bot()) // => true

		name, version = ua.Browser()
		fmt.Printf("Browser.Name: %v\n", name)       // => Googlebot
		fmt.Printf("Browser.version: %v\n", version) // => 2.1
	}
}
