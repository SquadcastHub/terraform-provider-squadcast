package api

import (
	"context"
	"fmt"
	"net/http"
)

type Organization struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func (client *Client) GetCurrentOrganization(ctx context.Context) (*Organization, error) {
	url := fmt.Sprintf("%s/organization", client.BaseURLV3)

	return Request[any, Organization](http.MethodGet, url, client, ctx, nil)
}
