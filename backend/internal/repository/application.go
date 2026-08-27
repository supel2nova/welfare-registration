package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/supel2nova/welfare-registration/backend/internal/domain"
	"github.com/supel2nova/welfare-registration/backend/internal/domain/appno"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
)

const ConstraintActivePerYear = "uq_app_active_per_year"

type MemberParams struct {
	Relation       string
	NationalIDHash *string
	FullName       string
	BirthYear      *int
	AnnualIncome   *int64
}

type ApplicationParams struct {
	ApplicationNo      string
	CitizenID          uuid.UUID
	FiscalYear         int
	IsFarmer           bool
	MaritalStatus      *string
	ExpenseToOthers    int64
	HasCreditCard      bool
	Snapshot           []byte
	RegistrationUnitID uuid.UUID
	CreatedByUserID    *uuid.UUID
	CreatedByClientID  *uuid.UUID
	SubmissionChannel  string
}

type VerificationParams struct {
	CitizenID    uuid.UUID
	Method       domain.VerificationMethod
	Verified     bool
	Note         *string
	VerifiedBy   *uuid.UUID
	ProviderCode *string
	ReferenceNo  *string
}

func (r *Repo) NextAppNo(ctx context.Context, tx pgx.Tx, fiscalYear int) (string, error) {
	var seq int64
	if err := tx.QueryRow(ctx, `SELECT nextval('app_no_seq')`).Scan(&seq); err != nil {
		return "", err
	}
	return appno.Format(fiscalYear, seq), nil
}

func (r *Repo) InsertApplication(ctx context.Context, tx pgx.Tx, p ApplicationParams) (uuid.UUID, time.Time, error) {
	const q = `
INSERT INTO applications (application_no, citizen_id, fiscal_year, is_farmer, marital_status,
                          expense_to_others, has_credit_card, applicant_snapshot,
                          registration_unit_id, created_by_user_id, created_by_client_id,
                          submission_channel)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
RETURNING id, submitted_at`

	var (
		id          uuid.UUID
		submittedAt time.Time
	)
	err := tx.QueryRow(ctx, q, p.ApplicationNo, p.CitizenID, p.FiscalYear, p.IsFarmer, p.MaritalStatus,
		p.ExpenseToOthers, p.HasCreditCard, p.Snapshot, p.RegistrationUnitID,
		p.CreatedByUserID, p.CreatedByClientID, p.SubmissionChannel).Scan(&id, &submittedAt)
	return id, submittedAt, err
}

func (r *Repo) InsertChildren(ctx context.Context, tx pgx.Tx, appID uuid.UUID, members []MemberParams, f dto.Financial) error {
	for _, m := range members {
		const q = `
INSERT INTO household_members (application_id, relation, national_id_hash, full_name, birth_year, annual_income)
VALUES ($1,$2,$3,$4,$5,$6)`
		if _, err := tx.Exec(ctx, q, appID, m.Relation, m.NationalIDHash, m.FullName, m.BirthYear, m.AnnualIncome); err != nil {
			return err
		}
	}

	for _, s := range f.IncomeSources {
		const q = `INSERT INTO income_sources (application_id, source_type, annual_amount) VALUES ($1,$2,$3)`
		if _, err := tx.Exec(ctx, q, appID, s.SourceType, s.AnnualAmount); err != nil {
			return err
		}
	}

	for _, a := range f.Assets {
		const q = `
INSERT INTO assets (application_id, asset_type, amount, unit, joint_account_holders, is_minor_account)
VALUES ($1,$2,$3::numeric,$4,COALESCE($5, 1),$6)`
		if _, err := tx.Exec(ctx, q, appID, a.AssetType, a.Amount.String(), a.Unit, a.JointAccountHolders, a.IsMinorAccount); err != nil {
			return err
		}
	}

	for _, l := range f.Liabilities {
		const q = `INSERT INTO liabilities (application_id, liability_type, credit_limit) VALUES ($1,$2,$3)`
		if _, err := tx.Exec(ctx, q, appID, l.LiabilityType, l.CreditLimit); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) InsertVerification(ctx context.Context, tx pgx.Tx, p VerificationParams) error {
	const q = `
INSERT INTO identity_verifications (citizen_id, method, verified, note, verified_by, provider_code, reference_no)
VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := tx.Exec(ctx, q, p.CitizenID, p.Method, p.Verified, p.Note, p.VerifiedBy, p.ProviderCode, p.ReferenceNo)
	return err
}

func (r *Repo) InsertStatusHistory(ctx context.Context, tx pgx.Tx, appID uuid.UUID, to domain.Status, actor domain.Actor) error {
	const q = `
INSERT INTO application_status_history (application_id, from_status, to_status, actor_type, actor_id, actor_role)
VALUES ($1, NULL, $2, $3, $4, $5)`
	_, err := tx.Exec(ctx, q, appID, to, actor.Type, actor.ActorID(), actor.Role)
	return err
}

func (r *Repo) FindActiveByHash(ctx context.Context, hash string, fiscalYear int) (*dto.DuplicateInfo, error) {
	const q = `
SELECT a.application_no, a.submitted_at, o.name_th
FROM applications a
JOIN citizens c      ON c.id = a.citizen_id
JOIN organizations o ON o.id = a.registration_unit_id
WHERE c.national_id_hash = $1 AND a.fiscal_year = $2 AND a.status <> 'CANCELLED'`

	return r.scanDuplicate(ctx, q, hash, fiscalYear)
}

func (r *Repo) FindActiveByCitizen(ctx context.Context, citizenID uuid.UUID, fiscalYear int) (*dto.DuplicateInfo, error) {
	const q = `
SELECT a.application_no, a.submitted_at, o.name_th
FROM applications a
JOIN organizations o ON o.id = a.registration_unit_id
WHERE a.citizen_id = $1 AND a.fiscal_year = $2 AND a.status <> 'CANCELLED'`

	return r.scanDuplicate(ctx, q, citizenID, fiscalYear)
}

func (r *Repo) scanDuplicate(ctx context.Context, q string, args ...any) (*dto.DuplicateInfo, error) {
	var (
		appNo       string
		submittedAt time.Time
		unit        string
	)
	err := r.pool.QueryRow(ctx, q, args...).Scan(&appNo, &submittedAt, &unit)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &dto.DuplicateInfo{
		RegisteredAt:   submittedAt.Format(time.DateOnly),
		RegisteredUnit: unit,
		ApplicationNo:  appNo,
		CanAppeal:      false,
	}, nil
}
