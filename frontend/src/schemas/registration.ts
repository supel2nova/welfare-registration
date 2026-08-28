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
      id_verify_note: z.string().nullable().optional(),
      title: z.string().optional(),
      first_name: z.string().optional(),
      last_name: z.string().optional(),
      birth_year_be: z.number().optional(),
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

  if (!p.national_id || !nationalId.test(p.national_id)) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'national_id'],
      message: 'เลขประจำตัวประชาชนไม่ถูกต้อง',
    })
  }
  if (!p.title) {
    ctx.addIssue({ code: 'custom', path: ['personal', 'title'], message: 'เลือกคำนำหน้า' })
  }
  if (!p.first_name || !isThaiName(p.first_name)) {
    ctx.addIssue({ code: 'custom', path: ['personal', 'first_name'], message: 'ชื่อไม่ถูกต้อง' })
  }
  if (!p.last_name || !isThaiName(p.last_name)) {
    ctx.addIssue({ code: 'custom', path: ['personal', 'last_name'], message: 'นามสกุลไม่ถูกต้อง' })
  }
  if (p.birth_year_be == null) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'birth_year_be'],
      message: 'ปีเกิดไม่ถูกต้อง',
    })
  }
  if (!p.phone || !(mobilePhone.test(p.phone) || landlinePhone.test(p.phone))) {
    ctx.addIssue({ code: 'custom', path: ['personal', 'phone'], message: 'หมายเลขโทรศัพท์ไม่ถูกต้อง' })
  }

  if (p.no_laser) {
    const note = (p.id_verify_note ?? '').trim()
    if (note.length < 10) {
      ctx.addIssue({
        code: 'custom',
        path: ['personal', 'id_verify_note'],
        message: 'ต้องระบุเหตุผลอย่างน้อย 10 ตัวอักษร',
      })
    }
  } else if (!p.laser_id || !laser.test(p.laser_id)) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'laser_id'],
      message: 'รหัสหลังบัตรไม่ถูกต้อง',
    })
  }

  const a = p.address
  if (!a?.house_no?.trim()) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'address', 'house_no'],
      message: 'กรอกบ้านเลขที่',
    })
  }
  if (!a?.province_code) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'address', 'province_code'],
      message: 'เลือกจังหวัด',
    })
  }
  if (!a?.district_code) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'address', 'district_code'],
      message: 'เลือกอำเภอ',
    })
  }
  if (!a?.subdistrict_code) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'address', 'subdistrict_code'],
      message: 'เลือกตำบล',
    })
  }
  if (!a?.postal_code || !/^\d{5}$/.test(a.postal_code)) {
    ctx.addIssue({
      code: 'custom',
      path: ['personal', 'address', 'postal_code'],
      message: 'รหัสไปรษณีย์ไม่ถูกต้อง',
    })
  }
})

export type RegistrationForm = z.infer<typeof registrationBaseSchema>
