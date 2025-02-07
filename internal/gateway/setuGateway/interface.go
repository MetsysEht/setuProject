package setuGateway

import (
	"context"
)

type ISetuGateway interface {
	VerifyPan(context.Context, *PANRequest) (*PANResponse, error)
}

type PANRequest struct {
	PAN     string `json:"pan"`
	Consent string `json:"consent"`
	Reason  string `json:"reason"`
}

type PANResponse struct {
	Data         Data   `json:"data"`
	Message      string `json:"message"`
	Verification string `json:"verification"`
	TraceID      string `json:"traceId"`
}

type Data struct {
	AadhaarSeedingStatus string `json:"aadhaar_seeding_status"`
	Category             string `json:"category"`
	FullName             string `json:"full_name"`
}
