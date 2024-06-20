package client

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/url"
	"path"
)

func (c *Client) get(_path string) ([]byte, error) {
	u, err := url.Parse(c.Gateway)
	if err != nil {
		return nil, err
	}

	u.Path = path.Join(u.Path, _path)

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

func (c *Client) post(_path string, payload []byte) (int, error) {
	u, err := url.Parse(c.Gateway)
	if err != nil {
		return -1, err
	}

	u.Path = path.Join(u.Path, _path)
	log.Println(u.String())
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
