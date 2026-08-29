import axios, { AxiosError, type AxiosInstance } from 'axios'
import type { ApiError, Envelope, FieldError } from '@/types/api'

export class ApiRequestError extends Error {
  readonly status: number
  readonly code: string
  readonly fieldErrors: FieldError[]
  readonly data: unknown

  constructor(err: ApiError) {
    super(err.message)
    this.name = 'ApiRequestError'
    this.status = err.status
    this.code = err.code
    this.fieldErrors = err.fieldErrors
    this.data = err.data
  }
}

type HeaderProvider = () => Record<string, string>

let headerProvider: HeaderProvider = () => ({})

export function setRequestHeaders(provider: HeaderProvider) {
  headerProvider = provider
}

function isEnvelope(v: unknown): v is Envelope<unknown> {
  return typeof v === 'object' && v !== null && 'statusCode' in v && 'statusDescription' in v
}

function toApiError(error: unknown): ApiError {
  if (error instanceof ApiRequestError) {
    return {
      status: error.status,
      code: error.code,
      message: error.message,
      fieldErrors: error.fieldErrors,
      data: error.data,
    }
  }
  if (error instanceof AxiosError) {
    const status = error.response?.status ?? 0
    const body = error.response?.data
    if (isEnvelope(body)) {
      return {
        status,
        code: body.errorCode ?? 'SYS001',
        message: body.errorMessage ?? body.statusDescription ?? 'ระบบขัดข้อง',
        fieldErrors: body.fieldErrors ?? [],
        data: body.data,
      }
    }
    return {
      status,
      code: 'SYS001',
      message: error.message || 'ระบบขัดข้อง',
      fieldErrors: [],
      data: null,
    }
  }
  if (error instanceof Error) {
    return { status: 0, code: 'SYS001', message: error.message, fieldErrors: [], data: null }
  }
  return { status: 0, code: 'SYS001', message: 'ระบบขัดข้อง', fieldErrors: [], data: null }
}

export function createClient(): AxiosInstance {
  const client = axios.create({
    baseURL: '/',
    timeout: 30000,
    headers: { 'Content-Type': 'application/json' },
  })

  client.interceptors.request.use((config) => {
    for (const [k, v] of Object.entries(headerProvider())) {
      if (v) config.headers.set(k, v)
    }
    return config
  })

  client.interceptors.response.use(
    (response) => {
      const body = response.data
      if (isEnvelope(body)) {
        if (body.statusCode === '0') {
          response.data = body.data
          return response
        }
        throw new ApiRequestError({
          status: response.status,
          code: body.errorCode ?? 'SYS001',
          message: body.errorMessage ?? body.statusDescription,
          fieldErrors: body.fieldErrors ?? [],
          data: body.data,
        })
      }
      return response
    },
    (error: unknown) => {
      throw new ApiRequestError(toApiError(error))
    },
  )

  return client
}

export const api = createClient()
