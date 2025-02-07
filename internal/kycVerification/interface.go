package kycVerification

import (
	"context"
)

type IManager interface {
	VerifyPan(context.Context, *PANVerification) (*PANVerification, error)
	CreateRPD(ctx context.Context, rpd *RPD) (*RPD, error)
	RPDWebhook(ctx context.Context, rpd *RPD, success bool) (*RPD, error)
}

type IRepo interface {
	SaveKYCVerification(_ context.Context, verification *PANVerification) error
	SaveRPDVerification(_ context.Context, verification *RPD) error
	GetRPDFromTraceID(_ context.Context, traceId string) (*RPD, error)
	GetKYCVerifiedUser(_ context.Context, rpd *RPD) bool
	UpdateRPDVerificationStatus(_ context.Context, rpd *RPD) error
}
