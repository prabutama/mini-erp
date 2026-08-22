import { createContext, useCallback, useContext, useMemo, useState, type ReactNode } from 'react'
import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'
import type { AuthTokens, CurrentUser } from './types'

type LoginInput = {
  email: string
  password: string
}

type SignupInput = {
  business_name: string
  admin_name: string
  email: string
  password: string
}

type AuthContextValue = {
  token: string | null
  user: CurrentUser | null
  isAuthenticated: boolean
  setToken: (token: string | null) => void
  loadCurrentUser: () => Promise<CurrentUser | null>
  login: (input: LoginInput) => Promise<CurrentUser | null>
  signup: (input: SignupInput) => Promise<CurrentUser | null>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | null>(null)

function extractAccessToken(payload: AuthTokens) {
  return payload.access_token || payload.token || null
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<CurrentUser | null>(null)

  const loadCurrentUser = useCallback(async () => {
    if (!token) {
      setUser(null)
      return null
    }

    const currentUser = await apiRequest<CurrentUser>(endpoints.me, { token })
    setUser(currentUser)
    return currentUser
  }, [token])

  const login = useCallback(async (input: LoginInput) => {
    const payload = await apiRequest<AuthTokens>(endpoints.login, {
      method: 'POST',
      body: input,
    })
    const accessToken = extractAccessToken(payload)
    setToken(accessToken)

    if (!accessToken) return null

    const currentUser = await apiRequest<CurrentUser>(endpoints.me, { token: accessToken })
    setUser(currentUser)
    return currentUser
  }, [])

  const signup = useCallback(async (input: SignupInput) => {
    const payload = await apiRequest<AuthTokens>(endpoints.signup, {
      method: 'POST',
      body: input,
    })
    const accessToken = extractAccessToken(payload)
    setToken(accessToken)

    if (!accessToken) return null
    const currentUser = await apiRequest<CurrentUser>(endpoints.me, { token: accessToken })
    setUser(currentUser)
    return currentUser
  }, [])

  const logout = useCallback(() => {
    setToken(null)
    setUser(null)
  }, [])

  const value = useMemo<AuthContextValue>(
    () => ({
      token,
      user,
      isAuthenticated: Boolean(token),
      setToken,
      loadCurrentUser,
      login,
      signup,
      logout,
    }),
    [loadCurrentUser, login, logout, signup, token, user],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider')
  }
  return context
}
