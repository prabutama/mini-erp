import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { useAuth } from '~/features/auth/auth-context'
import { getCurrentBusiness, getServiceOrdersSummary } from '~/features/dashboard/dashboard-api'
import { DashboardCard } from '~/features/dashboard/dashboard-card'
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

  return (
    <div>
      <div className="mb-8">
        <h1 className="display-md">Dashboard</h1>
        <p className="mt-3 text-body">
          {businessQuery.data?.name || 'Business workspace'}: first branch setup, active orders, and stock signals appear here.
        </p>
      </div>
      <div className="grid gap-6 md:grid-cols-3">
        <DashboardCard label="Orders" value={summaryQuery.isLoading ? '...' : String(total)} detail="From /api/v1/service-orders/summary." />
        <DashboardCard label="Open" value={String(summary?.open ?? summary?.counts?.open ?? 0)} detail="Orders waiting for work." />
        <DashboardCard label="In progress" value={String(summary?.in_progress ?? summary?.counts?.in_progress ?? 0)} detail="Active branch work." />
      </div>
    </div>
  )
}
