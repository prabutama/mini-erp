import { clsx } from 'clsx'
import type { ReactNode } from 'react'

export function PageHeader({ eyebrow, title, description, actions }: { eyebrow?: string; title: string; description?: string; actions?: ReactNode }) {
  return (
    <div className="flex flex-col gap-5 border-b border-hairline pb-8 lg:flex-row lg:items-end lg:justify-between">
      <div>
        {eyebrow ? <p className="mb-3 text-xs font-semibold uppercase tracking-[0.18em] text-muted">{eyebrow}</p> : null}
        <h1 className="display-md max-w-3xl">{title}</h1>
        {description ? <p className="mt-3 max-w-2xl text-body">{description}</p> : null}
      </div>
      {actions ? <div className="flex flex-wrap gap-3">{actions}</div> : null}
    </div>
  )
}

export function SectionHeader({ title, description, className }: { title: string; description?: string; className?: string }) {
  return (
    <div className={clsx('mb-5', className)}>
      <h2 className="text-lg font-semibold tracking-[-0.02em] text-ink">{title}</h2>
      {description ? <p className="mt-1 text-sm text-muted">{description}</p> : null}
    </div>
  )
}

export function EmptyState({ title, description, action }: { title: string; description?: string; action?: ReactNode }) {
  return (
    <div className="rounded-xl border border-dashed border-hairline bg-surface-soft p-10 text-center">
      <p className="text-base font-semibold text-ink">{title}</p>
      {description ? <p className="mx-auto mt-2 max-w-md text-sm text-muted">{description}</p> : null}
      {action ? <div className="mt-5 flex justify-center">{action}</div> : null}
    </div>
  )
}

export function MetricCard({ label, value, detail, dark }: { label: string; value: string | number; detail?: string; dark?: boolean }) {
  return (
    <div className={clsx('rounded-xl p-5', dark ? 'bg-surface-dark text-white' : 'border border-hairline bg-canvas shadow-soft')}>
      <p className={clsx('text-sm font-medium', dark ? 'text-on-dark-soft' : 'text-muted')}>{label}</p>
      <p className="mt-3 text-3xl font-semibold tracking-[-0.04em]">{value}</p>
      {detail ? <p className={clsx('mt-2 text-sm', dark ? 'text-on-dark-soft' : 'text-body')}>{detail}</p> : null}
    </div>
  )
}
