package main

import (
	"fmt"

	md "github.com/houseme/mobiledetect"
)

func main() {
	detect := md.New(nil, nil)

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
