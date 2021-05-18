package util

import (
	"fmt"
	neturl "net/url"
	"strconv"
	"strings"
)

func StringMapGetDefault(m map[string]string, k string, def string) string {
	v, ok := m[k]
	if !ok {
		v = def
	}
	return v
}

func StringFromJson(m map[string]interface{}, k string) (string, error) {

	sInterface, ok := m[k]
	if !ok {
		return "", fmt.Errorf("Failed to find key %v", k)
	}

	s, ok := sInterface.(string)
	if !ok {
		return "", fmt.Errorf("Failed to convert key %v to string", k)
	}

	return s, nil
}

func MapInterfaceFromJson(m map[string]interface{}, k string) (map[string]interface{}, error) {

	m2Interface, ok := m[k]
	if !ok {
		return nil, fmt.Errorf("Failed to find key %v", k)
	}

	m2, ok := m2Interface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Failed to convert key %v to string", k)
	}

	return m2, nil
}

func SanitizeUrl(url string) (string, error) {
	// Remove any leading and trailing spaces
	url = strings.TrimSpace(url)

	// See if the user put any quotes around it
	uqs, err := strconv.Unquote(url)
	if err == nil {
		url = uqs
	}

	// Make sure the URL isn't horibbly misformatted
	u, err := neturl.Parse(url)
	if err != nil {
		return "", fmt.Errorf("Error parsing url %v: %v", url, err)
	}

	// In case user forgets to put http:// on the front of the url
	if u.Scheme != "http" && u.Scheme != "https" {
		url = "http://" + url
	}

	return url, nil
}
