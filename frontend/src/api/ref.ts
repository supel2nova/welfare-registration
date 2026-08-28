import { api } from '@/api/client'
import type { DevUser, RefItem } from '@/types/api'

export async function fetchProvinces(): Promise<RefItem[]> {
  const { data } = await api.get<RefItem[]>('/api/v1/ref/provinces')
  return data
}

export async function fetchDistricts(provinceCode: string): Promise<RefItem[]> {
  const { data } = await api.get<RefItem[]>('/api/v1/ref/districts', {
    params: { province_code: provinceCode },
  })
  return data
}

export async function fetchSubdistricts(districtCode: string): Promise<RefItem[]> {
  const { data } = await api.get<RefItem[]>('/api/v1/ref/subdistricts', {
    params: { district_code: districtCode },
  })
  return data
}

export async function fetchDevUsers(): Promise<DevUser[]> {
  const { data } = await api.get<DevUser[]>('/api/v1/dev/users')
  return data
}
