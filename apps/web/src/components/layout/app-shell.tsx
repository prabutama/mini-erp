import { Link, Navigate, Outlet, useLocation } from '@tanstack/react-router'
import { clsx } from 'clsx'
import { useAuth } from '~/features/auth/auth-context'

const navItems = [
  { section: 'Workspace', items: [{ to: '/app', label: 'Dashboard' }] },
  { section: 'Operations', items: [{ to: '/app/workflows', label: 'Workflows' }, { to: '/app/services', label: 'Services' }, { to: '/app/service-orders', label: 'Service orders' }, { to: '/app/resources', label: 'Resources' }] },
  { section: 'People', items: [{ to: '/app/branches', label: 'Branches' }, { to: '/app/users', label: 'Users & access' }] },
  { section: 'Insight', items: [{ to: '/app/reports', label: 'Reports' }] },
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
    <div className="min-h-screen bg-surface-soft text-ink">
      <aside className="fixed inset-y-0 left-0 hidden w-72 border-r border-hairline bg-canvas p-5 lg:block">
        <div className="mb-8 rounded-2xl bg-surface-dark p-4 text-white">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-white text-sm font-bold text-ink">ME</div>
            <div>
              <p className="font-semibold tracking-[-0.03em]">Mini ERP</p>
              <p className="text-xs text-on-dark-soft">Service business workspace</p>
            </div>
          </div>
          <div className="mt-4 rounded-xl bg-surface-dark-elevated px-3 py-2 text-xs text-on-dark-soft">
            {auth.user?.role || 'Business workspace'}
          </div>
        </div>
        <nav className="grid gap-6">
          {navItems.map((group) => (
            <div key={group.section}>
              <p className="mb-2 px-3 text-[11px] font-semibold uppercase tracking-[0.18em] text-muted">{group.section}</p>
              <div className="grid gap-1">
                {group.items.map((item) => {
                  const active = location.pathname === item.to
                  return (
                    <Link
                      key={item.to}
                      to={item.to}
                      className={clsx(
                        'relative rounded-lg px-3 py-2.5 text-sm font-semibold transition-colors',
                        active ? 'bg-surface-card text-ink' : 'text-muted hover:bg-surface-soft hover:text-ink',
                      )}
                    >
                      {active ? <span className="absolute left-0 top-2 h-6 w-1 rounded-r-full bg-ink" /> : null}
                      <span className="pl-2">{item.label}</span>
                    </Link>
                  )
                })}
              </div>
            </div>
          ))}
        </nav>
      </aside>
      <main className="lg:pl-72">
        <header className="sticky top-0 z-10 flex h-16 items-center justify-between border-b border-hairline bg-canvas/90 px-4 backdrop-blur sm:px-6 lg:px-8">
          <div>
            <p className="text-xs font-semibold uppercase tracking-[0.18em] text-muted">Workspace</p>
            <p className="text-sm font-semibold text-ink">{auth.user?.role || 'Business workspace'}</p>
          </div>
          <button className="rounded-md border border-hairline bg-canvas px-4 py-2 text-sm font-semibold text-ink shadow-soft" onClick={auth.logout}>Logout</button>
        </header>
        <div className="mx-auto max-w-[1280px] px-4 py-8 sm:px-6 lg:px-8">
          <Outlet />
        </div>
      </main>
    </div>
  )
}
