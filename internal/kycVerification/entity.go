package kycVerification

import (
	"regexp"

	"github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PANVerification struct {
	UserID         string
	PAN            string
	Consent        bool
	Reason         string
	Success        string
	ResponseReason string
}

func (p *PANVerification) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.UserID, validation.Required),
		validation.Field(&p.PAN, validation.Required, validation.Match(regexp.MustCompile("[A-Z]{5}[0-9]{4}[A-Z]"))),
		validation.Field(&p.Consent, validation.Required),
		validation.Field(&p.Reason, validation.Required, validation.Length(20, 100)),
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

type RPD struct {
	UserID    string `json:"user_id"`
	ShortURL  string `json:"short_url"`
	UPILink   string `json:"upi_link"`
	RPDStatus string `json:"rpd_status"`
	TraceID   string `json:"traceId"`
	UpiBillID string `json:"upiBillId"`
}

func (r *RPD) Validate() error {
	err := validation.ValidateStruct(r,
		validation.Field(&r.UserID, validation.Required),
	)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

type KYCStatistics struct {
	TotalKYCAttempted              int64 `json:"total_kyc_attempted"`
	TotalKYCSuccessful             int64 `json:"total_kyc_successful"`
	TotalKYCFailed                 int64 `json:"total_kyc_failed"`
	TotalKYCFailedDueToPAN         int64 `json:"total_kyc_failed_due_to_pan"`
	TotalKYCFailedDueToBankAccount int64 `json:"total_kyc_failed_due_to_bank_account"`
	TotalKYCFailedDueToPANAndBank  int64 `json:"total_kyc_failed_due_to_pan_and_bank_account"`
}
