package weather

import (
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
