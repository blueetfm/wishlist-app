package embed

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

	"github.com/blueetfm/wishlist-app/backend/internal/apperr"
)

const (
	fetchTimeout  = 5 * time.Second
	maxBodyBytes  = 2 << 20 // 2MB, avoids reading unbounded/huge responses
)

// Service scrapes Open Graph metadata from a caller-supplied URL.
type Service struct {
	client *http.Client
}

// NewService constructs a Service with a hardened HTTP client: bounded
// timeout, and redirect targets are validated the same way as the original
// URL to prevent SSRF via redirect.
func NewService() *Service {
	return &Service{
		client: &http.Client{
			Timeout: fetchTimeout,
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
func (s *Service) Fetch(ctx context.Context, rawURL string) (*EmbedResponse, error) {
	// initial checks and error handling for our embed service's http client
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return nil, apperr.ErrInvalidURL
	}
	if err := validatePublicURL(parsed); err != nil {
		return nil, apperr.ErrInvalidURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsed.String(), nil)
	if err != nil {
		return nil, apperr.ErrInvalidURL
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; WishlistBot/1.0)")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", apperr.ErrUpstreamUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%w: status %d", apperr.ErrUpstreamUnavailable, resp.StatusCode)
	}

	// opengraph.New creates a new OpenGraph struct with specified URL
	// use this instead of opengraph.Fetch so that we can impose a size limit to the fetch
	og := opengraph.New(parsed.String())
	og.Intent.Strict = true // only trust <meta property="og:*">, to match documented ErrNoMetadata contract
	if err := og.Parse(io.LimitReader(resp.Body, maxBodyBytes)); err != nil {
		return nil, fmt.Errorf("%w: %v", apperr.ErrUpstreamUnavailable, err)
	}
	_ = og.ToAbs() // best-effort: resolve relative og:image URLs against the page URL

	embed := &EmbedResponse{
		Title:       og.Title,
		Description: og.Description,
	}
	if len(og.Image) > 0 {
		embed.ImageURL = og.Image[0].URL
	}
	if embed.Title == "" && embed.Description == "" && embed.ImageURL == "" {
		return nil, apperr.ErrNoMetadata
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
