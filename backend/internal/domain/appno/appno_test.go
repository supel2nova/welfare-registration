package appno_test

import (
	"testing"

	"github.com/supel2nova/welfare-registration/backend/internal/domain/appno"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		name string
		year int
		seq  int64
		want string
	}{
		{"ใบแรกของปี", 2026, 1, "WC-2026-0000001"},
		{"เลขกลางๆ", 2026, 42, "WC-2026-0000042"},
		{"เต็ม 7 หลักพอดี", 2026, 9999999, "WC-2026-9999999"},
		{"เกิน 7 หลัก ยาวขึ้นเองไม่ตัดทิ้ง", 2026, 12345678, "WC-2026-12345678"},
		{"คนละปีงบประมาณ", 2027, 1, "WC-2027-0000001"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := appno.Format(c.year, c.seq)
			if got != c.want {
				t.Fatalf("Format = %q want %q", got, c.want)
			}
			if len(got) > 20 {
				t.Fatalf("ยาว %d เกิน varchar(20)", len(got))
			}
			if !appno.IsValid(got) {
				t.Fatalf("สร้างเองแล้วตรวจไม่ผ่าน: %q", got)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"รูปแบบถูก", "WC-2026-0000001", true},
		{"ไม่มี prefix", "2026-0000001", false},
		{"prefix ผิด", "XX-2026-0000001", false},
		{"running สั้นไป", "WC-2026-000001", false},
		{"ปีสั้นไป", "WC-226-0000001", false},
		{"ว่าง", "", false},
		{"มีตัวอักษรใน running", "WC-2026-000000A", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := appno.IsValid(c.in); got != c.want {
				t.Errorf("IsValid(%q) = %v want %v", c.in, got, c.want)
			}
		})
	}
}
