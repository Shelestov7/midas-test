package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func NewHTTPClient(ctx context.Context, apiURL string, httpClient http.Client) *HTTPClient {
	return &HTTPClient{
		apiURL:     apiURL,
		httpClient: httpClient,
		ctx:        ctx,
	}
}

type HTTPClient struct {
	apiURL     string
	httpClient http.Client
	ctx        context.Context
}

func (client *HTTPClient) UserSimpleProtocolList(id string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/v1/user/simple_protocol_list?id=%s", client.apiURL, id)
	req, err := http.NewRequestWithContext(client.ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request error: %s", err.Error())
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %s", err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body error: %s", err.Error())
	}
	return body, nil
}
