import { Card } from '~/components/ui/card'

export function DashboardCard({ label, value, detail }: { label: string; value: string; detail: string }) {
  return (
    <Card>
      <p className="text-sm font-medium text-muted">{label}</p>
      <p className="mt-3 text-3xl font-semibold tracking-[-0.03em] text-ink">{value}</p>
      <p className="mt-2 text-sm text-body">{detail}</p>
    </Card>
  )
}
