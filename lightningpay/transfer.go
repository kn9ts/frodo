package lightningpay

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidSats         = errors.New("amount of sats must be greater than zero")
	ErrEmptyPaymentRequest = errors.New("payment request is empty")
)

type Invoice struct {
	PaymentRequest string
	AmountSats     int64
	Memo           string
}

type Issuer interface {
	CreateInvoice(ctx context.Context, amountSats int64, memo string) (Invoice, error)
}

type Payer interface {
	PayInvoice(ctx context.Context, paymentRequest string) (string, error)
}

type Receipt struct {
	PaymentRequest string
	PaymentID      string
}

func TransferSats(ctx context.Context, sender Payer, receiver Issuer, amountSats int64, memo string) (Receipt, error) {
	if ctx == nil {
		return Receipt{}, errors.New("context is required")
	}
	if amountSats <= 0 {
		return Receipt{}, ErrInvalidSats
	}
	if sender == nil {
		return Receipt{}, errors.New("sender is required")
	}
	if receiver == nil {
		return Receipt{}, errors.New("receiver is required")
	}

	invoice, err := receiver.CreateInvoice(ctx, amountSats, strings.TrimSpace(memo))
	if err != nil {
		return Receipt{}, fmt.Errorf("create invoice: %w", err)
	}
	if strings.TrimSpace(invoice.PaymentRequest) == "" {
		return Receipt{}, ErrEmptyPaymentRequest
	}
	if invoice.AmountSats != 0 && invoice.AmountSats != amountSats {
		return Receipt{}, fmt.Errorf("invoice amount mismatch: expected %d sats, got %d sats", amountSats, invoice.AmountSats)
	}

	paymentID, err := sender.PayInvoice(ctx, invoice.PaymentRequest)
	if err != nil {
		return Receipt{}, fmt.Errorf("pay invoice: %w", err)
	}

	return Receipt{PaymentRequest: invoice.PaymentRequest, PaymentID: paymentID}, nil
}
