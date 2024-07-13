package catapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
	"net/http"
	"net/url"
	"strings"
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

	if len(body) == 0 {
		return "", apperrors.InvalidCatBreed(breed, "not found")
	}

	if len(body) > 1 { // name is ambiguous
		const limit = 5

		otherBreeds := make([]string, 0, limit)
		for i := 0; i < limit && i < len(body); i++ {
			otherBreeds = append(otherBreeds, body[i].Name)
		}

		return "", apperrors.InvalidCatBreed(breed,
			fmt.Sprintf("is ambiguous, other variants: %s", strings.Join(otherBreeds, ", ")),
		)
	}

	if strings.ToLower(breed) != strings.ToLower(body[0].Name) {
		return "", apperrors.InvalidCatBreed(breed,
			fmt.Sprintf("maybe you meant '%s'?", body[0].Name),
		)
	}

	return body[0].Name, nil
}
