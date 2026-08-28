export type FieldError = {
  field: string
  code: string
  message: string
}

export type Envelope<T> = {
  data: T | null
  statusCode: string
  statusDescription: string
  errorCode?: string
  errorMessage?: string
  fieldErrors?: FieldError[]
}

export type CreateApplicationResponse = {
  application_id: string
  application_no: string
  status: string
  registration_unit: string
  submitted_at: string
}

export type DuplicateInfo = {
  registered_at: string
  registered_unit: string
  application_no: string
  can_appeal: boolean
  appeal_deadline: string | null
}

export type RefItem = {
  code: string
  name_th: string
}

export type DevUser = {
  id: string
  username: string
  role: string
  org_name: string
}

export type ApiError = {
  status: number
  code: string
  message: string
  fieldErrors: FieldError[]
  data: unknown
}

export type SubmitState =
  | { status: 'idle' }
  | { status: 'submitting' }
  | { status: 'success'; application: CreateApplicationResponse }
  | { status: 'invalid'; fieldErrors: FieldError[] }
  | { status: 'duplicate'; info: DuplicateInfo }
  | { status: 'error'; message: string }

export type AddressInput = {
  house_no: string
  moo: string | null
  road: string | null
  province_code: string
  district_code: string
  subdistrict_code: string
  postal_code: string
}

export type CreateApplicationRequest = {
  fiscal_year: number
  personal: {
    national_id: string
    laser_id: string | null
    id_verify_method: string
    id_verify_note: string | null
    title: string
    first_name: string
    last_name: string
    birth_year: number
    birth_month: number | null
    birth_day: number | null
    birth_precision: string
    phone: string
    is_farmer: boolean
    address: AddressInput
  }
  family: {
    marital_status: string | null
    members: Array<{
      relation: string
      national_id: string | null
      full_name: string
      birth_year: number | null
      annual_income: number | null
    }>
  } | null
  financial: {
    income_sources: Array<{ source_type: string; annual_amount: number }>
    expense_to_others: number
    assets: Array<{
      asset_type: string
      amount: number | string
      unit: string
      joint_account_holders: number | null
      is_minor_account: boolean
    }>
    liabilities: Array<{ liability_type: string; credit_limit: number }>
    has_credit_card: boolean
  }
}

export type AddressOption = {
  subdistrict_code: string
  subdistrict_name: string
  subdistrict_kind: string
  district_code: string
  district_name: string
  district_kind: string
  province_code: string
  province_name: string
  postal_code: string
}
