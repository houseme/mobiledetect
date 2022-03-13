package ua

import (
	"bytes"
	"fmt"
	"strings"
)

// New parses the given User-Agent string and get the resulting UserAgent object.
//
// Returns a UserAgent object that has been initialized after parsing
// the given User-Agent string.
func New(uas string) *UserAgent {
	ua := &UserAgent{
		ua: uas,
	}
	ua.initialize()
	tokens := parse(uas)

	// check is there URL
	for k := range tokens {
		if strings.HasPrefix(k, "http://") || strings.HasPrefix(k, "https://") {
			ua.url = k
			delete(tokens, k)
			break
		}
	}

	// OS lookup
	switch {
	case tokens.exists("Android"):
		ua.shortOS = Android
		ua.osVersion = tokens[Android]
		for s := range tokens {
			if strings.HasSuffix(s, "Build") {
				ua.device = strings.TrimSpace(s[:len(s)-5])
				ua.tablet = strings.Contains(strings.ToLower(ua.device), "tablet")
			}
		}

	case tokens.exists("iPhone"):
		ua.shortOS = IOS
		ua.osVersion = tokens.findMacOSVersion()
		ua.device = "iPhone"
		ua.mobile = true

	case tokens.exists("iPad"):
		ua.shortOS = IOS
		ua.osVersion = tokens.findMacOSVersion()
		ua.device = "iPad"
		ua.tablet = true

	case tokens.exists("Windows NT"):
		ua.shortOS = Windows
		ua.osVersion = tokens["Windows NT"]
		ua.desktop = true

	case tokens.exists("Windows Phone OS"):
		ua.shortOS = WindowsPhone
		ua.osVersion = tokens["Windows Phone OS"]
		ua.mobile = true

	case tokens.exists("Macintosh"):
		ua.shortOS = MacOS
		ua.osVersion = tokens.findMacOSVersion()
		ua.desktop = true

	case tokens.exists("Linux"):
		ua.shortOS = Linux
		ua.osVersion = tokens[Linux]
		ua.desktop = true

	case tokens.exists("CrOS"):
		ua.shortOS = ChromeOS
		ua.osVersion = tokens[ChromeOS]
		ua.desktop = true

	}

	switch {

	case tokens.exists("Googlebot"):
		ua.name = Googlebot
		ua.version = tokens[Googlebot]
		ua.bot = true
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens.exists("Applebot"):
		ua.name = Applebot
		ua.version = tokens[Applebot]
		ua.bot = true
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")
		ua.os = ""

	case tokens["Opera Mini"] != "":
		ua.name = OperaMini
		ua.version = tokens[OperaMini]
		ua.mobile = true

	case tokens["OPR"] != "":
		ua.name = Opera
		ua.version = tokens["OPR"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens["OPT"] != "":
		ua.name = OperaTouch
		ua.version = tokens["OPT"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	// Opera on iOS
	case tokens["OPiOS"] != "":
		ua.name = Opera
		ua.version = tokens["OPiOS"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	// Chrome on iOS
	case tokens["CriOS"] != "":
		ua.name = Chrome
		ua.version = tokens["CriOS"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	// Firefox on iOS
	case tokens["FxiOS"] != "":
		ua.name = Firefox
		ua.version = tokens["FxiOS"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens["Firefox"] != "":
		ua.name = Firefox
		ua.version = tokens[Firefox]
		_, ua.mobile = tokens["mobile"]
		_, ua.tablet = tokens["tablet"]

	case tokens["Vivaldi"] != "":
		ua.name = Vivaldi
		ua.version = tokens[Vivaldi]

	case tokens.exists("MSIE"):
		ua.name = InternetExplorer
		ua.version = tokens["MSIE"]

	case tokens["EdgiOS"] != "":
		ua.name = Edge
		ua.version = tokens["EdgiOS"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens["Edge"] != "":
		ua.name = Edge
		ua.version = tokens["Edge"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens["Edg"] != "":
		ua.name = Edge
		ua.version = tokens["Edg"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens["EdgA"] != "":
		ua.name = Edge
		ua.version = tokens["EdgA"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens["bingbot"] != "":
		ua.name = "Bingbot"
		ua.version = tokens["bingbot"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens["SamsungBrowser"] != "":
		ua.name = "Samsung Browser"
		ua.version = tokens["SamsungBrowser"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	// if chrome and Safari defined, find any other token sent descr
	case tokens.exists(Chrome) && tokens.exists(Safari):
		name := tokens.findBestMatch(true)
		if name != "" {
			ua.name = name
			ua.version = tokens[name]
			break
		}
		fallthrough

	case tokens.exists("Chrome"):
		ua.name = Chrome
		ua.version = tokens["Chrome"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens.exists("Brave Chrome"):
		ua.name = Chrome
		ua.version = tokens["Brave Chrome"]
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	case tokens.exists("Safari"):
		ua.name = Safari
		if v, ok := tokens["Version"]; ok {
			ua.version = v
		} else {
			ua.version = tokens["Safari"]
		}
		ua.mobile = tokens.existsAny("mobile", "mobile Safari")

	default:
		if ua.os == "Android" && tokens["Version"] != "" {
			ua.name = "Android browser"
			ua.version = tokens["Version"]
			ua.mobile = true
		} else {
			if name := tokens.findBestMatch(false); name != "" {
				ua.name = name
				ua.version = tokens[name]
			} else {
				ua.name = ua.ua
			}
			ua.bot = strings.Contains(strings.ToLower(ua.name), "bot")
			ua.mobile = tokens.existsAny("mobile", "mobile Safari")
		}
	}

	// if tablet, switch mobile to off
	if ua.tablet {
		ua.mobile = false
	}

	// if not already bot, check some popular bots and weather URL is set
	if !ua.bot {
		ua.bot = ua.url != ""
	}

	if !ua.bot {
		switch ua.name {
		case Twitterbot, FacebookExternalHit:
			ua.bot = true
		}
	}
	ua.Parse(uas)
	return ua
}

func parse(userAgent string) (clients properties) {
	clients = make(map[string]string)
	slash := false
	isURL := false
	var buff, val bytes.Buffer
	addToken := func() {
		if buff.Len() != 0 {
			s := strings.TrimSpace(buff.String())
			if _, ign := ignore[s]; !ign {
				if isURL {
					s = strings.TrimPrefix(s, "+")
				}

				if val.Len() == 0 { // only if value don't exists
					var ver string
					s, ver = checkVer(s) // determin version string and split
					clients[s] = ver
				} else {
					clients[s] = strings.TrimSpace(val.String())
				}
			}
		}
		buff.Reset()
		val.Reset()
		slash = false
		isURL = false
	}

	parOpen := false

	bua := []byte(userAgent)
	for i, c := range bua {

		switch {
		case c == 41: // )
			addToken()
			parOpen = false

		case parOpen && c == 59: // ;
			addToken()

		case c == 40: // (
			addToken()
			parOpen = true

		case slash && c == 32:
			addToken()

		case slash:
			val.WriteByte(c)

		case c == 47 && !isURL: //   /
			if i != len(bua)-1 && bua[i+1] == 47 && (bytes.HasSuffix(buff.Bytes(), []byte("http:")) || bytes.HasSuffix(buff.Bytes(), []byte("https:"))) {
				buff.WriteByte(c)
				isURL = true
			} else {
				slash = true
			}

		default:
			buff.WriteByte(c)
		}
	}
	addToken()

	return clients
}

func checkVer(s string) (name, v string) {
	i := strings.LastIndex(s, " ")
	if i == -1 {
		return s, ""
	}

	switch s[:i] {
	case "Linux", "Windows NT", "Windows Phone OS", "MSIE", "Android":
		return s[:i], s[i+1:]
	case "CrOS x86_64", "CrOS aarch64":
		j := strings.LastIndex(s[:i], " ")
		return s[:j], s[j+1 : i]
	default:
		return s, ""
	}
}

type properties map[string]string

func (p properties) exists(key string) bool {
	_, ok := p[key]
	return ok
}

func (p properties) existsAny(keys ...string) bool {
	for _, k := range keys {
		if _, ok := p[k]; ok {
			return true
		}
	}
	return false
}

func (p properties) findMacOSVersion() string {
	for k, v := range p {
		if strings.Contains(k, "OS") {
			if ver := findVersion(v); ver != "" {
				return ver
			} else if ver = findVersion(k); ver != "" {
				return ver
			}
		}
	}
	return ""
}

// findBestMatch from the rest of the bunch
//  in the first cycle only return key with version value
// if withVerValue is false, do another cycle and return any token
func (p properties) findBestMatch(withVerOnly bool) string {
	n := 2
	if withVerOnly {
		n = 1
	}
	for i := 0; i < n; i++ {
		for k, v := range p {
			switch k {
			case Chrome, Firefox, Safari, "Version", "mobile", "mobile Safari", "Mozilla", "AppleWebKit", "Windows NT", "Windows Phone OS", Android, "Macintosh", Linux, "GSA", ChromeOS, "Ubuntu":
			default:
				if i == 0 {
					if v != "" { // in first check, only return keys with value
						return k
					}
				} else {
					return k
				}
			}
		}
	}
	return ""
}

func findVersion(s string) string {
	if ver := rxMacOSVer.FindString(s); ver != "" {
		return strings.Replace(ver, "_", ".", -1)
	}
	return ""
}

// IsWindows shorthand function to check if OS == Windows
func (ua *UserAgent) IsWindows() bool {
	return ua.os == Windows
}

// IsAndroid shorthand function to check if OS == Android
func (ua *UserAgent) IsAndroid() bool {
	return ua.os == Android
}

// IsMacOS shorthand function to check if OS == MacOS
func (ua *UserAgent) IsMacOS() bool {
	return ua.os == MacOS
}

// IsIOS shorthand function to check if OS == IOS
func (ua *UserAgent) IsIOS() bool {
	return ua.os == IOS
}

// IsLinux shorthand function to check if OS == Linux
func (ua *UserAgent) IsLinux() bool {
	return ua.os == Linux
}

// IsOpera shorthand function to check if name == Opera
func (ua *UserAgent) IsOpera() bool {
	return ua.name == Opera
}

// IsOperaMini shorthand function to check if name == Opera Mini
func (ua *UserAgent) IsOperaMini() bool {
	return ua.name == OperaMini
}

// IsChrome shorthand function to check if name == Chrome
func (ua *UserAgent) IsChrome() bool {
	return ua.name == Chrome
}

// IsFirefox shorthand function to check if name == Firefox
func (ua *UserAgent) IsFirefox() bool {
	return ua.name == Firefox
}

// IsInternetExplorer shorthand function to check if name == Internet Explorer
func (ua *UserAgent) IsInternetExplorer() bool {
	return ua.name == InternetExplorer
}

// IsSafari shorthand function to check if name == Safari
func (ua *UserAgent) IsSafari() bool {
	return ua.name == Safari
}

// IsEdge shorthand function to check if name == Edge
func (ua *UserAgent) IsEdge() bool {
	return ua.name == Edge
}

// IsGoogleBot shorthand function to check if name == Googlebot
func (ua *UserAgent) IsGoogleBot() bool {
	return ua.name == Googlebot
}

// IsTwitterBot shorthand function to check if name == Twitterbot
func (ua *UserAgent) IsTwitterBot() bool {
	return ua.name == Twitterbot
}

// IsFacebookBot shorthand function to check if name == FacebookExternalHit
func (ua *UserAgent) IsFacebookBot() bool {
	return ua.name == FacebookExternalHit
}

// detectDevice some properties of the device from the given section.
func (ua *UserAgent) detectDevice(s section) {
	if !ua.mobile {
		return
	}
	if ua.platform == "iPhone" || ua.platform == "iPad" {
		ua.device = ua.platform
		return
	}
	// Android model
	if s.name == "Mozilla" && ua.platform == "Linux" && len(s.comment) > 2 {
		mostAndroidModel := s.comment[2]
		if strings.Contains(mostAndroidModel, "Android") || strings.Contains(mostAndroidModel, "Linux") {
			mostAndroidModel = s.comment[len(s.comment)-1]
		}
		tmp := strings.Split(mostAndroidModel, "Build")
		if len(tmp) > 0 {
			ua.device = strings.Trim(tmp[0], " ")
			return
		}
	}
	// traverse all item
	for _, v := range s.comment {
		if strings.Contains(v, "Build") {
			tmp := strings.Split(v, "Build")
			ua.device = strings.Trim(tmp[0], " ")
		}
	}
}

// Get the name of the bot from the website that may be in the given comment.
// If there is no website in the comment, then an empty string is returned.
func getFromSite(comment []string) string {
	if len(comment) == 0 {
		return ""
	}

	// Where we should check the website.
	idx := 2
	if len(comment) < 3 {
		idx = 0
	} else if len(comment) == 4 {
		idx = 3
	}

	// Pick the site.
	results := botFromSiteRegexp.FindStringSubmatch(comment[idx])
	if len(results) == 1 {
		// If it's a simple comment, just return the name of the site.
		if idx == 0 {
			return results[0]
		}

		// This is a large comment, usually the name will be in the previous
		// field of the comment.
		return strings.TrimSpace(comment[idx-1])
	}
	return ""
}

// Returns true if the info that we currently have corresponds to the Google or Bing mobile bot.
// This function also modifies some attributes in the receiver accordingly.
func (ua *UserAgent) googleOrBingBot() bool {
	// This is a hackish way to detect
	// Google's mobile bot (Googlebot, AdsBot-Google-Mobile, etc.)
	// (See https://support.google.com/webmasters/answer/1061943)
	// and Bing's mobile bot
	// (See https://www.bing.com/webmaster/help/which-crawlers-does-bing-use-8c184ec0)
	if strings.Contains(ua.ua, "Google") || strings.Contains(ua.ua, "bingbot") {
		ua.platform = ""
		ua.undecided = true
	}
	return ua.undecided
}

// Returns true if we think that it is iMessage-Preview.
// This function also modifies some attributes in the receiver accordingly.
func (ua *UserAgent) iMessagePreview() bool {
	// iMessage-Preview doesn't advertise itself.
	// We have a to rely on a hack
	// to detect it: it impersonates both facebook and Twitter bots.
	// See https://medium.com/@siggi/apples-imessage-impersonates-twitter-facebook-bots-when-scraping-cef85b2cbb7d
	if !strings.Contains(ua.ua, "facebookexternalhit") {
		return false
	}
	if !strings.Contains(ua.ua, "Twitterbot") {
		return false
	}
	ua.bot = true
	ua.browser.Name = "iMessage-Preview"
	ua.browser.Engine = ""
	ua.browser.EngineVersion = ""
	// We don't set the mobile flag because iMessage can be on iOS (mobile) or macOS (not mobile).
	return true
}

// Set the attributes of the receiver as given by the parameters.
// All the other parameters are set to empty.
func (ua *UserAgent) setSimple(name, version string, bot bool) {
	ua.bot = bot
	if !bot {
		ua.mozilla = ""
	}
	ua.browser.Name = name
	ua.browser.Version = version
	ua.browser.Engine = ""
	ua.browser.EngineVersion = ""
	ua.os = ""
	ua.localization = ""
}

// Fix some values for some weird browsers.
func (ua *UserAgent) fixOther(sections []section) {
	if len(sections) > 0 {
		ua.browser.Name = sections[0].name
		ua.browser.Version = sections[0].version
		ua.mozilla = ""
	}
}

// Check if we're dealing with a bot or with some weird browser.
// If that is the case, the receiver will be modified accordingly.
func (ua *UserAgent) checkBot(sections []section) {
	// If there's only one element, and it doesn't have the Mozilla string,
	// check whether this is a bot or not.
	if len(sections) == 1 && sections[0].name != "Mozilla" {
		ua.mozilla = ""

		// Check whether the name has some suspicious "bot" or "crawler" in his name.
		if botRegex.Match([]byte(sections[0].name)) {
			ua.setSimple(sections[0].name, "", true)
			return
		}

		// Tough luck, let's try to see if it has a website in his comment.
		if name := getFromSite(sections[0].comment); name != "" {
			// First, this is a bot.
			// Moreover, since it doesn't have the Mozilla string,
			// we can assume that the name and the version are
			// the ones from the first section.
			ua.setSimple(sections[0].name, sections[0].version, true)
			return
		}

		// At this point, we are sure that this is not a bot, but some weirdo.
		ua.setSimple(sections[0].name, sections[0].version, false)
	} else {
		// Let's iterate over the available comments and check for a website.
		for _, v := range sections {
			if name := getFromSite(v.comment); name != "" {
				// Ok, we've got a bot name.
				results := strings.SplitN(name, "/", 2)
				version := ""
				if len(results) == 2 {
					version = results[1]
				}
				ua.setSimple(results[0], version, true)
				return
			}
		}

		// We will assume that this is some other weird browser.
		ua.fixOther(sections)
	}
}

// Name of the browser.
func (ua *UserAgent) Name() string {
	return ua.name
}

// Version of the browser.
func (ua *UserAgent) Version() string {
	return ua.version
}

// OS returns a string containing the name of the Operating System.
func (ua *UserAgent) OS() string {
	return ua.os
}

// OSVersion returns a string containing the OS version.
func (ua *UserAgent) OSVersion() string {
	return ua.osVersion
}

// Device The Device returns a string containing the Phone Model like "Nexus 5X".
func (ua *UserAgent) Device() string {
	return ua.device
}

// Mobile returns true if it's a mobile device, false otherwise.
func (ua *UserAgent) Mobile() bool {
	return ua.mobile
}

// Tablet returns true if it's a tablet device, false otherwise.
func (ua *UserAgent) Tablet() bool {
	return ua.tablet
}

// Desktop returns true if it's a desktop device, false otherwise.
func (ua *UserAgent) Desktop() bool {
	return ua.desktop
}

// Bot returns true if it's a bot, false otherwise.
func (ua *UserAgent) Bot() bool {
	return ua.bot
}

// Mozilla returns the mozilla version (it's how the User Agent string begins:
// "Mozilla/5.0 ...", unless we're dealing with Opera, of course).
func (ua *UserAgent) Mozilla() string {
	return ua.mozilla
}

// URL returns the URL of the website that the user agent belongs to.
func (ua *UserAgent) URL() string {
	return ua.url
}

// UA returns the original given user agent.
func (ua *UserAgent) UA() string {
	return ua.ua
}

// Platform returns a string containing the platform.
func (ua *UserAgent) Platform() string {
	return ua.platform
}

// UserAgentBrowser returns a string containing the name of the browser.
func (ua *UserAgent) UserAgentBrowser() Browser {
	return ua.browser
}

// Localization returns a string containing the localization.
func (ua *UserAgent) Localization() string {
	return ua.localization
}

// ShortOS returns the short name of the operating system.
func (ua *UserAgent) ShortOS() string {
	return ua.shortOS
}

// Parse the given User-Agent string.
// After calling this function,
// the receiver will be set up with all the information that we've extracted.
func (ua *UserAgent) Parse(str string) {
	var sections []section

	ua.ua = str
	for index, limit := 0, len(str); index < limit; {
		s := parseSection(str, &index)
		if !ua.mobile && s.name == "Mobile" {
			ua.mobile = true
		}
		sections = append(sections, s)
	}
	if len(sections) > 0 {
		if sections[0].name == "Mozilla" {
			ua.mozilla = sections[0].version
		}

		ua.detectBrowser(sections)
		ua.detectOS(sections[0])
		ua.detectDevice(sections[0])
		if ua.undecided {
			ua.checkBot(sections)
		}
	}
}

// Read from the given string until the given delimiter or the
// end of the string have been reached.
//
// The first argument is the user agent string being parsed.
// The second argument is a reference pointing to the current index of the user agent string.
// The delimiter argument specifies which character is the delimiter
// and the cat argument determines whether nested '(' should be ignored or not.
//
// Returns an array of bytes containing what has been read.
func readUntil(ua string, index *int, delimiter byte, cat bool) []byte {
	var buffer []byte

	i := *index
	catalan := 0
	for ; i < len(ua); i = i + 1 {
		if ua[i] == delimiter {
			if catalan == 0 {
				*index = i + 1
				return buffer
			}
			catalan--
		} else if cat && ua[i] == '(' {
			catalan++
		}
		buffer = append(buffer, ua[i])
	}
	*index = i + 1
	return buffer
}

// Parse the given product, that is, just a name or a string
// formatted as Name/Version.
//
// It returns two strings.
// The first string is the name of the product,
// and the second string contains the version of the product.
func parseProduct(product []byte) (string, string) {
	prod := strings.SplitN(string(product), "/", 2)
	if len(prod) == 2 {
		return prod[0], prod[1]
	}
	return string(product), ""
}

// Parse a section.
// A section is typically formatted as follows "Name/Version (comment)".
// Both the comments and the version are optional.
//
// The first argument is the user agent string being parsed.
// The second argument is a reference pointing to the current index of the user agent string.
//
// Returns a section containing the information that we could extract
// from the last parsed section.
func parseSection(ua string, index *int) (s section) {
	var buffer []byte

	// Check for empty products
	if *index < len(ua) && ua[*index] != '(' && ua[*index] != '[' {
		buffer = readUntil(ua, index, ' ', false)
		s.name, s.version = parseProduct(buffer)
	}

	if *index < len(ua) && ua[*index] == '(' {
		*index++
		buffer = readUntil(ua, index, ')', true)
		s.comment = strings.Split(string(buffer), "; ")
		*index++
	}

	// Discards any trailing data within square brackets
	if *index < len(ua) && ua[*index] == '[' {
		*index++
		_ = readUntil(ua, index, ']', true)
		*index++
	}
	return s
}

// Initialize the parser.
func (ua *UserAgent) initialize() {
	ua.ua = ""
	ua.mozilla = ""
	ua.platform = ""
	ua.os = ""
	ua.localization = ""
	ua.device = ""
	ua.browser.Engine = ""
	ua.browser.EngineVersion = ""
	ua.browser.Name = ""
	ua.browser.Version = ""
	ua.bot = false
	ua.mobile = false
	ua.undecided = false
}

// Beautify the given string.
// Internal: beautify the UserAgent reference into a string, so it can be tested later on.
//
// ua - a UserAgent reference.
//
// Returns a string that contains the beautified representation.
func (ua *UserAgent) Beautify() (s string) {
	if len(ua.Mozilla()) > 0 {
		s += "Mozilla:" + ua.Mozilla() + " "
	}
	if len(ua.Platform()) > 0 {
		s += "Platform:" + ua.Platform() + " "
	}
	if len(ua.OS()) > 0 {
		s += "OS:" + ua.OS() + " "
	}
	if len(ua.Localization()) > 0 {
		s += "Localization:" + ua.Localization() + " "
	}
	if len(ua.Device()) > 0 {
		s += "Device:" + ua.Device() + " "
	}
	str1, str2 := ua.Browser()
	if len(str1) > 0 {
		s += "Browser:" + str1
		if len(str2) > 0 {
			s += "-" + str2 + " "
		} else {
			s += " "
		}
	}
	str1, str2 = ua.Engine()
	if len(str1) > 0 {
		s += "Engine:" + str1
		if len(str2) > 0 {
			s += "-" + str2 + " "
		} else {
			s += " "
		}
	}
	s += "Bot:" + fmt.Sprintf("%v", ua.Bot()) + " "
	s += "Mobile:" + fmt.Sprintf("%v", ua.Mobile())
	return s
}
