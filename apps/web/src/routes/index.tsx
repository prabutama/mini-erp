import { createFileRoute } from '@tanstack/react-router'
import { MarketingNav } from '~/components/layout/marketing-nav'
import { OperationsMockup } from '~/components/product-mockups/operations-mockup'
import { ButtonLink } from '~/components/ui/button'
import { Card } from '~/components/ui/card'

export const Route = createFileRoute('/')({
  component: LandingPage,
})

function LandingPage() {
  return (
    <div className="min-h-screen bg-canvas text-ink">
      <MarketingNav />
      <main>
        <section className="mx-auto grid max-w-content gap-10 px-4 py-20 sm:px-6 lg:grid-cols-[1.1fr_0.9fr] lg:px-8 lg:py-24">
          <div className="flex flex-col justify-center">
            <p className="mb-5 w-fit rounded-full bg-surface-card px-4 py-1.5 text-sm font-medium">Multi-tenant service ERP</p>
            <h1 className="display-xl max-w-3xl">Run branch operations without spreadsheet drift.</h1>
            <p className="mt-6 max-w-2xl text-lg leading-8 text-body">
              Mini ERP connects signup, branch access, workflows, service orders, resources, and reports through one API Gateway.
            </p>
            <div className="mt-8 flex flex-col gap-3 sm:flex-row">
              <ButtonLink to="/signup">Start tenant signup</ButtonLink>
              <ButtonLink to="/login" variant="secondary">Login</ButtonLink>
            </div>
          </div>
          <OperationsMockup />
        </section>

        <section id="features" className="border-t border-hairline-soft py-20">
          <div className="mx-auto max-w-content px-4 sm:px-6 lg:px-8">
            <h2 className="display-lg max-w-2xl">MVP flow already mapped to backend contracts.</h2>
            <div className="mt-10 grid gap-6 md:grid-cols-3">
              {[
                ['Tenant setup', 'Public signup creates one business and first Business Admin.'],
                ['Branch operations', 'Managers and Staff stay scoped to assigned branches.'],
                ['Work orders', 'Service orders move through fixed MVP transitions and assignments.'],
              ].map(([title, description]) => (
                <Card key={title}>
                  <h3 className="text-lg font-semibold">{title}</h3>
                  <p className="mt-3 text-body">{description}</p>
                </Card>
              ))}
            </div>
          </div>
        </section>
      </main>
      <footer className="bg-surface-dark px-4 py-16 text-on-dark-soft sm:px-6 lg:px-8">
        <div className="mx-auto max-w-content">
          <p className="text-lg font-semibold text-white">Mini ERP</p>
          <p className="mt-4 max-w-xl text-sm">Cal.com-like interface direction: white canvas, black CTAs, light cards, and real product UI fragments.</p>
        </div>
      </footer>
    </div>
  )
}
