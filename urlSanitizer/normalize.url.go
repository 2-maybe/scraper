package urlsanitizer

import (
	"fmt"
	"net"
	"net/url"
	"strings"
)

type protocol int

const (
	https protocol = iota
	http
)

func detectProtocol(data string) protocol {
	if data == "https" {
		return https
	}
	return http
}

func NormalizeUrl(preUrl string) (string, error) {
	urlMetaData, err := url.Parse(preUrl)
	if err != nil {
		return "", fmt.Errorf("err parsing url %w", err)
	}

	if urlMetaData.Scheme != "https" {
		return "", fmt.Errorf("only https allowed")
	}

	// 2. DNS Lookup (Note: This makes your sanitizer SLOW because it hits the network)
	// later cache dns
	_, err = net.LookupHost(urlMetaData.Host)
	if err != nil {
		return "", fmt.Errorf("malformed or unreachable host")
	}

	// 3. Keep Host + Path, drop Fragment (#) and Query (?) if you want unique pages
	// Lowercase the host, but keep the path case-sensitive (most servers care)
	cleanUrl := urlMetaData.Scheme + "://" + strings.ToLower(urlMetaData.Host) + urlMetaData.Path
	cleanUrl = strings.TrimSuffix(cleanUrl, "/") // consistency for 'site.com' vs '://site.com'

	return cleanUrl, nil
}
