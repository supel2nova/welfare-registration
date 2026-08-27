package dto_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/supel2nova/welfare-registration/backend/internal/dto"
)

func TestRequestHasNoActorFields(t *testing.T) {
	banned := map[string]bool{
		"registration_unit_id": true,
		"created_by_user_id":   true,
		"created_by_client_id": true,
		"submission_channel":   true,
		"status":               true,
		"application_no":       true,
		"org_id":               true,
		"organization_id":      true,
	}
	for _, tag := range jsonTags(reflect.TypeOf(dto.CreateApplicationRequest{}), map[reflect.Type]bool{}) {
		if banned[tag] {
			t.Errorf("DTO มี %q — org/user/client ต้องมาจาก Actor เท่านั้น (INVARIANT 4)", tag)
		}
	}
}

func TestRequestNeverUsesGinBinding(t *testing.T) {
	for _, tag := range structTags(reflect.TypeOf(dto.CreateApplicationRequest{}), map[reflect.Type]bool{}) {
		if strings.Contains(tag, `binding:"`) {
			t.Errorf("เจอ %q — binding หยุดที่ error ตัวแรก ทำให้เก็บ fieldErrors ไม่ครบ (INVARIANT 2)", tag)
		}
	}
}

func jsonTags(t reflect.Type, seen map[reflect.Type]bool) []string {
	var out []string
	for _, tag := range structTags(t, seen) {
		if v, ok := reflect.StructTag(tag).Lookup("json"); ok {
			out = append(out, strings.Split(v, ",")[0])
		}
	}
	return out
}

func structTags(t reflect.Type, seen map[reflect.Type]bool) []string {
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct || seen[t] {
		return nil
	}
	seen[t] = true
	var out []string
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		out = append(out, string(f.Tag))
		out = append(out, structTags(f.Type, seen)...)
	}
	return out
}
