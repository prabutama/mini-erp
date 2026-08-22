import { createFileRoute } from '@tanstack/react-router'
import { WorkflowsPage } from '~/features/operations/workflows-page'

export const Route = createFileRoute('/app/workflows')({
  component: WorkflowsPage,
})
