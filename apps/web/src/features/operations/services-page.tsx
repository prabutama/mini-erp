import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { FormEvent, useMemo } from 'react'
import { Badge } from '~/components/ui/badge'
import { Button } from '~/components/ui/button'
import { DataTable, type ColumnDef } from '~/components/ui/data-table'
import { Panel } from '~/components/ui/card'
import { FormField, Input } from '~/components/ui/input'
import { useAuth } from '~/features/auth/auth-context'
import { operationsApi, type ServiceDefinition } from './operations-api'

export function ServicesPage() {
  const token = useAuth().token || ''
  const queryClient = useQueryClient()
  const services = useQuery({ queryKey: ['service-definitions'], queryFn: () => operationsApi.services(token), enabled: Boolean(token) })
  const createService = useMutation({ mutationFn: (body: unknown) => operationsApi.createService(token, body), onSuccess: async () => queryClient.invalidateQueries({ queryKey: ['service-definitions'] }) })
  const columns = useMemo<ColumnDef<ServiceDefinition>[]>(() => [{ header: 'Name', accessorKey: 'name' }, { header: 'Code', accessorKey: 'code' }, { header: 'Description', accessorKey: 'description' }, { header: 'Status', cell: ({ row }) => <Badge>{row.original.status}</Badge> }], [])
  async function submit(event: FormEvent<HTMLFormElement>) { event.preventDefault(); const form = new FormData(event.currentTarget); await createService.mutateAsync({ name: String(form.get('name') || ''), description: String(form.get('description') || '') }); event.currentTarget.reset() }
  return <div className="grid gap-8 xl:grid-cols-[1fr_380px]"><section><h1 className="display-md">Service definitions</h1><p className="mt-3 text-body">Business-scoped catalog of services offered.</p><div className="mt-8"><DataTable data={services.data?.service_definitions || []} columns={columns} empty="No service definitions yet." /></div></section><Panel><h2 className="text-lg font-semibold">Create service</h2><form className="mt-5 grid gap-4" onSubmit={submit}><FormField label="Name"><Input name="name" required /></FormField><FormField label="Description"><Input name="description" /></FormField><Button>Create service</Button></form></Panel></div>
}
