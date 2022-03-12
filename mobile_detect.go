package mobiledetect

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

// @Project: mobile-detect
// @Author: houseme
// @Description:
// @File: mobile_detect
// @Version: 1.0.0
// @Date: 2021/3/13 20:27
// @Package mobiledetect

const (
	// A frequently used regular expression to extract version #s.
	verRegex = `([\w._\+]+)`

	// MobileGradeA .
	MobileGradeA = "A"
	// MobileGradeB .
	MobileGradeB = "B"
	// MobileGradeC .
	MobileGradeC = "C"
)

// Device Vars returns the route variables for the current request, if any.
func Device(r *http.Request) string {
	if rv := r.Context().Value("Device"); rv != nil {
		return rv.(string)
	}
	return ""
}

// DeviceHandler .
type DeviceHandler interface {
	Mobile(w http.ResponseWriter, r *http.Request, m *MobileDetect)
	Tablet(w http.ResponseWriter, r *http.Request, m *MobileDetect)
	Desktop(w http.ResponseWriter, r *http.Request, m *MobileDetect)
}

// Handler .
func Handler(h DeviceHandler, rules *rules) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := New(r, rules)
		if m.IsTablet() {
			h.Tablet(w, r, m)
		} else if m.IsMobile() {
			h.Mobile(w, r, m)
		} else {
			h.Desktop(w, r, m)
		}
	})
}

// HandlerMux .
func HandlerMux(s *http.ServeMux, rules *rules) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := New(r, rules)
		if m.IsTablet() {
			r = r.WithContext(context.WithValue(r.Context(), "Device", "Tablet"))
		} else if m.IsMobile() {
			r = r.WithContext(context.WithValue(r.Context(), "Device", "Mobile"))
		} else {
			r = r.WithContext(context.WithValue(r.Context(), "Device", "Desktop"))
		}
		s.ServeHTTP(w, r)
	})
}

// MobileDetect holds the structure to figure out a browser from a UserAgent string and methods necessary to make it happen
type MobileDetect struct {
	rules              *rules
	userAgent          string
	httpHeaders        map[string]string
	compiledRegexRules map[string]*regexp.Regexp
	*properties
}

// New creates the MobileDetect object
func New(r *http.Request, rules *rules) *MobileDetect {
	if nil == rules {
		rules = NewRules()
	}
	md := &MobileDetect{
		rules:              rules,
		userAgent:          r.UserAgent(),
		httpHeaders:        httpHeaders(r),
		compiledRegexRules: make(map[string]*regexp.Regexp, len(rules.mobileDetectionRules())),
		properties:         newProperties(),
	}
	return md
}

func httpHeaders(r *http.Request) map[string]string {
	httpHeaders := map[string]string{
		"SERVER_SOFTWARE":  r.Header.Get("SERVER_SOFTWARE"),
		"REQUEST_METHOD":   r.Method,
		"HOST":             r.Host,
		"X_REAL_IP":        r.Header.Get("X_REAL_IP"),
		"X_FORWARDED_FOR":  r.Header.Get("X_FORWARDED_FOR"),
		"CONNECTION":       r.Header.Get("CONNECTION"),
		"USER-AGENT":       r.UserAgent(),
		"ACCEPT":           r.Header.Get("ACCEPT"),
		"ACCEPT-LANGUAGE":  r.Header.Get("ACCEPT-LANGUAGE"),
		"ACCEPT-ENCODING":  r.Header.Get("ACCEPT-ENCODING"),
		"X_REQUESTED_WITH": r.Header.Get("X_REQUESTED_WITH"),
		"REFERER":          r.Referer(),
		"PRAGMA":           r.Header.Get("PRAGMA"),
		"CACHE_CONTROL":    r.Header.Get("CACHE_CONTROL"),
		"REMOTE_ADDR":      r.RemoteAddr,
		"REQUEST_TIME":     r.Header.Get("REQUEST_TIME"),
	}

	return httpHeaders
}

// PreCompileRegexRules .
func (md *MobileDetect) PreCompileRegexRules() *MobileDetect {
	for _, ruleValue := range md.rules.mobileDetectionRules() {
		md.match(ruleValue)
	}
	return md
}

// SetUserAgent .
func (md *MobileDetect) SetUserAgent(userAgent string) *MobileDetect {
	md.userAgent = userAgent
	return md
}

// SetHTTPHeaders .
func (md *MobileDetect) SetHTTPHeaders(httpHeaders map[string]string) *MobileDetect {
	md.httpHeaders = httpHeaders
	return md
}

// IsMobile is a specific case to detect only mobile browsers.
func (md *MobileDetect) IsMobile() bool {
	if md.CheckHTTPHeadersForMobile() {
		return true
	}
	return md.matchDetectionRulesAgainstUA()
}

// IsTablet is a specific case of detect only tablet browsers on tablets. Do not overlap with IsTablet
func (md *MobileDetect) IsTablet() bool {
	for _, ruleValue := range md.rules.tabletDevices {
		if md.match(ruleValue) {
			return true
		}
	}
	return false
}

// IsKey Is compared the detected browser with a "rule" from the existing rules list
func (md *MobileDetect) IsKey(key int) bool {
	return md.matchUAAgainstKey(key)
}

// Is It is recommended to use IsKey instead
func (md *MobileDetect) Is(key interface{}) bool {
	switch key.(type) {
	case string:
		name := strings.ToLower(key.(string))
		ruleKey, ok := md.rules.nameToKey(name)
		if !ok {
			return false
		}
		return md.matchUAAgainstKey(ruleKey)
	case int:
		ruleKey := key.(int)
		return md.IsKey(ruleKey)
	}
	return false
}

// VersionFloatKey VersionFloat does the same as Version, but returns a float number good for version comparison
func (md *MobileDetect) VersionFloatKey(propertyVal int) float64 {
	return md.properties.versionFloat(propertyVal, md.userAgent)
}

// VersionKey Version detects the browser version returning as string
func (md *MobileDetect) VersionKey(propertyVal int) string {
	return md.properties.version(propertyVal, md.userAgent)
}

// VersionFloat It is recommended to use VersionFloatKey instead
func (md *MobileDetect) VersionFloat(propertyName interface{}) float64 {
	switch propertyName.(type) {
	case string:
		return md.properties.versionFloatName(propertyName.(string), md.userAgent)
	case int:
		return md.VersionFloatKey(propertyName.(int))
	}
	return 0.0
}

// Version It is recommended to use VersionKey instead
func (md *MobileDetect) Version(propertyName interface{}) string {
	switch propertyName.(type) {
	case string:
		return md.properties.versionByName(propertyName.(string), md.userAgent)
	case int:
		return md.VersionKey(propertyName.(int))
	}
	return ""
}

// Search for a certain key in the rules array.
// If the key is found the try to match the corresponding regex agains the User-Agent.
func (md *MobileDetect) matchUAAgainstKey(key int) bool {
	ret := false
	rules := md.rules.mobileDetectionRules()
	for ruleKey, ruleValue := range rules {
		if key == ruleKey {
			ret = md.match(ruleValue)
			break
		}
	}

	return ret
}

// Find a detection rule that matches the current User-agent.
func (md *MobileDetect) matchDetectionRulesAgainstUA() bool {
	for _, ruleValue := range md.rules.mobileDetectionRules() {
		if "" != ruleValue {
			if md.match(ruleValue) {
				return true
			}
		}
	}

	return false
}

// Some detection rules are relative (not standard),because of the diversity of devices, vendors and
// their conventions in representing the User-Agent or the HTTP headers.
// This method will be used to check custom regexes against the User-Agent string.
// @todo: search in the HTTP headers too.
func (md *MobileDetect) match(ruleValue string) bool {
	// Escape the special character which is the delimiter
	// rule = strings.Replace(rule, `\`, `\/`, -1)
	ruleValue = `(?is)` + ruleValue
	var re *regexp.Regexp
	re = md.compiledRegexRules[ruleValue]
	if nil == re {
		md.compiledRegexRules[ruleValue] = regexp.MustCompile(ruleValue)
	}
	re = md.compiledRegexRules[ruleValue]
	ret := re.MatchString(md.userAgent)
	return ret
}

// CheckHTTPHeadersForMobile looks for mobile rules to confirm if the browser is a mobile browser
func (md *MobileDetect) CheckHTTPHeadersForMobile() bool {
	for _, mobileHeader := range md.mobileHeaders() {
		if headerString, ok := md.httpHeaders[mobileHeader]; ok {
			mobileHeaderMatches := md.mobileHeaderMatches()
			if matches, ok := mobileHeaderMatches[mobileHeader]; ok {
				for _, match := range matches {
					if -1 != strings.Index(headerString, match) {
						return true
					}
				}
				return false
			}
			return true
		}
	}
	return false
}

func (md *MobileDetect) mobileHeaders() []string {
	return []string{
		"HTTP_ACCEPT",
		"HTTP_X_WAP_PROFILE",
		"HTTP_X_WAP_CLIENTID",
		"HTTP_WAP_CONNECTION",
		"HTTP_PROFILE",
		// Reported by Opera on Nokia devices (eg. C3).
		"HTTP_X_OPERAMINI_PHONE_UA",
		"HTTP_X_NOKIA_GATEWAY_ID",
		"HTTP_X_ORANGE_ID",
		"HTTP_X_VODAFONE_3GPDPCONTEXT",
		"HTTP_X_HUAWEI_USERID",
		// Reported by Windows Smartphones.
		"HTTP_UA_OS",
		// Reported by Verizon, Vodafone proxy system.
		"HTTP_X_MOBILE_GATEWAY",
		// Seend this on HTC Sensation. @ref: SensationXE_Beats_Z715e.
		"HTTP_X_ATT_DEVICEID",
		// Seen this on a HTC.
		"HTTP_UA_CPU",
	}
}

func (md *MobileDetect) mobileHeaderMatches() map[string][]string {
	return map[string][]string{
		"HTTP_ACCEPT": []string{
			// Opera Mini; @reference: http://dev.opera.com/articles/view/opera-binary-markup-language/
			"application/x-obml2d",
			// BlackBerry devices.
			"application/vnd.rim.html",
			"text/vnd.wap.wml",
			"application/vnd.wap.xhtml+xml",
		},
		"HTTP_UA_CPU": []string{"ARM"},
	}
}

// MobileGrade returns a graduation similar to jQuery's Graded Browse Support
func (md *MobileDetect) MobileGrade() string {
	isMobile := md.IsMobile()

	if md.isMobileGradeA(isMobile) {
		return MobileGradeA
	}
	if md.isMobileGradeB() {
		return MobileGradeB
	}
	if md.isMobileGradeC(isMobile) {
		return MobileGradeC
	}
	return MobileGradeC
}

func (md *MobileDetect) isMobileGradeA(isMobile bool) bool {
	if md.VersionFloat("iPad") >= 4.3 || md.VersionFloat("iPhone") >= 4.3 || md.VersionFloat("iPod") >= 4.3 ||
		(md.VersionFloat("Android") > 2.1 && md.Is("Webkit")) ||
		md.VersionFloat("Windows Phone OS") >= 7.5 ||
		md.Is("BlackBerry") && md.VersionFloat("BlackBerry") >= 6.0 ||
		md.match("Playbook.*Tablet") ||
		(md.VersionFloat("webOS") >= 1.4 && md.match("Palm|Pre|Pixi")) ||
		md.match("hp.*TouchPad") ||
		(md.Is("Firefox") && md.VersionFloat("Firefox") >= 18) ||
		(md.Is("Chrome") && md.Is("AndroidOS") && md.VersionFloat("Android") >= 4.0) ||
		(md.Is("Skyfire") && md.VersionFloat("Skyfire") >= 4.1 && md.Is("AndroidOS") && md.VersionFloat("Android") >= 2.3) ||
		(md.Is("Opera") && md.VersionFloat("Opera Mobi") > 11.5 && md.Is("AndroidOS")) ||
		md.Is("MeeGoOS") ||
		md.Is("SailfishOS") ||
		md.Is("Tizen") ||
		md.Is("Dolfin") && md.VersionFloat("Bada") >= 2.0 ||
		((md.Is("UC Browser") || md.Is("Dolfin")) && md.VersionFloat("Android") >= 2.3) ||
		(md.match("Kindle Fire") || md.Is("Kindle") && md.VersionFloat("Kindle") >= 3.0) ||
		(md.Is("AndroidOS") && md.Is("NookTablet")) ||
		(md.VersionFloat("Chrome") >= 16 && isMobile) ||
		(md.VersionFloat("Safari") >= 5.0 && isMobile) ||
		(md.VersionFloat("Firefox") >= 10.0 && isMobile) ||
		(md.VersionFloat("MSIE") >= 7.0 && isMobile) ||
		(md.VersionFloat("Opera") >= 10 && isMobile) {
		return true
	}
	return false
}

func (md *MobileDetect) isMobileGradeB() bool {
	if (md.VersionFloat("iPad") < 4.3 || md.VersionFloat("iPhone") < 4.3 || md.VersionFloat("iPod") < 4.3 || md.Is("Blackberry") && md.VersionFloat("BlackBerry") >= 5 && md.VersionFloat("BlackBerry") < 6) ||
		(md.VersionFloat("Opera Mini") >= 5.0 && md.VersionFloat("Opera Mini") <= 7.0 && (md.VersionFloat(
			"Android") >= 2.3 || md.Is("iOS"))) ||
		md.match("NokiaN8|NokiaC7|N97.*Series60|Symbian/3") ||
		(md.VersionFloat("Opera Mobi") >= 11 && md.Is("SymbianOS")) {
		return true
	}
	return false
}

func (md *MobileDetect) isMobileGradeC(isMobile bool) bool {
	if md.VersionFloat("Blackberry") <= 5.0 || md.match("MSIEMobile|Windows CE.*Mobile") || md.VersionFloat(
		"Windows Mobile") <= 5.2 || md.VersionFloat("iPad") <= 3.2 || md.VersionFloat("iPhone") <= 3.2 || md.
		VersionFloat("iPod") <= 3.2 || md.VersionFloat("IE") <= 7.0 && !isMobile {
		return true
	}
	return false
}
