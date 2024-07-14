package catapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
)

const BaseURL = "https://api.thecatapi.com"

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c Client) CheckBreed(ctx context.Context, breed string) (formattedBreed string, err error) {
	path := "/v1/breeds/search?attach_image=0&q=" + url.QueryEscape(breed)
	requestURL := BaseURL + path

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
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

	if resp.StatusCode != http.StatusOK {
		errBody, _ := io.ReadAll(resp.Body)
		return "", apperrors.Internal(errors.New(string(errBody))).
			Wrap("status code != 200").WithMetadata("status", resp.StatusCode)
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
