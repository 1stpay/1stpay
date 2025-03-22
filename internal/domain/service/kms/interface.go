package kms

type WalletData struct {
	Address    string
	PrivateKey string
}

type WalletProvider interface {
	Validate(address string) (bool, error)
	Create() (WalletData, error)
}
