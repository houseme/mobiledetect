// Copyright (C) 2012-2021 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package ua

import (
	"strings"
)

// Normalize the name of the operating system.
// By now, this just affects Windows NT.
//
// Returns a string containing the normalized name for the Operating System.
func normalizeOS(name string) string {
	sp := strings.SplitN(name, " ", 3)
	if len(sp) != 3 || sp[1] != "NT" {
		return name
	}

	switch sp[2] {
	case "5.0":
		return "Windows 2000"
	case "5.01":
		return "Windows 2000, Service Pack 1 (SP1)"
	case "5.1":
		return "Windows XP"
	case "5.2":
		return "Windows XP x64 Edition"
	case "6.0":
		return "Windows Vista"
	case "6.1":
		return "Windows 7"
	case "6.2":
		return "Windows 8"
	case "6.3":
		return "Windows 8.1"
	case "10.0":
		return "Windows 10"
	}
	return name
}

// Guess the OS, the localization and if this is a mobile device for a
// Webkit-powered browser.
//
// The first argument p is a reference to the current UserAgent,
// and the second argument is a slice of strings containing the comment.
func webkit(ua *UserAgent, comment []string) {
	if ua.platform == "webOS" {
		ua.browser.Name = ua.platform
		ua.os = "Palm"
		if len(comment) > 2 {
			ua.localization = comment[2]
		}
		ua.mobile = true
	} else if ua.platform == "Symbian" {
		ua.mobile = true
		ua.browser.Name = ua.platform
		ua.os = comment[0]
	} else if ua.platform == "Linux" {
		ua.mobile = true
		if ua.browser.Name == "Safari" {
			ua.browser.Name = "Android"
		}
		if len(comment) > 1 {
			if comment[1] == "U" || comment[1] == "arm_64" {
				if len(comment) > 2 {
					ua.os = comment[2]
				} else {
					ua.mobile = false
					ua.os = comment[0]
				}
			} else {
				ua.os = comment[1]
			}
		}
		if len(comment) > 3 {
			ua.localization = comment[3]
		} else if len(comment) == 3 {
			_ = ua.googleOrBingBot()
		}
	} else if len(comment) > 0 {
		if len(comment) > 3 {
			ua.localization = comment[3]
		}
		if strings.HasPrefix(comment[0], "Windows NT") {
			ua.os = normalizeOS(comment[0])
		} else if len(comment) < 2 {
			ua.localization = comment[0]
		} else if len(comment) < 3 {
			if !ua.googleOrBingBot() && !ua.iMessagePreview() {
				ua.os = normalizeOS(comment[1])
			}
		} else {
			ua.os = normalizeOS(comment[2])
		}
		if ua.platform == "BlackBerry" {
			ua.browser.Name = ua.platform
			if ua.os == "Touch" {
				ua.os = ua.platform
			}
		}
	}

	// Special case for Firefox on iPad, where the platform is advertised as Macintosh instead of iPad
	if ua.platform == "Macintosh" && ua.browser.Engine == "AppleWebKit" && ua.browser.Name == "Firefox" {
		ua.platform = "iPad"
		ua.mobile = true
	}
}

// Guess the OS, the localization and if this is a mobile device
// for a Gecko-powered browser.
//
// The first argument p is a reference to the current UserAgent,
// and the second argument is a slice of strings containing the comment.
func gecko(ua *UserAgent, comment []string) {
	if len(comment) > 1 {
		if comment[1] == "U" || comment[1] == "arm_64" {
			if len(comment) > 2 {
				ua.os = normalizeOS(comment[2])
			} else {
				ua.os = normalizeOS(comment[1])
			}
		} else {
			if strings.Contains(ua.platform, "Android") {
				ua.mobile = true
				ua.platform, ua.os = normalizeOS(comment[1]), ua.platform
			} else if comment[0] == "Mobile" || comment[0] == "Tablet" {
				ua.mobile = true
				ua.os = "FirefoxOS"
			} else {
				if ua.os == "" {
					ua.os = normalizeOS(comment[1])
				}
			}
		}
		// Only parse 4th comment as localization if it doesn't start with rv: .
		// For example, Firefox on Ubuntu contains "rv:XX.X" in this field.
		if len(comment) > 3 && !strings.HasPrefix(comment[3], "rv:") {
			ua.localization = comment[3]
		}
	}
}

// Guess the OS, the localization and if this is a mobile device for Internet Explorer.
//
// The first argument p is a reference to the current UserAgent,
// and the second argument is a slice of strings containing the comment.
func trident(ua *UserAgent, comment []string) {
	// Internet Explorer only runs on Windows.
	ua.platform = "Windows"

	// The OS can be set before to handle a new case in IE11.
	if ua.os == "" {
		if len(comment) > 2 {
			ua.os = normalizeOS(comment[2])
		} else {
			ua.os = "Windows NT 4.0"
		}
	}

	// Last but not least, let's detect if it comes from a mobile device.
	for _, v := range comment {
		if strings.HasPrefix(v, "IEMobile") {
			ua.mobile = true
			return
		}
	}
}

// Guess the OS, the localization and if this is a mobile device for Opera.
//
// The first argument p is a reference to the current UserAgent,
// and the second argument is a slice of strings containing the comment.
func opera(ua *UserAgent, comment []string) {
	slen := len(comment)

	if strings.HasPrefix(comment[0], "Windows") {
		ua.platform = "Windows"
		ua.os = normalizeOS(comment[0])
		if slen > 2 {
			if slen > 3 && strings.HasPrefix(comment[2], "MRA") {
				ua.localization = comment[3]
			} else {
				ua.localization = comment[2]
			}
		}
	} else {
		if strings.HasPrefix(comment[0], "Android") {
			ua.mobile = true
		}
		ua.platform = comment[0]
		if slen > 1 {
			ua.os = comment[1]
			if slen > 3 {
				ua.localization = comment[3]
			}
		} else {
			ua.os = comment[0]
		}
	}
}

// Guess the OS.
// Android browsers send Dalvik as the user agent in the request header.
//
// The first argument p is a reference to the current UserAgent,
// and the second argument is a slice of strings containing the comment.
func dalvik(ua *UserAgent, comment []string) {
	slen := len(comment)

	if strings.HasPrefix(comment[0], "Linux") {
		ua.platform = comment[0]
		if slen > 2 {
			ua.os = comment[2]
		}
		ua.mobile = true
	}
}

// Given the comment of the first section of the UserAgent string,
// get the platform.
func getPlatform(comment []string) string {
	if len(comment) > 0 {
		if comment[0] != "compatible" {
			if strings.HasPrefix(comment[0], "Windows") {
				return "Windows"
			} else if strings.HasPrefix(comment[0], "Symbian") {
				return "Symbian"
			} else if strings.HasPrefix(comment[0], "webOS") {
				return "webOS"
			} else if comment[0] == "BB10" {
				return "BlackBerry"
			}
			return comment[0]
		}
	}
	return ""
}

// Detect some properties of the OS from the given section.
func (ua *UserAgent) detectOS(s section) {
	ua.os = ""
	if s.name == "Mozilla" {
		// Get the platform here.
		// Be aware that IE11 provides a new format
		// that is not backwards-compatible with previous versions of IE.
		ua.platform = getPlatform(s.comment)
		if ua.platform == "Windows" && len(s.comment) > 0 {
			ua.os = normalizeOS(s.comment[0])
		}

		// And finally, get the OS depending on the engine.
		switch ua.browser.Engine {
		case "":
			ua.undecided = true
		case "Gecko":
			gecko(ua, s.comment)
		case "AppleWebKit":
			webkit(ua, s.comment)
		case "Trident":
			trident(ua, s.comment)
		}
	} else if s.name == "Opera" {
		if len(s.comment) > 0 {
			opera(ua, s.comment)
		}
	} else if s.name == "Dalvik" {
		if len(s.comment) > 0 {
			dalvik(ua, s.comment)
		}
	} else if s.name == "okhttp" {
		ua.mobile = true
		ua.browser.Name = "OkHttp"
		ua.browser.Version = s.version
	} else {
		// Check whether this is a bot or just a weird browser.
		ua.undecided = true
	}
}

// Return OS name and version from a slice of strings created from the full name of the OS.
func osName(osSplit []string) (name, version string) {
	if len(osSplit) == 1 {
		name = osSplit[0]
		version = ""
	} else {
		// Assume a version is stored in the last part of the array.
		nameSplit := osSplit[:len(osSplit)-1]
		version = osSplit[len(osSplit)-1]

		// Nicer looking Mac OS X
		if len(nameSplit) >= 2 && nameSplit[0] == "Intel" && nameSplit[1] == "Mac" {
			nameSplit = nameSplit[1:]
		}
		name = strings.Join(nameSplit, " ")

		if strings.Contains(version, "x86") || strings.Contains(version, "i686") {
			// x86_64 and i868 are not Linux versions but architectures
			version = ""
		} else if version == "X" && name == "Mac OS" {
			// X is not a version for macOS.
			name = name + " " + version
			version = ""
		}
	}
	return name, version
}

// OSInfo returns combined information for the operating system.
func (ua *UserAgent) OSInfo() OSInfo {
	// Special case for iPhone weirdness
	os := strings.Replace(ua.os, "like Mac OS X", "", 1)
	os = strings.Replace(os, "CPU", "", 1)
	os = strings.Trim(os, " ")

	osSplit := strings.Split(os, " ")

	// Special case for x64 edition of Windows
	if os == "Windows XP x64 Edition" {
		osSplit = osSplit[:len(osSplit)-2]
	}

	name, version := osName(osSplit)

	// Special case for names that contain a forward slash version separator.
	if strings.Contains(name, "/") {
		s := strings.Split(name, "/")
		name = s[0]
		version = s[1]
	}

	// Special case for versions that use underscores
	version = strings.Replace(version, "_", ".", -1)

	return OSInfo{
		FullName: ua.os,
		Name:     name,
		Version:  version,
	}
}
