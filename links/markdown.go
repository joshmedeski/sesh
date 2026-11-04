package links

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"codeberg.org/readeck/go-readability/v2"
	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
)

const (
	fetchTimeout = 10 * time.Second
	maxPageBytes = 5 << 20
)

var client = &http.Client{Timeout: fetchTimeout}

// Markdown fetches an http(s) page and returns its main content as markdown.
func Markdown(pageURL string) (string, error) {
	u, err := url.Parse(pageURL)
	if err != nil {
		return "", err
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", fmt.Errorf("unsupported scheme %q", u.Scheme)
	}

	resp, err := client.Get(pageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetching %s: %s", pageURL, resp.Status)
	}
	if !strings.Contains(resp.Header.Get("Content-Type"), "text/html") {
		return "", errors.New("not an HTML page")
	}

	article, err := readability.FromReader(io.LimitReader(resp.Body, maxPageBytes), resp.Request.URL)
	if err != nil {
		return "", err
	}
	if article.Node == nil {
		return "", errors.New("no readable content")
	}
	body, err := htmltomarkdown.ConvertNode(article.Node, converter.WithDomain(resp.Request.URL.String()))
	if err != nil {
		return "", err
	}

	md := strings.TrimSpace(string(body))
	if title := article.Title(); title != "" {
		md = "# " + title + "\n\n" + md
	}
	return md, nil
}
