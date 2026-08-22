import { apiBaseUrl } from '~/lib/config/env'
import { ApiError, type ApiErrorBody } from './errors'

type ApiOptions = Omit<RequestInit, 'body'> & {
  token?: string | null
  body?: unknown
}

export async function apiRequest<T>(path: string, options: ApiOptions = {}): Promise<T> {
  const { token, body, headers, ...init } = options
  const requestHeaders = new Headers(headers)

  if (body !== undefined) {
    requestHeaders.set('Content-Type', 'application/json')
  }

  if (token) {
    requestHeaders.set('Authorization', `Bearer ${token}`)
  }

  const response = await fetch(`${apiBaseUrl()}${path}`, {
    ...init,
    headers: requestHeaders,
    body: body === undefined ? undefined : JSON.stringify(body),
  })

  const text = await response.text()
  const payload = text ? JSON.parse(text) : null

  if (!response.ok) {
    throw new ApiError(response.status, (payload || {}) as ApiErrorBody)
  }

  return payload as T
}
