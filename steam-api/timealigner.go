package steam

import (
	"encoding/json"
	"net/url"
	"time"
	"github.com/lucasmbaia/goskins/steam-api/request"
)

const (
	alignTime   = "https://api.steampowered.com/ITwoFactorService/QueryTime/v0001?"
)

type timeSyncResponse struct {
	Response struct {
		ServerTime                 timestamp `json:"server_time"`
		SkewTolerence              seconds   `json:"skew_tolerance_seconds"`
		LargeTimeJink              seconds   `json:"large_time_jink"`
		ProbeFrequency             seconds   `json:"probe_frequency_seconds"`
		AdjustedTimeProbeFrequency seconds   `json:"adjusted_time_probe_frequency_seconds"`
		HintProbeFrequency         seconds   `json:"hint_probe_frequency_seconds"`
		SyncTimeout                seconds   `json:"sync_timeout"`
		RetryDelay                 seconds   `json:"try_again_seconds"`
		MaxAttempts                int       `json:"max_attempts"`
	} `json:"response"`
}

var (
	aligned		bool
	timeDifference	time.Duration
)

func (s *Session) GetSteamTime() time.Time {
	if !aligned {
		s.AlignTime()
	}

	return time.Now().Add(timeDifference)
}

func (s *Session) AlignTime() (err error) {
	var (
		resp	request.Response
		tr	timeSyncResponse
		params	url.Values
	)

	params = url.Values{
		"steamid": []string{"0"},
	}

	if resp, err = s.client.Request(request.GET, alignTime + params.Encode(), request.Options{
		Headers:    map[string]string{
			"User-Agent":	    httpUserAgentValue,
			"Accept":	    httpAcceptValue,
		},
	}); err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body, tr); err != nil {
		return
	}

	timeDifference = tr.Response.ServerTime.Sub(time.Now())
	aligned = true

	return
}
