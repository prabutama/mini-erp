import { clsx } from 'clsx'
import type { HTMLAttributes } from 'react'

type CardProps = HTMLAttributes<HTMLDivElement> & {
  variant?: 'soft' | 'canvas' | 'dark'
}

export function Card({ className, variant = 'soft', ...props }: CardProps) {
  return (
    <div
      className={clsx(
        'rounded-xl p-6',
        variant === 'soft' && 'bg-surface-card',
        variant === 'canvas' && 'border border-hairline bg-canvas shadow-soft',
        variant === 'dark' && 'bg-surface-dark text-on-dark',
        className,
      )}
      {...props}
    />
  )
}

export function Panel({ className, ...props }: HTMLAttributes<HTMLDivElement>) {
  return <div className={clsx('rounded-xl border border-hairline bg-canvas p-6 shadow-soft', className)} {...props} />
}
