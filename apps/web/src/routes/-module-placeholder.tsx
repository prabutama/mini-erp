import { Badge } from '~/components/ui/badge'
import { Panel } from '~/components/ui/card'

export function ModulePlaceholder({ title, endpoint, note }: { title: string; endpoint: string; note?: string }) {
  return (
    <div>
      <div className="mb-8 flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="display-md">{title}</h1>
          <p className="mt-3 text-body">MVP screen scaffold ready for backend integration.</p>
        </div>
        <Badge>{endpoint}</Badge>
      </div>
      <Panel>
        <div className="rounded-lg border border-dashed border-hairline bg-surface-soft p-8 text-center">
          <p className="text-lg font-semibold text-ink">Data table and forms go here.</p>
          <p className="mx-auto mt-3 max-w-xl text-body">
            Next pass wires TanStack Query and TanStack Table to the API Gateway endpoint above.
          </p>
          {note ? <p className="mt-4 text-sm text-muted">{note}</p> : null}
        </div>
      </Panel>
    </div>
  )
}
