package ua

import (
	"regexp"
)

// Constants for browsers and operating systems for easier comparison
const (
	Windows      = "Windows"
	WindowsPhone = "Windows Phone"
	Android      = "Android"
	MacOS        = "macOS"
	IOS          = "iOS"
	Linux        = "Linux"
	ChromeOS     = "CrOS"

	Opera            = "Opera"
	OperaMini        = "Opera Mini"
	OperaTouch       = "Opera Touch"
	Chrome           = "Chrome"
	Firefox          = "Firefox"
	InternetExplorer = "Internet Explorer"
	Safari           = "Safari"
	Edge             = "Edge"
	Vivaldi          = "Vivaldi"

	Googlebot           = "Googlebot"
	Twitterbot          = "Twitterbot"
	FacebookExternalHit = "facebookexternalhit"
	Applebot            = "Applebot"
)

var (
	rxMacOSVer        = regexp.MustCompile(`[_\d.]+`)
	botFromSiteRegexp = regexp.MustCompile(`http[s]?://.+\.\w+`)
	botRegex          = regexp.MustCompile("(?i)(bot|crawler|sp([iy])der|search|worm|fetch|nutch)")
	ie11Regexp        = regexp.MustCompile("^rv:(.+)$")
)

var ignore = map[string]struct{}{
	"KHTML, like Gecko": {},
	"U":                 {},
	"compatible":        {},
	"Mozilla":           {},
	"WOW64":             {},
}

// A section contains the name of the product, its version and
// an optional comment.
type section struct {
	name    string
	version string
	comment []string
}

// UserAgent struct containing all data extracted from parsed user-agent string
type UserAgent struct {
	name         string
	version      string
	os           string
	osVersion    string
	shortOS      string
	device       string
	mobile       bool
	tablet       bool
	desktop      bool
	bot          bool
	url          string
	ua           string
	platform     string
	browser      Browser
	mozilla      string
	localization string
	undecided    bool
}

// Browser The browser is a struct containing all the information that we might be
// interested in the browser.
type Browser struct {
	// The name of the browser's engine.
	Engine string

	// The version of the browser's engine.
	EngineVersion string

	// The name of the browser.
	Name string

	// The version of the browser.
	Version string
}

// OSInfo represents full information on the operating system extracted from the user agent.
type OSInfo struct {
	// Full name of the operating system. This is identical to the output of ua.OS()
	FullName string

	// Name of the operating system. This is sometimes a shorter version of the
	// operating system name, e.g. "Mac OS X" instead of "Intel Mac OS X"
	Name string

	// Operating system version, e.g. 7 for Windows 7 or 10.8 for Max OS X Mountain Lion
	Version string
}
