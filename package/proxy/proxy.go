package proxy

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"golang.org/x/exp/rand"
)

// Proxy represents a single proxy entry from the JSON file
type Proxy struct {
	IP        string   `json:"ip"`
	Port      string   `json:"port"`
	Protocols []string `json:"protocols"`
}

// LoadProxies loads proxy data from the JSON file
func loadProxies() ([]Proxy, error) {
	file, err := os.ReadFile("proxys/proxys.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var proxies []Proxy
	err = json.Unmarshal(file, &proxies)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return proxies, nil
}

// GetRandomProxy returns a random proxy URL from the list of proxies loaded from the JSON file.
//
// If there are no proxies available, it returns nil.
//
// If there is an error loading the proxies, it will print the error to the console and return nil.
func GetRandomProxy() *url.URL {
	var proxyList, err = loadProxies()

	if err != nil {
		fmt.Printf("Error loading proxies: %v\n", err)
		return nil
	}

	if len(proxyList) == 0 {
		return nil
	}

	proxy := proxyList[rand.Intn(len(proxyList))]

	proxyURL, err := url.Parse(fmt.Sprintf("%s://%s:%s", proxy.Protocols[0], proxy.IP, proxy.Port))

	if err != nil {
		fmt.Printf("Error parsing proxy URL: %v\n", err)
		return nil
	}
	return proxyURL
}
