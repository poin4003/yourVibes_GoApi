package request

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/poin4003/yourVibes_GoApi/internal/application/advertise/command"
)

type ConfirmPaymentRequest struct {
	PartnerCode  string `form:"partnerCode"`
	OrderId      string `form:"orderId"`
	RequestId    string `form:"requestId"`
	Amount       string `form:"amount"`
	OrderInfo    string `form:"orderInfo"`
	OrderType    string `form:"orderType"`
	TransId      string `form:"transId"`
	ResultCode   string `form:"resultCode"`
	Message      string `form:"message"`
	PayType      string `form:"payType"`
	ResponseTime string `form:"responseTime"`
	ExtraData    string `form:"extraData"`
	Signature    string `form:"signature"`
}

func ValidateConfirmPaymentRequest(input interface{}) error {
	query, ok := input.(*ConfirmPaymentRequest)
	if !ok {
		return fmt.Errorf("validate ConfirmPaymentRequest failed")
	}

	return validation.ValidateStruct(query)
}

func (req *ConfirmPaymentRequest) ToConfirmPaymentCommand() (*command.ConfirmPaymentCommand, error) {
	if req == nil {
		return nil, fmt.Errorf("request in confirm payment request is nil")
	}

	var billId uuid.UUID
	if req.OrderInfo != "" {
		parseBillId, err := uuid.Parse(req.OrderInfo)
		if err != nil {
			return nil, err
		}
		billId = parseBillId
	}

	return &command.ConfirmPaymentCommand{
		BillId: &billId,
	}, nil
}
