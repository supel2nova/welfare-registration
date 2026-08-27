package nationalid_test

import (
	"testing"

	"github.com/supel2nova/welfare-registration/backend/internal/domain/nationalid"
)

func TestIsValid(t *testing.T) {
	cases := []struct {
		name string
		id   string
		want bool
	}{
		{"เลขถูกต้อง (fixture ผู้ยื่น)", "1234567890121", true},
		{"เลขถูกต้อง (fixture คู่สมรส)", "1234567890139", true},
		{"12 หลัก สั้นไป", "123456789012", false},
		{"14 หลัก ยาวไป", "12345678901234", false},
		{"checksum ผิด", "1111111111111", false},
		{"มีตัวอักษรปน", "12345678901a1", false},
		{"หลักสุดท้ายเป็นตัวอักษร", "123456789012x", false},
		{"ว่าง", "", false},
		{"เว้นวรรคแทนเลข", "1234567890 21", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nationalid.IsValid(c.id); got != c.want {
				t.Errorf("IsValid(%q) = %v want %v", c.id, got, c.want)
			}
		})
	}
}

func TestHash(t *testing.T) {
	const id = "1234567890121"

	t.Run("เลขเดิม pepper เดิม ได้ค่าเดิมเสมอ", func(t *testing.T) {
		if nationalid.Hash("p", id) != nationalid.Hash("p", id) {
			t.Fatal("hash ไม่ deterministic")
		}
	})
	t.Run("ยาว 64 ตัวพอดี ลงคอลัมน์ varchar(64)", func(t *testing.T) {
		if got := len(nationalid.Hash("p", id)); got != 64 {
			t.Fatalf("len = %d want 64", got)
		}
	})
	t.Run("คนละเลข ได้คนละค่า", func(t *testing.T) {
		if nationalid.Hash("p", id) == nationalid.Hash("p", "1234567890139") {
			t.Fatal("ชนกัน")
		}
	})
	t.Run("★ pepper ต่างกัน ได้คนละค่า", func(t *testing.T) {
		if nationalid.Hash("p1", id) == nationalid.Hash("p2", id) {
			t.Fatal("pepper ไม่มีผล = brute-force ได้")
		}
	})
	t.Run("★ ผลลัพธ์ต้องไม่มีเลขบัตรอยู่ในนั้น", func(t *testing.T) {
		if h := nationalid.Hash("p", id); len(h) == 0 || contains(h, id) {
			t.Fatal("เลขบัตรโผล่ใน hash")
		}
	})
}

func TestMask(t *testing.T) {
	cases := []struct{ name, id, want string }{
		{"รูปแบบมาตรฐาน", "1234567890121", "1-2345-xxxxx-xx-1"},
		{"ความยาวผิด คืนค่าว่าง", "123", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := nationalid.Mask(c.id); got != c.want {
				t.Errorf("Mask(%q) = %q want %q", c.id, got, c.want)
			}
		})
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
