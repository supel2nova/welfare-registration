import { api } from '@/api/client'
import type { CreateApplicationRequest, CreateApplicationResponse } from '@/types/api'

export async function createApplication(
  body: CreateApplicationRequest,
): Promise<CreateApplicationResponse> {
  const { data } = await api.post<CreateApplicationResponse>('/api/v1/applications', body)
  return data
}
