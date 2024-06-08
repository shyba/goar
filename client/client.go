package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/liteseed/goar/types"
)

// arweave HTTP API: https://docs.arweave.org/developers/server/http-api

type Client struct {
	Client  *http.Client
	Gateway string
}

func New(gateway string) *Client {
	return &Client{
		Client:  &http.Client{Timeout: time.Second * 10},
		Gateway: gateway,
	}
}

func (c *Client) GetTransaction(id string) (*types.Transaction, error) {
	body, err := c.get(fmt.Sprintf("tx/%s", id))
	if err != nil {
		return nil, err
	}
	t := &types.Transaction{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Client) GetTransactionStatus(id string) (*types.TransactionStatus, error) {
	body, err := c.get(fmt.Sprintf("tx/%s/status", id))
	if err != nil {
		return nil, err
	}

	t := &types.TransactionStatus{}
	err = json.Unmarshal(body, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (c *Client) GetTransactionField(id string, field string) (string, error) {
	body, err := c.get(fmt.Sprintf("tx/%s/%s", id, field))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) GetTransactionData(id string) ([]byte, error) {
	body, err := c.get(id)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (c *Client) GetTransactionPrice(size int, target string) (string, error) {
	url := fmt.Sprintf("price/%d/%s", size, target)
	body, err := c.get(url)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func (c *Client) SubmitTransaction(t *types.Transaction) ([]byte, int, error) {
	b, err := json.Marshal(t)
	if err != nil {
		return nil, -1, err
	}

	body, statusCode, err := c.httpPost("tx", b)
	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil
}

func (c *Client) GetWalletBalance(address string) (string, error) {
	body, err := c.get(fmt.Sprintf("wallet/%s/balance", address))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) GetLastTransactionID(address string) (string, error) {
	body, err := c.get(fmt.Sprintf("wallet/%s/last_tx", address))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) GetBlockByID(id string) (*types.Block, error) {
	body, err := c.get(fmt.Sprintf("block/hash/%s", id))
	if err != nil {
		return nil, err
	}
	b := &types.Block{}
	err = json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Client) GetBlockByHeight(height string) (*types.Block, error) {
	body, err := c.get(fmt.Sprintf("block/hash/%s", height))
	if err != nil {
		return nil, err
	}
	b := &types.Block{}
	err = json.Unmarshal(body, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (c *Client) GetNetworkInfo() (*types.NetworkInfo, error) {
	body, err := c.get("info")
	if err != nil {
		return nil, err
	}
	n := &types.NetworkInfo{}
	err = json.Unmarshal(body, n)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (c *Client) UploadChunk(chunk *types.GetChunkResult) ([]byte, int, error) {
	b, err := json.Marshal(chunk)
	if err != nil {
		return nil, -1, err
	}
	body, statusCode, err := c.httpPost("tx", b)
	if err != nil {
		return nil, -1, err
	}

	return body, statusCode, nil
}
