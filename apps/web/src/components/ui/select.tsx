import { clsx } from 'clsx'
import type { SelectHTMLAttributes } from 'react'

export function Select({ className, ...props }: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={clsx('h-10 w-full rounded-md border border-hairline bg-canvas px-3.5 text-base text-ink outline-none focus:border-ink', className)}
      {...props}
    />
  )
}
