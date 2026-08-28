package dto

import "time"

type CreateApplicationResponse struct {
	ApplicationID    string    `json:"application_id"`
	ApplicationNo    string    `json:"application_no"`
	Status           string    `json:"status"`
	RegistrationUnit string    `json:"registration_unit"`
	SubmittedAt      time.Time `json:"submitted_at"`
}

type DuplicateInfo struct {
	RegisteredAt   string  `json:"registered_at"`
	RegisteredUnit string  `json:"registered_unit"`
	ApplicationNo  string  `json:"application_no"`
	CanAppeal      bool    `json:"can_appeal"`
	AppealDeadline *string `json:"appeal_deadline"`
}

type RefItem struct {
	Code   string `json:"code"`
	NameTH string `json:"name_th"`
}

type AddressOption struct {
	SubdistrictCode string `json:"subdistrict_code"`
	SubdistrictName string `json:"subdistrict_name"`
	SubdistrictKind string `json:"subdistrict_kind"`
	DistrictCode    string `json:"district_code"`
	DistrictName    string `json:"district_name"`
	DistrictKind    string `json:"district_kind"`
	ProvinceCode    string `json:"province_code"`
	ProvinceName    string `json:"province_name"`
	PostalCode      string `json:"postal_code"`
}

type DevUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	OrgName  string `json:"org_name"`
}
