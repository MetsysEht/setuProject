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
