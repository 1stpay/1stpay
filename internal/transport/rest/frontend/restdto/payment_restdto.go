package restdto

import "time"

type PaymentDetailResponseRestDTO struct {
	RequestedAmount    float64                           `json:"requested_amount"`
	Email              *string                           `json:"email"`
	ExpiresAt          *time.Time                        `json:"expires_at"`
	Merchant           MerchatDetailResponseDTO          `json:"merchant"`
	PaymentAddressList []PaymentAddressDetailResponseDTO `json:"payment_address_list"`
}

type PaymentAddressDetailResponseDTO struct {
	PublicKey       string             `json:"public_key"`
	RequestedAmount string             `json:"requested_amount"`
	Token           TokenDetailRestDTO `json:"token"`
}
