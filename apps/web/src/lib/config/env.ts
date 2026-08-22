import { z } from 'zod'

const clientEnvSchema = z.object({
  VITE_API_BASE_URL: z.string().default(''),
})

export const clientEnv = clientEnvSchema.parse(import.meta.env)

export function apiBaseUrl() {
  return clientEnv.VITE_API_BASE_URL.replace(/\/$/, '')
}
