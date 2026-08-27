package httpx_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/supel2nova/welfare-registration/backend/pkg/apperror"
	"github.com/supel2nova/welfare-registration/backend/pkg/httpx"
)

func marshal(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestOK(t *testing.T) {
	got := marshal(t, httpx.OK(map[string]string{"application_no": "WC-2026-0000001"}))
	want := `{"data":{"application_no":"WC-2026-0000001"},"statusCode":"0","statusDescription":"Success"}`
	if got != want {
		t.Errorf("got  %s\nwant %s", got, want)
	}
}

func TestFail(t *testing.T) {
	t.Run("400 พร้อม fieldErrors และ data เป็น null", func(t *testing.T) {
		status, env := httpx.Fail(apperror.Validation([]apperror.FieldError{
			apperror.Field("personal.national_id", apperror.CodeNationalID),
		}))
		if status != 400 {
			t.Errorf("status = %d want 400", status)
		}
		got := marshal(t, env)
		if !strings.Contains(got, `"data":null`) {
			t.Errorf(`ไม่มี "data":null: %s`, got)
		}
		if !strings.Contains(got, `"field":"personal.national_id"`) || !strings.Contains(got, `"code":"VAL001"`) {
			t.Errorf("fieldErrors ไม่ครบ: %s", got)
		}
	})

	t.Run("409 คืนข้อมูลใบเดิมใน data ไม่ใช่ null", func(t *testing.T) {
		status, env := httpx.Fail(apperror.Duplicate(map[string]any{"application_no": "WC-2026-0000042"}))
		if status != 409 {
			t.Errorf("status = %d want 409", status)
		}
		got := marshal(t, env)
		if strings.Contains(got, `"data":null`) {
			t.Errorf("409 ต้องคืนข้อมูลใบเดิม: %s", got)
		}
		if !strings.Contains(got, `"errorCode":"DUP001"`) {
			t.Errorf("errorCode ผิด: %s", got)
		}
	})

	t.Run("500 ไม่หลุดรายละเอียดภายในออกไป", func(t *testing.T) {
		status, env := httpx.Fail(apperror.Internal(errSecret{}))
		if status != 500 {
			t.Errorf("status = %d want 500", status)
		}
		if got := marshal(t, env); strings.Contains(got, "pgx: connection refused to 10.0.0.5") {
			t.Errorf("รายละเอียดภายในหลุด: %s", got)
		}
	})
}

type errSecret struct{}

func (errSecret) Error() string { return "pgx: connection refused to 10.0.0.5" }
