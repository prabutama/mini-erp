import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'

export type CurrentBusiness = {
  id?: string
  business_id?: string
  name?: string
  status?: string
}

export type ServiceOrdersSummary = {
  total?: number
  open?: number
  in_progress?: number
  completed?: number
  cancelled?: number
  counts?: Record<string, number>
}

export function getCurrentBusiness(token: string) {
  return apiRequest<CurrentBusiness>(endpoints.currentBusiness, { token })
}

export function getServiceOrdersSummary(token: string) {
  return apiRequest<ServiceOrdersSummary>(endpoints.serviceOrdersSummary, { token })
}
