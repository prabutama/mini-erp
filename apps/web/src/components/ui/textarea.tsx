import { clsx } from 'clsx'
import type { TextareaHTMLAttributes } from 'react'

export function Textarea({ className, ...props }: TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return <textarea className={clsx('min-h-24 w-full rounded-md border border-hairline bg-canvas px-3.5 py-2 text-base text-ink outline-none focus:border-ink', className)} {...props} />
}
