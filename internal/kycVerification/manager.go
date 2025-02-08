package kycVerification

import (
	"context"

	"github.com/MetsysEht/setuProject/internal/gateway/setuGateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Manager struct {
	repo        IRepo
	setuGateway setuGateway.ISetuGateway
}

func NewManager(kycVerificationRepo IRepo, gateway setuGateway.ISetuGateway) IManager {
	return &Manager{
		kycVerificationRepo,
		gateway,
	}
}

func (m *Manager) VerifyPan(ctx context.Context, panDetails *PANVerification) (*PANVerification, error) {
	consent := "N"
	if panDetails.Consent {
		consent = "Y"
	}
	resp, err := m.setuGateway.VerifyPan(ctx, &setuGateway.PANRequest{
		PAN:     panDetails.PAN,
		Consent: consent,
		Reason:  panDetails.Reason,
	})
	if err != nil {
		return nil, err
	}
	panDetails.Success = resp.Verification
	panDetails.ResponseReason = resp.Message
	err = m.repo.SaveKYCVerification(ctx, panDetails)
	if err != nil {
		return nil, err
	}
	return panDetails, nil
}

func (m *Manager) CreateRPD(ctx context.Context, rpd *RPD) (*RPD, error) {
	if !m.repo.GetKYCVerifiedUser(ctx, rpd) {
		return nil, status.Error(codes.InvalidArgument, "User not verified PAN")
	}
	resp, err := m.setuGateway.CreateRPD(ctx, &setuGateway.RPDPayload{
		RedirectionConfig: setuGateway.RedirectionConfig{},
		AdditionalData: map[string]string{
			"user_id": rpd.UserID,
		},
	})
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, status.Error(codes.Internal, "Failed to create RPD")
	}
	if resp.Status == "BAV_REVERSE_PENNY_DROP_CREATED" {
		rpd.RPDStatus = "RPD Created"
	}
	rpd.ShortURL = resp.ShortURL
	rpd.TraceID = resp.TraceID
	rpd.UpiBillID = resp.UpiBillID
	rpd.UPILink = resp.UpiLink
	err = m.repo.SaveRPDVerification(ctx, rpd)
	if err != nil {
		return nil, err
	}
	return rpd, nil
}

func (m *Manager) RPDWebhook(ctx context.Context, rpd *RPD, success bool) (*RPD, error) {
	rpd, err := m.repo.GetRPDFromTraceID(ctx, rpd.TraceID)
	if err != nil {
		return nil, err
	}
	if success {
		rpd.RPDStatus = "Success"
	} else {
		rpd.RPDStatus = "Failed"
	}
	err = m.repo.UpdateRPDVerificationStatus(ctx, rpd)
	if err != nil {
		return nil, err
	}
	return rpd, nil
}

func (m *Manager) GetStats(ctx context.Context) (*KYCStatistics, error) {
	totalKYCAttempt, err := m.repo.GetTotalKYCAttempts(ctx)
	if err != nil {
		return nil, err
	}
	totalKYCSuccess, err := m.repo.GetTotalKYCSuccess(ctx)
	if err != nil {
		return nil, err
	}
	totalRPDFailed, err := m.repo.GetTotalRPDKYCFailed(ctx)
	if err != nil {
		return nil, err
	}
	totalPANFailed, err := m.repo.GetTotalPANKYCFailed(ctx)
	if err != nil {
		return nil, err
	}
	totalKYCFailed := totalPANFailed + totalRPDFailed
	return &KYCStatistics{
		TotalKYCAttempted:              totalKYCAttempt,
		TotalKYCSuccessful:             totalKYCSuccess,
		TotalKYCFailed:                 totalKYCFailed,
		TotalKYCFailedDueToPAN:         totalPANFailed,
		TotalKYCFailedDueToBankAccount: totalRPDFailed,
		TotalKYCFailedDueToPANAndBank:  -1,
	}, nil
}
