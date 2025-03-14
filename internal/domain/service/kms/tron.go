package kms

type TronProvider struct{}

func (p TronProvider) Validate(address string) (bool, error) {
	return true, nil
}

func (p TronProvider) Create() (WalletData, error) {
	return WalletData{
		Address:    "TExampleAddress",
		PrivateKey: "ExamplePrivateKey",
	}, nil
}
