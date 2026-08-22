import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent, useMemo } from 'react'
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
import { operationsApi } from '~/features/operations/operations-api'

type Resource = { resource_id: string; branch_id: string; name: string; code: string; unit: string; type: string; status: string }
type Order = { service_order_id: string; title: string }

export function ResourcesPage() {
  const token = useAuth().token || ''
  const queryClient = useQueryClient()
  const resources = useQuery({ queryKey: ['resources'], queryFn: () => apiRequest<{ resources: Resource[] }>(endpoints.resources, { token }), enabled: Boolean(token) })
  const branches = useQuery({ queryKey: ['branches'], queryFn: async () => normalizeBranches(await listBranches(token)), enabled: Boolean(token) })
  const orders = useQuery({ queryKey: ['service-orders'], queryFn: () => operationsApi.orders(token), enabled: Boolean(token) })
  const createResource = useMutation({ mutationFn: (body: unknown) => apiRequest(endpoints.resources, { method: 'POST', token, body }), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['resources'] }) })
  const stock = useMutation({ mutationFn: ({ resourceId, body }: { resourceId: string; body: unknown }) => apiRequest(endpoints.stockMovements(resourceId), { method: 'POST', token, body }), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['resources'] }) })
  const usage = useMutation({ mutationFn: ({ orderId, body }: { orderId: string; body: unknown }) => apiRequest(endpoints.serviceOrderUsage(orderId), { method: 'POST', token, body }) })
  const columns = useMemo<ColumnDef<Resource>[]>(() => [{ header: 'Name', accessorKey: 'name' }, { header: 'Code', accessorKey: 'code' }, { header: 'Unit', accessorKey: 'unit' }, { header: 'Type', accessorKey: 'type' }, { header: 'Status', cell: ({ row }) => <Badge>{row.original.status}</Badge> }], [])
  const list = resources.data?.resources || []
  async function submitResource(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await createResource.mutateAsync({ branch_id: String(form.get('branch_id') || ''), name: String(form.get('name') || ''), unit: String(form.get('unit') || 'pcs'), type: String(form.get('type') || 'stock') }); event.currentTarget.reset() }
  async function submitStock(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await stock.mutateAsync({ resourceId: String(form.get('resource_id') || ''), body: { movement_type: String(form.get('movement_type') || 'in'), quantity: Number(form.get('quantity') || 0), reason: String(form.get('reason') || ''), service_order_id: String(form.get('service_order_id') || '') } }); event.currentTarget.reset() }
  async function submitUsage(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await usage.mutateAsync({ orderId: String(form.get('order_id') || ''), body: { resource_id: String(form.get('resource_id') || ''), quantity: Number(form.get('quantity') || 0), reason: String(form.get('reason') || '') } }); event.currentTarget.reset() }
  return <div className="grid gap-8 xl:grid-cols-[1fr_380px]"><section><h1 className="display-md">Resources</h1><p className="mt-3 text-body">Branch-scoped resources, stock movement, and order usage.</p><div className="mt-8"><DataTable data={list} columns={columns} empty="No resources yet." /></div></section><aside className="grid gap-6"><Panel><h2 className="text-lg font-semibold">Create resource</h2><form className="mt-5 grid gap-4" onSubmit={submitResource}><FormField label="Branch"><Select name="branch_id" required>{(branches.data || []).map((branch) => <option key={branch.branch_id || branch.id} value={branch.branch_id || branch.id}>{branch.name}</option>)}</Select></FormField><FormField label="Name"><Input name="name" required /></FormField><FormField label="Unit"><Input name="unit" defaultValue="pcs" /></FormField><FormField label="Type"><Input name="type" defaultValue="stock" /></FormField><Button>Create resource</Button></form></Panel><Panel><h2 className="text-lg font-semibold">Stock movement</h2><form className="mt-5 grid gap-4" onSubmit={submitStock}><FormField label="Resource"><Select name="resource_id" required>{list.map((resource) => <option key={resource.resource_id} value={resource.resource_id}>{resource.name}</option>)}</Select></FormField><FormField label="Movement type"><Select name="movement_type"><option value="in">in</option><option value="out">out</option></Select></FormField><FormField label="Quantity"><Input name="quantity" type="number" step="0.01" required /></FormField><FormField label="Reason"><Input name="reason" /></FormField><FormField label="Service order"><Select name="service_order_id"><option value="">None</option>{((orders.data?.service_orders || []) as Order[]).map((order) => <option key={order.service_order_id} value={order.service_order_id}>{order.title}</option>)}</Select></FormField><Button>Record stock</Button></form></Panel><Panel><h2 className="text-lg font-semibold">Order usage</h2><form className="mt-5 grid gap-4" onSubmit={submitUsage}><FormField label="Order"><Select name="order_id" required>{((orders.data?.service_orders || []) as Order[]).map((order) => <option key={order.service_order_id} value={order.service_order_id}>{order.title}</option>)}</Select></FormField><FormField label="Resource"><Select name="resource_id" required>{list.map((resource) => <option key={resource.resource_id} value={resource.resource_id}>{resource.name}</option>)}</Select></FormField><FormField label="Quantity"><Input name="quantity" type="number" step="0.01" required /></FormField><FormField label="Reason"><Input name="reason" /></FormField><Button>Record usage</Button></form></Panel></aside></div>
}
