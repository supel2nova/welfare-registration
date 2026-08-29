import { inject, provide, type InjectionKey, type Ref } from 'vue'

export type FormContext = {
  fieldErrors: Ref<Record<string, string>>
  isSubmitting: Ref<boolean>
  clear: (path: string) => void
  clearPrefix: (prefix: string) => void
  clearAll: () => void
  setError: (path: string, message: string) => void
}

const key: InjectionKey<FormContext> = Symbol('form-context')

export function provideFormContext(ctx: FormContext) {
  provide(key, ctx)
}

export function useFormContext(): FormContext {
  const ctx = inject(key)
  if (!ctx) {
    throw new Error('useFormContext ต้องอยู่ภายใต้ provideFormContext')
  }
  return ctx
}

export function useFieldError(path: () => string) {
  const { fieldErrors } = useFormContext()
  return () => fieldErrors.value[path()] ?? ''
}
