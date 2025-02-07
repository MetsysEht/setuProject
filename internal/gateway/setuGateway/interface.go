package setuGateway

import (
	"context"
)

type ISetuGateway interface {
	VerifyPan(context.Context, *PANRequest) (*PANResponse, error)
	CreateRPD(context.Context, *RPDPayload) (*RPDResponse, error)
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

// RPDPayload represents the structure of the JSON payload
type RPDPayload struct {
	RedirectionConfig RedirectionConfig `json:"redirectionConfig"`
	AdditionalData    map[string]string `json:"additionalData"`
}

// RedirectionConfig represents the nested "redirectionConfig" field
type RedirectionConfig struct {
	RedirectURL string `json:"redirectUrl"`
	Timeout     int    `json:"timeout"`
}

type RPDResponse struct {
	ID        string         `json:"id"`
	ShortURL  string         `json:"shortUrl"`
	Status    string         `json:"status"`
	TraceID   string         `json:"traceId"`
	UpiBillID string         `json:"upiBillId"`
	UpiLink   string         `json:"upiLink"`
	ValidUpto string         `json:"validUpto"`
	Error     *ErrorResponse `json:"error"`
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail represents the nested "error" field
type ErrorDetail struct {
	Code    string `json:"code"`
	Detail  string `json:"detail"`
	TraceID string `json:"traceId"`
}
