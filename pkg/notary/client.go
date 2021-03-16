package notary

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"time"
	"vax/pkg/model/dto"
)

var (
	NOTARY_HOST              string
	NOTARY_PORT              string
	httpClient               *http.Client
	ErrEmptyHash             = errors.New("empty hash")
	ErrEmptyId               = errors.New("empty id")
	ErrInvalidServerResponse = errors.New("invalid server response")
)

const (
	eventsPath = "/event"
)

type createEventRequest struct {
	Hash string `json:"hash"`
}

func init() {
	if NOTARY_HOST = os.Getenv("NOTARY_HOST"); NOTARY_HOST == "" {
		NOTARY_HOST = "localhost"
	}
	if NOTARY_PORT = os.Getenv("NOTARY_PORT"); NOTARY_PORT == "" {
		NOTARY_PORT = "8080"
	}

	httpClient = &http.Client{
		Timeout: 5 * time.Second,
	}
}

func CreateEvent(hash string) (*dto.Event, error) {
	if hash == "" {
		return nil, ErrEmptyHash
	}

	body := createEventRequest{
		Hash: hash,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, getBaseEventsPath(), bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var event dto.Event
	if err := json.NewDecoder(res.Body).Decode(&event); err != nil {
		return nil, err
	}

	return &event, nil
}

func GetEvent(id string) (*dto.Event, error) {
	if id == "" {
		return nil, ErrEmptyId
	}

	req, err := http.NewRequest(http.MethodGet, getBaseEventsPath()+"/"+url.PathEscape(id), nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, ErrInvalidServerResponse
	}

	var event dto.Event
	if err := json.NewDecoder(res.Body).Decode(&event); err != nil {
		return nil, err
	}

	return &event, nil
}

func getBaseEventsPath() string {
	return "http://" + NOTARY_HOST + ":" + NOTARY_PORT + "/api/v1" + eventsPath
}
