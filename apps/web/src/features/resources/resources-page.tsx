import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent } from 'react'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { DataTable, type ColumnDef } from '~/components/ui/data-table'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { MetricCard, PageHeader, SectionHeader } from '~/components/ui/page'
import { Select } from '~/components/ui/select'
import { useAuth } from '~/features/auth/auth-context'
import { listBranches } from '~/features/branches/branch-api'
import { normalizeBranches } from '~/features/branches/types'
import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'
import { operationsApi } from '~/features/operations/operations-api'

type Resource = { resource_id: string; branch_id: string; name: string; code: string; unit: string; type: string; status: string }
type Order = { service_order_id: string; title: string }

function shortId(value: string) {
  return value ? value.slice(0, 8) : '-'
}

export function ResourcesPage() {
  const token = useAuth().token || ''
  const queryClient = useQueryClient()
  const resources = useQuery({ queryKey: ['resources'], queryFn: () => apiRequest<{ resources: Resource[] }>(endpoints.resources, { token }), enabled: Boolean(token) })
  const branches = useQuery({ queryKey: ['branches'], queryFn: async () => normalizeBranches(await listBranches(token)), enabled: Boolean(token) })
  const orders = useQuery({ queryKey: ['service-orders'], queryFn: () => operationsApi.orders(token), enabled: Boolean(token) })
  const createResource = useMutation({ mutationFn: (body: unknown) => apiRequest(endpoints.resources, { method: 'POST', token, body }), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['resources'] }) })
  const stock = useMutation({ mutationFn: ({ resourceId, body }: { resourceId: string; body: unknown }) => apiRequest(endpoints.stockMovements(resourceId), { method: 'POST', token, body }), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['resources'] }) })
  const usage = useMutation({ mutationFn: ({ orderId, body }: { orderId: string; body: unknown }) => apiRequest(endpoints.serviceOrderUsage(orderId), { method: 'POST', token, body }) })
  const list = resources.data?.resources || []
  const stockCount = list.filter((resource) => resource.type === 'stock').length

  const columns: ColumnDef<Resource>[] = [
    {
      header: 'Resource',
      cell: ({ row }) => (
        <div>
          <p className="font-semibold text-ink">{row.original.name}</p>
          <p className="mt-1 text-xs text-muted">{row.original.code || shortId(row.original.resource_id)} · {row.original.unit}</p>
        </div>
      ),
    },
    { header: 'Type', cell: ({ row }) => <Badge>{row.original.type}</Badge> },
    { header: 'Status', cell: ({ row }) => <Badge tone={row.original.status === 'active' ? 'success' : 'neutral'}>{row.original.status}</Badge> },
    { header: 'Branch', cell: ({ row }) => <span className="font-mono text-xs text-muted">{shortId(row.original.branch_id)}</span> },
  ]

  async function submitResource(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    await createResource.mutateAsync({ branch_id: String(form.get('branch_id') || ''), name: String(form.get('name') || ''), unit: String(form.get('unit') || 'pcs'), type: String(form.get('type') || 'stock') })
    event.currentTarget.reset()
  }

  async function submitStock(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    await stock.mutateAsync({ resourceId: String(form.get('resource_id') || ''), body: { movement_type: String(form.get('movement_type') || 'in'), quantity: Number(form.get('quantity') || 0), reason: String(form.get('reason') || ''), service_order_id: String(form.get('service_order_id') || '') } })
    event.currentTarget.reset()
  }

  async function submitUsage(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    await usage.mutateAsync({ orderId: String(form.get('order_id') || ''), body: { resource_id: String(form.get('resource_id') || ''), quantity: Number(form.get('quantity') || 0), reason: String(form.get('reason') || '') } })
    event.currentTarget.reset()
  }

  return (
    <div className="grid gap-8">
      <PageHeader eyebrow="Inventory" title="Resources" description="Manage branch-scoped resources, stock movements, and usage against service orders." />

      <div className="grid gap-4 md:grid-cols-3">
        <MetricCard dark label="Total resources" value={list.length} detail="Available in this business." />
        <MetricCard label="Stock items" value={stockCount} detail="Trackable resource inventory." />
        <MetricCard label="Orders loaded" value={orders.data?.service_orders?.length || 0} detail="Can receive resource usage." />
      </div>

      <div className="grid gap-8 xl:grid-cols-[1fr_400px]">
        <section>
          <SectionHeader title="Resource catalog" description="Inventory and non-stock resources organized by branch." />
          <DataTable data={list} columns={columns} empty="No resources yet." />
        </section>

        <aside className="grid gap-6">
          <Panel>
            <SectionHeader title="Create resource" description="Add a resource under a branch." />
            <form className="grid gap-4" onSubmit={submitResource}>
              <FormField label="Branch"><Select name="branch_id" required>{(branches.data || []).map((branch) => <option key={branch.branch_id || branch.id} value={branch.branch_id || branch.id}>{branch.name}</option>)}</Select></FormField>
              <FormField label="Name"><Input name="name" required /></FormField>
              <FormField label="Unit"><Input name="unit" defaultValue="pcs" /></FormField>
              <FormField label="Type"><Input name="type" defaultValue="stock" /></FormField>
              <Button disabled={createResource.isPending}>Create resource</Button>
            </form>
          </Panel>

          <Panel>
            <SectionHeader title="Stock movement" description="Record inbound or outbound stock." />
            <form className="grid gap-4" onSubmit={submitStock}>
              <FormField label="Resource"><Select name="resource_id" required>{list.map((resource) => <option key={resource.resource_id} value={resource.resource_id}>{resource.name}</option>)}</Select></FormField>
              <FormField label="Movement type"><Select name="movement_type"><option value="in">in</option><option value="out">out</option></Select></FormField>
              <FormField label="Quantity"><Input name="quantity" type="number" step="0.01" required /></FormField>
              <FormField label="Reason"><Input name="reason" /></FormField>
              <FormField label="Service order"><Select name="service_order_id"><option value="">None</option>{((orders.data?.service_orders || []) as Order[]).map((order) => <option key={order.service_order_id} value={order.service_order_id}>{order.title}</option>)}</Select></FormField>
              <Button disabled={stock.isPending}>Record stock</Button>
            </form>
          </Panel>

          <Panel>
            <SectionHeader title="Order usage" description="Attach resource consumption to work orders." />
            <form className="grid gap-4" onSubmit={submitUsage}>
              <FormField label="Order"><Select name="order_id" required>{((orders.data?.service_orders || []) as Order[]).map((order) => <option key={order.service_order_id} value={order.service_order_id}>{order.title}</option>)}</Select></FormField>
              <FormField label="Resource"><Select name="resource_id" required>{list.map((resource) => <option key={resource.resource_id} value={resource.resource_id}>{resource.name}</option>)}</Select></FormField>
              <FormField label="Quantity"><Input name="quantity" type="number" step="0.01" required /></FormField>
              <FormField label="Reason"><Input name="reason" /></FormField>
              <Button disabled={usage.isPending}>Record usage</Button>
            </form>
          </Panel>
        </aside>
      </div>
    </div>
  )
}
