# Go/Golang package for parsing user agent strings
[![GoDoc](https://pkg.go.dev/badge/github.com/houseme/go-mobile-detect.svg)](https://pkg.go.dev/github.com/houseme/go-mobile-detect)

Package `ua.Parse(userAgent string)` function parses browser's and bot's user agents strings and determins:
+ User agent name and version (Chrome, Firefox, Googlebot, etc.)
+ Operating system name and version  (Windows, Android, iOS etc.)
+ Device type (mobile, desktop, tablet, bot)
+ Device name if available (iPhone, iPad, Huawei VNS-L21)
+ URL provided by the bot (http://www.google.com/bot.html etc.)

## Status

Still need some work on detecting Andorid device names.

Fill free to report an issue for any User-Agent string not recognized or misinterpreted.

## Installation <a id="installation"></a>
```shell
go get -u -v github.com/houseme/mobiledetect
```

## Example<a id="example"></a>

```go
package main

import (
    "fmt"
    "strings"

    "github.com/houseme/mobiledetect/ua"
)

func main() {
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
    }

    for _, s := range userAgents {
        ua := ua.ParseUA(s)
        fmt.Println()
        fmt.Println(ua.UA())
        fmt.Println(strings.Repeat("=", len(ua.UA())))
        fmt.Println("Name:", ua.Name(), "v", ua.Version())
        fmt.Println("OS:", ua.OS(), "v", ua.OSVersion())
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
        fmt.Printf("%v\n", ua.Mobile())   // => true
        fmt.Printf("%v\n", ua.Bot())      // => false
        fmt.Printf("%v\n", ua.Mozilla())  // => "5.0"
        fmt.Printf("%v\n", ua.Model())    // => "Nexus One"
    
        fmt.Printf("%v\n", ua.Platform()) // => "Linux"
        fmt.Printf("%v\n", ua.OS())       // => "Android 2.3.7"
    
        name, version := ua.Engine()
        fmt.Printf("%v\n", name)          // => "AppleWebKit"
        fmt.Printf("%v\n", version)       // => "533.1"
    
        name, version = ua.Browser()
        fmt.Printf("%v\n", name)          // => "Android"
        fmt.Printf("%v\n", version)       // => "4.0"
    
        // Let's see an example with a bot.
    
        ua.Parse("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
    
        fmt.Printf("%v\n", ua.Bot())      // => true
    
        name, version = ua.Browser()
        fmt.Printf("%v\n", name)          // => Googlebot
        fmt.Printf("%v\n", version)       // => 2.1
    }
}


```

### Shorthand functions

Beside `UserAgent{}` struct and its properties returned by `ua.Parse()`, there is a bunch of shorthand functions for most popular browsers and operating systems, so this code:

```go
    ua := ua.Parse(userAgentString)
    if ua.OS == "Android" && ua.Name == "Chrome" {
        // do something
    }
```
can be also written on this way:
```go
    ua := ua.Parse(userAgentString)
    if ua.IsAndroid() && ua.IsChrome() {
        // do something
    }
```

### Notice

+ Opera and Opera Mini are two browsers, since they operate on very different ways.
+ If Googlebot (or any other bot) is detected and it is using its mobile crawler, both `bot` and `mobile` flags will be set to `true`.


### License

Go Mobile Detect is an open-source script released under [MIT License](http://www.opensource.org/licenses/mit-license.php). thanks for [mileusna](https://github.com/mileusna/useragent) and [mssola](https://github.com/mssola/user_agent).
