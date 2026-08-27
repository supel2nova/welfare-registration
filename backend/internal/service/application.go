package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/supel2nova/welfare-registration/backend/internal/domain"
	"github.com/supel2nova/welfare-registration/backend/internal/domain/nationalid"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
	"github.com/supel2nova/welfare-registration/backend/internal/verifier"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
	"github.com/supel2nova/welfare-registration/backend/pkg/idcrypto"
)

type ApplicationService struct {
	repo     *repository.Repo
	verifier verifier.Verifier
	cipher   idcrypto.Cipher
	pepper   string
	now      func() time.Time

	insertChildren func(ctx context.Context, tx pgx.Tx, appID uuid.UUID, members []repository.MemberParams, f dto.Financial) error
}

func NewApplicationService(repo *repository.Repo, v verifier.Verifier, cipher idcrypto.Cipher, pepper string) *ApplicationService {
	s := &ApplicationService{
		repo:     repo,
		verifier: v,
		cipher:   cipher,
		pepper:   pepper,
		now:      time.Now,
	}
	s.insertChildren = repo.InsertChildren
	return s
}

func (s *ApplicationService) Create(ctx context.Context, actor domain.Actor, req dto.CreateApplicationRequest) (*dto.CreateApplicationResponse, error) {
	if fieldErrs := validateFormat(req, s.now()); len(fieldErrs) > 0 {
		return nil, apperror.Validation(fieldErrs)
	}

	addr, err := s.repo.ResolveAddress(ctx,
		req.Personal.Address.ProvinceCode,
		req.Personal.Address.DistrictCode,
		req.Personal.Address.SubdistrictCode,
		req.Personal.Address.PostalCode,
	)
	if errors.Is(err, repository.ErrAddressNotFound) {
		return nil, apperror.BadRequest(apperror.CodeAddress,
			apperror.Field("personal.address.subdistrict_code", apperror.CodeAddress))
	}
	if err != nil {
		return nil, apperror.Internal(err)
	}
	if !addr.PostalMatches {
		return nil, apperror.BadRequest(apperror.CodeAddress,
			apperror.Field("personal.address.postal_code", apperror.CodeAddress))
	}

	hash := nationalid.Hash(s.pepper, req.Personal.NationalID)
	if dup, err := s.repo.FindActiveByHash(ctx, hash, req.FiscalYear); err != nil {
		return nil, apperror.Internal(err)
	} else if dup != nil {
		return nil, apperror.Duplicate(dup)
	}

	method, providerCode, referenceNo, err := s.verifyIdentity(ctx, req)
	if err != nil {
		return nil, err
	}

	enc, err := s.cipher.Seal(req.Personal.NationalID, "")
	if err != nil {
		return nil, apperror.Internal(err)
	}

	snap, err := buildSnapshot(req, addr)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	members := memberParams(s.pepper, req)

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return nil, apperror.Internal(err)
	}
	defer tx.Rollback(ctx)

	addrID, err := s.repo.InsertAddress(ctx, tx, repository.AddressParams{
		HouseNo:         req.Personal.Address.HouseNo,
		Moo:             req.Personal.Address.Moo,
		Road:            req.Personal.Address.Road,
		SubdistrictCode: req.Personal.Address.SubdistrictCode,
		DistrictCode:    req.Personal.Address.DistrictCode,
		ProvinceCode:    req.Personal.Address.ProvinceCode,
		PostalCode:      req.Personal.Address.PostalCode,
		SubdistrictName: addr.SubdistrictName,
		DistrictName:    addr.DistrictName,
		ProvinceName:    addr.ProvinceName,
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}

	citizenID, err := s.repo.UpsertCitizen(ctx, tx, repository.CitizenParams{
		NationalIDHash: hash,
		NationalIDEnc:  enc,
		Title:          req.Personal.Title,
		FirstName:      req.Personal.FirstName,
		LastName:       req.Personal.LastName,
		BirthYear:      req.Personal.BirthYear,
		BirthMonth:     req.Personal.BirthMonth,
		BirthDay:       req.Personal.BirthDay,
		BirthPrecision: req.Personal.BirthPrecision,
		Phone:          req.Personal.Phone,
		AddressID:      addrID,
	})
	if err != nil {
		return nil, apperror.Internal(err)
	}

	appNo, err := s.repo.NextAppNo(ctx, tx, req.FiscalYear)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	var marital *string
	if req.Family != nil {
		marital = req.Family.MaritalStatus
	}

	appID, submittedAt, err := s.repo.InsertApplication(ctx, tx, repository.ApplicationParams{
		ApplicationNo:      appNo,
		CitizenID:          citizenID,
		FiscalYear:         req.FiscalYear,
		IsFarmer:           req.Personal.IsFarmer,
		MaritalStatus:      marital,
		ExpenseToOthers:    req.Financial.ExpenseToOthers,
		HasCreditCard:      req.Financial.HasCreditCard,
		Snapshot:           snap,
		RegistrationUnitID: actor.OrgID,
		CreatedByUserID:    actor.UserID,
		CreatedByClientID:  actor.ClientID,
		SubmissionChannel:  string(domain.ChannelWalkIn),
	})
	if repository.IsUniqueViolation(err, repository.ConstraintActivePerYear) {
		_ = tx.Rollback(ctx)
		return nil, s.duplicateAfterConflict(ctx, hash, req.FiscalYear)
	}
	if err != nil {
		return nil, apperror.Internal(err)
	}

	if err := s.insertChildren(ctx, tx, appID, members, req.Financial); err != nil {
		return nil, apperror.Internal(err)
	}

	var note *string
	if method == domain.MethodManual {
		note = req.Personal.IDVerifyNote
	}
	verified := method != domain.MethodPendingVerification
	var verifiedBy *uuid.UUID
	if method == domain.MethodManual {
		verifiedBy = actor.UserID
	}
	if err := s.repo.InsertVerification(ctx, tx, repository.VerificationParams{
		CitizenID:    citizenID,
		Method:       method,
		Verified:     verified,
		Note:         note,
		VerifiedBy:   verifiedBy,
		ProviderCode: providerCode,
		ReferenceNo:  referenceNo,
	}); err != nil {
		return nil, apperror.Internal(err)
	}

	if err := s.repo.InsertStatusHistory(ctx, tx, appID, domain.StatusSubmitted, actor); err != nil {
		return nil, apperror.Internal(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, apperror.Internal(err)
	}

	unit, err := s.repo.OrgName(ctx, actor.OrgID)
	if err != nil {
		return nil, apperror.Internal(err)
	}

	return &dto.CreateApplicationResponse{
		ApplicationID:    appID.String(),
		ApplicationNo:    appNo,
		Status:           string(domain.StatusSubmitted),
		RegistrationUnit: unit,
		SubmittedAt:      submittedAt,
	}, nil
}

func (s *ApplicationService) verifyIdentity(ctx context.Context, req dto.CreateApplicationRequest) (
	domain.VerificationMethod, *string, *string, error,
) {
	switch domain.IDVerifyMethod(req.Personal.IDVerifyMethod) {
	case domain.VerifyManualCardCheck:
		return domain.MethodManual, nil, nil, nil
	case domain.VerifyChipRead:
		return domain.MethodChipRead, nil, nil, nil
	case domain.VerifyLaserCode:
		laser := ""
		if req.Personal.LaserID != nil {
			laser = *req.Personal.LaserID
		}
		res, err := s.verifier.Verify(ctx, verifier.Request{
			NationalID: req.Personal.NationalID,
			LaserID:    laser,
		})
		if errors.Is(err, verifier.ErrUnavailable) {
			return domain.MethodPendingVerification, nil, nil, nil
		}
		if err != nil {
			return "", nil, nil, apperror.Internal(err)
		}
		if !res.Matched {
			return "", nil, nil, apperror.BadRequest(apperror.CodeKYC,
				apperror.Field("personal.laser_id", apperror.CodeKYC))
		}
		pc, rn := res.ProviderCode, res.ReferenceNo
		return domain.MethodLaserCode, &pc, &rn, nil
	default:
		return "", nil, nil, apperror.BadRequest(apperror.CodeEnum,
			apperror.Field("personal.id_verify_method", apperror.CodeEnum))
	}
}

func (s *ApplicationService) duplicateAfterConflict(ctx context.Context, hash string, fiscalYear int) error {
	citizenID, err := s.repo.FindCitizenIDByHash(ctx, hash)
	if err != nil {
		return apperror.Internal(err)
	}
	existing, err := s.repo.FindActiveByCitizen(ctx, citizenID, fiscalYear)
	if err != nil {
		return apperror.Internal(err)
	}
	if existing == nil {
		return apperror.Internal(errors.New("unique violation แต่ไม่พบใบเดิม"))
	}
	return apperror.Duplicate(existing)
}

func memberParams(pepper string, req dto.CreateApplicationRequest) []repository.MemberParams {
	if req.Family == nil {
		return nil
	}
	out := make([]repository.MemberParams, 0, len(req.Family.Members))
	for _, m := range req.Family.Members {
		p := repository.MemberParams{
			Relation:     m.Relation,
			FullName:     m.FullName,
			BirthYear:    m.BirthYear,
			AnnualIncome: m.AnnualIncome,
		}
		if m.NationalID != nil && *m.NationalID != "" {
			h := nationalid.Hash(pepper, *m.NationalID)
			p.NationalIDHash = &h
		}
		out = append(out, p)
	}
	return out
}
