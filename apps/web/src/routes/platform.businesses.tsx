import { createFileRoute } from '@tanstack/react-router'
import { PlatformBusinessesPage } from '~/features/platform/platform-businesses-page'

export const Route = createFileRoute('/platform/businesses')({
  component: PlatformBusinessesPage,
})
