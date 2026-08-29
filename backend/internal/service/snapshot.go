package service

import (
	"encoding/json"

	"github.com/supel2nova/welfare-registration/backend/internal/domain/nationalid"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/internal/repository"
)

type snapshotPersonal struct {
	NationalIDMask string      `json:"national_id_mask"`
	IDVerifyMethod string      `json:"id_verify_method"`
	IDVerifyNote   *string     `json:"id_verify_note"`
	Title          string      `json:"title"`
	FirstName      string      `json:"first_name"`
	LastName       string      `json:"last_name"`
	BirthYear      int         `json:"birth_year"`
	BirthMonth     *int        `json:"birth_month"`
	BirthDay       *int        `json:"birth_day"`
	BirthPrecision string      `json:"birth_precision"`
	Phone          string      `json:"phone"`
	IsFarmer       bool        `json:"is_farmer"`
	Address        dto.Address `json:"address"`
}

type snapshotAddressNames struct {
	SubdistrictName string `json:"subdistrict_name"`
	DistrictName    string `json:"district_name"`
	ProvinceName    string `json:"province_name"`
}

type snapshotMember struct {
	Relation       string  `json:"relation"`
	NationalIDMask *string `json:"national_id_mask"`
	FullName       string  `json:"full_name"`
	BirthYear      *int    `json:"birth_year"`
	AnnualIncome   *int64  `json:"annual_income"`
}

type snapshotFamily struct {
	MaritalStatus *string          `json:"marital_status"`
	Members       []snapshotMember `json:"members"`
}

type applicantSnapshot struct {
	FiscalYear   int                  `json:"fiscal_year"`
	Personal     snapshotPersonal     `json:"personal"`
	AddressNames snapshotAddressNames `json:"address_names"`
	Family       *snapshotFamily      `json:"family"`
	Financial    dto.Financial        `json:"financial"`
}

func buildSnapshot(req dto.CreateApplicationRequest, addr repository.ResolvedAddress) ([]byte, error) {
	s := applicantSnapshot{
		FiscalYear: req.FiscalYear,
		Personal: snapshotPersonal{
			NationalIDMask: nationalid.Mask(req.Personal.NationalID),
			IDVerifyMethod: req.Personal.IDVerifyMethod,
			IDVerifyNote:   req.Personal.IDVerifyNote,
			Title:          req.Personal.Title,
			FirstName:      req.Personal.FirstName,
			LastName:       req.Personal.LastName,
			BirthYear:      req.Personal.BirthYear,
			BirthMonth:     req.Personal.BirthMonth,
			BirthDay:       req.Personal.BirthDay,
			BirthPrecision: req.Personal.BirthPrecision,
			Phone:          req.Personal.Phone,
			IsFarmer:       req.Personal.IsFarmer,
			Address:        req.Personal.Address,
		},
		AddressNames: snapshotAddressNames{
			SubdistrictName: addr.SubdistrictName,
			DistrictName:    addr.DistrictName,
			ProvinceName:    addr.ProvinceName,
		},
		Family:    maskFamily(req.Family),
		Financial: req.Financial,
	}
	return json.Marshal(s)
}

func maskFamily(f *dto.Family) *snapshotFamily {
	if f == nil {
		return nil
	}
	out := &snapshotFamily{
		MaritalStatus: f.MaritalStatus,
		Members:       make([]snapshotMember, 0, len(f.Members)),
	}
	for _, m := range f.Members {
		member := snapshotMember{
			Relation:     m.Relation,
			FullName:     m.FullName,
			BirthYear:    m.BirthYear,
			AnnualIncome: m.AnnualIncome,
		}
		if m.NationalID != nil {
			masked := nationalid.Mask(*m.NationalID)
			member.NationalIDMask = &masked
		}
		out.Members = append(out.Members, member)
	}
	return out
}
