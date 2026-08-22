import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { FormEvent, useState } from 'react'
import { Button } from '~/components/ui/button'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { useAuth } from '~/features/auth/auth-context'

export const Route = createFileRoute('/signup')({
  component: SignupPage,
})

function SignupPage() {
  const auth = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    const form = new FormData(event.currentTarget)
    const user = await auth.signup({
      business_name: String(form.get('business_name') || ''),
      admin_name: String(form.get('admin_name') || ''),
      email: String(form.get('email') || ''),
      password: String(form.get('password') || ''),
    }).catch((err: Error) => {
      setError(err.message)
      return null
    })

    if (!user) return
    await navigate({ to: '/app' })
  }

  return (
    <main className="grid min-h-screen place-items-center bg-surface-soft px-4 py-12">
      <Panel className="w-full max-w-lg">
        <h1 className="display-md">Create tenant</h1>
        <p className="mt-3 text-body">Signup creates one business and the first Business Admin. Branch comes next.</p>
        <form className="mt-8 grid gap-4" onSubmit={onSubmit}>
          <FormField label="Business name">
            <Input name="business_name" autoComplete="organization" required />
          </FormField>
          <FormField label="Admin name">
            <Input name="admin_name" autoComplete="name" required />
          </FormField>
          <FormField label="Email">
            <Input name="email" type="email" autoComplete="email" required />
          </FormField>
          <FormField label="Password">
            <Input name="password" type="password" autoComplete="new-password" required />
          </FormField>
          {error ? <p className="text-sm text-error">{error}</p> : null}
          <Button type="submit">Create business</Button>
        </form>
      </Panel>
    </main>
  )
}
