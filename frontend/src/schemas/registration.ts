import { z } from 'zod'

const thaiNameChars = /^[\u0E00-\u0E7F .-]+$/
const thaiConsonant = /[\u0E01-\u0E2E]/
const thaiDigit = /[\u0E50-\u0E59]/

function isThaiName(value: string): boolean {
  const s = value.trim()
  return s.length > 0 && thaiNameChars.test(s) && !thaiDigit.test(s) && thaiConsonant.test(s)
}
const mobilePhone = /^0[689]\d{8}$/
const landlinePhone = /^0[2-7]\d{7}$/
const nationalId = /^\d{13}$/
const laser = /^[A-Z]{2}\d-?\d{7}-?\d{2}$/

export const addressSchema = z.object({
  house_no: z.string().optional(),
  moo: z.string().nullable().optional(),
  road: z.string().nullable().optional(),
  province_code: z.string().optional(),
  district_code: z.string().optional(),
  subdistrict_code: z.string().optional(),
  postal_code: z.string().optional(),
})

export const memberSchema = z.object({
  relation: z.string().optional(),
  national_id: z.string().nullable().optional(),
  full_name: z.string().optional(),
  birth_year: z.number().nullable().optional(),
  annual_income: z.number().nullable().optional(),
})

export const incomeSchema = z.object({
  source_type: z.string().optional(),
  annual_amount: z.number().optional(),
})

export const assetSchema = z.object({
  asset_type: z.string().optional(),
  amount: z.union([z.number(), z.string()]).optional(),
  unit: z.string().optional(),
  joint_account_holders: z.number().nullable().optional(),
  is_minor_account: z.boolean().optional(),
})

export const liabilitySchema = z.object({
  liability_type: z.string().optional(),
  credit_limit: z.number().optional(),
})

export const registrationBaseSchema = z.object({
  fiscal_year: z.number().optional(),
  personal: z
    .object({
      national_id: z.string().optional(),
      laser_id: z.string().nullable().optional(),
      id_verify_method: z.string().optional(),
      id_verify_reason: z.string().optional(),
      id_verify_note: z.string().nullable().optional(),
      title: z.string().optional(),
      first_name: z.string().optional(),
      last_name: z.string().optional(),
      birth_year_be: z.number().nullable().optional(),
      birth_month: z.number().nullable().optional(),
      birth_day: z.number().nullable().optional(),
      birth_precision: z.string().optional(),
      phone: z.string().optional(),
      is_farmer: z.boolean().optional(),
      no_laser: z.boolean().optional(),
      address: addressSchema.optional(),
    })
    .optional(),
  family: z
    .object({
      marital_status: z.string().nullable().optional(),
      members: z.array(memberSchema).optional(),
    })
    .nullable()
    .optional(),
  financial: z
    .object({
      income_sources: z.array(incomeSchema).optional(),
      expense_to_others: z.number().optional(),
      assets: z.array(assetSchema).optional(),
      liabilities: z.array(liabilitySchema).optional(),
      has_credit_card: z.boolean().optional(),
    })
    .optional(),
})

export const registrationSubmitSchema = registrationBaseSchema.superRefine((val, ctx) => {
  if (val.fiscal_year == null || val.fiscal_year < 2020) {
    ctx.addIssue({ code: 'custom', path: ['fiscal_year'], message: 'ปีงบประมาณไม่ถูกต้อง' })
  }

  const p = val.personal
  if (!p) {
    ctx.addIssue({ code: 'custom', path: ['personal'], message: 'ข้อมูลส่วนตัวไม่ครบ' })
    return
  }

  const need = (path: string[], label: string) =>
    ctx.addIssue({ code: 'custom', path, message: `กรุณาระบุ${label}` })
  const pick = (path: string[], label: string) =>
    ctx.addIssue({ code: 'custom', path, message: `กรุณาเลือก${label}` })
  const wrong = (path: string[], message: string) =>
    ctx.addIssue({ code: 'custom', path, message })

  if (!p.national_id) need(['personal', 'national_id'], 'เลขประจำตัวประชาชน')
  else if (!nationalId.test(p.national_id))
    wrong(['personal', 'national_id'], 'เลขประจำตัวประชาชนต้องมี 13 หลัก')

  if (!p.title) pick(['personal', 'title'], 'คำนำหน้า')

  if (!p.first_name?.trim()) need(['personal', 'first_name'], 'ชื่อ')
  else if (!isThaiName(p.first_name)) wrong(['personal', 'first_name'], 'ชื่อต้องเป็นภาษาไทย')

  if (!p.last_name?.trim()) need(['personal', 'last_name'], 'นามสกุล')
  else if (!isThaiName(p.last_name)) wrong(['personal', 'last_name'], 'นามสกุลต้องเป็นภาษาไทย')

  if (p.birth_year_be == null) need(['personal', 'birth_year'], 'ปีเกิด')

  if (!p.phone) need(['personal', 'phone'], 'เบอร์โทร')
  else if (!(mobilePhone.test(p.phone) || landlinePhone.test(p.phone)))
    wrong(['personal', 'phone'], 'มือถือ 10 หลัก หรือเบอร์บ้าน 9 หลัก')

  if (p.no_laser) {
    if (!p.id_verify_reason) pick(['personal', 'id_verify_reason'], 'เหตุผลที่ไม่ใช้รหัสหลังบัตร')
    else if (p.id_verify_reason === 'OTHER') {
      const note = (p.id_verify_note ?? '').trim()
      if (!note) need(['personal', 'id_verify_note'], 'เหตุผล')
      else if (note.length < 10)
        wrong(['personal', 'id_verify_note'], 'ต้องระบุเหตุผลอย่างน้อย 10 ตัวอักษร')
    }
  } else if (!p.laser_id) {
    need(['personal', 'laser_id'], 'รหัสหลังบัตร')
  } else if (!laser.test(p.laser_id)) {
    wrong(['personal', 'laser_id'], 'รหัสหลังบัตรไม่ถูกต้อง เช่น JT8-1234567-89')
  }

  const a = p.address
  if (!a?.house_no?.trim()) need(['personal', 'address', 'house_no'], 'บ้านเลขที่')

  if (!a?.subdistrict_code || !a?.district_code || !a?.province_code || !a?.postal_code) {
    pick(['personal', 'address', 'subdistrict_code'], 'ตำบล/แขวง จากช่องค้นหา')
  }
})

export type RegistrationForm = z.infer<typeof registrationBaseSchema>
