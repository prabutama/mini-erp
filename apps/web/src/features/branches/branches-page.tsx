import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent, useState } from 'react'
import { Button } from '~/components/ui/button'
import { Badge } from '~/components/ui/badge'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { useAuth } from '~/features/auth/auth-context'
import { queryKeys } from '~/lib/api/query-keys'
import { createBranch, listBranches } from './branch-api'
import { normalizeBranches } from './types'

export function BranchesPage() {
  const auth = useAuth()
  const queryClient = useQueryClient()
  const [error, setError] = useState<string | null>(null)

  const branchesQuery = useQuery({
    queryKey: queryKeys.branches,
    queryFn: async () => normalizeBranches(await listBranches(auth.token || '')),
    enabled: Boolean(auth.token),
  })

  const createMutation = useMutation({
    mutationFn: (input: { name: string; code?: string; address?: string }) => createBranch(auth.token || '', input),
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: queryKeys.branches })
    },
  })

  async function onSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setError(null)
    const form = new FormData(event.currentTarget)
    const name = String(form.get('name') || '').trim()
    const code = String(form.get('code') || '').trim()
    const address = String(form.get('address') || '').trim()

    await createMutation.mutateAsync({ name, code: code || undefined, address: address || undefined }).catch((err: Error) => {
      setError(err.message)
    })

    if (!createMutation.isError) {
      event.currentTarget.reset()
    }
  }

  const branches = branchesQuery.data || []

  return (
    <div className="grid gap-8 xl:grid-cols-[1fr_380px]">
      <section>
        <div className="mb-8 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h1 className="display-md">Branches</h1>
            <p className="mt-3 text-body">Create first branch after tenant signup. API Gateway validates business scope.</p>
          </div>
          <Badge>GET /api/v1/branches</Badge>
        </div>

        <Panel className="overflow-hidden p-0">
          {branchesQuery.isLoading ? (
            <p className="p-6 text-body">Loading branches...</p>
          ) : branches.length === 0 ? (
            <div className="p-10 text-center">
              <p className="text-lg font-semibold">No branch yet.</p>
              <p className="mt-2 text-body">Create one to unlock placements, orders, and resources.</p>
            </div>
          ) : (
            <div className="divide-y divide-hairline">
              {branches.map((branch) => {
                const id = branch.id || branch.branch_id || branch.name || 'branch'
                return (
                  <div key={id} className="grid gap-2 p-5 sm:grid-cols-[1fr_auto] sm:items-center">
                    <div>
                      <p className="font-semibold text-ink">{branch.name || 'Unnamed branch'}</p>
                      <p className="mt-1 text-sm text-muted">{branch.address || 'No address set'}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      {branch.code ? <Badge>{branch.code}</Badge> : null}
                      {branch.status ? <Badge>{branch.status}</Badge> : null}
                    </div>
                  </div>
                )
              })}
            </div>
          )}
        </Panel>
      </section>

      <Panel>
        <h2 className="text-lg font-semibold">Create branch</h2>
        <p className="mt-2 text-sm text-body">Do not send `business_id`; backend derives it from authenticated context.</p>
        <form className="mt-6 grid gap-4" onSubmit={onSubmit}>
          <FormField label="Branch name">
            <Input name="name" required placeholder="Main Branch" />
          </FormField>
          <FormField label="Code">
            <Input name="code" placeholder="MAIN" />
          </FormField>
          <FormField label="Address">
            <Input name="address" placeholder="Street, city" />
          </FormField>
          {error ? <p className="text-sm text-error">{error}</p> : null}
          <Button type="submit" disabled={createMutation.isPending}>
            {createMutation.isPending ? 'Creating...' : 'Create branch'}
          </Button>
        </form>
      </Panel>
    </div>
  )
}
