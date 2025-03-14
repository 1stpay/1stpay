package kms

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"errors"

	"github.com/btcsuite/btcd/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type TronProvider struct{}

func (p TronProvider) Validate(address string) (bool, error) {
	return true, nil
}

func (p TronProvider) Create() (WalletData, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return WalletData{}, err
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)

	address, err := tronAddressFromPubKey(&privateKey.PublicKey)
	if err != nil {
		return WalletData{}, err
	}

	privateKeyStr := hexutil.Encode(privateKeyBytes)[2:]

	return WalletData{
		Address:    address,
		PrivateKey: privateKeyStr,
	}, nil
}

func tronAddressFromPubKey(pub *ecdsa.PublicKey) (string, error) {
	pubBytes := crypto.FromECDSAPub(pub)
	if len(pubBytes) != 65 {
		return "", errors.New("wrong pub key length")
	}
	pubBytes = pubBytes[1:]

	hash := crypto.Keccak256(pubBytes)
	if len(hash) < 20 {
		return "", errors.New("wrong hash length")
	}
	addrBytes := hash[len(hash)-20:]

	tronAddr := append([]byte{0x41}, addrBytes...)

	checksum := doubleSHA256(tronAddr)[:4]

	finalAddr := append(tronAddr, checksum...)

	tronAddress := base58.Encode(finalAddr)

	return tronAddress, nil
}

func doubleSHA256(data []byte) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:]
}
