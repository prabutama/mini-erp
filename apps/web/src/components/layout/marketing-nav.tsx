import { ButtonLink } from '~/components/ui/button'

export function MarketingNav() {
  return (
    <header className="sticky top-0 z-20 border-b border-hairline-soft bg-canvas/95 backdrop-blur">
      <div className="mx-auto flex h-16 max-w-content items-center justify-between px-4 sm:px-6 lg:px-8">
        <a href="/" className="text-lg font-semibold tracking-[-0.03em] text-ink">Mini ERP</a>
        <nav className="hidden items-center gap-7 text-sm font-medium text-muted md:flex">
          <a href="#features">Features</a>
          <a href="#workflow">Workflow</a>
          <a href="#deployment">Deployment</a>
        </nav>
        <div className="flex items-center gap-2">
          <ButtonLink to="/login" variant="ghost">Login</ButtonLink>
          <ButtonLink to="/signup">Sign up</ButtonLink>
        </div>
      </div>
    </header>
  )
}
