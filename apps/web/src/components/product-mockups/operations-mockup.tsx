import { Badge } from '~/components/ui/badge'

export function OperationsMockup() {
  return (
    <div className="rounded-xl border border-hairline bg-canvas p-4 shadow-card">
      <div className="mb-4 flex items-center justify-between border-b border-hairline pb-3">
        <div>
          <p className="text-sm font-semibold">Service orders</p>
          <p className="text-xs text-muted">Branch-scoped workflow</p>
        </div>
        <Badge className="bg-badge-emerald/20">Live</Badge>
      </div>
      <div className="grid gap-3">
        {[
          ['AC Repair', 'open', 'North Branch'],
          ['Generator Check', 'in_progress', 'Main Branch'],
          ['Water Pump Service', 'completed', 'Main Branch'],
        ].map(([name, status, branch]) => (
          <div key={name} className="rounded-lg border border-hairline bg-surface-soft p-3">
            <div className="flex items-center justify-between gap-3">
              <p className="text-sm font-semibold">{name}</p>
              <Badge>{status}</Badge>
            </div>
            <p className="mt-1 text-xs text-muted">{branch}</p>
          </div>
        ))}
      </div>
    </div>
  )
}
