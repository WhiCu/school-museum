package telemetry

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type UmamiClient struct {
	URL       string
	WebsiteID string
}

func NewUmami(url, websiteID string) *UmamiClient {
	return &UmamiClient{URL: url, WebsiteID: websiteID}
}

func (u *UmamiClient) Track(r *http.Request, title string) {
	go func() {
		// Извлекаем IP из заголовков (порядок важен)
		ip := r.Header.Get("X-Forwarded-For")
		if ip == "" {
			ip = r.Header.Get("X-Real-IP")
		}
		if ip == "" {
			ip = r.RemoteAddr
		}

		payload := map[string]interface{}{
			"type": "event",
			"payload": map[string]interface{}{
				"website":   u.WebsiteID,
				"url":       r.URL.Path,
				"hostname":  r.Host,
				"title":     title,
				"language":  r.Header.Get("Accept-Language"),
				"referrer":  r.Referer(),
				"screen":    "1920x1080",   // можно сделать конфигурируемым
				"ip":        ip,            // ← передаём IP в payload [citation:4]
				"userAgent": r.UserAgent(), // ← User-Agent в payload [citation:4]
			},
		}

		// Важно: НЕ передаём эти заголовки в HTTP-запросе!
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", u.URL+"/api/send", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		// User-Agent и X-Forwarded-For убираем!

		http.DefaultClient.Do(req)
	}()
}
