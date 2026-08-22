import { clsx } from 'clsx'
import type { HTMLAttributes } from 'react'

export function Badge({ className, ...props }: HTMLAttributes<HTMLSpanElement>) {
  return <span className={clsx('inline-flex rounded-full bg-surface-card px-3 py-1 text-[13px] font-medium text-ink', className)} {...props} />
}
