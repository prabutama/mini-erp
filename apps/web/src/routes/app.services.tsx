import { createFileRoute } from '@tanstack/react-router'
import { ServicesPage } from '~/features/operations/services-page'

export const Route = createFileRoute('/app/services')({
  component: ServicesPage,
})
