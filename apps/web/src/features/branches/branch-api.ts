import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'
import type { Branch, BranchListResponse } from './types'

export type CreateBranchInput = {
  name: string
  code?: string
  address?: string
}

export function listBranches(token: string) {
  return apiRequest<BranchListResponse>(endpoints.branches, { token })
}

export function createBranch(token: string, input: CreateBranchInput) {
  return apiRequest<Branch>(endpoints.branches, {
    method: 'POST',
    token,
    body: input,
  })
}
