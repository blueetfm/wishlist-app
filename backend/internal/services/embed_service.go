package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/otiai10/opengraph/v2"

	"github.com/blueetfm/wishlist-app/backend/internal/models"
)

const (
	embedFetchTimeout = 5 * time.Second
	embedMaxBodyBytes = 2 << 20 // 2MB, avoids reading unbounded/huge responses
)

// EmbedService scrapes Open Graph metadata from a caller-supplied URL.
type EmbedService struct {
	client *http.Client
}

// NewEmbedService constructs an EmbedService with a hardened HTTP client:
// bounded timeout, and redirect targets are validated the same way as the
// original URL to prevent SSRF via redirect.
func NewEmbedService() *EmbedService {
	return &EmbedService{
		client: &http.Client{
			Timeout: embedFetchTimeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 5 {
					return errors.New("too many redirects")
				}
				return validatePublicURL(req.URL)
			},
		},
	}
}

// Fetch retrieves rawURL and extracts its Open Graph tags.
// Returns ErrInvalidURL if rawURL isn't an allowed http(s) address, ErrUpstreamUnavailable
// if the target couldn't be reached, and ErrNoMetadata if it was reached but
// has no recognizable Open Graph tags.
func (s *EmbedService) Fetch(ctx context.Context, rawURL string) (*models.EmbedResponse, error) {
	// initial checks and error handling for our embedservice's http client
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, ErrInvalidURL
	}
	if err := validatePublicURL(parsed); err != nil {
		return nil, ErrInvalidURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, ErrInvalidURL
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WishlistBot/1.0)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUpstreamUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status %d", ErrUpstreamUnavailable, resp.StatusCode)
	}

	// opengraph.New creates a new OpenGraph struct with specified URL
	// use this instead of opengraph.Fetch so that we can impose a size limit to the fetch
	og := opengraph.New(parsed.String())
	og.Intent.Strict = true // only trust <meta property="og:*">, to match documented ErrNoMetadata contract
	if err := og.Parse(io.LimitReader(resp.Body, embedMaxBodyBytes)); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUpstreamUnavailable, err)
	}
	_ = og.ToAbs() // best-effort: resolve relative og:image URLs against the page URL

	embed := &models.EmbedResponse{
		Title:       og.Title,
		Description: og.Description,
	}
	if len(og.Image) > 0 {
		embed.ImageURL = og.Image[0].URL
	}
	if embed.Title == "" && embed.Description == "" && embed.ImageURL == "" {
		return nil, ErrNoMetadata
	}
	return embed, nil
}

// validatePublicURL prevents against SSRF
func validatePublicURL(u *url.URL) error {
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("scheme must be http or https")
	}
	host := u.Hostname()
	if host == "" {
		return errors.New("missing host")
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return err
	}
	for _, ip := range ips {
		if isDisallowedIP(ip) {
			return fmt.Errorf("host resolves to a disallowed address: %s", ip)
		}
	}
	return nil
}

func isDisallowedIP(ip net.IP) bool {
	return ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() || ip.IsUnspecified() || ip.IsMulticast()
}
