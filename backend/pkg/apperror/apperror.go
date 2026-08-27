package apperror

import "fmt"

const (
	CodeInvalidPayload = "VAL000"
	CodeNationalID     = "VAL001"
	CodeLaserID        = "VAL002"
	CodeName           = "VAL003"
	CodeBirthDate      = "VAL004"
	CodePhone          = "VAL005"
	CodeAddress        = "VAL006"
	CodeAddressMissing = "VAL007"
	CodeMemberDupSelf  = "VAL008"
	CodeSpouseCount    = "VAL009"
	CodeTooMany        = "VAL010"
	CodeNegative       = "VAL011"
	CodeEnum           = "VAL012"
	CodeAssetUnit      = "VAL013"
	CodeFiscalYear     = "VAL014"
	CodeIDVerification = "VAL015"
	CodeKYC            = "VAL_KYC"
	CodeDuplicate      = "DUP001"
	CodeForbidden      = "PERM001"
	CodeInternal       = "SYS001"
)

var messages = map[string]string{
	CodeInvalidPayload: "ข้อมูลไม่ถูกต้อง",
	CodeNationalID:     "เลขประจำตัวประชาชนไม่ถูกต้อง",
	CodeLaserID:        "รหัสหลังบัตรไม่ถูกต้อง",
	CodeName:           "ชื่อ-นามสกุลไม่ถูกต้อง",
	CodeBirthDate:      "วันเดือนปีเกิดไม่ถูกต้อง",
	CodePhone:          "หมายเลขโทรศัพท์ไม่ถูกต้อง",
	CodeAddress:        "ที่อยู่ไม่ถูกต้องหรือไม่สอดคล้องกัน",
	CodeAddressMissing: "ข้อมูลที่อยู่ไม่ครบ",
	CodeMemberDupSelf:  "เลขประจำตัวประชาชนของสมาชิกซ้ำกับผู้ยื่น",
	CodeSpouseCount:    "ระบุคู่สมรสได้ 1 คน",
	CodeTooMany:        "จำนวนรายการเกินกำหนด",
	CodeNegative:       "ตัวเลขไม่ถูกต้อง",
	CodeEnum:           "ค่าที่เลือกไม่ถูกต้อง",
	CodeAssetUnit:      "หน่วยไม่สอดคล้องกับประเภททรัพย์สิน",
	CodeFiscalYear:     "ปีงบประมาณไม่ถูกต้อง",
	CodeIDVerification: "ข้อมูลยืนยันตัวตนไม่สอดคล้องกัน",
	CodeKYC:            "ยืนยันตัวตนไม่ผ่าน",
	CodeDuplicate:      "เลขประจำตัวประชาชนนี้ลงทะเบียนแล้ว",
	CodeForbidden:      "ไม่มีสิทธิ์ดำเนินการ",
	CodeInternal:       "ระบบขัดข้อง กรุณาลองใหม่อีกครั้ง",
}

func Message(code string) string { return messages[code] }

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Field(path, code string) FieldError {
	return FieldError{Field: path, Code: code, Message: messages[code]}
}

func FieldMsg(path, code, message string) FieldError {
	return FieldError{Field: path, Code: code, Message: message}
}

type Error struct {
	HTTPStatus  int
	Code        string
	Message     string
	FieldErrors []FieldError
	Data        any
	cause       error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.cause }

func Validation(fieldErrors []FieldError) *Error {
	return &Error{
		HTTPStatus:  400,
		Code:        CodeInvalidPayload,
		Message:     messages[CodeInvalidPayload],
		FieldErrors: fieldErrors,
	}
}

func BadRequest(code string, fieldErrors ...FieldError) *Error {
	return &Error{
		HTTPStatus:  400,
		Code:        code,
		Message:     messages[code],
		FieldErrors: fieldErrors,
	}
}

func Duplicate(data any) *Error {
	return &Error{
		HTTPStatus: 409,
		Code:       CodeDuplicate,
		Message:    messages[CodeDuplicate],
		Data:       data,
	}
}

func Forbidden() *Error {
	return &Error{HTTPStatus: 403, Code: CodeForbidden, Message: messages[CodeForbidden]}
}

func Internal(cause error) *Error {
	return &Error{
		HTTPStatus: 500,
		Code:       CodeInternal,
		Message:    messages[CodeInternal],
		cause:      cause,
	}
}
