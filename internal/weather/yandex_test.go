package weather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/saintbyte/home-ctrl/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParamDefaults(t *testing.T) {
	task := config.Task{
		Params: map[string]any{
			"key": "key",
			"lat": "55.75",
			"lon": "37.62",
		},
	}
	w := NewYandexWeather(task)

	assert.True(t, w.Enabled())
	assert.Equal(t, "ru_RU", w.opts.lang)
	assert.Equal(t, 2, w.opts.limit)
	assert.True(t, w.opts.hours)
	assert.True(t, w.opts.extra)
}

func TestParamOverrides(t *testing.T) {
	task := config.Task{
		Params: map[string]any{
			"key":   "key",
			"lat":   "55.75",
			"lon":   "37.62",
			"lang":  "en_US",
			"limit": 5,
			"hours": false,
			"extra": false,
		},
	}
	w := NewYandexWeather(task)

	assert.Equal(t, "en_US", w.opts.lang)
	assert.Equal(t, 5, w.opts.limit)
	assert.False(t, w.opts.hours)
	assert.False(t, w.opts.extra)
}

func TestEnabled(t *testing.T) {
	assert.False(t, NewYandexWeather(config.Task{}).Enabled())
}

func TestBuildURL(t *testing.T) {
	task := config.Task{
		Params: map[string]any{
			"key":   "key",
			"lat":   "55.75",
			"lon":   "37.62",
			"lang":  "ru_RU",
			"limit": 2,
			"hours": true,
			"extra": true,
		},
	}
	w := NewYandexWeather(task)
	raw := w.buildURL()
	u, err := url.Parse(raw)
	require.NoError(t, err)
	assert.Equal(t, yandexWeatherURL, u.Scheme+"://"+u.Host+u.Path)
	q := u.Query()
	assert.Equal(t, "55.75", q.Get("lat"))
	assert.Equal(t, "37.62", q.Get("lon"))
	assert.Equal(t, "ru_RU", q.Get("lang"))
	assert.Equal(t, "2", q.Get("limit"))
	assert.Equal(t, "true", q.Get("hours"))
	assert.Equal(t, "true", q.Get("extra"))
}

func TestGetReturnsEmptySnapshot(t *testing.T) {
	w := NewYandexWeather(config.Task{})
	snap := w.Get()
	assert.True(t, snap.LastUpdate.IsZero())
	assert.Nil(t, snap.Data)
	assert.Empty(t, snap.Error)
}

func TestFetchSuccess(t *testing.T) {
	payload := map[string]any{"temp": -5, "condition": "cloudy"}
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		json.NewEncoder(rw).Encode(payload)
	}))
	defer srv.Close()

	task := config.Task{
		Params: map[string]any{
			"key":   "test-key",
			"lat":   "55.75",
			"lon":   "37.62",
		},
	}
	w := newTestYandexWeather(task, srv.URL)

	err := w.Fetch()
	require.NoError(t, err)

	snap := w.Get()
	assert.False(t, snap.LastUpdate.IsZero())
	assert.Empty(t, snap.Error)
	assert.NotNil(t, snap.Data)

	var got map[string]any
	require.NoError(t, json.Unmarshal(snap.Data, &got))
	assert.Equal(t, float64(-5), got["temp"])
	assert.Equal(t, "cloudy", got["condition"])
}

func TestFetchSetsErrorOnNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusUnauthorized)
		rw.Write([]byte("unauthorized"))
	}))
	defer srv.Close()

	task := config.Task{
		Params: map[string]any{
			"key":   "bad-key",
			"lat":   "55.75",
			"lon":   "37.62",
		},
	}
	w := newTestYandexWeather(task, srv.URL)

	err := w.Fetch()
	assert.Error(t, err)

	snap := w.Get()
	assert.Contains(t, snap.Error, "401")
}

func TestFetchSetsErrorOnInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("not json at all"))
	}))
	defer srv.Close()

	task := config.Task{
		Params: map[string]any{
			"key":   "test-key",
			"lat":   "55.75",
			"lon":   "37.62",
		},
	}
	w := newTestYandexWeather(task, srv.URL)

	err := w.Fetch()
	assert.Error(t, err)

	snap := w.Get()
	assert.Contains(t, snap.Error, "invalid json")
}

func TestFetchSetsErrorOnNetworkFailure(t *testing.T) {
	task := config.Task{
		Params: map[string]any{
			"key":   "test-key",
			"lat":   "55.75",
			"lon":   "37.62",
		},
	}
	w := newTestYandexWeather(task, "http://localhost:1")

	err := w.Fetch()
	assert.Error(t, err)

	snap := w.Get()
	assert.NotEmpty(t, snap.Error)
}

func TestFetchSendsAuthHeader(t *testing.T) {
	var gotKey string
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		gotKey = req.Header.Get("X-Yandex-Weather-Key")
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{}`))
	}))
	defer srv.Close()

	task := config.Task{
		Params: map[string]any{
			"key":   "my-secret-key",
			"lat":   "55.75",
			"lon":   "37.62",
		},
	}
	w := newTestYandexWeather(task, srv.URL)

	_ = w.Fetch()
	assert.Equal(t, "my-secret-key", gotKey)
}

func TestFetchSetsCorrectQueryParams(t *testing.T) {
	var gotQuery url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		gotQuery = req.URL.Query()
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte(`{}`))
	}))
	defer srv.Close()

	task := config.Task{
		Params: map[string]any{
			"key":   "key",
			"lat":   "61.25",
			"lon":   "73.38",
			"lang":  "en_US",
			"limit": 3,
			"hours": false,
			"extra": true,
		},
	}
	w := newTestYandexWeather(task, srv.URL)

	_ = w.Fetch()
	assert.Equal(t, "61.25", gotQuery.Get("lat"))
	assert.Equal(t, "73.38", gotQuery.Get("lon"))
	assert.Equal(t, "en_US", gotQuery.Get("lang"))
	assert.Equal(t, "3", gotQuery.Get("limit"))
	assert.Equal(t, "false", gotQuery.Get("hours"))
	assert.Equal(t, "true", gotQuery.Get("extra"))
}

func TestFetchOverwritesPreviousSnapshot(t *testing.T) {
	call := 0
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		call++
		rw.WriteHeader(http.StatusOK)
		if call == 1 {
			rw.Write([]byte(`{"temp": -10}`))
		} else {
			rw.Write([]byte(`{"temp": 5}`))
		}
	}))
	defer srv.Close()

	task := config.Task{
		Params: map[string]any{
			"key":   "key",
			"lat":   "55.75",
			"lon":   "37.62",
		},
	}
	w := newTestYandexWeather(task, srv.URL)

	_ = w.Fetch()
	snap1 := w.Get()
	var d1 map[string]any
	_ = json.Unmarshal(snap1.Data, &d1)
	assert.Equal(t, float64(-10), d1["temp"])

	_ = w.Fetch()
	snap2 := w.Get()
	var d2 map[string]any
	_ = json.Unmarshal(snap2.Data, &d2)
	assert.Equal(t, float64(5), d2["temp"])
}

// newTestYandexWeather creates a YandexWeather with a custom base URL for testing.
func newTestYandexWeather(task config.Task, baseURL string) *YandexWeather {
	w := NewYandexWeather(task)
	// Override buildURL by manipulating opts to return a URL pointing to the test server.
	// We do this by replacing the constant in buildURL's output.
	origClient := w.client
	_ = origClient
	w.client = &http.Client{}
	// Monkey-patch: we can't change the const, so we wrap via a custom transport.
	w.client.Transport = &testTransport{base: baseURL, orig: w.opts}
	return w
}

// testTransport rewrites requests to the test server URL, preserving query params.
type testTransport struct {
	base string
	orig options
}

func (t *testTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	u, err := url.Parse(t.base)
	if err != nil {
		return nil, err
	}
	u.RawQuery = req.URL.RawQuery
	req.URL = u
	return http.DefaultTransport.RoundTrip(req)
}
