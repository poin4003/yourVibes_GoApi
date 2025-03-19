package consts

type VoucherType bool

type VoucherStatus bool

const (
	PERCENTAGE       VoucherType   = true
	FIX_AMOUNT       VoucherType   = false
	VOUCHER_ACTIVE   VoucherStatus = true
	VOUCHER_INACTIVE VoucherStatus = false
)

var VoucherTypes = []interface{}{
	PERCENTAGE,
	FIX_AMOUNT,
}
