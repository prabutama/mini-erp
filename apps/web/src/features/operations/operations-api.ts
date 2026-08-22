import { apiRequest } from '~/lib/api/client'
import { endpoints } from '~/lib/api/endpoints'

export type Workflow = { workflow_id: string; name: string; description: string; status: string }
export type ServiceDefinition = { service_definition_id: string; name: string; code: string; description: string; status: string }
export type ServiceOrder = { service_order_id: string; branch_id: string; service_definition_id: string; title: string; description: string; status: string; priority: string }

export const operationsApi = {
  workflows: (token: string) => apiRequest<{ workflows: Workflow[] }>(endpoints.workflows, { token }),
  createWorkflow: (token: string, body: unknown) => apiRequest<Workflow>(endpoints.workflows, { method: 'POST', token, body }),
  createStatus: (token: string, workflowId: string, body: unknown) => apiRequest(endpoints.workflowStatuses(workflowId), { method: 'POST', token, body }),
  createTransition: (token: string, workflowId: string, body: unknown) => apiRequest(endpoints.workflowTransitions(workflowId), { method: 'POST', token, body }),
  services: (token: string) => apiRequest<{ service_definitions: ServiceDefinition[] }>(endpoints.serviceDefinitions, { token }),
  createService: (token: string, body: unknown) => apiRequest<ServiceDefinition>(endpoints.serviceDefinitions, { method: 'POST', token, body }),
  orders: (token: string) => apiRequest<{ service_orders: ServiceOrder[] }>(endpoints.serviceOrders, { token }),
  createOrder: (token: string, body: unknown) => apiRequest<ServiceOrder>(endpoints.serviceOrders, { method: 'POST', token, body }),
  assignOrder: (token: string, orderId: string, assignedUserId: string) => apiRequest(endpoints.assignServiceOrder(orderId), { method: 'POST', token, body: { assigned_user_id: assignedUserId } }),
  transitionOrder: (token: string, orderId: string, status: string) => apiRequest(endpoints.transitionServiceOrder(orderId), { method: 'POST', token, body: { status } }),
}
