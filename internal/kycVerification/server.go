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
		Message: resp.Reason,
	}, nil
}
