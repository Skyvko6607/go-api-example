package enums

type PaymentType int

const (
	PayPal PaymentType = iota
	Barion
	Stripe
	BankTransfer
)

func (p PaymentType) String() string {
	switch p {
	case PayPal:
		return "PayPal"
	case Barion:
		return "Barion"
	case Stripe:
		return "Stripe"
	case BankTransfer:
		return "BankTransfer"
	default:
		return "Unknown"
	}
}
