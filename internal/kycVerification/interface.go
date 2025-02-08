package kycVerification

import (
	"context"
)

type IManager interface {
	VerifyPan(context.Context, *PANVerification) (*PANVerification, error)
	CreateRPD(context.Context, *RPD) (*RPD, error)
	RPDWebhook(context.Context, *RPD, bool) (*RPD, error)
	GetStats(context.Context) (*KYCStatistics, error)
}

type IRepo interface {
	SaveKYCVerification(context.Context, *PANVerification) error
	SaveRPDVerification(context.Context, *RPD) error
	GetRPDFromTraceID(context.Context, string) (*RPD, error)
	GetKYCVerifiedUser(context.Context, *RPD) bool
	UpdateRPDVerificationStatus(context.Context, *RPD) error
	GetTotalKYCAttempts(context.Context) (int64, error)
	GetTotalKYCSuccess(context.Context) (int64, error)
	GetTotalRPDKYCFailed(context.Context) (int64, error)
	GetTotalPANKYCFailed(context.Context) (int64, error)
}
