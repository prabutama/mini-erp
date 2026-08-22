import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { FormEvent, useState } from 'react'
import { Button } from '~/components/ui/button'
import { FormField, Input } from '~/components/ui/input'
import { Panel } from '~/components/ui/card'
import { useAuth } from '~/features/auth/auth-context'

export const Route = createFileRoute('/login')({
  component: LoginPage,
})

function LoginPage() {
  const auth = useAuth()
  const navigate = useNavigate()
  const [error, setError] = useState<string | null>(null)

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    const form = new FormData(event.currentTarget)
    const user = await auth.login({
      email: String(form.get('email') || ''),
      password: String(form.get('password') || ''),
    }).catch((err: Error) => {
      setError(err.message)
      return null
    })

    if (!user) return
    await navigate({ to: user.role === 'Platform Admin' ? '/platform/businesses' : '/app' })
  }

  return (
    <main className="grid min-h-screen place-items-center bg-surface-soft px-4 py-12">
      <Panel className="w-full max-w-md">
        <h1 className="display-md">Login</h1>
        <p className="mt-3 text-body">Use your Mini ERP account to continue.</p>
        <form className="mt-8 grid gap-4" onSubmit={onSubmit}>
          <FormField label="Email">
            <Input name="email" type="email" autoComplete="email" required />
          </FormField>
          <FormField label="Password">
            <Input name="password" type="password" autoComplete="current-password" required />
          </FormField>
          {error ? <p className="text-sm text-error">{error}</p> : null}
          <Button type="submit">Login</Button>
        </form>
      </Panel>
    </main>
  )
}
