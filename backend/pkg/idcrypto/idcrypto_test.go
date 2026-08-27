package idcrypto_test

import (
	"strings"
	"testing"

	"github.com/supel2nova/welfare-registration/backend/pkg/idcrypto"
)

const (
	key      = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	otherKey = "fedcba9876543210fedcba9876543210fedcba9876543210fedcba9876543210"
	id       = "1234567890121"
	hashA    = "aaaa1111"
	hashB    = "bbbb2222"
)

func newCipher(t *testing.T, k string) idcrypto.Cipher {
	t.Helper()
	c, err := idcrypto.New(k)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestNewRejectsBadKey(t *testing.T) {
	cases := []struct{ name, key string }{
		{"ไม่ใช่ hex", "zzzz"},
		{"สั้นเกินไป", "0123456789abcdef"},
		{"ว่าง", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if _, err := idcrypto.New(c.key); err == nil {
				t.Error("ต้อง error")
			}
		})
	}
}

func TestSealOpen(t *testing.T) {
	c := newCipher(t, key)

	sealed, err := c.Seal(id, hashA)
	if err != nil {
		t.Fatal(err)
	}

	t.Run("ถอดกลับได้เลขเดิม", func(t *testing.T) {
		got, err := c.Open(sealed, hashA)
		if err != nil || got != id {
			t.Fatalf("Open = %q, %v", got, err)
		}
	})

	t.Run("ciphertext ต้องไม่มีเลขบัตรอยู่ในนั้น", func(t *testing.T) {
		if strings.Contains(string(sealed), id) {
			t.Fatal("เลขบัตรโผล่ใน ciphertext")
		}
	})

	t.Run("เลขเดิมเข้ารหัสสองครั้งได้คนละค่า", func(t *testing.T) {
		again, _ := c.Seal(id, hashA)
		if string(again) == string(sealed) {
			t.Fatal("nonce ซ้ำ = เดาได้ว่าสองแถวนี้คือคนเดียวกัน")
		}
	})

	t.Run("ย้าย ciphertext ไปแถวอื่นแล้วถอดไม่ออก", func(t *testing.T) {
		if _, err := c.Open(sealed, hashB); err == nil {
			t.Fatal("ก๊อป ciphertext ข้ามแถวได้")
		}
	})

	t.Run("key คนละดอกถอดไม่ออก", func(t *testing.T) {
		if _, err := newCipher(t, otherKey).Open(sealed, hashA); err == nil {
			t.Fatal("ถอดได้ทั้งที่ key ไม่ตรง")
		}
	})

	t.Run("แก้ไบต์เดียวก็ถอดไม่ออก", func(t *testing.T) {
		bad := append([]byte(nil), sealed...)
		bad[len(bad)-1] ^= 0xff
		if _, err := c.Open(bad, hashA); err == nil {
			t.Fatal("ciphertext ถูกแก้แล้วยังผ่าน")
		}
	})

	t.Run("ข้อมูลสั้นเกินไปไม่ panic", func(t *testing.T) {
		if _, err := c.Open([]byte{0x01, 0x02}, hashA); err == nil {
			t.Fatal("ต้อง error")
		}
	})
}
