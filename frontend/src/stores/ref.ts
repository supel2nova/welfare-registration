import { defineStore } from 'pinia'
import { ref } from 'vue'
import { fetchDistricts, fetchProvinces, fetchSubdistricts } from '@/api/ref'
import type { RefItem } from '@/types/api'

export const useRefStore = defineStore('ref', () => {
  const provinces = ref<RefItem[]>([])
  const districts = ref<RefItem[]>([])
  const subdistricts = ref<RefItem[]>([])

  async function loadProvinces() {
    if (provinces.value.length === 0) {
      provinces.value = await fetchProvinces()
    }
  }

  async function loadDistricts(provinceCode: string) {
    districts.value = provinceCode ? await fetchDistricts(provinceCode) : []
    subdistricts.value = []
  }

  async function loadSubdistricts(districtCode: string) {
    subdistricts.value = districtCode ? await fetchSubdistricts(districtCode) : []
  }

  return { provinces, districts, subdistricts, loadProvinces, loadDistricts, loadSubdistricts }
})
