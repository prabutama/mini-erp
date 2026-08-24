import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent } from 'react'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { DataTable, type ColumnDef } from '~/components/ui/data-table'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { EmptyState, MetricCard, PageHeader, SectionHeader } from '~/components/ui/page'
import { Select } from '~/components/ui/select'
import { useAuth } from '~/features/auth/auth-context'
import { listBranches } from '~/features/branches/branch-api'
import { normalizeBranches } from '~/features/branches/types'
import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'
import { operationsApi, type ServiceOrder } from './operations-api'

type User = { user_id: string; full_name: string; email: string }

function statusTone(status: string) {
  if (status === 'completed') return 'success'
  if (status === 'in_progress') return 'dark'
  if (status === 'cancelled') return 'warning'
  return 'neutral'
}

function shortId(value: string) {
  return value ? value.slice(0, 8) : '-'
}

export function OrdersPage() {
  const token = useAuth().token || ''
  const queryClient = useQueryClient()
  const orders = useQuery({ queryKey: ['service-orders'], queryFn: () => operationsApi.orders(token), enabled: Boolean(token) })
  const services = useQuery({ queryKey: ['service-definitions'], queryFn: () => operationsApi.services(token), enabled: Boolean(token) })
  const branches = useQuery({ queryKey: ['branches'], queryFn: async () => normalizeBranches(await listBranches(token)), enabled: Boolean(token) })
  const users = useQuery({ queryKey: ['users'], queryFn: () => apiRequest<{ users: User[] }>(endpoints.users, { token }), enabled: Boolean(token) })
  const createOrder = useMutation({ mutationFn: (body: unknown) => operationsApi.createOrder(token, body), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['service-orders'] }) })
  const assignOrder = useMutation({ mutationFn: ({ orderId, userId }: { orderId: string; userId: string }) => operationsApi.assignOrder(token, orderId, userId), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['service-orders'] }) })
  const transitionOrder = useMutation({ mutationFn: ({ orderId, status }: { orderId: string; status: string }) => operationsApi.transitionOrder(token, orderId, status), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['service-orders'] }) })
  const orderList = orders.data?.service_orders || []
  const openCount = orderList.filter((order) => order.status === 'open').length
  const activeCount = orderList.filter((order) => order.status === 'in_progress').length

  const columns: ColumnDef<ServiceOrder>[] = [
    {
      header: 'Order',
      cell: ({ row }) => (
        <div>
          <p className="font-semibold text-ink">{row.original.title}</p>
          <p className="mt-1 text-xs text-muted">#{shortId(row.original.service_order_id)} · service {shortId(row.original.service_definition_id)}</p>
        </div>
      ),
    },
    { header: 'Priority', cell: ({ row }) => <Badge>{row.original.priority || 'normal'}</Badge> },
    { header: 'Status', cell: ({ row }) => <Badge tone={statusTone(row.original.status)}>{row.original.status}</Badge> },
    { header: 'Branch', cell: ({ row }) => <span className="font-mono text-xs text-muted">{shortId(row.original.branch_id)}</span> },
  ]

  async function submitCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    await createOrder.mutateAsync({
      branch_id: String(form.get('branch_id') || ''),
      service_definition_id: String(form.get('service_definition_id') || ''),
      title: String(form.get('title') || ''),
      description: String(form.get('description') || ''),
      priority: String(form.get('priority') || 'normal'),
    })
    event.currentTarget.reset()
  }

  async function submitAction(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    const form = new FormData(event.currentTarget)
    const orderId = String(form.get('order_id') || '')
    const assignedUserId = String(form.get('assigned_user_id') || '')
    const status = String(form.get('status') || '')
    if (assignedUserId) await assignOrder.mutateAsync({ orderId, userId: assignedUserId })
    if (status) await transitionOrder.mutateAsync({ orderId, status })
  }

  return (
    <div className="grid gap-8">
      <PageHeader
        eyebrow="Operations"
        title="Service orders"
        description="Create, assign, and move internal work orders through the branch queue."
      />

      <div className="grid gap-4 md:grid-cols-3">
        <MetricCard dark label="Total queue" value={orderList.length} detail="All visible work orders." />
        <MetricCard label="Open" value={openCount} detail="Ready for assignment." />
        <MetricCard label="In progress" value={activeCount} detail="Currently being handled." />
      </div>

      <div className="grid gap-8 xl:grid-cols-[1fr_400px]">
        <section>
          <SectionHeader title="Work queue" description="Operational queue with status, branch, and service references." />
          {orders.isLoading ? <EmptyState title="Loading service orders" /> : <DataTable data={orderList} columns={columns} empty="No service orders yet." />}
        </section>

        <aside className="grid gap-6">
          <Panel>
            <SectionHeader title="Create order" description="Start new branch work from a service definition." />
            <form className="grid gap-4" onSubmit={submitCreate}>
              <FormField label="Branch">
                <Select name="branch_id" required>{(branches.data || []).map((branch) => <option key={branch.branch_id || branch.id} value={branch.branch_id || branch.id}>{branch.name}</option>)}</Select>
              </FormField>
              <FormField label="Service">
                <Select name="service_definition_id" required>{(services.data?.service_definitions || []).map((service) => <option key={service.service_definition_id} value={service.service_definition_id}>{service.name}</option>)}</Select>
              </FormField>
              <FormField label="Title"><Input name="title" required /></FormField>
              <FormField label="Description"><Input name="description" /></FormField>
              <FormField label="Priority"><Input name="priority" defaultValue="normal" /></FormField>
              <Button disabled={createOrder.isPending}>{createOrder.isPending ? 'Creating...' : 'Create order'}</Button>
            </form>
          </Panel>

          <Panel>
            <SectionHeader title="Assign / transition" description="Route work to a user or move status forward." />
            <form className="grid gap-4" onSubmit={submitAction}>
              <FormField label="Order">
                <Select name="order_id" required>{orderList.map((order) => <option key={order.service_order_id} value={order.service_order_id}>{order.title}</option>)}</Select>
              </FormField>
              <FormField label="Assign user">
                <Select name="assigned_user_id"><option value="">No assignment</option>{(users.data?.users || []).map((user) => <option key={user.user_id} value={user.user_id}>{user.full_name || user.email}</option>)}</Select>
              </FormField>
              <FormField label="Transition">
                <Select name="status"><option value="">No transition</option><option value="in_progress">in_progress</option><option value="completed">completed</option><option value="cancelled">cancelled</option></Select>
              </FormField>
              <Button disabled={assignOrder.isPending || transitionOrder.isPending}>Apply update</Button>
            </form>
          </Panel>
        </aside>
      </div>
    </div>
  )
}
