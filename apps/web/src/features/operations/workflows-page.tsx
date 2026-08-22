import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent, useMemo } from 'react'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { DataTable, type ColumnDef } from '~/components/ui/data-table'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { Select } from '~/components/ui/select'
import { useAuth } from '~/features/auth/auth-context'
import { operationsApi, type Workflow } from './operations-api'

export function WorkflowsPage() {
  const token = useAuth().token || ''
  const queryClient = useQueryClient()
  const workflows = useQuery({ queryKey: ['workflows'], queryFn: () => operationsApi.workflows(token), enabled: Boolean(token) })
  const createWorkflow = useMutation({ mutationFn: (body: unknown) => operationsApi.createWorkflow(token, body), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['workflows'] }) })
  const createStatus = useMutation({ mutationFn: ({ workflowId, body }: { workflowId: string; body: unknown }) => operationsApi.createStatus(token, workflowId, body) })
  const createTransition = useMutation({ mutationFn: ({ workflowId, body }: { workflowId: string; body: unknown }) => operationsApi.createTransition(token, workflowId, body) })
  const columns = useMemo<ColumnDef<Workflow>[]>(() => [
    { header: 'Name', accessorKey: 'name' },
    { header: 'Description', accessorKey: 'description' },
    { header: 'Status', cell: ({ row }) => <Badge>{row.original.status}</Badge> },
  ], [])
  const list = workflows.data?.workflows || []

  async function submitWorkflow(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await createWorkflow.mutateAsync({ name: String(form.get('name') || ''), description: String(form.get('description') || '') }); event.currentTarget.reset() }
  async function submitStatus(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await createStatus.mutateAsync({ workflowId: String(form.get('workflow_id') || ''), body: { code: String(form.get('code') || ''), name: String(form.get('name') || ''), category: String(form.get('category') || 'active'), sort_order: Number(form.get('sort_order') || 0), is_initial: form.get('is_initial') === 'on', is_terminal: form.get('is_terminal') === 'on' } }); event.currentTarget.reset() }
  async function submitTransition(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await createTransition.mutateAsync({ workflowId: String(form.get('workflow_id') || ''), body: { from_status_code: String(form.get('from_status_code') || ''), to_status_code: String(form.get('to_status_code') || '') } }); event.currentTarget.reset() }

  return (
    <div className="grid gap-8 xl:grid-cols-[1fr_380px]"><section><h1 className="display-md">Workflows</h1><p className="mt-3 text-body">Define statuses and transitions for business workflows.</p><div className="mt-8"><DataTable data={list} columns={columns} empty="No workflows yet." /></div></section><aside className="grid gap-6"><Panel><h2 className="text-lg font-semibold">Create workflow</h2><form className="mt-5 grid gap-4" onSubmit={submitWorkflow}><FormField label="Name"><Input name="name" required /></FormField><FormField label="Description"><Input name="description" /></FormField><Button>Create workflow</Button></form></Panel><Panel><h2 className="text-lg font-semibold">Add status</h2><form className="mt-5 grid gap-4" onSubmit={submitStatus}><FormField label="Workflow"><Select name="workflow_id" required>{list.map((workflow) => <option key={workflow.workflow_id} value={workflow.workflow_id}>{workflow.name}</option>)}</Select></FormField><FormField label="Code"><Input name="code" placeholder="open" required /></FormField><FormField label="Name"><Input name="name" placeholder="Open" required /></FormField><FormField label="Category"><Input name="category" placeholder="active" /></FormField><FormField label="Sort order"><Input name="sort_order" type="number" defaultValue="10" /></FormField><label className="text-sm"><input name="is_initial" type="checkbox" /> Initial</label><label className="text-sm"><input name="is_terminal" type="checkbox" /> Terminal</label><Button>Add status</Button></form></Panel><Panel><h2 className="text-lg font-semibold">Add transition</h2><form className="mt-5 grid gap-4" onSubmit={submitTransition}><FormField label="Workflow"><Select name="workflow_id" required>{list.map((workflow) => <option key={workflow.workflow_id} value={workflow.workflow_id}>{workflow.name}</option>)}</Select></FormField><FormField label="From"><Input name="from_status_code" placeholder="open" required /></FormField><FormField label="To"><Input name="to_status_code" placeholder="in_progress" required /></FormField><Button>Add transition</Button></form></Panel></aside></div>
  )
}
