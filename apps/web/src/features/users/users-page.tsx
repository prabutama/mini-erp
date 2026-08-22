import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent, useMemo, useState } from 'react'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { DataTable, type ColumnDef } from '~/components/ui/data-table'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { Select } from '~/components/ui/select'
import { useAuth } from '~/features/auth/auth-context'
import { listBranches } from '~/features/branches/branch-api'
import { normalizeBranches } from '~/features/branches/types'
import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'

type User = { user_id: string; email: string; full_name: string; status: string }
type Role = { name: string; scope: string }

export function UsersPage() {
  const auth = useAuth()
  const queryClient = useQueryClient()
  const [message, setMessage] = useState<string | null>(null)
  const token = auth.token || ''

  const usersQuery = useQuery({ queryKey: ['users'], queryFn: () => apiRequest<{ users: User[] }>(endpoints.users, { token }), enabled: Boolean(token) })
  const rolesQuery = useQuery({ queryKey: ['roles'], queryFn: () => apiRequest<{ roles: Role[] }>(endpoints.roles, { token }), enabled: Boolean(token) })
  const branchesQuery = useQuery({ queryKey: ['branches'], queryFn: async () => normalizeBranches(await listBranches(token)), enabled: Boolean(token) })

  const createUserMutation = useMutation({
    mutationFn: (body: unknown) => apiRequest<User>(endpoints.users, { method: 'POST', token, body }),
    onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['users'] }),
  })
  const roleMutation = useMutation({ mutationFn: ({ userId, role }: { userId: string; role: string }) => apiRequest(endpoints.userRoles(userId), { method: 'POST', token, body: { role } }) })
  const placementMutation = useMutation({
    mutationFn: (input: { userId: string; branch_id: string; position: string; employment_type: string }) =>
      apiRequest(endpoints.userPlacements(input.userId), { method: 'POST', token, body: { branch_id: input.branch_id, position: input.position, employment_type: input.employment_type } }),
  })

  const columns = useMemo<ColumnDef<User>[]>(() => [
    { header: 'Name', accessorKey: 'full_name' },
    { header: 'Email', accessorKey: 'email' },
    { header: 'Status', cell: ({ row }) => <Badge>{row.original.status}</Badge> },
  ], [])

  async function submitCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setMessage(null)
    const form = new FormData(event.currentTarget)
    await createUserMutation.mutateAsync({
      full_name: String(form.get('full_name') || ''),
      email: String(form.get('email') || ''),
      password: String(form.get('password') || ''),
      role: String(form.get('role') || 'Staff'),
    })
    event.currentTarget.reset()
    setMessage('User created.')
  }

  async function submitAccess(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setMessage(null)
    const form = new FormData(event.currentTarget)
    const userId = String(form.get('user_id') || '')
    const role = String(form.get('role') || '')
    const branchId = String(form.get('branch_id') || '')
    if (role) await roleMutation.mutateAsync({ userId, role })
    if (branchId) await placementMutation.mutateAsync({ userId, branch_id: branchId, position: String(form.get('position') || 'Staff'), employment_type: String(form.get('employment_type') || 'full_time') })
    setMessage('Access updated.')
  }

  return (
    <div className="grid gap-8 xl:grid-cols-[1fr_380px]">
      <section>
        <h1 className="display-md">Users and access</h1>
        <p className="mt-3 text-body">Create tenant users, assign fixed MVP roles, and place users into branches.</p>
        <div className="mt-8"><DataTable data={usersQuery.data?.users || []} columns={columns} empty="No users returned." /></div>
      </section>
      <aside className="grid gap-6">
        <Panel>
          <h2 className="text-lg font-semibold">Create user</h2>
          <form className="mt-5 grid gap-4" onSubmit={submitCreate}>
            <FormField label="Full name"><Input name="full_name" required /></FormField>
            <FormField label="Email"><Input name="email" type="email" required /></FormField>
            <FormField label="Password"><Input name="password" type="password" required /></FormField>
            <FormField label="Role"><Select name="role">{(rolesQuery.data?.roles || []).map((role) => <option key={role.name}>{role.name}</option>)}</Select></FormField>
            <Button disabled={createUserMutation.isPending}>Create user</Button>
          </form>
        </Panel>
        <Panel>
          <h2 className="text-lg font-semibold">Role and placement</h2>
          <form className="mt-5 grid gap-4" onSubmit={submitAccess}>
            <FormField label="User"><Select name="user_id" required>{(usersQuery.data?.users || []).map((user) => <option key={user.user_id} value={user.user_id}>{user.full_name || user.email}</option>)}</Select></FormField>
            <FormField label="Role"><Select name="role"><option value="">No role change</option>{(rolesQuery.data?.roles || []).map((role) => <option key={role.name}>{role.name}</option>)}</Select></FormField>
            <FormField label="Branch"><Select name="branch_id"><option value="">No placement</option>{(branchesQuery.data || []).map((branch) => <option key={branch.branch_id || branch.id} value={branch.branch_id || branch.id}>{branch.name}</option>)}</Select></FormField>
            <FormField label="Position"><Input name="position" placeholder="Technician" /></FormField>
            <FormField label="Employment type"><Input name="employment_type" placeholder="full_time" /></FormField>
            <Button disabled={roleMutation.isPending || placementMutation.isPending}>Save access</Button>
          </form>
          {message ? <p className="mt-3 text-sm text-success">{message}</p> : null}
        </Panel>
      </aside>
    </div>
  )
}
