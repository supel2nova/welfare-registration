package laserid_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/supel2nova/welfare-registration/backend/internal/domain/laserid"
)

func TestIsValidFormat(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want bool
	}{
		{"มีขีดครบตามบัตร", "JT8-1234567-89", true},
		{"ไม่มีขีดเลย", "JT8123456789", true},
		{"ขีดตัวแรกตัวเดียว", "JT8-123456789", true},
		{"ขีดตัวหลังตัวเดียว", "JT81234567-89", true},
		{"สั้นเกินไป", "ABC123", false},
		{"ตัวอักษรพิมพ์เล็ก", "jt8-1234567-89", false},
		{"ตัวอักษร 3 ตัวนำ", "JTX-1234567-89", false},
		{"ตัวเลขเกิน", "JT8-12345678-89", false},
		{"ว่าง", "", false},
		{"มีช่องว่างท้าย", "JT8-1234567-89 ", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := laserid.IsValidFormat(c.in); got != c.want {
				t.Errorf("IsValidFormat(%q) = %v want %v", c.in, got, c.want)
			}
		})
	}
}

func TestNoStoreFunction(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	banned := []string{"func Hash", "func Store", "func Save", "func Encrypt"}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		src, err := os.ReadFile(f)
		if err != nil {
			t.Fatal(err)
		}
		for _, b := range banned {
			if strings.Contains(string(src), b) {
				t.Errorf("%s มี %q — laser_id ห้ามถูกเก็บ (INVARIANT 5)", f, b)
			}
		}
	}
}
