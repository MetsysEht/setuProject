package kycVerification

import (
	"context"

	"github.com/MetsysEht/setuProject/internal/kycVerification/model"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) IRepo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Save(_ context.Context, verification *PANVerification) error {
	status := false
	if verification.Success == "SUCCESS" {
		status = true
	}
	mod := model.PANVerification{
		UserID:        verification.UserID,
		PAN:           verification.PAN,
		Consent:       verification.Consent,
		RequestReason: verification.Reason,
		Status:        status,
	}
	tx := r.db.Save(&mod)
	return tx.Error
}
