import { api } from '@/api/client'
import type { AddressOption, DevUser } from '@/types/api'

export async function fetchDevUsers(): Promise<DevUser[]> {
  const { data } = await api.get<DevUser[]>('/api/v1/dev/users')
  return data
}

export async function searchAddress(q: string, signal?: AbortSignal): Promise<AddressOption[]> {
  const { data } = await api.get<AddressOption[]>('/api/v1/ref/address-search', {
    params: { q },
    signal,
    headers: { 'Cache-Control': 'no-cache' },
  })
  return data
}
