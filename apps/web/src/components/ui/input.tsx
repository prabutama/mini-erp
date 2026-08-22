import { clsx } from 'clsx'
import type { InputHTMLAttributes, LabelHTMLAttributes, ReactNode } from 'react'

export function Input({ className, ...props }: InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input
      className={clsx(
        'h-10 w-full rounded-md border border-hairline bg-canvas px-3.5 text-base text-ink outline-none focus:border-ink',
        className,
      )}
      {...props}
    />
  )
}

export function Label({ className, ...props }: LabelHTMLAttributes<HTMLLabelElement>) {
  return <label className={clsx('text-sm font-medium text-ink', className)} {...props} />
}

export function FormField({ label, children }: { label: string; children: ReactNode }) {
  return (
    <label className="grid gap-2">
      <span className="text-sm font-medium text-ink">{label}</span>
      {children}
    </label>
  )
}
