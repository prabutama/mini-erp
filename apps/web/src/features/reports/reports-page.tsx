import { useQuery } from '@tanstack/react-query'
import { useMemo } from 'react'
import { DataTable, type ColumnDef } from '~/components/ui/data-table'
import { Panel } from '~/components/ui/card'
import { useAuth } from '~/features/auth/auth-context'
import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'

type AuditEvent = { event_id: string; event_type: string; producer: string; occurred_at: string; actor_id: string }

export function ReportsPage() {
  const token = useAuth().token || ''
  const events = useQuery({ queryKey: ['reports', 'audit-events'], queryFn: () => apiRequest<{ events: AuditEvent[] }>(endpoints.auditEvents, { token }), enabled: Boolean(token) })
  const summary = useQuery({ queryKey: ['reports', 'operations-summary'], queryFn: () => apiRequest(endpoints.operationsSummary, { token }), enabled: Boolean(token) })
  const columns = useMemo<ColumnDef<AuditEvent>[]>(() => [{ header: 'Event', accessorKey: 'event_type' }, { header: 'Producer', accessorKey: 'producer' }, { header: 'Actor', accessorKey: 'actor_id' }, { header: 'Occurred', accessorKey: 'occurred_at' }], [])
  return <div><h1 className="display-md">Reports</h1><p className="mt-3 text-body">Reporting reads from reporting_db. Empty audit data is expected until NATS ingestion is wired.</p><div className="mt-8 grid gap-6 lg:grid-cols-3"><Panel><p className="text-sm text-muted">Audit events</p><p className="mt-3 text-3xl font-semibold">{events.data?.events?.length || 0}</p></Panel><Panel className="lg:col-span-2"><p className="text-sm text-muted">Operations summary</p><pre className="mt-3 overflow-auto rounded-md bg-surface-soft p-3 text-xs">{JSON.stringify(summary.data || {}, null, 2)}</pre></Panel></div><div className="mt-8"><DataTable data={events.data?.events || []} columns={columns} empty="No audit events projected yet." /></div></div>
}
