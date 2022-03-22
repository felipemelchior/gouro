package main

import (
	"bufio"
	"net/url"
	"os"
)

var extensions = []string{"js", "css", "png", "jpg", "jpeg", "svg", "ico", "webp",
	"ttf", "otf", "woff", "gif", "pdf", "bmp", "eot", "mp3",
	"mp4", "woff2", "avi"}

var urlmap = make(map[string]string)
var known_params []string

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()

		parsed, _ := url.Parse(line)
		host := parsed.Scheme + "://" + parsed.Host
		if _, found := urlmap[host]; !found {
			urlmap[host] = ""
		}
		params := parsed.Query()

		for param, _ := range params {
			if !contains(known_params, param) {
				known_params = append(known_params, param)
			}
		}
	}
}
