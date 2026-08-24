import { useQuery } from '@tanstack/react-query'
import { Badge } from '~/components/ui/badge'
import { Panel } from '~/components/ui/card'
import { EmptyState, MetricCard, PageHeader, SectionHeader } from '~/components/ui/page'
import { useAuth } from '~/features/auth/auth-context'
import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'

type AuditEvent = { event_id: string; event_type: string; producer: string; occurred_at: string; actor_id: string; data: string }
type OperationsSummary = { open_orders?: number; in_progress_orders?: number; completed_orders?: number; cancelled_orders?: number; resources_used?: number }

function formatDate(value: string) {
  if (!value) return '-'
  return new Intl.DateTimeFormat('en', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value))
}

export function ReportsPage() {
  const token = useAuth().token || ''
  const events = useQuery({ queryKey: ['reports', 'audit-events'], queryFn: () => apiRequest<{ events: AuditEvent[] }>(endpoints.auditEvents, { token }), enabled: Boolean(token) })
  const summary = useQuery({ queryKey: ['reports', 'operations-summary'], queryFn: () => apiRequest<OperationsSummary>(endpoints.operationsSummary, { token }), enabled: Boolean(token) })
  const eventList = events.data?.events || []
  const summaryData = summary.data || {}

  return (
    <div className="grid gap-8">
      <PageHeader eyebrow="Reporting" title="Reports" description="Audit events now flow from NATS JetStream into reporting projections." />

      <div className="grid gap-4 md:grid-cols-4">
        <MetricCard dark label="Audit events" value={eventList.length} detail="Latest projected activity." />
        <MetricCard label="Open" value={summaryData.open_orders ?? 0} detail="Open snapshot count." />
        <MetricCard label="In progress" value={summaryData.in_progress_orders ?? 0} detail="Active snapshot count." />
        <MetricCard label="Resources used" value={summaryData.resources_used ?? 0} detail="Current summary value." />
      </div>

      <Panel>
        <SectionHeader title="Audit timeline" description="Domain events consumed by Reporting with idempotent event IDs." />
        {eventList.length === 0 ? (
          <EmptyState title="No audit events yet" description="Create or update service orders to generate NATS-backed audit activity." />
        ) : (
          <div className="divide-y divide-hairline">
            {eventList.map((event) => (
              <div key={event.event_id} className="grid gap-3 py-4 sm:grid-cols-[1fr_auto] sm:items-center">
                <div>
                  <div className="flex flex-wrap items-center gap-2">
                    <Badge tone="dark">{event.event_type}</Badge>
                    <span className="text-xs font-medium text-muted">{event.producer}</span>
                  </div>
                  <p className="mt-2 text-sm text-body">Actor {event.actor_id ? event.actor_id.slice(0, 8) : 'system'} · event {event.event_id.slice(0, 8)}</p>
                </div>
                <p className="text-sm font-medium text-muted">{formatDate(event.occurred_at)}</p>
              </div>
            ))}
          </div>
        )}
      </Panel>
    </div>
  )
}
