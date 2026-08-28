import { computed, nextTick, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { createApplication } from '@/api/application'
import { ApiRequestError } from '@/api/client'
import { provideFormContext } from '@/composables/useFormContext'
import { emptyForm } from '@/constants/options'
import { registrationSubmitSchema } from '@/schemas/registration'
import type { CreateApplicationRequest, DuplicateInfo, SubmitState } from '@/types/api'
import type { RegistrationForm } from '@/types/form'
import { toCE } from '@/utils/buddhist'

export function useRegistrationForm() {
  const router = useRouter()

  const form = reactive<RegistrationForm>(emptyForm())
  const fieldErrors = ref<Record<string, string>>({})
  const isSubmitting = ref(false)
  const submitState = ref<SubmitState>({ status: 'idle' })
  const duplicateOpen = ref(false)
  const duplicateInfo = ref<DuplicateInfo | null>(null)
  const formError = ref('')

  provideFormContext({
    fieldErrors,
    isSubmitting,
    clear: (path) => {
      if (fieldErrors.value[path]) {
        const next = { ...fieldErrors.value }
        delete next[path]
        fieldErrors.value = next
      }
    },
    clearAll: () => {
      fieldErrors.value = {}
    },
    setError: (path, message) => {
      fieldErrors.value = { ...fieldErrors.value, [path]: message }
    },
  })

  const isDirty = computed(() => JSON.stringify(form) !== JSON.stringify(emptyForm()))

  function reset() {
    Object.assign(form, emptyForm())
    fieldErrors.value = {}
    formError.value = ''
    submitState.value = { status: 'idle' }
    duplicateOpen.value = false
    duplicateInfo.value = null
  }

  function setErrorsFromPaths(issues: Array<{ path: PropertyKey[]; message: string }>) {
    const map: Record<string, string> = {}
    for (const issue of issues) {
      const path = issue.path.map(String).join('.')
      if (path && !map[path]) map[path] = issue.message
    }
    fieldErrors.value = map
  }

  function buildPayload(): CreateApplicationRequest {
    const p = form.personal
    const precision = p.birth_precision
    let birth_month = p.birth_month
    let birth_day = p.birth_day
    if (precision === 'YEAR_ONLY') {
      birth_month = null
      birth_day = null
    } else if (precision === 'YEAR_MONTH') {
      birth_day = null
    }

    return {
      fiscal_year: form.fiscal_year,
      personal: {
        national_id: p.national_id,
        laser_id: p.no_laser ? null : p.laser_id || null,
        id_verify_method: p.no_laser ? 'MANUAL_CARD_CHECK' : 'LASER_CODE',
        id_verify_note: p.no_laser ? p.id_verify_note : null,
        title: p.title,
        first_name: p.first_name,
        last_name: p.last_name,
        birth_year: p.birth_year_be == null ? 0 : toCE(p.birth_year_be),
        birth_month,
        birth_day,
        birth_precision: precision,
        phone: p.phone,
        is_farmer: p.is_farmer,
        address: { ...p.address },
      },
      family: form.includeFamily
        ? {
            marital_status: form.family.marital_status,
            members: form.family.members.map((m) => ({
              relation: m.relation,
              national_id: m.national_id || null,
              full_name: m.full_name,
              birth_year: m.birth_year,
              annual_income: m.annual_income,
            })),
          }
        : null,
      financial: {
        income_sources: form.financial.income_sources.map((s) => ({
          source_type: s.source_type,
          annual_amount: Number(s.annual_amount) || 0,
        })),
        expense_to_others: Number(form.financial.expense_to_others) || 0,
        assets: form.financial.assets.map((a) => ({
          asset_type: a.asset_type,
          amount: a.amount,
          unit: a.unit,
          joint_account_holders: a.joint_account_holders,
          is_minor_account: a.is_minor_account,
        })),
        liabilities: form.financial.liabilities.map((l) => ({
          liability_type: l.liability_type,
          credit_limit: Number(l.credit_limit) || 0,
        })),
        has_credit_card: form.financial.has_credit_card,
      },
    }
  }

  function scrollToCenter(el: HTMLElement, duration = 450) {
    const rect = el.getBoundingClientRect()
    const maxTop = document.documentElement.scrollHeight - window.innerHeight
    const target = Math.max(
      0,
      Math.min(maxTop, window.scrollY + rect.top - (window.innerHeight - rect.height) / 2),
    )
    const start = window.scrollY
    const delta = target - start
    if (Math.abs(delta) < 2) return

    if (window.matchMedia('(prefers-reduced-motion: reduce)').matches) {
      window.scrollTo(0, target)
      return
    }

    const startedAt = performance.now()
    const step = (now: number) => {
      const p = Math.min((now - startedAt) / duration, 1)
      const eased = p < 0.5 ? 2 * p * p : 1 - Math.pow(-2 * p + 2, 2) / 2
      window.scrollTo(0, start + delta * eased)
      if (p < 1) requestAnimationFrame(step)
    }
    requestAnimationFrame(step)
  }

  async function focusFirstError() {
    await nextTick()
    for (let frame = 0; frame < 10; frame++) {
      const el = document.querySelector<HTMLElement>('.field-input-error')
      if (el) {
        el.focus({ preventScroll: true })
        scrollToCenter(el)
        return
      }
      await new Promise(requestAnimationFrame)
    }
  }

  function applyFieldErrors(entries: Array<{ field: string; message: string }>) {
    const map: Record<string, string> = {}
    for (const fe of entries) map[fe.field] = fe.message
    fieldErrors.value = map
  }

  async function submit() {
    formError.value = ''
    fieldErrors.value = {}
    submitState.value = { status: 'submitting' }
    isSubmitting.value = true

    const parsed = registrationSubmitSchema.safeParse(form)
    if (!parsed.success) {
      setErrorsFromPaths(parsed.error.issues)
      submitState.value = {
        status: 'invalid',
        fieldErrors: Object.entries(fieldErrors.value).map(([field, message]) => ({
          field,
          code: 'VAL000',
          message,
        })),
      }
      isSubmitting.value = false
      void focusFirstError()
      return
    }

    try {
      const res = await createApplication(buildPayload())
      submitState.value = { status: 'success', application: res }
      await router.push({
        name: 'success',
        query: {
          application_no: res.application_no,
          registration_unit: res.registration_unit,
          submitted_at: res.submitted_at,
          status: res.status,
        },
      })
    } catch (err: unknown) {
      if (err instanceof ApiRequestError) {
        if (err.status === 409 && err.data && typeof err.data === 'object') {
          const info = err.data as DuplicateInfo
          duplicateInfo.value = info
          duplicateOpen.value = true
          submitState.value = { status: 'duplicate', info }
        } else if (err.fieldErrors.length > 0) {
          applyFieldErrors(err.fieldErrors)
          submitState.value = { status: 'invalid', fieldErrors: err.fieldErrors }
          void focusFirstError()
        } else {
          formError.value = err.message
          submitState.value = { status: 'error', message: err.message }
        }
      } else {
        const message = err instanceof Error ? err.message : 'ระบบขัดข้อง'
        formError.value = message
        submitState.value = { status: 'error', message }
      }
    } finally {
      isSubmitting.value = false
    }
  }

  return {
    form,
    fieldErrors,
    isSubmitting,
    submitState,
    duplicateOpen,
    duplicateInfo,
    formError,
    isDirty,
    reset,
    submit,
  }
}
