package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"
)

// UmamiClient отправляет события аналитики в Umami.
// Использует ограниченный пул воркеров, чтобы не расти неограниченно под нагрузкой.
type UmamiClient struct {
	url       string
	websiteID string
	client    *http.Client
	log       *slog.Logger
	queue     chan trackEvent
}

type trackEvent struct {
	path      string
	hostname  string
	title     string
	language  string
	referrer  string
	userAgent string
	ip        string
}

// UmamiOpts содержит параметры для создания клиента Umami.
type UmamiOpts struct {
	URL       string
	WebsiteID string
	Username  string
	Password  string
	Domain    string
}

const (
	// Максимальное количество событий в очереди.
	umamiBufSize = 1024
	// Количество воркеров, отправляющих события.
	umamiWorkers = 4
)

// NewUmamiClient создаёт клиент для Umami.
// Если WebsiteID пуст — автоматически ищет или создаёт website через API.
func NewUmamiClient(ctx context.Context, opts UmamiOpts, log *slog.Logger) (*UmamiClient, error) {
	c := &UmamiClient{
		url: opts.URL,
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        umamiWorkers * 2,
				MaxIdleConnsPerHost: umamiWorkers * 2,
				IdleConnTimeout:     30 * time.Second,
			},
		},
		log:   log,
		queue: make(chan trackEvent, umamiBufSize),
	}

	websiteID := opts.WebsiteID
	if websiteID == "" {
		id, err := c.ensureWebsite(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("umami: auto-create website: %w", err)
		}
		websiteID = id
	}
	c.websiteID = websiteID

	c.client.Timeout = 5 * time.Second

	// Запускаем пул воркеров
	for range umamiWorkers {
		go c.worker(ctx)
	}

	return c, nil
}

// worker читает события из канала и отправляет их в Umami.
func (u *UmamiClient) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-u.queue:
			if !ok {
				return
			}
			u.send(ev)
		}
	}
}

// ensureWebsite логинится в Umami, ищет сайт по домену; если не найден — создаёт.
// Повторяет попытки, т.к. Umami может быть ещё не готов при старте.
func (u *UmamiClient) ensureWebsite(ctx context.Context, opts UmamiOpts) (string, error) {
	const maxRetries = 10
	var lastErr error

	for i := range maxRetries {
		if i > 0 {
			u.log.Info("umami: waiting for service...", slog.Int("attempt", i+1))
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(3 * time.Second):
			}
		}

		token, err := u.login(ctx, opts.Username, opts.Password)
		if err != nil {
			lastErr = fmt.Errorf("login: %w", err)
			continue
		}

		id, err := u.findWebsite(ctx, token, opts.Domain)
		if err != nil {
			lastErr = fmt.Errorf("find website: %w", err)
			continue
		}
		if id != "" {
			u.log.Info("umami: found existing website", slog.String("id", id), slog.String("domain", opts.Domain))
			return id, nil
		}

		id, err = u.createWebsite(ctx, token, opts.Domain)
		if err != nil {
			lastErr = fmt.Errorf("create website: %w", err)
			continue
		}
		u.log.Info("umami: created website", slog.String("id", id), slog.String("domain", opts.Domain))
		return id, nil
	}

	return "", fmt.Errorf("after %d attempts: %w", maxRetries, lastErr)
}

func (u *UmamiClient) login(ctx context.Context, username, password string) (string, error) {
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.url+"/api/auth/login", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := u.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Token, nil
}

func (u *UmamiClient) findWebsite(ctx context.Context, token, domain string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.url+"/api/websites", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := u.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		Data []struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, w := range result.Data {
		if w.Domain == domain {
			return w.ID, nil
		}
	}
	return "", nil
}

func (u *UmamiClient) createWebsite(ctx context.Context, token, domain string) (string, error) {
	body, _ := json.Marshal(map[string]string{"domain": domain, "name": domain})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.url+"/api/websites", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := u.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ID, nil
}

// Track ставит событие в очередь для асинхронной отправки.
// Если очередь полна — событие отбрасывается (back-pressure).
func (u *UmamiClient) Track(r *http.Request, eventName string) {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.Header.Get("X-Real-IP")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}

	hostname := r.Host
	if h, _, err := net.SplitHostPort(hostname); err == nil {
		hostname = h
	}

	ev := trackEvent{
		path:      r.URL.Path,
		hostname:  hostname,
		title:     eventName,
		language:  r.Header.Get("Accept-Language"),
		referrer:  r.Referer(),
		userAgent: r.UserAgent(),
		ip:        ip,
	}

	select {
	case u.queue <- ev:
	default:
		// Очередь полна — сбрасываем событие, чтобы не блокировать запрос
	}
}

// send выполняет HTTP-отправку одного события в Umami.
func (u *UmamiClient) send(ev trackEvent) {
	payload := map[string]any{
		"type": "event",
		"payload": map[string]any{
			"website":  u.websiteID,
			"url":      ev.path,
			"hostname": ev.hostname,
			"title":    ev.title,
			"language": ev.language,
			"referrer": ev.referrer,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		u.log.Warn("umami: failed to marshal payload", slog.String("error", err.Error()))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.url+"/api/send", bytes.NewReader(body))
	if err != nil {
		u.log.Warn("umami: failed to create request", slog.String("error", err.Error()))
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", ev.userAgent)
	if ev.ip != "" {
		req.Header.Set("X-Forwarded-For", ev.ip)
	}

	resp, err := u.client.Do(req)
	if err != nil {
		u.log.Warn("umami: failed to send event", slog.String("error", err.Error()))
		return
	}
	resp.Body.Close()
	if resp.StatusCode >= 300 {
		u.log.Warn("umami: unexpected response",
			slog.Int("status", resp.StatusCode),
			slog.String("url", ev.path))
	}
}

// Middleware возвращает HTTP middleware, которое автоматически трекает каждый запрос в Umami.
func (u *UmamiClient) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u.Track(r, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
