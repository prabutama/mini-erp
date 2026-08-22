export type Branch = {
  id?: string
  branch_id?: string
  name?: string
  code?: string
  address?: string
  status?: string
  created_at?: string
}

export type BranchListResponse = Branch[] | { branches?: Branch[]; items?: Branch[]; data?: Branch[] }

export function normalizeBranches(payload: BranchListResponse): Branch[] {
  if (Array.isArray(payload)) return payload
  return payload.branches || payload.items || payload.data || []
}
