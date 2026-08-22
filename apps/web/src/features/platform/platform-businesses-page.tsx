import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent, useMemo } from 'react'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { DataTable, type ColumnDef } from '~/components/ui/data-table'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { Select } from '~/components/ui/select'
import { Textarea } from '~/components/ui/textarea'
import { useAuth } from '~/features/auth/auth-context'
import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'

type Business = { business_id: string; name: string; code: string; status: string; plan: string; platform_notes?: string }

export function PlatformBusinessesPage() {
  const token = useAuth().token || ''
  const queryClient = useQueryClient()
  const businesses = useQuery({ queryKey: ['platform', 'businesses'], queryFn: () => apiRequest<{ businesses: Business[] }>(endpoints.platformBusinesses, { token }), enabled: Boolean(token) })
  const update = useMutation({ mutationFn: ({ businessId, body }: { businessId: string; body: unknown }) => apiRequest(endpoints.platformBusiness(businessId), { method: 'PATCH', token, body }), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['platform', 'businesses'] }) })
  const list = businesses.data?.businesses || []
  const columns = useMemo<ColumnDef<Business>[]>(() => [{ header: 'Name', accessorKey: 'name' }, { header: 'Code', accessorKey: 'code' }, { header: 'Plan', accessorKey: 'plan' }, { header: 'Status', cell: ({ row }) => <Badge>{row.original.status}</Badge> }], [])
  async function submit(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await update.mutateAsync({ businessId: String(form.get('business_id') || ''), body: { status: String(form.get('status') || ''), plan: String(form.get('plan') || ''), platform_notes: String(form.get('platform_notes') || ''), suspended_at: String(form.get('suspended_at') || '') } }) }
  return <div className="grid gap-8 xl:grid-cols-[1fr_380px]"><section><h1 className="display-md">Platform businesses</h1><p className="mt-3 text-body">Oversight-only tenant records. No branch, user, workflow, resource, or report management here.</p><div className="mt-8"><DataTable data={list} columns={columns} empty="No businesses returned." /></div></section><Panel><h2 className="text-lg font-semibold">Update platform fields</h2><form className="mt-5 grid gap-4" onSubmit={submit}><FormField label="Business"><Select name="business_id" required>{list.map((business) => <option key={business.business_id} value={business.business_id}>{business.name}</option>)}</Select></FormField><FormField label="Status"><Input name="status" placeholder="active" /></FormField><FormField label="Plan"><Input name="plan" placeholder="mvp" /></FormField><FormField label="Suspended at"><Input name="suspended_at" placeholder="RFC3339 or blank" /></FormField><FormField label="Platform notes"><Textarea name="platform_notes" /></FormField><Button>Update tenant</Button></form></Panel></div>
}
