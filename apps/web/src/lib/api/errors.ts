export type ApiErrorBody = {
  code?: string
  message?: string
  error?: string
}

export class ApiError extends Error {
  status: number
  code: string

  constructor(status: number, body: ApiErrorBody) {
    super(body.message || body.error || `Request failed with status ${status}`)
    this.name = 'ApiError'
    this.status = status
    this.code = body.code || 'REQUEST_FAILED'
  }
}
