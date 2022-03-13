package ua

import (
	"regexp"
	"strings"
)

var ie11Regexp = regexp.MustCompile("^rv:(.+)$")

// Browser is a struct containing all the information that we might be
// interested from the browser.
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

// Extract all the information that we can get from the User-Agent string
// about the browser and update the receiver with this information.
//
// The function receives just one argument "sections", that contains the
// sections from the User-Agent string after being parsed.
func (ua *UserAgent) detectBrowser(sections []section) {
	slen := len(sections)

	if sections[0].name == "Opera" {
		ua.browser.Name = "Opera"
		ua.browser.Version = sections[0].version
		ua.browser.Engine = "Presto"
		if slen > 1 {
			ua.browser.EngineVersion = sections[1].version
		}
	} else if sections[0].name == "Dalvik" {
		// When Dalvik VM is in use, there is no browser info attached to ua.
		// Although browser is still a Mozilla/5.0 compatible.
		ua.mozilla = "5.0"
	} else if slen > 1 {
		engine := sections[1]
		ua.browser.Engine = engine.name
		ua.browser.EngineVersion = engine.version
		if slen > 2 {
			sectionIndex := 2
			// The version after the engine comment is empty on e.g. Ubuntu
			// platforms so if this is the case, let's use the next in line.
			if sections[2].version == "" && slen > 3 {
				sectionIndex = 3
			}
			ua.browser.Version = sections[sectionIndex].version
			if engine.name == "AppleWebKit" {
				for _, comment := range engine.comment {
					if len(comment) > 5 &&
						(strings.HasPrefix(comment, "Googlebot") || strings.HasPrefix(comment, "bingbot")) {
						ua.undecided = true
						break
					}
				}
				switch sections[slen-1].name {
				case "Edge":
					ua.browser.Name = "Edge"
					ua.browser.Version = sections[slen-1].version
					ua.browser.Engine = "EdgeHTML"
					ua.browser.EngineVersion = ""
				case "Edg":
					if !ua.undecided {
						ua.browser.Name = "Edge"
						ua.browser.Version = sections[slen-1].version
						ua.browser.Engine = "AppleWebKit"
						ua.browser.EngineVersion = sections[slen-2].version
					}
				case "OPR":
					ua.browser.Name = "Opera"
					ua.browser.Version = sections[slen-1].version
				case "mobile":
					ua.browser.Name = "Mobile App"
					ua.browser.Version = ""
				default:
					switch sections[slen-3].name {
					case "YaBrowser":
						ua.browser.Name = "YaBrowser"
						ua.browser.Version = sections[slen-3].version
					case "coc_coc_browser":
						ua.browser.Name = "Coc Coc"
						ua.browser.Version = sections[slen-3].version
					default:
						switch sections[slen-2].name {
						case "Electron":
							ua.browser.Name = "Electron"
							ua.browser.Version = sections[slen-2].version
						case "DuckDuckGo":
							ua.browser.Name = "DuckDuckGo"
							ua.browser.Version = sections[slen-2].version
						default:
							switch sections[sectionIndex].name {
							case "Chrome", "CriOS":
								ua.browser.Name = "Chrome"
							case "HeadlessChrome":
								ua.browser.Name = "Headless Chrome"
							case "Chromium":
								ua.browser.Name = "Chromium"
							case "GSA":
								ua.browser.Name = "Google App"
							case "FxiOS":
								ua.browser.Name = "Firefox"
							default:
								ua.browser.Name = "Safari"
							}
						}
					}
					// It's possible the google-bot emulates these now
					for _, comment := range engine.comment {
						if len(comment) > 5 &&
							(strings.HasPrefix(comment, "Googlebot") || strings.HasPrefix(comment, "bingbot")) {
							ua.undecided = true
							break
						}
					}
				}
			} else if engine.name == "Gecko" {
				name := sections[2].name
				if name == "MRA" && slen > 4 {
					name = sections[4].name
					ua.browser.Version = sections[4].version
				}
				ua.browser.Name = name
			} else if engine.name == "like" && sections[2].name == "Gecko" {
				// This is the new user agent from Internet Explorer 11.
				ua.browser.Engine = "Trident"
				ua.browser.Name = "Internet Explorer"
				for _, c := range sections[0].comment {
					version := ie11Regexp.FindStringSubmatch(c)
					if len(version) > 0 {
						ua.browser.Version = version[1]
						return
					}
				}
				ua.browser.Version = ""
			}
		}
	} else if slen == 1 && len(sections[0].comment) > 1 {
		comment := sections[0].comment
		if comment[0] == "compatible" && strings.HasPrefix(comment[1], "MSIE") {
			ua.browser.Engine = "Trident"
			ua.browser.Name = "Internet Explorer"
			// The MSIE version may be reported as the compatibility version.
			// For IE 8 through 10, the Trident token is more accurate.
			// http://msdn.microsoft.com/en-us/library/ie/ms537503(v=vs.85).aspx#VerToken
			for _, v := range comment {
				if strings.HasPrefix(v, "Trident/") {
					switch v[8:] {
					case "4.0":
						ua.browser.Version = "8.0"
					case "5.0":
						ua.browser.Version = "9.0"
					case "6.0":
						ua.browser.Version = "10.0"
					}
					break
				}
			}
			// If the Trident token is not provided, fall back to MSIE token.
			if ua.browser.Version == "" {
				ua.browser.Version = strings.TrimSpace(comment[1][4:])
			}
		}
	}
}

// Engine returns two strings.
// The first string is the name of the engine and the
// second one is the version of the engine.
func (ua *UserAgent) Engine() (string, string) {
	return ua.browser.Engine, ua.browser.EngineVersion
}

// Browser returns two strings. The first string is the name of the browser and the
// second one is the version of the browser.
func (ua *UserAgent) Browser() (string, string) {
	return ua.browser.Name, ua.browser.Version
}
