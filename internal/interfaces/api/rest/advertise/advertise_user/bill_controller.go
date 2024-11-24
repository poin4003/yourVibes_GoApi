package advertise_user

import "context"

type cBill struct {
}

func NewBillController() *cBill {
	return &cBill{}
}

func (c *cBill) ConfirmPayment(ctx context.Context) {
	//var confirmRequest
}
