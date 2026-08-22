import { Link, type LinkProps } from '@tanstack/react-router'
import { clsx } from 'clsx'
import type { ButtonHTMLAttributes, ComponentProps, ReactNode } from 'react'

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: 'primary' | 'secondary' | 'ghost'
}

const variants = {
  primary: 'bg-primary text-white hover:bg-primary-active',
  secondary: 'border border-hairline bg-canvas text-ink',
  ghost: 'bg-transparent text-ink',
}

export function Button({ className, variant = 'primary', ...props }: ButtonProps) {
  return (
    <button
      className={clsx(
        'inline-flex h-10 items-center justify-center rounded-md px-5 text-sm font-semibold transition-colors disabled:bg-primary-disabled disabled:text-muted',
        variants[variant],
        className,
      )}
      {...props}
    />
  )
}

type ButtonLinkProps = ComponentProps<typeof Link> & {
  variant?: ButtonProps['variant']
  children: ReactNode
}

export function ButtonLink({ className, variant = 'primary', children, ...props }: ButtonLinkProps) {
  return (
    <Link
      className={clsx(
        'inline-flex h-10 items-center justify-center rounded-md px-5 text-sm font-semibold transition-colors',
        variants[variant],
        className,
      )}
      {...props}
    >
      {children}
    </Link>
  )
}
