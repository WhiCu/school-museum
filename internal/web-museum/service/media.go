package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	errUnsupportedMediaURL = errors.New("unsupported media URL")

	rePostDataJSON = regexp.MustCompile(`window\.postDataJSON=\"([\s\S]*?)\"</script>`)
	reOgVideo      = regexp.MustCompile(`<meta\s+property=\"og:video\"[^>]*content=\"([^\"]+)\"`)
	reOgImage      = regexp.MustCompile(`<meta\s+property=\"og:image\"[^>]*content=\"([^\"]+)\"`)
	reTwitterImage = regexp.MustCompile(`<meta\s+name=\"twitter:image\"[^>]*content=\"([^\"]+)\"`)
)

type imgurPostData struct {
	Cover struct {
		URL      string `json:"url"`
		MimeType string `json:"mime_type"`
	} `json:"cover"`
	Media []struct {
		URL      string `json:"url"`
		MimeType string `json:"mime_type"`
	} `json:"media"`
}

func (s *Service) ResolveExternalMedia(ctx context.Context, rawURL string) (string, string, error) {
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed == nil {
		return "", "", fmt.Errorf("invalid url")
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", "", fmt.Errorf("unsupported url scheme")
	}

	host := strings.ToLower(strings.TrimPrefix(parsed.Hostname(), "www."))
	if host != "imgur.com" {
		return "", "", errUnsupportedMediaURL
	}

	return s.resolveImgur(ctx, parsed.String())
}

func (s *Service) resolveImgur(ctx context.Context, sourceURL string) (string, string, error) {
	html, err := s.fetchPage(ctx, sourceURL)
	if err != nil {
		return "", "", err
	}

	if u, t, ok := extractFromPostDataJSON(html); ok {
		return u, t, nil
	}

	if u, ok := firstSubmatch(reOgVideo, html); ok {
		return u, "video", nil
	}

	if u, ok := firstSubmatch(reTwitterImage, html); ok {
		return strings.ReplaceAll(u, "?fbplay", ""), "image", nil
	}

	if u, ok := firstSubmatch(reOgImage, html); ok {
		return strings.ReplaceAll(u, "?fbplay", ""), "image", nil
	}

	return "", "", fmt.Errorf("failed to resolve imgur media")
}

func (s *Service) fetchPage(ctx context.Context, sourceURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "school-museum/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("external source returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 3*1024*1024))
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func extractFromPostDataJSON(html string) (string, string, bool) {
	m := rePostDataJSON.FindStringSubmatch(html)
	if len(m) < 2 {
		return "", "", false
	}

	decoded, err := strconv.Unquote(`"` + m[1] + `"`)
	if err != nil {
		return "", "", false
	}

	var data imgurPostData
	if err := json.Unmarshal([]byte(decoded), &data); err != nil {
		return "", "", false
	}

	if len(data.Media) > 0 && data.Media[0].URL != "" {
		return data.Media[0].URL, classifyMedia(data.Media[0].URL, data.Media[0].MimeType), true
	}

	if data.Cover.URL != "" {
		return data.Cover.URL, classifyMedia(data.Cover.URL, data.Cover.MimeType), true
	}

	return "", "", false
}

func classifyMedia(sourceURL, mimeType string) string {
	mt := strings.ToLower(mimeType)
	if strings.HasPrefix(mt, "video/") {
		return "video"
	}

	ext := strings.ToLower(path.Ext(sourceURL))
	switch ext {
	case ".mp4", ".webm", ".mov", ".m4v", ".ogv", ".ogg":
		return "video"
	default:
		return "image"
	}
}

func firstSubmatch(re *regexp.Regexp, body string) (string, bool) {
	m := re.FindStringSubmatch(body)
	if len(m) < 2 {
		return "", false
	}
	return m[1], true
}
