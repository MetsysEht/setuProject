package kycVerification

import (
	"context"
)

type IManager interface {
	VerifyPan(context.Context, *PANVerification) (*PANVerification, error)
	CreateRPD(ctx context.Context, rpd *RPD) (*RPD, error)
}

type IRepo interface {
	SaveKYCVerification(_ context.Context, verification *PANVerification) error
	SaveRPDVerification(_ context.Context, verification *RPD) error
}
