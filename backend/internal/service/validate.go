package service

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/supel2nova/welfare-registration/backend/internal/domain"
	"github.com/supel2nova/welfare-registration/backend/internal/domain/birthdate"
	"github.com/supel2nova/welfare-registration/backend/internal/domain/laserid"
	"github.com/supel2nova/welfare-registration/backend/internal/domain/nationalid"
	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
)

var (
	thaiName      = regexp.MustCompile(`^[\p{Thai} .-]+$`)
	thaiConsonant = regexp.MustCompile(`[\x{0E01}-\x{0E2E}]`)
	thaiDigit     = regexp.MustCompile(`[\x{0E50}-\x{0E59}]`)
	mobile        = regexp.MustCompile(`^0[689]\d{8}$`)
	landline      = regexp.MustCompile(`^0[2-7]\d{7}$`)
	amount        = regexp.MustCompile(`^\d{1,12}(\.\d{1,2})?$`)
	digits2       = regexp.MustCompile(`^\d{2}$`)
	digits4       = regexp.MustCompile(`^\d{4}$`)
	digits5       = regexp.MustCompile(`^\d{5}$`)
	digits6       = regexp.MustCompile(`^\d{6}$`)
	maxNameLen    = 100
)

type errs []apperror.FieldError

func (e *errs) add(path, code string) {
	*e = append(*e, apperror.Field(path, code))
}

func (e *errs) addMsg(path, code, msg string) {
	*e = append(*e, apperror.FieldMsg(path, code, msg))
}

func validateFormat(req dto.CreateApplicationRequest, today time.Time) []apperror.FieldError {
	var e errs

	if req.FiscalYear < domain.MinFiscalYear || req.FiscalYear > today.Year()+1 {
		e.add("fiscal_year", apperror.CodeFiscalYear)
	}

	validatePersonal(&e, req.Personal, today)
	validateFamily(&e, req.Family, req.Personal.NationalID)
	validateFinancial(&e, req.Financial)

	return e
}

func validatePersonal(e *errs, p dto.Personal, today time.Time) {
	if !nationalid.IsValid(p.NationalID) {
		e.add("personal.national_id", apperror.CodeNationalID)
	}

	validateIDVerification(e, p)

	if !domain.IsValidTitle(domain.Title(p.Title)) {
		e.add("personal.title", apperror.CodeEnum)
	}
	if !isThaiName(p.FirstName, maxNameLen) {
		e.add("personal.first_name", apperror.CodeName)
	}
	if !isThaiName(p.LastName, maxNameLen) {
		e.add("personal.last_name", apperror.CodeName)
	}

	bd := birthdate.BirthDate{
		Year:      p.BirthYear,
		Month:     p.BirthMonth,
		Day:       p.BirthDay,
		Precision: birthdate.Precision(p.BirthPrecision),
	}
	for _, field := range bd.Validate(today) {
		e.add("personal."+field, apperror.CodeBirthDate)
	}

	if !isValidPhone(p.Phone) {
		e.add("personal.phone", apperror.CodePhone)
	}

	validateAddress(e, p.Address)
}

func validateIDVerification(e *errs, p dto.Personal) {
	method := domain.IDVerifyMethod(p.IDVerifyMethod)
	if !domain.IsValidIDVerifyMethod(method) {
		e.add("personal.id_verify_method", apperror.CodeEnum)
		return
	}

	hasLaser := p.LaserID != nil && *p.LaserID != ""

	switch method {
	case domain.VerifyLaserCode:
		if !hasLaser {
			e.addMsg("personal.laser_id", apperror.CodeIDVerification,
				"ไม่มีรหัสหลังบัตร ต้องเลือกวิธียืนยันตัวตนแบบอื่น")
			return
		}
		if !laserid.IsValidFormat(*p.LaserID) {
			e.add("personal.laser_id", apperror.CodeLaserID)
		}
	case domain.VerifyManualCardCheck:
		if hasLaser {
			e.addMsg("personal.laser_id", apperror.CodeIDVerification,
				"มีรหัสหลังบัตรแล้ว ต้องยืนยันด้วยรหัสหลังบัตร")
		}
		note := ""
		if p.IDVerifyNote != nil {
			note = strings.TrimSpace(*p.IDVerifyNote)
		}
		if utf8.RuneCountInString(note) < domain.MinManualNoteLen {
			e.addMsg("personal.id_verify_note", apperror.CodeIDVerification,
				fmt.Sprintf("ต้องระบุเหตุผลอย่างน้อย %d ตัวอักษร", domain.MinManualNoteLen))
		}
	}
}

func validateAddress(e *errs, a dto.Address) {
	house := strings.TrimSpace(a.HouseNo)
	if house == "" || utf8.RuneCountInString(house) > 50 {
		e.add("personal.address.house_no", apperror.CodeAddressMissing)
	}
	if !digits2.MatchString(a.ProvinceCode) {
		e.add("personal.address.province_code", apperror.CodeAddress)
	}
	if !digits4.MatchString(a.DistrictCode) {
		e.add("personal.address.district_code", apperror.CodeAddress)
	}
	if !digits6.MatchString(a.SubdistrictCode) {
		e.add("personal.address.subdistrict_code", apperror.CodeAddress)
	}
	if !digits5.MatchString(a.PostalCode) {
		e.add("personal.address.postal_code", apperror.CodeAddress)
	}
}

func validateFamily(e *errs, f *dto.Family, applicantID string) {
	if f == nil {
		return
	}
	if f.MaritalStatus != nil && !domain.IsValidMaritalStatus(domain.MaritalStatus(*f.MaritalStatus)) {
		e.add("family.marital_status", apperror.CodeEnum)
	}
	if len(f.Members) > domain.MaxHouseholdMembers {
		e.addMsg("family.members", apperror.CodeTooMany,
			fmt.Sprintf("ระบุสมาชิกได้ไม่เกิน %d คน", domain.MaxHouseholdMembers))
		return
	}

	spouse := 0
	for i, m := range f.Members {
		path := fmt.Sprintf("family.members[%d]", i)

		if !domain.IsValidRelation(domain.Relation(m.Relation)) {
			e.add(path+".relation", apperror.CodeEnum)
		}
		if domain.Relation(m.Relation) == domain.RelationSpouse {
			spouse++
			if spouse > domain.MaxSpouse {
				e.add(path+".relation", apperror.CodeSpouseCount)
			}
		}
		if !isThaiName(m.FullName, 200) {
			e.add(path+".full_name", apperror.CodeName)
		}
		if m.NationalID != nil && *m.NationalID != "" {
			switch {
			case !nationalid.IsValid(*m.NationalID):
				e.add(path+".national_id", apperror.CodeNationalID)
			case *m.NationalID == applicantID:
				e.add(path+".national_id", apperror.CodeMemberDupSelf)
			}
		}
		if m.AnnualIncome != nil && *m.AnnualIncome < 0 {
			e.add(path+".annual_income", apperror.CodeNegative)
		}
	}
}

func validateFinancial(e *errs, f dto.Financial) {
	if len(f.IncomeSources) > domain.MaxIncomeSources {
		e.addMsg("financial.income_sources", apperror.CodeTooMany,
			fmt.Sprintf("ระบุแหล่งรายได้ได้ไม่เกิน %d รายการ", domain.MaxIncomeSources))
	} else {
		for i, s := range f.IncomeSources {
			path := fmt.Sprintf("financial.income_sources[%d]", i)
			if !domain.IsValidIncomeType(domain.IncomeSourceType(s.SourceType)) {
				e.add(path+".source_type", apperror.CodeEnum)
			}
			if s.AnnualAmount < 0 {
				e.add(path+".annual_amount", apperror.CodeNegative)
			}
		}
	}

	if f.ExpenseToOthers < 0 {
		e.add("financial.expense_to_others", apperror.CodeNegative)
	}

	if len(f.Assets) > domain.MaxAssets {
		e.addMsg("financial.assets", apperror.CodeTooMany,
			fmt.Sprintf("ระบุทรัพย์สินได้ไม่เกิน %d รายการ", domain.MaxAssets))
	} else {
		for i, a := range f.Assets {
			validateAsset(e, fmt.Sprintf("financial.assets[%d]", i), a)
		}
	}

	if len(f.Liabilities) > domain.MaxLiabilities {
		e.addMsg("financial.liabilities", apperror.CodeTooMany,
			fmt.Sprintf("ระบุหนี้สินได้ไม่เกิน %d รายการ", domain.MaxLiabilities))
	} else {
		for i, l := range f.Liabilities {
			path := fmt.Sprintf("financial.liabilities[%d]", i)
			if !domain.IsValidLiabilityType(domain.LiabilityType(l.LiabilityType)) {
				e.add(path+".liability_type", apperror.CodeEnum)
			}
			if l.CreditLimit < 0 {
				e.add(path+".credit_limit", apperror.CodeNegative)
			}
		}
	}
}

func validateAsset(e *errs, path string, a dto.Asset) {
	want, known := domain.UnitFor(domain.AssetType(a.AssetType))
	if !known {
		e.add(path+".asset_type", apperror.CodeEnum)
	} else if domain.Unit(a.Unit) != want {
		e.addMsg(path+".unit", apperror.CodeAssetUnit,
			fmt.Sprintf("หน่วยของ %s ต้องเป็น %s", a.AssetType, want))
	}
	if !amount.MatchString(a.Amount.String()) {
		e.add(path+".amount", apperror.CodeNegative)
	}
	if a.JointAccountHolders != nil && *a.JointAccountHolders < 1 {
		e.add(path+".joint_account_holders", apperror.CodeNegative)
	}
}

func isValidPhone(s string) bool {
	return mobile.MatchString(s) || landline.MatchString(s)
}

func isThaiName(s string, max int) bool {
	s = strings.TrimSpace(s)
	n := utf8.RuneCountInString(s)
	if n < 1 || n > max {
		return false
	}
	if !thaiName.MatchString(s) || thaiDigit.MatchString(s) {
		return false
	}
	return thaiConsonant.MatchString(s)
}
