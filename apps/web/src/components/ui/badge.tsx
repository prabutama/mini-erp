import { clsx } from 'clsx'
import type { HTMLAttributes } from 'react'

type BadgeProps = HTMLAttributes<HTMLSpanElement> & {
  tone?: 'neutral' | 'success' | 'warning' | 'dark'
}

export function Badge({ className, tone = 'neutral', ...props }: BadgeProps) {
  return (
    <span
      className={clsx(
        'inline-flex rounded-full px-3 py-1 text-[13px] font-semibold leading-none',
        tone === 'neutral' && 'bg-surface-card text-ink',
        tone === 'success' && 'bg-emerald-50 text-emerald-700',
        tone === 'warning' && 'bg-amber-50 text-amber-700',
        tone === 'dark' && 'bg-ink text-white',
        className,
      )}
      {...props}
    />
  )
}
