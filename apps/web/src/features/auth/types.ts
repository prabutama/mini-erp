export type RoleName = 'Platform Admin' | 'Business Admin' | 'Manager' | 'Staff'

export type CurrentUser = {
  user_id?: string
  id?: string
  email?: string
  name?: string
  role?: RoleName | string
  business_id?: string
  assigned_branch_ids?: string[]
  permissions?: string[]
}

export type AuthTokens = {
  access_token?: string
  refresh_token?: string
  token?: string
}
