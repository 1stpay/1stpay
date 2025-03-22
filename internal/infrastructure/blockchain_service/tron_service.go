package blockchain_service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type TronService struct {
	rpcUrl string
}

func NewTronService(rpcURL string) (*TronService, error) {
	return &TronService{rpcUrl: rpcURL}, nil
}

type TronAccountResponse struct {
	Data    []TronAccountData `json:"data"`
	Success bool              `json:"success"`
	Meta    json.RawMessage   `json:"meta"`
}

type TronAccountData struct {
	Address string              `json:"address"`
	Balance int64               `json:"balance"`
	TRC20   []map[string]string `json:"trc20"`
}

func (s TronService) fetchAccountData(address string) (*TronAccountResponse, error) {
	url := fmt.Sprintf("%s/v1/accounts/%s", strings.TrimRight(s.rpcUrl, "/"), address)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http get error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body error: %w", err)
	}

	var accountResp TronAccountResponse
	if err := json.Unmarshal(body, &accountResp); err != nil {
		return nil, fmt.Errorf("json unmarshal error: %w", err)
	}

	if !accountResp.Success || len(accountResp.Data) == 0 {
		return nil, errors.New("no account data returned or success false")
	}

	return &accountResp, nil
}

func (s TronService) GetNativeBalance(address string) (*big.Int, error) {

	accountResp, err := s.fetchAccountData(address)
	if err != nil {
		return nil, err
	}

	if !accountResp.Success || len(accountResp.Data) == 0 {
		return nil, errors.New("no account data returned or success false")
	}

	balance := big.NewInt(accountResp.Data[0].Balance)
	return balance, nil
}

func (s TronService) GetTokenBalance(address, tokenAddress string) (*big.Int, error) {
	accountResp, err := s.fetchAccountData(address)
	if err != nil {
		return nil, err
	}

	if !accountResp.Success || len(accountResp.Data) == 0 {
		return nil, errors.New("no account data returned or success false")
	}

	for _, trc20Obj := range accountResp.Data[0].TRC20 {
		if val, ok := trc20Obj[tokenAddress]; ok {
			balance := new(big.Int)
			_, ok := balance.SetString(val, 10)
			if !ok {
				return nil, fmt.Errorf("failed to parse token balance: %s", val)
			}
			return balance, nil
		}
	}

	return nil, fmt.Errorf("token %s not found in account %s", tokenAddress, address)
}
