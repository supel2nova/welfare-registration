export function digitsOnly(value: string): string {
  return value.replace(/\D/g, '')
}

export function formatNationalId(value: string): string {
  const d = digitsOnly(value).slice(0, 13)
  if (d.length <= 1) return d
  if (d.length <= 5) return `${d.slice(0, 1)}-${d.slice(1)}`
  if (d.length <= 10) return `${d.slice(0, 1)}-${d.slice(1, 5)}-${d.slice(5)}`
  if (d.length <= 12) return `${d.slice(0, 1)}-${d.slice(1, 5)}-${d.slice(5, 10)}-${d.slice(10)}`
  return `${d.slice(0, 1)}-${d.slice(1, 5)}-${d.slice(5, 10)}-${d.slice(10, 12)}-${d.slice(12)}`
}

export function parseNationalId(value: string): string {
  return digitsOnly(value).slice(0, 13)
}

export function isNationalIdComplete(value: string): boolean {
  return parseNationalId(value).length === 13
}

export function parseLaserId(value: string): string {
  const cleaned = value.toUpperCase().replace(/[^A-Z0-9]/g, '')
  let letters = ''
  let digits = ''
  for (const c of cleaned) {
    if (letters.length < 2) {
      if (/[A-Z]/.test(c)) letters += c
      continue
    }
    if (digits.length < 10 && /\d/.test(c)) digits += c
  }
  return letters + digits
}

export function formatLaserId(value: string): string {
  const s = parseLaserId(value)
  if (s.length <= 3) return s
  if (s.length <= 10) return `${s.slice(0, 3)}-${s.slice(3)}`
  return `${s.slice(0, 3)}-${s.slice(3, 10)}-${s.slice(10)}`
}

export function isLaserIdComplete(value: string): boolean {
  return parseLaserId(value).length === 12
}

const mobilePattern = /^0[689]\d{8}$/
const landlinePattern = /^0[2-7]\d{7}$/

export function isMobilePhone(digits: string): boolean {
  return /^0[689]/.test(digits)
}

export function isLandlinePhone(digits: string): boolean {
  return /^0[2-7]/.test(digits)
}

export function parsePhone(value: string): string {
  const d = digitsOnly(value)
  return isLandlinePhone(d) ? d.slice(0, 9) : d.slice(0, 10)
}

export function formatPhone(value: string): string {
  const d = parsePhone(value)

  if (d.startsWith('02')) {
    if (d.length <= 2) return d
    if (d.length <= 5) return `${d.slice(0, 2)}-${d.slice(2)}`
    return `${d.slice(0, 2)}-${d.slice(2, 5)}-${d.slice(5, 9)}`
  }

  if (isLandlinePhone(d)) {
    if (d.length <= 3) return d
    if (d.length <= 6) return `${d.slice(0, 3)}-${d.slice(3)}`
    return `${d.slice(0, 3)}-${d.slice(3, 6)}-${d.slice(6, 9)}`
  }

  if (d.length <= 3) return d
  if (d.length <= 6) return `${d.slice(0, 3)}-${d.slice(3)}`
  return `${d.slice(0, 3)}-${d.slice(3, 6)}-${d.slice(6, 10)}`
}

export function isPhoneValid(value: string): boolean {
  const d = parsePhone(value)
  return mobilePattern.test(d) || landlinePattern.test(d)
}
