import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { ButtonLink } from '~/components/ui/button'
import { Card, Panel } from '~/components/ui/card'
import { MetricCard, PageHeader, SectionHeader } from '~/components/ui/page'
import { useAuth } from '~/features/auth/auth-context'
import { getCurrentBusiness, getServiceOrdersSummary } from '~/features/dashboard/dashboard-api'
import { queryKeys } from '~/lib/api/query-keys'

export const Route = createFileRoute('/app/')({
  component: DashboardPage,
})

function DashboardPage() {
  const auth = useAuth()
  const businessQuery = useQuery({
    queryKey: queryKeys.currentBusiness,
    queryFn: () => getCurrentBusiness(auth.token || ''),
    enabled: Boolean(auth.token),
  })
  const summaryQuery = useQuery({
    queryKey: queryKeys.serviceOrdersSummary,
    queryFn: () => getServiceOrdersSummary(auth.token || ''),
    enabled: Boolean(auth.token),
  })

  const summary = summaryQuery.data
  const total = summary?.total ?? Object.values(summary?.counts || {}).reduce((sum, value) => sum + value, 0)
  const businessName = businessQuery.data?.name || 'Business workspace'

  return (
    <div className="grid gap-8">
      <PageHeader
        eyebrow="Command center"
        title={businessName}
        description="Track branch setup, work order movement, and operational signals from one focused workspace."
        actions={<ButtonLink to="/app/service-orders">Open work queue</ButtonLink>}
      />

      <div className="grid gap-4 md:grid-cols-4">
        <MetricCard dark label="Total orders" value={summaryQuery.isLoading ? '...' : total} detail="All internal work orders." />
        <MetricCard label="Open" value={summary?.open ?? summary?.counts?.open ?? 0} detail="Waiting for assignment." />
        <MetricCard label="In progress" value={summary?.in_progress ?? summary?.counts?.in_progress ?? 0} detail="Active branch work." />
        <MetricCard label="Completed" value={summary?.completed ?? summary?.counts?.completed ?? 0} detail="Finished orders." />
      </div>

      <div className="grid gap-6 lg:grid-cols-[1.2fr_0.8fr]">
        <Panel>
          <SectionHeader title="Next best actions" description="Set up the operational graph in the right order." />
          <div className="grid gap-3 sm:grid-cols-2">
            <ButtonLink variant="secondary" to="/app/branches" className="justify-start">Create branch</ButtonLink>
            <ButtonLink variant="secondary" to="/app/services" className="justify-start">Define services</ButtonLink>
            <ButtonLink variant="secondary" to="/app/workflows" className="justify-start">Configure workflow</ButtonLink>
            <ButtonLink variant="secondary" to="/app/resources" className="justify-start">Add resources</ButtonLink>
          </div>
        </Panel>
        <Card variant="dark">
          <p className="text-sm font-medium text-on-dark-soft">MVP status</p>
          <p className="mt-4 text-2xl font-semibold tracking-[-0.04em]">Backend contracts connected.</p>
          <p className="mt-3 text-sm leading-6 text-on-dark-soft">Signup, branches, access, service orders, resources, NATS audit events, and reporting run through API Gateway.</p>
        </Card>
      </div>
    </div>
  )
}
