package detect

import (
	"regexp"
	"strconv"
	"strings"
)

// @Project: go-mobile-detect
// @Author: houseme
// @Description:
// @File: properties.go
// @Version: 1.0.0
// @Date: 2021/3/13 21:14
// @Package go_mobile_detect

const (
	PropMobile = iota
	PropBuild
	PropVersion
	PropVendorid
	PropIpad
	PropIphone
	PropIpod
	PropKindle
	PropChrome
	PropCoast
	PropDolfin
	PropFirefox
	PropFennec
	PropIe
	PropNetfront
	PropNokiabrowser
	PropOpera
	PropOperaMini
	PropOperaMobi
	PropUcBrowser
	PropMqqbrowser
	PropMicromessenger
	PropBaiduboxapp
	PropBaidubrowser
	PropSafari
	PropSkyfire
	PropTizen
	PropWebkit
	PropGecko
	PropTrident
	PropPresto
	PropIos
	PropAndroid
	PropBlackberry
	PropBrew
	PropJava
	PropWindowsPhoneOs
	PropWindowsPhone
	PropWindowsCe
	PropWindowsNt
	PropSymbian
	PropWebos
)

var (
	propertiesNameToVal = map[string]int{
		"mobile":           PropMobile,
		"build":            PropBuild,
		"version":          PropVersion,
		"vendorid":         PropVendorid,
		"ipad":             PropIpad,
		"iphone":           PropIphone,
		"ipod":             PropIpod,
		"kindle":           PropKindle,
		"chrome":           PropChrome,
		"coast":            PropCoast,
		"dolfin":           PropDolfin,
		"firefox":          PropFirefox,
		"fennec":           PropFennec,
		"ie":               PropIe,
		"netfront":         PropNetfront,
		"nokiabrowser":     PropNokiabrowser,
		"opera":            PropOpera,
		"opera mini":       PropOperaMini,
		"opera mobi":       PropOperaMobi,
		"uc browser":       PropUcBrowser,
		"mqqbrowser":       PropMqqbrowser,
		"micromessenger":   PropMicromessenger,
		"baiduboxapp":      PropBaiduboxapp,
		"baidubrowser":     PropBaidubrowser,
		"safari":           PropSafari,
		"skyfire":          PropSkyfire,
		"tizen":            PropTizen,
		"webkit":           PropWebkit,
		"gecko":            PropGecko,
		"trident":          PropTrident,
		"presto":           PropPresto,
		"ios":              PropIos,
		"android":          PropAndroid,
		"blackberry":       PropBlackberry,
		"brew":             PropBrew,
		"java":             PropJava,
		"windows phone os": PropWindowsPhoneOs,
		"windows phone":    PropWindowsPhone,
		"windows ce":       PropWindowsCe,
		"windows nt":       PropWindowsNt,
		"symbian":          PropSymbian,
		"webos":            PropWebos,
	}

	// Properties helps parsing User Agent string, extracting useful segments of text.
	//VER refers to the regular expression defined in the constant self::VER.
	props = [...][]string{
		// Build
		//MOBILE:PROP_
		[]string{`Mobile/[VER]`},
		//PROP_BUILD:
		[]string{`Build/[VER]`},
		//PROP_VERSION:
		[]string{`Version/[VER]`},
		//PROP_VENDORID:
		[]string{`VendorID/[VER]`},
		// Devices
		//PROP_IPAD:
		[]string{`iPad.*CPU[a-z ]+[VER]`},
		//PROP_IPHONE:
		[]string{`iPhone.*CPU[a-z ]+[VER]`},
		//PROP_IPOD:
		[]string{`iPod.*CPU[a-z ]+[VER]`},
		//`BlackBerry`    : array(`BlackBerry[VER]`, `BlackBerry [VER];`),
		//PROP_KINDLE:
		[]string{`Kindle/[VER]`},
		// Browser
		//PROP_CHROME:
		[]string{`Chrome/[VER]`, `CriOS/[VER]`, `CrMo/[VER]`},
		//PROP_COAST:
		[]string{`Coast/[VER]`},
		//PROP_DOLFIN:
		[]string{`Dolfin/[VER]`},
		// @reference: https://developer.mozilla.org/en-US/docs/User_Agent_Strings_Reference
		//PROP_FIREFOX:
		[]string{`Firefox/[VER]`},
		//PROP_FENNEC:
		[]string{`Fennec/[VER]`},
		// @reference: http://msdn.microsoft.com/en-us/library/ms537503(v=vs.85).aspx
		//PROP_IE:
		[]string{`IEMobile/[VER];`, `IEMobile [VER]`, `MSIE [VER];`},
		// http://en.wikipedia.org/wiki/NetFront
		//PROP_NETFRONT:
		[]string{`NetFront/[VER]`},
		//PROP_NOKIABROWSER:
		[]string{`NokiaBrowser/[VER]`},
		//PROP_OPERA:
		[]string{` OPR/[VER]`, `Opera Mini/[VER]`, `Version/[VER]`},
		//PROP_OPERA_MINI:
		[]string{`Opera Mini/[VER]`},
		//PROP_OPERA_MOBI:
		[]string{`Version/[VER]`},
		//PROP_UC_BROWSER:
		[]string{`UC Browser[VER]`},
		//PROP_MQQBROWSER:
		[]string{`MQQBrowser/[VER]`},
		//PROP_MICROMESSENGER:
		[]string{`MicroMessenger/[VER]`},
		//PROP_BAIDUBOXAPP
		[]string{`baiduboxapp/[VER]`},
		//PROP_BAIDUBROWSER
		[]string{`baidubrowser/[VER]`},
		// @note: Safari 7534.48.3 is actually Version 5.1.
		// @note: On BlackBerry the Version is overwriten by the OS.
		//PROP_SAFARI:
		[]string{`Version/[VER]`, `Safari/[VER]`},
		//PROP_SKYFIRE:
		[]string{`Skyfire/[VER]`},
		//PROP_TIZEN:
		[]string{`Tizen/[VER]`},
		//PROP_WEBKIT:
		[]string{`webkit[ /][VER]`},
		// Engine
		//PROP_GECKO:
		[]string{`Gecko/[VER]`},
		//PROP_TRIDENT:
		[]string{`Trident/[VER]`},
		//PROP_PRESTO:
		[]string{`Presto/[VER]`},
		// OS
		//PROP_IOS:
		[]string{` \bOS\b [VER] `},
		//PROP_ANDROID:
		[]string{`Android [VER]`},
		//PROP_BLACKBERRY:
		[]string{`BlackBerry[\w]+/[VER]`, `BlackBerry.*Version/[VER]`, `Version/[VER]`},
		//PROP_BREW:
		[]string{`BREW [VER]`},
		//PROP_JAVA:
		[]string{`Java/[VER]`},
		// @reference: http://windowsteamblog.com/windows_phone/b/wpdev/archive/2011/08/29/introducing-the-ie9-on-windows-phone-mango-user-agent-string.aspx
		// @reference: http://en.wikipedia.org/wiki/Windows_NT#Releases
		//PROP_WINDOWS_PHONE_OS:
		[]string{`Windows Phone OS [VER]`, `Windows Phone [VER]`},
		//PROP_WINDOWS_PHONE:
		[]string{`Windows Phone [VER]`},
		//PROP_WINDOWS_CE:
		[]string{`Windows CE/[VER]`},
		// http://social.msdn.microsoft.com/Forums/en-US/windowsdeveloperpreviewgeneral/thread/6be392da-4d2f-41b4-8354-8dcee20c85cd
		//PROP_WINDOWS_NT:
		[]string{`Windows NT [VER]`},
		//PROP_SYMBIAN:
		[]string{`SymbianOS/[VER]`, `Symbian/[VER]`},
		//PROP_WEBOS:
		[]string{`webOS/[VER]`, `hpwOS/[VER];`},
	}
)

type properties struct {
	cache map[string]*regexp.Regexp
}

func newProperties() *properties {
	p := &properties{}
	p.cache = make(map[string]*regexp.Regexp)
	p.preCompile()
	return p
}

func (p *properties) preCompile() {
	for _, property := range props {
		for _, pattern := range property {
			p.compiledRegexByPattern(pattern)
		}
	}
}

func (p *properties) compiledRegexByPattern(propertyPattern string) *regexp.Regexp {
	re, ok := p.cache[propertyPattern]
	if !ok {
		p.cache[propertyPattern] = regexp.MustCompile(propertyPattern)
	}
	re = p.cache[propertyPattern]
	return re
}

func (p *properties) version(propertyVal int, userAgent string) string {
	if len(props) >= propertyVal {
		for _, propertyMatchString := range props[propertyVal] {
			propertyPattern := `(?is)` + strings.Replace(string(propertyMatchString), `[VER]`, verRegex, -1)

			// Escape the special character which is the delimiter.
			//propertyPattern = strings.Replace(propertyPattern, `/`, `\/`, -1)

			// Identify and extract the version.
			re := p.compiledRegexByPattern(propertyPattern)
			match := re.FindStringSubmatch(userAgent)
			if len(match) > 0 {
				return match[1]
			}
		}
	}
	return ""
}

func (p *properties) nameToKey(propertyName string) int {
	propertyName = strings.ToLower(propertyName)
	propertyVal, ok := propertiesNameToVal[propertyName]
	if !ok {
		return -1
	}
	return propertyVal
}

func (p *properties) versionByName(propertyName, userAgent string) string {
	if "" != propertyName {
		propertyVal := p.nameToKey(propertyName)
		if -1 != propertyVal {
			return p.version(propertyVal, userAgent)
		}
	}
	return ""
}

func (p *properties) versionFloatName(propertyName, userAgent string) float64 {
	propertyVal := p.nameToKey(propertyName)
	if -1 != propertyVal {
		return p.versionFloat(propertyVal, userAgent)
	}
	return 0.0
}

func (p *properties) versionFloat(propertyVal int, userAgent string) float64 {
	version := p.version(propertyVal, userAgent)
	replacer := strings.NewReplacer(`_`, `.`, `/`, `.`)
	version = replacer.Replace(version)

	versionNumbers := strings.Split(version, `.`)

	versionNumbersLength := len(versionNumbers)
	if versionNumbersLength > 1 {
		firstNumber := versionNumbers[0]
		retVersion := make([]string, versionNumbersLength-1)
		for i := 1; i < versionNumbersLength; i++ {
			retVersion[(i - 1)] = strings.Replace(versionNumbers[i], `.`, ``, -1)
		}

		version = firstNumber + `.` + strings.Join(retVersion, ``)
	}
	versionFloat, err := strconv.ParseFloat(version, 64)

	if nil != err {
		return 0.0
	}
	return versionFloat
}
