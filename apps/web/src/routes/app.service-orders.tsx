import { createFileRoute } from '@tanstack/react-router'
import { OrdersPage } from '~/features/operations/orders-page'

export const Route = createFileRoute('/app/service-orders')({
  component: OrdersPage,
})
