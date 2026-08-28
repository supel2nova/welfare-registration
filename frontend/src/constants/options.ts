import type {
  AssetRow,
  IncomeRow,
  LiabilityRow,
  MemberRow,
  RegistrationForm,
  SelectOption,
} from '@/types/form'

export const titleOptions: SelectOption[] = [
  { value: 'นาย', label: 'นาย' },
  { value: 'นาง', label: 'นาง' },
  { value: 'นางสาว', label: 'นางสาว' },
]

export const precisionOptions: SelectOption[] = [
  { value: 'FULL', label: 'วัน เดือน ปี' },
  { value: 'YEAR_MONTH', label: 'เดือน ปี' },
  { value: 'YEAR_ONLY', label: 'ปีอย่างเดียว' },
]

export const maritalOptions: SelectOption[] = [
  { value: 'SINGLE', label: 'โสด' },
  { value: 'MARRIED', label: 'สมรส' },
  { value: 'DIVORCED', label: 'หย่าร้าง' },
  { value: 'WIDOWED', label: 'หม้าย' },
]

export const relationOptions: SelectOption[] = [
  { value: 'SPOUSE', label: 'คู่สมรส' },
  { value: 'CHILD', label: 'บุตร' },
  { value: 'PARENT', label: 'บิดา/มารดา' },
  { value: 'OTHER', label: 'อื่นๆ' },
]

export const incomeTypeOptions: SelectOption[] = [
  { value: 'SALARY', label: 'เงินเดือน' },
  { value: 'AGRI', label: 'เกษตร' },
  { value: 'TRADE', label: 'ค้าขาย' },
  { value: 'RENT', label: 'ค่าเช่า' },
  { value: 'OTHER', label: 'อื่นๆ' },
]

export const liabilityTypeOptions: SelectOption[] = [
  { value: 'LOAN_HOME', label: 'บ้าน' },
  { value: 'LOAN_VEHICLE', label: 'ยานพาหนะ' },
  { value: 'LOAN_PERSONAL', label: 'ส่วนบุคคล' },
  { value: 'LOAN_AGRI', label: 'เกษตร' },
  { value: 'OTHER', label: 'อื่นๆ' },
]

export const assetUnits: Record<string, string> = {
  DEPOSIT: 'THB',
  LOTTERY: 'THB',
  BOND: 'THB',
  SECURITIES: 'THB',
  LAND_AGRI: 'RAI',
  LAND_RESIDENTIAL: 'SQ_WA',
  CONDO: 'SQ_M',
  VEHICLE_CAR: 'COUNT',
  VEHICLE_MOTORCYCLE: 'COUNT',
  VEHICLE_TRICYCLE: 'COUNT',
  VEHICLE_FARM: 'COUNT',
}

export const assetTypeOptions: SelectOption[] = [
  { value: 'DEPOSIT', label: 'เงินฝาก' },
  { value: 'LOTTERY', label: 'สลากออมทรัพย์' },
  { value: 'BOND', label: 'พันธบัตร' },
  { value: 'SECURITIES', label: 'หลักทรัพย์' },
  { value: 'LAND_AGRI', label: 'ที่ดินเกษตร' },
  { value: 'LAND_RESIDENTIAL', label: 'ที่ดินที่อยู่อาศัย' },
  { value: 'CONDO', label: 'ห้องชุด' },
  { value: 'VEHICLE_CAR', label: 'รถยนต์' },
  { value: 'VEHICLE_MOTORCYCLE', label: 'รถจักรยานยนต์' },
  { value: 'VEHICLE_TRICYCLE', label: 'รถสามล้อ' },
  { value: 'VEHICLE_FARM', label: 'รถเพื่อการเกษตร' },
]

export const unitLabels: Record<string, string> = {
  THB: 'บาท',
  RAI: 'ไร่',
  SQ_WA: 'ตารางวา',
  SQ_M: 'ตารางเมตร',
  COUNT: 'คัน',
}

export function emptyMember(): MemberRow {
  return { relation: 'CHILD', national_id: null, full_name: '', birth_year: null, annual_income: null }
}

export function emptyIncome(): IncomeRow {
  return { source_type: 'OTHER', annual_amount: 0 }
}

export function emptyAsset(): AssetRow {
  return {
    asset_type: 'DEPOSIT',
    amount: 0,
    unit: 'THB',
    joint_account_holders: 1,
    is_minor_account: false,
  }
}

export function emptyLiability(): LiabilityRow {
  return { liability_type: 'LOAN_PERSONAL', credit_limit: 0 }
}

export function emptyForm(): RegistrationForm {
  return {
    fiscal_year: 2026,
    personal: {
      national_id: '',
      laser_id: '',
      id_verify_note: '',
      title: 'นาย',
      first_name: '',
      last_name: '',
      birth_year_be: 2528,
      birth_month: 3,
      birth_day: 12,
      birth_precision: 'FULL',
      phone: '',
      is_farmer: true,
      no_laser: false,
      address: {
        house_no: '',
        moo: null,
        road: null,
        province_code: '',
        district_code: '',
        subdistrict_code: '',
        postal_code: '',
      },
    },
    includeFamily: false,
    family: {
      marital_status: 'MARRIED',
      members: [
        { relation: 'SPOUSE', national_id: '', full_name: '', birth_year: null, annual_income: null },
      ],
    },
    financial: {
      income_sources: [{ source_type: 'AGRI', annual_amount: 0 }],
      expense_to_others: 0,
      assets: [emptyAsset()],
      liabilities: [],
      has_credit_card: false,
    },
  }
}
