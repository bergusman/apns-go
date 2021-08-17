package apns

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrClientNotificationNil = errors.New("client: notification is nil")
	ErrClientTokenNil        = errors.New("client: token is nil")
)

// Token-based client for sending remote notifications.
type Client struct {
	Token      *Token
	HTTPClient *http.Client
}

// NewClient creates client with token and http.Client client,
// pass nil for client to use http.DefaultClient.
func NewClient(token *Token, httpClient *http.Client) *Client {
	return &Client{
		Token:      token,
		HTTPClient: httpClient,
	}
}

// Push sends remote notification.
func (c *Client) Push(n *Notification) (*Response, error) {
	return c.PushWithContext(context.Background(), n)
}

// Push sends remote notification with context.
func (c *Client) PushWithContext(ctx context.Context, n *Notification) (*Response, error) {
	if n == nil {
		return nil, ErrClientNotificationNil
	}

	req, err := n.BuildRequestWithContext(ctx)
	if err != nil {
		return nil, err
	}

	if c.Token == nil {
		return nil, ErrClientTokenNil
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
