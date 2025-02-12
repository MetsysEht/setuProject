package kycVerification

import (
	"context"

	kyc_verificationv1 "github.com/MetsysEht/setuProject/rpc/kycVerification"
)

type Server struct {
	kyc_verificationv1.UnimplementedKYCVerificationServiceServer
	manager IManager
}

func NewServer(manager IManager) *Server {
	return &Server{
		manager: manager,
	}
}

func (s *Server) VerifyPan(ctx context.Context, req *kyc_verificationv1.VerifyPanRequest) (*kyc_verificationv1.VerifyPanResponse, error) {
	request := &PANVerification{
		UserID:  req.GetUserId(),
		PAN:     req.GetPan(),
		Consent: req.GetConsent(),
		Reason:  req.GetReason(),
	}
	err := request.Validate()
	if err != nil {
		return nil, err
	}
	resp, err := s.manager.VerifyPan(ctx, request)
	if err != nil {
		return nil, err
	}
	return &kyc_verificationv1.VerifyPanResponse{
		Success: resp.Success,
		Message: resp.ResponseReason,
	}, nil
}

func (s *Server) VerifyRPD(ctx context.Context, req *kyc_verificationv1.RPDRequest) (*kyc_verificationv1.RpdResponse, error) {
	rpd := &RPD{
		UserID: req.GetUserId(),
	}
	err := rpd.Validate()
	if err != nil {
		return nil, err
	}
	rpd, err = s.manager.CreateRPD(ctx, rpd)
	if err != nil {
		return nil, err
	}
	resp := &kyc_verificationv1.RpdResponse{
		ShortUrl: rpd.ShortURL,
		Status:   rpd.RPDStatus,
		UpiLink:  rpd.UPILink,
	}
	return resp, nil
}

func (s *Server) RPDWebhook(ctx context.Context, req *kyc_verificationv1.RPDWebhookRequest) (*kyc_verificationv1.RPDWebhookResponse, error) {
	rpd := &RPD{
		TraceID: req.TraceId,
	}
	rpd, err := s.manager.RPDWebhook(ctx, rpd, req.Data.Rpd.Success)
	if err != nil {
		return nil, err
	}
	return &kyc_verificationv1.RPDWebhookResponse{}, nil
}

func (s *Server) GetStats(ctx context.Context, _ *kyc_verificationv1.Empty) (*kyc_verificationv1.KYCStatistics, error) {
	stats, err := s.manager.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	resp := &kyc_verificationv1.KYCStatistics{
		TotalKycAttempted:                    int32(stats.TotalKYCAttempted),
		TotalKycSuccessful:                   int32(stats.TotalKYCSuccessful),
		TotalKycFailed:                       int32(stats.TotalKYCFailed),
		TotalKycFailedDueToPan:               int32(stats.TotalKYCFailedDueToPAN),
		TotalKycFailedDueToBankAccount:       int32(stats.TotalKYCFailedDueToBankAccount),
		TotalKycFailedDueToPanAndBankAccount: -1,
	}
	return resp, nil
}
