package kycVerification

import (
	"context"

	"github.com/MetsysEht/setuProject/internal/gateway/setuGateway"
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
	err = m.repo.Save(ctx, panDetails)
	if err != nil {
		return nil, err
	}
	return panDetails, nil
}
