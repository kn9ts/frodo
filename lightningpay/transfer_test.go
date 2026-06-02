package lightningpay

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type stubIssuer struct {
	called     bool
	amountSats int64
	memo       string
	invoice    Invoice
	err        error
}

func (s *stubIssuer) CreateInvoice(_ context.Context, amountSats int64, memo string) (Invoice, error) {
	s.called = true
	s.amountSats = amountSats
	s.memo = memo
	return s.invoice, s.err
}

type stubPayer struct {
	called         bool
	paymentRequest string
	paymentID      string
	err            error
}

func (s *stubPayer) PayInvoice(_ context.Context, paymentRequest string) (string, error) {
	s.called = true
	s.paymentRequest = paymentRequest
	return s.paymentID, s.err
}

func TestTransferSatsPaysGeneratedInvoice(t *testing.T) {
	issuer := &stubIssuer{invoice: Invoice{PaymentRequest: "lnbc2500n1p...", AmountSats: 2500, Memo: "coffee"}}
	payer := &stubPayer{paymentID: "payment-123"}

	receipt, err := TransferSats(context.Background(), payer, issuer, 2500, " coffee ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !issuer.called {
		t.Fatal("expected receiver invoice creation to be called")
	}
	if issuer.amountSats != 2500 {
		t.Fatalf("expected amount 2500, got %d", issuer.amountSats)
	}
	if issuer.memo != "coffee" {
		t.Fatalf("expected trimmed memo %q, got %q", "coffee", issuer.memo)
	}
	if !payer.called {
		t.Fatal("expected sender payment to be called")
	}
	if payer.paymentRequest != "lnbc2500n1p..." {
		t.Fatalf("expected invoice %q, got %q", "lnbc2500n1p...", payer.paymentRequest)
	}
	if !reflect.DeepEqual(receipt, Receipt{PaymentRequest: "lnbc2500n1p...", PaymentID: "payment-123"}) {
		t.Fatalf("unexpected receipt: %#v", receipt)
	}
}

func TestTransferSatsRejectsInvalidAmount(t *testing.T) {
	_, err := TransferSats(context.Background(), &stubPayer{}, &stubIssuer{}, 0, "memo")
	if !errors.Is(err, ErrInvalidSats) {
		t.Fatalf("expected ErrInvalidSats, got %v", err)
	}
}

func TestTransferSatsRejectsEmptyInvoice(t *testing.T) {
	_, err := TransferSats(context.Background(), &stubPayer{}, &stubIssuer{invoice: Invoice{AmountSats: 100}}, 100, "memo")
	if !errors.Is(err, ErrEmptyPaymentRequest) {
		t.Fatalf("expected ErrEmptyPaymentRequest, got %v", err)
	}
}

func TestTransferSatsWrapsSenderError(t *testing.T) {
	issuer := &stubIssuer{invoice: Invoice{PaymentRequest: "lnbc1...", AmountSats: 1000}}
	payer := &stubPayer{err: errors.New("insufficient balance")}

	_, err := TransferSats(context.Background(), payer, issuer, 1000, "memo")
	if err == nil || !errors.Is(err, payer.err) {
		t.Fatalf("expected wrapped sender error, got %v", err)
	}
}
