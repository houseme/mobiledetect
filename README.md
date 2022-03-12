# Go-MobileDetect

### Library for detecting mobile devices and tablets
### Go Mobile Detect is a lightweight Go package imported from PHP for detecting mobile devices (including tablets).

[![Go Reference](https://pkg.go.dev/badge/github.com/houseme/go-mobile-detect.svg)](https://pkg.go.dev/github.com/houseme/go-mobile-detect)
[![Go](https://github.com/houseme/go-mobile-detect/actions/workflows/go.yml/badge.svg)](https://github.com/housemecn/go-mobile-detect/actions/workflows/go.yml)
![GitHub](https://img.shields.io/github/license/houseme/go-mobile-detect?style=plastic)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/houseme/go-mobile-detect/main?style=flat-square)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/houseme/go-mobile-detect?style=flat-square)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/houseme/go-mobile-detect?style=flat-square)

### What is it?

#### It uses the User-Agent string combined with specific HTTP headers to detect the mobile environment.
#### The package is imported from [MobileDetect](http://mobiledetect.net/) which was originally written in PHP.

### Installation 

     go get -u -v github.com/houseme/go-mobiledetect 

### Updates 

#### Version 1.1.2


### Why is it useful?

There are different ways of using the package: 

- [Basic usage](examples/app.go) 
- [Basic router implementation](examples/router/main.go)
- [Handler interface implementation](examples/handler/main.go)
- [Mux interface implementation](examples/mux/main.go)

### License

Go Mobile Detect is an open-source script released under [MIT License](http://www.opensource.org/licenses/mit-license.php). thanks for [Shaked](https://github.com/Shaked/gomobiledetect) and [serbanghita](https://github.com/serbanghita/Mobile-Detect).