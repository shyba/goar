package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/liteseed/goar/tx"
)

// arweave HTTP API: https://docs.arweave.org/developers/server/http-api

type Client struct {
	client *http.Client
	url    string
}

func New(node string) *Client {
	httpClient := http.DefaultClient
	return &Client{client: httpClient, url: node}
}

func (c *Client) GetTransaction(id string) (*tx.Transaction, error) {
	body, err := c.get(fmt.Sprintf("transaction/%s", id))
	if err != nil {
		return nil, err
	}

	t := &tx.Transaction{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Client) GetTransactionStatus(id string) (*TransactionStatus, error) {
	body, err := c.get(fmt.Sprintf("transaction/%s/status", id))
	if err != nil {
		return nil, err
	}

	t := &TransactionStatus{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Client) GetTransactionField(id string, field string) (string, error) {
	body, err := c.get(fmt.Sprintf("transaction/%s/%s", id, field))
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) GetTransactionData(id string) ([]byte, error) {
	body, err := c.get(fmt.Sprintf("transaction/%s", id))
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) GetTransactionTags(id string) ([]tx.Tag, error) {
	jsTags, err := c.GetTransactionField(id, "tags")
	if err != nil {
		return nil, err
	}

	tags := make([]tx.Tag, 0)
	if err := json.Unmarshal([]byte(jsTags), &tags); err != nil {
		return nil, err
	}
	return tags, nil
}

func (c *Client) GetTransactionPrice(size int, target *string) (int64, error) {
	url := fmt.Sprintf("price/%d", size)
	if target != nil {
		url = fmt.Sprintf("%v/%v", url, *target)
	}

	body, err := c.get(url)
	if err != nil {
		return 0, err
	}
	price, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		return 0, err
	}
	return price, nil
}

func (c *Client) GetTransactionAnchor() (string, error) {
	body, err := c.get("transaction_anchor")
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) SubmitTransaction(transaction *tx.Transaction) (status string, code int, err error) {
	b, err := json.Marshal(transaction)
	if err != nil {
		return
	}

	body, statusCode, err := c.httpPost("transaction", b)
	status = string(body)
	code = statusCode
	return
}
