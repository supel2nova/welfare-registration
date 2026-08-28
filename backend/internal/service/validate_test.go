package service

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/supel2nova/welfare-registration/backend/internal/dto"
	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
)

var today = time.Date(2026, 8, 27, 15, 30, 0, 0, time.UTC)

func s(v string) *string { return &v }
func i(v int) *int       { return &v }
func i64(v int64) *int64 { return &v }

func validRequest() dto.CreateApplicationRequest {
	return dto.CreateApplicationRequest{
		FiscalYear: 2026,
		Personal: dto.Personal{
			NationalID:     "1234567890121",
			LaserID:        s("JT8-1234567-89"),
			IDVerifyMethod: "LASER_CODE",
			Title:          "นาย",
			FirstName:      "สมชาย",
			LastName:       "ใจดี",
			BirthYear:      1985,
			BirthMonth:     i(3),
			BirthDay:       i(12),
			BirthPrecision: "FULL",
			Phone:          "0812345678",
			IsFarmer:       true,
			Address: dto.Address{
				HouseNo:         "12/34",
				Moo:             s("5"),
				ProvinceCode:    "50",
				DistrictCode:    "5001",
				SubdistrictCode: "500101",
				PostalCode:      "50200",
			},
		},
		Family: &dto.Family{
			MaritalStatus: s("MARRIED"),
			Members: []dto.Member{
				{Relation: "SPOUSE", NationalID: s("1234567890139"), FullName: "นางสมหญิง ใจดี", BirthYear: i(1988), AnnualIncome: i64(40000)},
			},
		},
		Financial: dto.Financial{
			IncomeSources:   []dto.IncomeSource{{SourceType: "AGRI", AnnualAmount: 45000}},
			ExpenseToOthers: 0,
			Assets: []dto.Asset{
				{AssetType: "DEPOSIT", Amount: json.Number("15000"), Unit: "THB", JointAccountHolders: i(1)},
				{AssetType: "LAND_AGRI", Amount: json.Number("8"), Unit: "RAI"},
			},
			Liabilities:   []dto.Liability{{LiabilityType: "LOAN_PERSONAL", CreditLimit: 30000}},
			HasCreditCard: false,
		},
	}
}

func has(fes []apperror.FieldError, field, code string) bool {
	for _, fe := range fes {
		if fe.Field == field && fe.Code == code {
			return true
		}
	}
	return false
}

func TestValidateFormat(t *testing.T) {
	cases := []struct {
		name  string
		mut   func(*dto.CreateApplicationRequest)
		field string
		code  string
	}{
		{"#1 เลขบัตร 12 หลัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.NationalID = "123456789012"
		}, "personal.national_id", apperror.CodeNationalID},

		{"#2 เลขบัตร checksum ผิด", func(r *dto.CreateApplicationRequest) {
			r.Personal.NationalID = "1111111111111"
		}, "personal.national_id", apperror.CodeNationalID},

		{"#3 รหัสหลังบัตรผิดรูปแบบ", func(r *dto.CreateApplicationRequest) {
			r.Personal.LaserID = s("ABC123")
		}, "personal.laser_id", apperror.CodeLaserID},

		{"#4 ชื่อมีตัวเลข", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "สมชาย123"
		}, "personal.first_name", apperror.CodeName},

		{"#6 วันเกิดพรุ่งนี้", func(r *dto.CreateApplicationRequest) {
			r.Personal.BirthYear, r.Personal.BirthMonth, r.Personal.BirthDay = 2026, i(8), i(28)
		}, "personal.birth_day", apperror.CodeBirthDate},

		{"#7 วันที่ 31 กุมภาพันธ์", func(r *dto.CreateApplicationRequest) {
			r.Personal.BirthMonth, r.Personal.BirthDay = i(2), i(31)
		}, "personal.birth_day", apperror.CodeBirthDate},

		{"#9 YEAR_ONLY แต่ส่งเดือนมาด้วย", func(r *dto.CreateApplicationRequest) {
			r.Personal.BirthPrecision, r.Personal.BirthDay = "YEAR_ONLY", nil
		}, "personal.birth_month", apperror.CodeBirthDate},

		{"#10 เบอร์โทรผิดรูปแบบ", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "123456"
		}, "personal.phone", apperror.CodePhone},

		{"มือถือ 9 หลัก ขาดไปหนึ่งตัว", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "081234567"
		}, "personal.phone", apperror.CodePhone},

		{"เบอร์บ้าน 8 หลัก สั้นไป", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "02123456"
		}, "personal.phone", apperror.CodePhone},

		{"เบอร์บ้าน 10 หลัก ยาวไป", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "0212345678"
		}, "personal.phone", apperror.CodePhone},

		{"ส่งมาพร้อมขีดคั่น", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "099-119-2231"
		}, "personal.phone", apperror.CodePhone},

		{"ขึ้นต้นด้วย 1", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "1234567890"
		}, "personal.phone", apperror.CodePhone},

		{"ขึ้นต้นด้วย 9 ไม่ใช่ 0", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "9912345678"
		}, "personal.phone", apperror.CodePhone},

		{"ส่งมาเป็นรูปแบบสากล +66", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "+66812345678"
		}, "personal.phone", apperror.CodePhone},

		{"ตัด 0 ออกแล้วใส่ 66 นำหน้า", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "66812345678"
		}, "personal.phone", apperror.CodePhone},

		{"มือถือขึ้นต้น 07 (ไม่มีจริง)", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "0712345678"
		}, "personal.phone", apperror.CodePhone},

		{"#11 รายได้ติดลบ", func(r *dto.CreateApplicationRequest) {
			r.Financial.IncomeSources[0].AnnualAmount = -100
		}, "financial.income_sources[0].annual_amount", apperror.CodeNegative},

		{"#13 LAND_AGRI แต่หน่วยเป็นบาท", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets[1].Unit = "THB"
		}, "financial.assets[1].unit", apperror.CodeAssetUnit},

		{"#15 คู่สมรส 2 คน", func(r *dto.CreateApplicationRequest) {
			r.Family.Members = append(r.Family.Members, dto.Member{Relation: "SPOUSE", FullName: "นางสาวสมศรี ใจดี"})
		}, "family.members[1].relation", apperror.CodeSpouseCount},

		{"#16 เลขบัตรสมาชิกซ้ำกับผู้ยื่น", func(r *dto.CreateApplicationRequest) {
			r.Family.Members[0].NationalID = s("1234567890121")
		}, "family.members[0].national_id", apperror.CodeMemberDupSelf},

		{"#19 MANUAL แต่ไม่ระบุเหตุผล", func(r *dto.CreateApplicationRequest) {
			r.Personal.LaserID, r.Personal.IDVerifyMethod = nil, "MANUAL_CARD_CHECK"
		}, "personal.id_verify_note", apperror.CodeIDVerification},

		{"#19b MANUAL แต่เหตุผลสั้นเกินไป", func(r *dto.CreateApplicationRequest) {
			r.Personal.LaserID, r.Personal.IDVerifyMethod = nil, "MANUAL_CARD_CHECK"
			r.Personal.IDVerifyNote = s("เลือน")
		}, "personal.id_verify_note", apperror.CodeIDVerification},

		{"#20 ไม่มีรหัสหลังบัตร แต่เลือก LASER_CODE", func(r *dto.CreateApplicationRequest) {
			r.Personal.LaserID = nil
		}, "personal.laser_id", apperror.CodeIDVerification},

		{"มีรหัสหลังบัตร แต่เลือก MANUAL", func(r *dto.CreateApplicationRequest) {
			r.Personal.IDVerifyMethod = "MANUAL_CARD_CHECK"
			r.Personal.IDVerifyNote = s("บัตรตลอดชีพ รหัสเลือนอ่านไม่ออก")
		}, "personal.laser_id", apperror.CodeIDVerification},

		{"วิธียืนยันตัวตนไม่รู้จัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.IDVerifyMethod = "PENDING_VERIFICATION"
		}, "personal.id_verify_method", apperror.CodeEnum},

		{"คำนำหน้าไม่รู้จัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.Title = "คุณ"
		}, "personal.title", apperror.CodeEnum},

		{"ปีงบประมาณเก่าเกินไป", func(r *dto.CreateApplicationRequest) {
			r.FiscalYear = 2019
		}, "fiscal_year", apperror.CodeFiscalYear},

		{"ปีงบประมาณล้ำไปสองปี", func(r *dto.CreateApplicationRequest) {
			r.FiscalYear = 2028
		}, "fiscal_year", apperror.CodeFiscalYear},

		{"บ้านเลขที่ว่าง", func(r *dto.CreateApplicationRequest) {
			r.Personal.Address.HouseNo = "  "
		}, "personal.address.house_no", apperror.CodeAddressMissing},

		{"รหัสตำบลไม่ครบ 6 หลัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.Address.SubdistrictCode = "5001"
		}, "personal.address.subdistrict_code", apperror.CodeAddress},

		{"สมาชิกเกิน 20 คน", func(r *dto.CreateApplicationRequest) {
			for len(r.Family.Members) <= 20 {
				r.Family.Members = append(r.Family.Members, dto.Member{Relation: "CHILD", FullName: "เด็กชายสมหมาย ใจดี"})
			}
		}, "family.members", apperror.CodeTooMany},

		{"ประเภททรัพย์สินไม่รู้จัก", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets[0].AssetType = "SPACESHIP"
		}, "financial.assets[0].asset_type", apperror.CodeEnum},

		{"จำนวนทรัพย์สินติดลบ", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets[0].Amount = json.Number("-15000")
		}, "financial.assets[0].amount", apperror.CodeNegative},

		{"จำนวนทรัพย์สินทศนิยมเกิน 2 ตำแหน่ง", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets[1].Amount = json.Number("8.555")
		}, "financial.assets[1].amount", apperror.CodeNegative},

		{"เจ้าของบัญชีร่วมเป็นศูนย์", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets[0].JointAccountHolders = i(0)
		}, "financial.assets[0].joint_account_holders", apperror.CodeNegative},

		{"ค่าใช้จ่ายเลี้ยงดูติดลบ", func(r *dto.CreateApplicationRequest) {
			r.Financial.ExpenseToOthers = -1
		}, "financial.expense_to_others", apperror.CodeNegative},

		{"ประเภทหนี้สินไม่รู้จัก", func(r *dto.CreateApplicationRequest) {
			r.Financial.Liabilities[0].LiabilityType = "LOAN_MOON"
		}, "financial.liabilities[0].liability_type", apperror.CodeEnum},

		{"สถานภาพสมรสไม่รู้จัก", func(r *dto.CreateApplicationRequest) {
			r.Family.MaritalStatus = s("COMPLICATED")
		}, "family.marital_status", apperror.CodeEnum},

		{"ชื่อสมาชิกมีตัวเลข", func(r *dto.CreateApplicationRequest) {
			r.Family.Members[0].FullName = "สมหญิง 2"
		}, "family.members[0].full_name", apperror.CodeName},

		{"ชื่อเป็นจุดล้วน ไม่มีตัวอักษรสักตัว", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "..."
		}, "personal.first_name", apperror.CodeName},

		{"ชื่อเป็นไม้ยมกล้วน", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "ๆๆๆ"
		}, "personal.first_name", apperror.CodeName},

		{"ชื่อมีเลขไทย", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "สมชาย๑๒๓"
		}, "personal.first_name", apperror.CodeName},

		{"ชื่อเป็นอักษรโรมัน", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "John"
		}, "personal.first_name", apperror.CodeName},

		{"ชื่อไทยปนอังกฤษ", func(r *dto.CreateApplicationRequest) {
			r.Personal.LastName = "ใจดี Smith"
		}, "personal.last_name", apperror.CodeName},

		{"ชื่อว่างเปล่า", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "   "
		}, "personal.first_name", apperror.CodeName},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := validRequest()
			c.mut(&req)
			got := validateFormat(req, today)
			if !has(got, c.field, c.code) {
				t.Errorf("ไม่พบ %s/%s\nได้: %+v", c.field, c.code, got)
			}
		})
	}
}

func TestValidateFormatPasses(t *testing.T) {
	cases := []struct {
		name string
		mut  func(*dto.CreateApplicationRequest)
	}{
		{"ใบสมัครที่ถูกต้องครบถ้วน", func(r *dto.CreateApplicationRequest) {}},

		{"#5 วันเกิดเป็นวันนี้", func(r *dto.CreateApplicationRequest) {
			r.Personal.BirthYear, r.Personal.BirthMonth, r.Personal.BirthDay = 2026, i(8), i(27)
		}},

		{"#8 ไม่ทราบวันและเดือนเกิด", func(r *dto.CreateApplicationRequest) {
			r.Personal.BirthPrecision = "YEAR_ONLY"
			r.Personal.BirthMonth, r.Personal.BirthDay = nil, nil
		}},

		{"ทราบแค่ปีและเดือนเกิด", func(r *dto.CreateApplicationRequest) {
			r.Personal.BirthPrecision, r.Personal.BirthDay = "YEAR_MONTH", nil
		}},

		{"#12 ไม่มีรายได้เลย", func(r *dto.CreateApplicationRequest) {
			r.Financial.IncomeSources = nil
		}},

		{"มือถือ AIS 10 หลัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "0991192231"
		}},

		{"มือถือ 06 นำหน้า", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "0645556677"
		}},

		{"เบอร์บ้านกรุงเทพ 9 หลัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "021234567"
		}},

		{"เบอร์บ้านต่างจังหวัด 9 หลัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "053123456"
		}},

		{"เบอร์บ้านภาคใต้ 9 หลัก", func(r *dto.CreateApplicationRequest) {
			r.Personal.Phone = "074123456"
		}},

		{"#14 ไม่กรอกข้อมูลครอบครัว", func(r *dto.CreateApplicationRequest) {
			r.Family = nil
		}},

		{"#18 บัตรตลอดชีพ ยืนยันด้วยตาพร้อมเหตุผล", func(r *dto.CreateApplicationRequest) {
			r.Personal.LaserID, r.Personal.IDVerifyMethod = nil, "MANUAL_CARD_CHECK"
			r.Personal.IDVerifyNote = s("บัตรตลอดชีพ รหัสหลังบัตรเลือนอ่านไม่ออก ตรวจกับทะเบียนบ้านแล้ว")
		}},

		{"อ่านชิปแทน ไม่ต้องมีรหัสหลังบัตร", func(r *dto.CreateApplicationRequest) {
			r.Personal.LaserID, r.Personal.IDVerifyMethod = nil, "CHIP_READ"
		}},

		{"ที่ดิน 8.5 ไร่", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets[1].Amount = json.Number("8.50")
		}},

		{"ไม่มีทรัพย์สินและหนี้สินเลย", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets, r.Financial.Liabilities = nil, nil
		}},

		{"ไม่ระบุจำนวนเจ้าของบัญชีร่วม", func(r *dto.CreateApplicationRequest) {
			r.Financial.Assets[0].JointAccountHolders = nil
		}},

		{"สมาชิก 20 คนพอดี", func(r *dto.CreateApplicationRequest) {
			for len(r.Family.Members) < 20 {
				r.Family.Members = append(r.Family.Members, dto.Member{Relation: "CHILD", FullName: "เด็กชายสมหมาย ใจดี"})
			}
		}},

		{"ปีงบประมาณล่วงหน้าหนึ่งปี", func(r *dto.CreateApplicationRequest) {
			r.FiscalYear = 2027
		}},

		{"ชื่อมีจุดและเว้นวรรค", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "สมชาย ณ ป้อมเพชร"
		}},

		{"นามสกุลมีขีดคั่น (ชื่อไทยมุสลิม)", func(r *dto.CreateApplicationRequest) {
			r.Personal.LastName = "อัล-ฮาซัน"
		}},

		{"ชื่อมีคำนำหน้าราชสกุลย่อ", func(r *dto.CreateApplicationRequest) {
			r.Personal.FirstName = "ม.ร.ว.สมชาย"
		}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := validRequest()
			c.mut(&req)
			if got := validateFormat(req, today); len(got) != 0 {
				t.Errorf("ต้องผ่าน แต่ได้: %+v", got)
			}
		})
	}
}

func TestValidate_CollectsAllErrors(t *testing.T) {
	req := validRequest()
	req.Personal.NationalID = "1111111111111"
	req.Personal.FirstName = "สมชาย123"
	req.Personal.Phone = "123456"

	got := validateFormat(req, today)

	if len(got) != 3 {
		t.Fatalf("ผิด 3 ช่อง ต้องได้ 3 รายการ ได้ %d: %+v", len(got), got)
	}
	for _, want := range []struct{ field, code string }{
		{"personal.national_id", apperror.CodeNationalID},
		{"personal.first_name", apperror.CodeName},
		{"personal.phone", apperror.CodePhone},
	} {
		if !has(got, want.field, want.code) {
			t.Errorf("ขาด %s/%s", want.field, want.code)
		}
	}
}

func TestValidateNeverEchoesInput(t *testing.T) {
	req := validRequest()
	req.Personal.NationalID = "9999999999999"
	req.Personal.LaserID = s("XY9-7654321-01")
	req.Personal.Phone = "0999999999x"

	for _, fe := range validateFormat(req, today) {
		for _, secret := range []string{"9999999999999", "XY9-7654321-01", "0999999999x"} {
			if contains(fe.Message, secret) {
				t.Errorf("message สะท้อนค่าที่ผู้ใช้กรอกกลับไป: %q", fe.Message)
			}
		}
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
