package client

import (
	"bytes"
	"io"
	"net/url"
	"path"
)

func (c *Client) get(_path string) ([]byte, error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, _path)

	resp, err := c.client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	return body, nil
}

func (c *Client) httpPost(_path string, payload []byte) (body []byte, statusCode int, err error) {
	u, err := url.Parse(c.url)
	if err != nil {
		return
	}

	u.Path = path.Join(u.Path, _path)

	resp, err := c.client.Post(u.String(), "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode
	body, err = io.ReadAll(resp.Body)
	return
}
