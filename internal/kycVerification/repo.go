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

func (r *Repo) SaveKYCVerification(_ context.Context, verification *PANVerification) error {
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

func (r *Repo) SaveRPDVerification(_ context.Context, rpd *RPD) error {
	mod := model.RPDVerification{
		UserID:  rpd.UserID,
		TraceID: rpd.TraceID,
		Status:  rpd.RPDStatus,
	}
	tx := r.db.Save(&mod)
	return tx.Error
}

func (r *Repo) UpdateRPDVerificationStatus(_ context.Context, rpd *RPD) error {
	mod := model.RPDVerification{
		UserID:  rpd.UserID,
		TraceID: rpd.TraceID,
		Status:  rpd.RPDStatus,
	}
	tx := r.db.Model(&mod).Where("trace_id", rpd.TraceID).Update("status", rpd.RPDStatus)
	return tx.Error
}

func (r *Repo) GetRPDFromTraceID(_ context.Context, traceId string) (*RPD, error) {
	mod := &model.RPDVerification{}
	tx := r.db.Where("trace_id = ?", traceId).First(&mod)
	if tx.Error != nil {
		return nil, tx.Error
	}
	rpd := &RPD{
		UserID:    mod.UserID,
		RPDStatus: mod.Status,
		TraceID:   mod.TraceID,
	}
	return rpd, tx.Error
}

func (r *Repo) GetKYCVerifiedUser(_ context.Context, rpd *RPD) bool {
	mod := &model.PANVerification{}
	tx := r.db.Where("user_id = ? AND status = ?", rpd.UserID, true).First(&mod)
	if tx.Error != nil {
		return false
	}
	return true
}

func (r *Repo) GetTotalKYCAttempts(_ context.Context) (int64, error) {
	var count int64

	tx := r.db.Model(&model.PANVerification{}).Count(&count)
	if tx.Error != nil {
		return 0, nil
	}
	return count, nil
}

func (r *Repo) GetTotalKYCSuccess(_ context.Context) (int64, error) {
	var count int64

	tx := r.db.Model(&model.RPDVerification{}).Where("status = Success").Count(&count)
	if tx.Error != nil {
		return 0, nil
	}
	return count, nil
}

func (r *Repo) GetTotalRPDKYCFailed(_ context.Context) (int64, error) {
	var count int64

	tx := r.db.Model(&model.RPDVerification{}).Where("status = Failed").Count(&count)
	if tx.Error != nil {
		return 0, nil
	}
	return count, nil
}

func (r *Repo) GetTotalPANKYCFailed(_ context.Context) (int64, error) {
	var count int64

	tx := r.db.Model(&model.PANVerification{}).Where("status = ?", false).Count(&count)
	if tx.Error != nil {
		return 0, nil
	}
	return count, nil
}
