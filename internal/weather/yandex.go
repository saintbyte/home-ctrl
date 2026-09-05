package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/saintbyte/home-ctrl/internal/config"
)

const yandexWeatherURL = "https://api.weather.yandex.ru/v2/forecast"

// Snapshot is an in-memory copy of the last successful weather response
type Snapshot struct {
	Data       json.RawMessage `json:"data"`
	LastUpdate time.Time       `json:"last_update"`
	Error      string          `json:"error,omitempty"`
}

// options holds the resolved Yandex Weather parameters for a task
type options struct {
	key   string
	lat   string
	lon   string
	lang  string
	limit int
	hours bool
	extra bool
}

// NewYandexWeather downloads weather data from the Yandex Weather API and
// keeps the latest response in memory.
type YandexWeather struct {
	client *http.Client
	opts   options

	mu   sync.RWMutex
	snap Snapshot
}

// NewYandexWeather creates a new Yandex Weather service for the given task.
// The task params may specify key, lat, lon (required) and lang (default
// ru_RU), limit (default 2), hours (default true), extra (default true).
func NewYandexWeather(task config.Task) *YandexWeather {
	return &YandexWeather{
		client: &http.Client{Timeout: 30 * time.Second},
		opts:   resolveOptions(task.Params),
	}
}

func resolveOptions(params map[string]any) options {
	opts := options{
		lang:  "ru_RU",
		limit: 2,
		hours: true,
		extra: true,
	}
	if v := config.StringParam(params, "key"); v != "" {
		opts.key = v
	}
	if v := config.StringParam(params, "lat"); v != "" {
		opts.lat = v
	}
	if v := config.StringParam(params, "lon"); v != "" {
		opts.lon = v
	}
	if v := config.StringParam(params, "lang"); v != "" {
		opts.lang = v
	}
	if v := config.IntParam(params, "limit"); v > 0 {
		opts.limit = v
	}
	if v, ok := config.BoolParam(params, "hours"); ok {
		opts.hours = v
	}
	if v, ok := config.BoolParam(params, "extra"); ok {
		opts.extra = v
	}
	return opts
}

// Enabled reports whether weather downloading is configured.
func (w *YandexWeather) Enabled() bool {
	return w.opts.key != ""
}

// BuildURL returns the request URL for the current configuration.
func (w *YandexWeather) buildURL() string {
	q := url.Values{}
	q.Set("lat", w.opts.lat)
	q.Set("lon", w.opts.lon)
	q.Set("lang", w.opts.lang)
	if w.opts.limit > 0 {
		q.Set("limit", strconv.Itoa(w.opts.limit))
	}
	q.Set("hours", strconv.FormatBool(w.opts.hours))
	q.Set("extra", strconv.FormatBool(w.opts.extra))
	return yandexWeatherURL + "?" + q.Encode()
}

// Fetch downloads weather data and stores it in memory.
func (w *YandexWeather) Fetch() error {
	data, err := w.fetch()
	w.mu.Lock()
	defer w.mu.Unlock()
	if err != nil {
		w.snap.Error = err.Error()
		return err
	}
	w.snap = Snapshot{
		Data:       json.RawMessage(data),
		LastUpdate: time.Now(),
	}
	return nil
}

func (w *YandexWeather) fetch() ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, w.buildURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("yandex weather: build request: %w", err)
	}
	req.Header.Set("X-Yandex-Weather-Key", w.opts.key)

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yandex weather: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yandex weather: read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("yandex weather: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if !json.Valid(body) {
		return nil, fmt.Errorf("yandex weather: invalid json response")
	}

	return body, nil
}

// Get returns a copy of the latest weather snapshot.
func (w *YandexWeather) Get() Snapshot {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.snap
}
