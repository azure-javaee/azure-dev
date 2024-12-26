package appdetect

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func download(requestUrl string) ([]byte, error) {
	parsedUrl, err := url.ParseRequestURI(requestUrl)
	if err != nil {
		return nil, err
	}
	if !isAllowedHost(parsedUrl.Host) {
		return nil, fmt.Errorf("invalid host")
	}
	timeOut := 30 * time.Second
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		timeOut = 60 * time.Second
	}
	client := &http.Client{
		Timeout: timeOut,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	resp, err := client.Get(requestUrl)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("failed to close http response body")
		}
	}(resp.Body)
	return io.ReadAll(resp.Body)
}

func isAllowedHost(host string) bool {
	return host == "repo.maven.apache.org"
}
