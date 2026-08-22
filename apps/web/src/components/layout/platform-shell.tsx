import { Link, Navigate, Outlet, useLocation } from '@tanstack/react-router'
import { useAuth } from '~/features/auth/auth-context'

export function PlatformShell() {
  const auth = useAuth()
  const location = useLocation()

  if (!auth.isAuthenticated) {
    return <Navigate to="/login" search={{ redirect: location.href }} />
  }

  if (auth.user?.role && auth.user.role !== 'Platform Admin') {
    return <Navigate to="/app" />
  }

  return (
    <div className="min-h-screen bg-canvas">
      <header className="border-b border-hairline bg-canvas">
        <div className="mx-auto flex h-16 max-w-content items-center justify-between px-4 sm:px-6 lg:px-8">
          <div className="text-lg font-semibold tracking-[-0.03em]">Mini ERP Platform</div>
          <nav className="flex items-center gap-5 text-sm font-medium">
            <Link to="/platform/businesses">Businesses</Link>
            <button onClick={auth.logout}>Logout</button>
          </nav>
        </div>
      </header>
      <main className="mx-auto max-w-content px-4 py-8 sm:px-6 lg:px-8">
        <Outlet />
      </main>
    </div>
  )
}
