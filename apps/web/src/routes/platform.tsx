import { createFileRoute } from '@tanstack/react-router'
import { PlatformShell } from '~/components/layout/platform-shell'

export const Route = createFileRoute('/platform')({
  component: PlatformShell,
})
