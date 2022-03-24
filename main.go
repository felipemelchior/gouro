package main

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type urlmap_t struct {
	host string
	path string
}

var urlmap = make(map[urlmap_t]url.Values)

var extensions = []string{"js", "css", "png", "jpg", "jpeg", "svg", "ico", "webp",
	"ttf", "otf", "woff", "gif", "pdf", "bmp", "eot", "mp3",
	"mp4", "woff2", "avi"}
var known_params []string
var known_patterns []string
var re_content, _ = regexp.Compile(`(post|blog)s?|docs|support/|/(\d{4}|pages?)/\d+/`)
var re_int, _ = regexp.Compile(`/\d+([?/]|$)`)

func isContent(path string) bool {
	for _, part := range strings.Split(path, "/") {
		if strings.Count(part, "-") > 3 {
			return true
		}
	}
	return false
}

func hasBadExtensions(path string) bool {
	if path == "/" {
		return false
	}

	for _, bad_ext := range extensions {
		if strings.HasSuffix(path, "."+bad_ext) {
			return true
		}
	}
	return false
}

func addNewParams(params url.Values) bool {
	var has_new_params bool
	var new_param bool

	if len(known_params) == 0 {
		for unknown_param := range params {
			known_params = append(known_params, unknown_param)
		}
		return true
	}

	has_new_params = false
	for unknown_param := range params {
		new_param = true
		for _, known_param := range known_params {
			if unknown_param == known_param {
				new_param = false
			}
		}

		if new_param {
			known_params = append(known_params, unknown_param)
			has_new_params = true
		}
	}
	return has_new_params
}

func createPattern(path string) {
	var new_parts []string

	for _, part := range strings.Split(regexp.QuoteMeta(path), "/") {
		if _, err := strconv.Atoi(part); err == nil {
			new_parts = append(new_parts, "\\d+")
		} else {
			new_parts = append(new_parts, part)
		}
	}

	addPattern := true
	new_pattern := strings.Join(new_parts, "/")
	for _, pattern := range known_patterns {
		if pattern == new_pattern {
			addPattern = false
			break
		}
	}

	if addPattern {
		known_patterns = append(known_patterns, new_pattern)
	}
}

func matchesPattern(path string) bool {
	for _, pattern := range known_patterns {
		if matched, _ := regexp.MatchString(pattern, path); matched {
			return true
		}
	}
	return false
}

func hostExists(urlmap_to_test urlmap_t) bool {
	for urlmap_aux := range urlmap {
		if urlmap_aux.host == urlmap_to_test.host && urlmap_aux.path == urlmap_to_test.path {
			return true
		}
	}
	return false
}

func compareParams(new_params url.Values, og_params url.Values) bool {
	for param := range new_params {
		if !og_params.Has(param) {
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
		params := parsed.Query()
		has_new_params := addNewParams(params)
		if hasBadExtensions(parsed.Path) || re_content.MatchString(parsed.Path) || isContent(parsed.Path) {
			continue
		}
		urlmap_aux := urlmap_t{
			host: host,
			path: parsed.Path,
		}
		if ((len(params) == 0) || has_new_params) && re_int.MatchString(parsed.Path) {
			if matchesPattern(parsed.Path) {
				continue
			} else {
				createPattern(parsed.Path)
				urlmap[urlmap_aux] = params
			}
		} else if !hostExists(urlmap_aux) {
			urlmap[urlmap_aux] = params
		} else if has_new_params || compareParams(params, urlmap[urlmap_aux]) {
			for param := range params {
				if !urlmap[urlmap_aux].Has(param) {
					urlmap[urlmap_aux].Add(param, params.Get(param))
				}
			}
		}
	}

	for urlmap_aux := range urlmap {
		if len(urlmap[urlmap_aux]) > 0 {
			fmt.Println(urlmap_aux.host + urlmap_aux.path + "?" + urlmap[urlmap_aux].Encode())
		} else {
			fmt.Println(urlmap_aux.host + urlmap_aux.path)
		}
	}
}
