package apns

import (
	"context"
	"errors"
	"net/http"
)

type Client struct {
	Token      *Token
	HTTPClient *http.Client
}

func NewClient(token *Token, httpClient *http.Client) *Client {
	return &Client{
		Token:      token,
		HTTPClient: httpClient,
	}
}

func (c *Client) Push(n *Notification) (*Response, error) {
	return c.PushWithContext(context.Background(), n)
}

func (c *Client) PushWithContext(ctx context.Context, n *Notification) (*Response, error) {
	if n == nil {
		return nil, errors.New("notification is nil")
	}

	req, err := n.BuildRequestWithContext(ctx)
	if err != nil {
		return nil, err
	}

	err = c.Token.SetAuthorization(req.Header)
	if err != nil {
		return nil, err
	}

	httpClient := c.HTTPClient
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return ParseResponse(res)
}
