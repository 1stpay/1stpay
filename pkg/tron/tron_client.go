package tron

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type TronClient struct {
	rpcURL     string
	httpClient *http.Client
}

func NewTronClient(rpcURL string) *TronClient {
	return &TronClient{
		rpcURL:     strings.TrimRight(rpcURL, "/"),
		httpClient: &http.Client{},
	}
}

func TronHexToBase58(hexAddr string) (string, error) {
	addrBytes, err := hex.DecodeString(hexAddr)
	if err != nil {
		return "", fmt.Errorf("failed to decode hex address: %w", err)
	}
	first := sha256.Sum256(addrBytes)
	second := sha256.Sum256(first[:])
	checksum := second[:4]
	fullBytes := append(addrBytes, checksum...)
	return base58.Encode(fullBytes), nil
}

func Base58ToTronHex(base58Addr string) (string, error) {
	decoded := base58.Decode(base58Addr)
	if len(decoded) != 25 {
		return "", fmt.Errorf("invalid address length")
	}
	checksum := decoded[len(decoded)-4:]
	payload := decoded[:len(decoded)-4]

	hash0 := sha256.Sum256(payload)
	hash1 := sha256.Sum256(hash0[:])

	if string(checksum) != string(hash1[:4]) {
		return "", fmt.Errorf("invalid checksum")
	}

	return fmt.Sprintf("%x", payload), nil
}

func ConstructTronTokenTxData(recipientHex string, amount *big.Int) (string, error) {
	parsedABI, err := abi.JSON(strings.NewReader(`[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`))
	if err != nil {
		return "", fmt.Errorf("error while txn build")
	}
	data, err := parsedABI.Pack("transfer", common.HexToAddress(recipientHex), amount)
	if err != nil {
		return "", fmt.Errorf("error while txn build")
	}
	return hex.EncodeToString(data[4:]), nil
}

func GetTronAddressFromPrivateKey(privateKey string) (string, error) {
	senderPrivateKey := strings.TrimPrefix(privateKey, "0x")
	privKey, err := crypto.HexToECDSA(senderPrivateKey)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}
	pubKey := privKey.Public()
	publicKeyECDSA, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("failed to cast public key to ECDSA")
	}
	pubBytes := crypto.FromECDSAPub(publicKeyECDSA)
	hash := crypto.Keccak256(pubBytes[1:])
	addrBytes := hash[len(hash)-20:]
	tronAddrBytes := append([]byte{0x41}, addrBytes...)
	hexAddr := hex.EncodeToString(tronAddrBytes)
	base58Addr, err := TronHexToBase58(hexAddr)
	if err != nil {
		return "", fmt.Errorf("failed to convert hex address to base58: %w", err)
	}
	return base58Addr, nil
}

func (c *TronClient) doRequest(ctx context.Context, method, path string, payload interface{}) ([]byte, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.rpcURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("received non-OK status: %d, url, %s, body: %s", resp.StatusCode, path, string(respBody))
	}
	return respBody, nil
}

func (c *TronClient) FetchAccountData(ctx context.Context, address string) (*AccountResponse, error) {
	respBody, err := c.doRequest(ctx, "GET", fmt.Sprintf("/v1/accounts/%s", address), nil)
	if err != nil {
		return nil, err
	}
	var accountResp AccountResponse
	if err := json.Unmarshal(respBody, &accountResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}
	if !accountResp.Success {
		return nil, errors.New("no account data returned or success false")
	}
	return &accountResp, nil
}

func (c *TronClient) CreateTransaction(ctx context.Context, fromAddress, toAddress string, amount *big.Int) (*Transaction, error) {
	reqPayload := TransactionCreateRequest{
		OwnerAddress: fromAddress,
		ToAddress:    toAddress,
		Amount:       int(amount.Int64()),
		Visible:      true,
	}
	respBody, err := c.doRequest(ctx, "POST", "/wallet/createtransaction", reqPayload)
	if err != nil {
		return nil, err
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
		return nil, errors.New(errResp.Error)
	}
	var tx Transaction
	if err := json.Unmarshal(respBody, &tx); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transaction response: %w", err)
	}
	if tx.TxID == "" {
		return nil, errors.New("transaction ID is empty in response")
	}
	return &tx, nil
}

func (c *TronClient) SignTransaction(tx *Transaction, privateKey string) error {
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}
	rawDataBytes, err := hexutil.Decode("0x" + tx.RawDataHex)
	if err != nil {
		return fmt.Errorf("unable to decode rawDataHex: %w", err)
	}
	hash := sha256.New()
	hash.Write(rawDataBytes)
	txHash := hash.Sum(nil)
	signature, err := crypto.Sign(txHash, privKey)
	if err != nil {
		return fmt.Errorf("unable to sign tx: %w", err)
	}
	tx.Signature = append(tx.Signature, hexutil.Encode(signature)[2:])
	return nil
}

func (c *TronClient) BroadcastTransaction(ctx context.Context, tx *Transaction) (common.Hash, error) {
	respBody, err := c.doRequest(ctx, "POST", "/wallet/broadcasttransaction", tx)
	if err != nil {
		return common.Hash{}, err
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
		return common.Hash{}, errors.New(errResp.Error)
	}
	var broadcastResp BroadcastTransactionResponse
	if err := json.Unmarshal(respBody, &broadcastResp); err != nil {
		return common.Hash{}, errors.New("failed to unmarshal broadcast response")
	}
	return common.HexToHash(broadcastResp.TxID), nil
}

// Returns wallet's TRX balance
func (c *TronClient) GetNativeBalance(ctx context.Context, address string) (*big.Int, error) {
	accountResp, err := c.FetchAccountData(ctx, address)
	if err != nil {
		return nil, err
	}
	if len(accountResp.Data) == 0 {
		return big.NewInt(0), nil
	}
	return big.NewInt(accountResp.Data[0].Balance), nil
}

// Returns wallet balance for given TRC20 contract.
func (c *TronClient) GetTokenBalance(ctx context.Context, address, tokenAddress string) (*big.Int, error) {
	accountResp, err := c.FetchAccountData(ctx, address)
	if err != nil {
		return nil, err
	}
	if !accountResp.Success || len(accountResp.Data) == 0 {
		return nil, errors.New("no account data returned or success false")
	}
	for _, token := range accountResp.Data[0].TRC20 {
		if val, ok := token[tokenAddress]; ok {
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

func (c *TronClient) TransferNative(
	ctx context.Context,
	senderPrivateKey,
	toAddress string,
	amount *big.Int,
) (common.Hash, error) {
	senderAddress, err := GetTronAddressFromPrivateKey(senderPrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get sender address: %w", err)
	}
	tx, err := c.CreateTransaction(ctx, senderAddress, toAddress, amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create transaction: %w", err)
	}
	if err = c.SignTransaction(tx, senderPrivateKey); err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}
	return c.BroadcastTransaction(ctx, tx)
}

func (c *TronClient) TransferToken(
	ctx context.Context,
	senderPrivateKey,
	tokenContractAddress,
	toAddress string,
	amount *big.Int,
) (common.Hash, error) {
	senderAddress, err := GetTronAddressFromPrivateKey(senderPrivateKey)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while getting address from private key")
	}

	toAddressHex, err := Base58ToTronHex(toAddress)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while address conversion")
	}
	parameterHex, err := ConstructTronTokenTxData(toAddressHex, amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("error while constructing txn data: %w", err)
	}

	triggerPayload := TriggerSmartContractRequest{
		OwnerAddress:     senderAddress,
		ContractAddress:  tokenContractAddress,
		FunctionSelector: "transfer(address,uint256)",
		Parameter:        parameterHex,
		FeeLimit:         10_000_000, // @TODO
		Visible:          true,
	}

	respBody, err := c.doRequest(ctx, "POST", "/wallet/triggersmartcontract", triggerPayload)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to trigger smart contract: %w", err)
	}

	var errResp ErrorResponse
	if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
		return common.Hash{}, fmt.Errorf("trigger smart contract error: %s", errResp.Error)
	}

	var resp TriggerSmartContractResponse
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return common.Hash{}, fmt.Errorf("failed to unmarshal transaction: %w", err)
	}
	if !resp.Result.Result {
		return common.Hash{}, errors.New("error while smart contract trigger")
	}
	tx := resp.Transaction
	if tx.TxID == "" {
		return common.Hash{}, errors.New("transaction id is empty in response")
	}

	if err = c.SignTransaction(&tx, senderPrivateKey); err != nil {
		return common.Hash{}, fmt.Errorf("failed to sign transaction: %w", err)
	}

	return c.BroadcastTransaction(ctx, &tx)
}
