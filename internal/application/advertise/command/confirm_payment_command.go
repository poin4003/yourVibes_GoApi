package command

import (
	"github.com/google/uuid"
)

type ConfirmPaymentCommand struct {
	BillId *uuid.UUID
}
