package catapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"net/http"
	"net/url"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c Client) CheckBreed(ctx context.Context, breed string) (formattedBreed string, err error) {

	requestURL := fmt.Sprintf("https://api.thecatapi.com/v1/breeds/search?q=%s&attach_image=0",
		url.QueryEscape(breed),
	)

	req, err := http.NewRequestWithContext(ctx,
		"GET", requestURL,
		nil,
	)

	if err != nil {
		return "", apperrors.Internal(err).Wrap("http: new request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", apperrors.Internal(err).Wrap("http: default client: do")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	var body []struct {
		Name string `json:"name"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", apperrors.Internal(err).Wrap("json: unmarshall body")
	}

	if len(body) != 1 { // 0 if breed not found, > 1 if name is ambiguous
		return "", apperrors.InvalidCatBreed(breed)
	}

	return body[0].Name, nil
}
