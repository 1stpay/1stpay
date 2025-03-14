package kms

type SolanaProvider struct{}

func (p SolanaProvider) Validate(address string) (bool, error) {
	return true, nil
}

func (p SolanaProvider) Create() (WalletData, error) {
	return WalletData{
		Address:    "SolExampleAddress",
		PrivateKey: "ExamplePrivateKey",
	}, nil
}
