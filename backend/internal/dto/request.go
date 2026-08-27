package dto

import "encoding/json"

type CreateApplicationRequest struct {
	FiscalYear int       `json:"fiscal_year"`
	Personal   Personal  `json:"personal"`
	Family     *Family   `json:"family"`
	Financial  Financial `json:"financial"`
}

type Personal struct {
	NationalID     string  `json:"national_id"`
	LaserID        *string `json:"laser_id"`
	IDVerifyMethod string  `json:"id_verify_method"`
	IDVerifyNote   *string `json:"id_verify_note"`
	Title          string  `json:"title"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	BirthYear      int     `json:"birth_year"`
	BirthMonth     *int    `json:"birth_month"`
	BirthDay       *int    `json:"birth_day"`
	BirthPrecision string  `json:"birth_precision"`
	Phone          string  `json:"phone"`
	IsFarmer       bool    `json:"is_farmer"`
	Address        Address `json:"address"`
}

type Address struct {
	HouseNo         string  `json:"house_no"`
	Moo             *string `json:"moo"`
	Road            *string `json:"road"`
	ProvinceCode    string  `json:"province_code"`
	DistrictCode    string  `json:"district_code"`
	SubdistrictCode string  `json:"subdistrict_code"`
	PostalCode      string  `json:"postal_code"`
}

type Family struct {
	MaritalStatus *string  `json:"marital_status"`
	Members       []Member `json:"members"`
}

type Member struct {
	Relation     string  `json:"relation"`
	NationalID   *string `json:"national_id"`
	FullName     string  `json:"full_name"`
	BirthYear    *int    `json:"birth_year"`
	AnnualIncome *int64  `json:"annual_income"`
}

type Financial struct {
	IncomeSources   []IncomeSource `json:"income_sources"`
	ExpenseToOthers int64          `json:"expense_to_others"`
	Assets          []Asset        `json:"assets"`
	Liabilities     []Liability    `json:"liabilities"`
	HasCreditCard   bool           `json:"has_credit_card"`
}

type IncomeSource struct {
	SourceType   string `json:"source_type"`
	AnnualAmount int64  `json:"annual_amount"`
}

type Asset struct {
	AssetType           string      `json:"asset_type"`
	Amount              json.Number `json:"amount"`
	Unit                string      `json:"unit"`
	JointAccountHolders *int        `json:"joint_account_holders"`
	IsMinorAccount      bool        `json:"is_minor_account"`
}

type Liability struct {
	LiabilityType string `json:"liability_type"`
	CreditLimit   int64  `json:"credit_limit"`
}
