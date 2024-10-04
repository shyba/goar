package client

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"path"
)

func (c *Client) get(route string) ([]byte, error) {
	u, err := url.Parse(c.Gateway)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, route)

	resp, err := c.Client.Get(u.String())
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) post(route string, payload []byte) (int, error) {
	u, err := url.Parse(c.Gateway)
	if err != nil {
		return -1, err
	}

	u.Path = path.Join(u.Path, route)
	resp, err := c.Client.Post(u.String(), "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return -1, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}
	code := resp.StatusCode
	if code >= 400 {
		return resp.StatusCode, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}
	return code, nil
}
