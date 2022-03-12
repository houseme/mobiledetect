package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/houseme/mobiledetect"
)

func handler(w http.ResponseWriter, r *http.Request) {
	detect := mobiledetect.New(r, nil)
	requestValue := r.URL.Query().Get("r")
	fmt.Fprintln(w, "isMobile?", detect.IsMobile())
	fmt.Fprintln(w, "isTablet?", detect.IsTablet())
	fmt.Fprintln(w, "is(request)?", requestValue, " ", detect.Is(requestValue))
	fmt.Fprintln(w, "isKey(request)?", requestValue, " ", detect.IsKey(mobiledetect.IPHONE))
	fmt.Fprintln(w, "Version: ", detect.Version(requestValue))
	fmt.Fprintln(w, "VersionKey: ", detect.Version(mobiledetect.PropIphone))
	fmt.Fprintln(w, "VersionFloat: ", detect.Version(requestValue))
	fmt.Fprintln(w, "VersionFloatKey: ", detect.Version(mobiledetect.PropIphone))
	// Any mobile device (phones or tablets).
	fmt.Println(detect.IsMobile())

	// Any tablet device.
	fmt.Println(detect.IsTablet())

	// Exclude tablets.
	fmt.Println(detect.IsMobile() && !detect.IsTablet())

	// Check for a specific platform with the help of the magic methods:
	fmt.Println(detect.Is("iPhone"))

	// Alternative method is() for checking specific properties.
	// WARNING: this method is in BETA, some keyword properties will change in the future.
	fmt.Println(detect.Is("Chrome"))
	fmt.Println(detect.Is("iOS"))
	fmt.Println(detect.Is("UC Browser"))

	// Get the version() of components.
	// WARNING: this method is in BETA, some keyword properties will change in the future.
	fmt.Println(detect.VersionFloat("Android"))
}

func main() {
	log.Println("Starting local server http://localhost:10001/check (cmd+click to open from terminal)")
	http.HandleFunc("/check", handler)
	http.ListenAndServe(":10001", nil)
}
