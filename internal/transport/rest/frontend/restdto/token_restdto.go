package restdto

type TokenDetailRestDTO struct {
	Name       string                  `json:"name"`
	Symbol     string                  `json:"symbol"`
	Logo       *string                 `json:"logo"`
	Blockchain BlockchainDetailRestDTO `json:"blockchain"`
}
