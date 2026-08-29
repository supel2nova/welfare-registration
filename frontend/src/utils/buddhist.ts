export function currentFiscalYear(now = new Date()): number {
  return now.getMonth() >= 9 ? now.getFullYear() + 1 : now.getFullYear()
}

export function toCE(buddhistYear: number): number {
  return buddhistYear - 543
}

export function toBE(christianYear: number): number {
  return christianYear + 543
}

export function formatBEDate(ceYear: number, month?: number | null, day?: number | null): string {
  const y = String(toBE(ceYear))
  if (month == null) return y
  const m = String(month).padStart(2, '0')
  if (day == null) return `${y}-${m}`
  return `${y}-${m}-${String(day).padStart(2, '0')}`
}

const thaiMonths = [
  'มกราคม', 'กุมภาพันธ์', 'มีนาคม', 'เมษายน', 'พฤษภาคม', 'มิถุนายน',
  'กรกฎาคม', 'สิงหาคม', 'กันยายน', 'ตุลาคม', 'พฤศจิกายน', 'ธันวาคม',
]

export function formatThaiDateTime(iso: string): string {
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  const time = `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
  return `${d.getDate()} ${thaiMonths[d.getMonth()]} ${toBE(d.getFullYear())} เวลา ${time} น.`
}
