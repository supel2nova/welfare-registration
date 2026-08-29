import type { AddressInput } from '@/types/api'

export type MemberRow = {
  relation: string
  national_id: string | null
  full_name: string
  birth_year: number | null
  annual_income: number | null
}

export type IncomeRow = {
  source_type: string
  annual_amount: number
}

export type AssetRow = {
  asset_type: string
  amount: number | string
  unit: string
  joint_account_holders: number | null
  is_minor_account: boolean
}

export type LiabilityRow = {
  liability_type: string
  credit_limit: number
}

export type PersonalForm = {
  national_id: string
  laser_id: string | null
  id_verify_reason: string
  id_verify_note: string | null
  title: string
  first_name: string
  last_name: string
  birth_year_be: number | null
  birth_month: number | null
  birth_day: number | null
  birth_precision: string
  phone: string
  is_farmer: boolean
  no_laser: boolean
  address: AddressInput
}

export type FamilyForm = {
  marital_status: string | null
  members: MemberRow[]
}

export type FinancialForm = {
  income_sources: IncomeRow[]
  expense_to_others: number
  assets: AssetRow[]
  liabilities: LiabilityRow[]
  has_credit_card: boolean
}

export type RegistrationForm = {
  fiscal_year: number
  personal: PersonalForm
  includeFamily: boolean
  family: FamilyForm
  financial: FinancialForm
}

export type SelectOption = {
  value: string
  label: string
}
