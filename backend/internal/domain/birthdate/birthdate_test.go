package birthdate_test

import (
	"testing"
	"time"

	"github.com/supel2nova/welfare-registration/backend/internal/domain/birthdate"
)

func p(v int) *int { return &v }

var today = time.Date(2026, 8, 27, 15, 30, 0, 0, time.UTC)

func TestValidate(t *testing.T) {
	cases := []struct {
		name string
		in   birthdate.BirthDate
		want []string
	}{
		{"ครบสามช่อง", birthdate.BirthDate{Year: 1985, Month: p(3), Day: p(12), Precision: birthdate.Full}, nil},
		{"วันเกิดวันนี้ ผ่าน", birthdate.BirthDate{Year: 2026, Month: p(8), Day: p(27), Precision: birthdate.Full}, nil},
		{"วันเกิดพรุ่งนี้ ไม่ผ่าน", birthdate.BirthDate{Year: 2026, Month: p(8), Day: p(28), Precision: birthdate.Full}, []string{"birth_day"}},
		{"31 ก.พ. ไม่มีจริง", birthdate.BirthDate{Year: 1985, Month: p(2), Day: p(31), Precision: birthdate.Full}, []string{"birth_day"}},
		{"29 ก.พ. ปีอธิกสุรทิน ผ่าน", birthdate.BirthDate{Year: 1984, Month: p(2), Day: p(29), Precision: birthdate.Full}, nil},
		{"29 ก.พ. ปีปกติ ไม่ผ่าน", birthdate.BirthDate{Year: 1985, Month: p(2), Day: p(29), Precision: birthdate.Full}, []string{"birth_day"}},
		{"ไม่ทราบวันและเดือน ผ่าน", birthdate.BirthDate{Year: 1950, Precision: birthdate.YearOnly}, nil},
		{"ทราบแค่เดือน ผ่าน", birthdate.BirthDate{Year: 1950, Month: p(7), Precision: birthdate.YearMonth}, nil},
		{"YEAR_ONLY แต่ส่งเดือนมา", birthdate.BirthDate{Year: 1950, Month: p(3), Precision: birthdate.YearOnly}, []string{"birth_month"}},
		{"YEAR_MONTH แต่ส่งวันมา", birthdate.BirthDate{Year: 1950, Month: p(3), Day: p(1), Precision: birthdate.YearMonth}, []string{"birth_day"}},
		{"FULL แต่ไม่ส่งวัน", birthdate.BirthDate{Year: 1950, Month: p(3), Precision: birthdate.Full}, []string{"birth_day"}},
		{"FULL แต่ไม่ส่งทั้งวันและเดือน", birthdate.BirthDate{Year: 1950, Precision: birthdate.Full}, []string{"birth_month", "birth_day"}},
		{"precision ไม่รู้จัก", birthdate.BirthDate{Year: 1950, Precision: "SOMETHING"}, []string{"birth_precision"}},
		{"ปีเก่าเกินไป", birthdate.BirthDate{Year: 1899, Precision: birthdate.YearOnly}, []string{"birth_year"}},
		{"ปีอนาคต", birthdate.BirthDate{Year: 2027, Precision: birthdate.YearOnly}, []string{"birth_year"}},
		{"เดือน 13", birthdate.BirthDate{Year: 1985, Month: p(13), Day: p(1), Precision: birthdate.Full}, []string{"birth_month"}},
		{"เดือน 0", birthdate.BirthDate{Year: 1985, Month: p(0), Day: p(1), Precision: birthdate.Full}, []string{"birth_month"}},
		{"ผิดทั้งปีและเดือน", birthdate.BirthDate{Year: 1800, Month: p(99), Day: p(1), Precision: birthdate.Full}, []string{"birth_year", "birth_month"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.in.Validate(today)
			if len(got) != len(c.want) {
				t.Fatalf("Validate = %v want %v", got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Fatalf("Validate = %v want %v", got, c.want)
				}
			}
		})
	}
}

func TestEarliest(t *testing.T) {
	cases := []struct {
		name string
		in   birthdate.BirthDate
		want time.Time
	}{
		{"ครบสามช่อง ใช้ตามจริง", birthdate.BirthDate{Year: 1985, Month: p(3), Day: p(12), Precision: birthdate.Full},
			time.Date(1985, 3, 12, 0, 0, 0, 0, time.UTC)},
		{"ไม่ทราบวัน นับวันที่ 1", birthdate.BirthDate{Year: 1985, Month: p(3), Precision: birthdate.YearMonth},
			time.Date(1985, 3, 1, 0, 0, 0, 0, time.UTC)},
		{"ไม่ทราบวันและเดือน นับ 1 ม.ค.", birthdate.BirthDate{Year: 1985, Precision: birthdate.YearOnly},
			time.Date(1985, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.in.Earliest(); !got.Equal(c.want) {
				t.Errorf("Earliest = %v want %v", got, c.want)
			}
		})
	}
}

func TestAgeAt(t *testing.T) {
	cases := []struct {
		name string
		in   birthdate.BirthDate
		want int
	}{
		{"เกิดวันนี้เมื่อ 41 ปีก่อน", birthdate.BirthDate{Year: 1985, Month: p(8), Day: p(27), Precision: birthdate.Full}, 41},
		{"อีกหนึ่งวันถึงวันเกิด ยังไม่ครบปี", birthdate.BirthDate{Year: 1985, Month: p(8), Day: p(28), Precision: birthdate.Full}, 40},
		{"★ ไม่ทราบวันเดือน ได้อายุมากสุด = เป็นคุณกับผู้ยื่น", birthdate.BirthDate{Year: 1985, Precision: birthdate.YearOnly}, 41},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.in.AgeAt(today); got != c.want {
				t.Errorf("AgeAt = %d want %d", got, c.want)
			}
		})
	}
}
