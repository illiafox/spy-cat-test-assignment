package catapi

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/illiafox/spy-cat-test-assignment/app/internal/apperrors"
)

//go:embed breeds.json
var breedsFile []byte

type Client struct {
	Breeds map[string]struct{}
}

func NewClient() (*Client, error) {
	var breeds []string

	err := json.Unmarshal(breedsFile, &breeds)
	if err != nil {
		return nil, fmt.Errorf("unmarshall breeds from file: %w", err)
	}

	client := &Client{}

	client.Breeds = make(map[string]struct{}, len(breeds))
	for _, breed := range breeds {
		client.Breeds[breed] = struct{}{}
	}

	return client, nil
}

func (c Client) CheckBreed(_ context.Context, breed string) (formattedBreed string, err error) {
	_, ok := c.Breeds[breed]
	if !ok {
		return "", apperrors.InvalidCatBreed(breed, "not found")
	}

	return breed, nil
}
