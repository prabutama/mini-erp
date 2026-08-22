import { Link, Navigate, Outlet, useLocation } from '@tanstack/react-router'
import { clsx } from 'clsx'
import { useAuth } from '~/features/auth/auth-context'

const navItems = [
  { to: '/app', label: 'Dashboard' },
  { to: '/app/branches', label: 'Branches' },
  { to: '/app/users', label: 'Users' },
  { to: '/app/workflows', label: 'Workflows' },
  { to: '/app/services', label: 'Services' },
  { to: '/app/service-orders', label: 'Orders' },
  { to: '/app/resources', label: 'Resources' },
  { to: '/app/reports', label: 'Reports' },
] as const

export function AppShell() {
  const location = useLocation()
  const auth = useAuth()

  if (!auth.isAuthenticated) {
    return <Navigate to="/login" search={{ redirect: location.href }} />
  }

  if (auth.user?.role === 'Platform Admin') {
    return <Navigate to="/platform/businesses" />
  }

  return (
    <div className="min-h-screen bg-canvas text-ink">
      <aside className="fixed inset-y-0 left-0 hidden w-64 border-r border-hairline bg-canvas p-4 lg:block">
        <div className="mb-8 text-lg font-semibold tracking-[-0.03em]">Mini ERP</div>
        <nav className="grid gap-1">
          {navItems.map((item) => (
            <Link
              key={item.to}
              to={item.to}
              className={clsx(
                'rounded-md px-3 py-2 text-sm font-medium',
                location.pathname === item.to ? 'bg-surface-card text-ink' : 'text-muted hover:text-ink',
              )}
            >
              {item.label}
            </Link>
          ))}
        </nav>
      </aside>
      <main className="lg:pl-64">
        <header className="flex h-16 items-center justify-between border-b border-hairline bg-canvas px-4 sm:px-6 lg:px-8">
          <div>
            <p className="text-sm text-muted">{auth.user?.role || 'Business workspace'}</p>
          </div>
          <button className="text-sm font-semibold text-ink" onClick={auth.logout}>Logout</button>
        </header>
        <div className="px-4 py-8 sm:px-6 lg:px-8">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
