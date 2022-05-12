package iotex

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/vignesh-innblockchain/blockatlas/pkg/blockatlas"
	"github.com/trustwallet/golibs/client"
	"github.com/trustwallet/golibs/types"
)

type Client struct {
	client.Request
}

func (c *Client) GetLatestBlock() (int64, error) {
	var chainMeta ChainMeta
	err := c.Get(&chainMeta, "chainmeta", nil)
	if err != nil {
		return 0, err
	}
	b, err := strconv.ParseInt(chainMeta.Height, 10, 64)
	if err != nil {
		return 0, err
	}
	return b, nil
}

func (c *Client) GetTxsInBlock(number int64) ([]*ActionInfo, error) {
	path := fmt.Sprintf("transfers/block/%d", number)
	var resp Response
	err := c.Get(&resp, path, nil)
	if err != nil {
		return nil, err
	}
	return resp.ActionInfo, nil
}

func (c *Client) GetTxsOfAddress(address string, start int64) (*Response, error) {
	var response Response
	err := c.Get(&response, "actions/addr/"+address, url.Values{
		"start": {strconv.FormatInt(start, 10)},
		"count": {strconv.Itoa(types.TxPerPage)},
	})

	if err != nil {
		return nil, blockatlas.ErrSourceConn
	}
	return &response, err
}

func (c *Client) GetAddressTotalTransactions(address string) (int64, error) {
	var account AccountInfo
	err := c.Get(&account, "accounts/"+address, nil)
	if err != nil {
		return 0, nil
	}
	numActions, err := strconv.ParseInt(account.AccountMeta.NumActions, 10, 64)
	if err != nil {
		return 0, nil
	}
	return numActions, nil
}

func (c *Client) GetValidators() (blockatlas.ValidatorPage, error) {
	var validators blockatlas.ValidatorPage
	err := c.Get(&validators, "staking/validators", nil)
	if err != nil {
		return nil, err
	}
	return validators, nil
}

func (c *Client) GetDelegations(address string) (blockatlas.DelegationsPage, error) {
	var delegations blockatlas.DelegationsPage
	err := c.Get(&delegations, "staking/delegations/"+address, nil)
	if err != nil {
		return nil, err
	}
	return delegations, nil
}

func (c *Client) GetAccount(address string) (*AccountInfo, error) {
	var account AccountInfo
	err := c.Get(&account, "accounts/"+address, nil)
	if err != nil {
		return nil, err
	}
	return &account, nil
}
