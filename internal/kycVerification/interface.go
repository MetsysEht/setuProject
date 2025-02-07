package kycVerification

import (
	"context"
)

type IManager interface {
	VerifyPan(context.Context, *PANVerification) (*PANVerification, error)
}

type IRepo interface {
	Save(_ context.Context, verification *PANVerification) error
}
