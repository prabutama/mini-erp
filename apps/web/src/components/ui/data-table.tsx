import type { ReactNode } from 'react'
import { EmptyState } from './page'

export type ColumnDef<T> = {
  header: string
  accessorKey?: keyof T
  cell?: (ctx: { row: { original: T } }) => ReactNode
}

export function DataTable<T>({ data, columns, empty }: { data: T[]; columns: ColumnDef<T>[]; empty: string }) {
  if (data.length === 0) {
    return <EmptyState title={empty} description="Records will appear here after activity is created." />
  }

  return (
    <div className="overflow-hidden rounded-xl border border-hairline bg-canvas shadow-soft">
      <table className="min-w-full divide-y divide-hairline text-sm">
        <thead className="bg-surface-soft/80 text-left text-muted">
          <tr>
            {columns.map((column) => (
              <th key={column.header} className="px-5 py-3 text-xs font-semibold uppercase tracking-[0.12em]">{column.header}</th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-hairline bg-canvas text-ink">
          {data.map((item, index) => (
            <tr key={index} className="transition-colors hover:bg-surface-soft">
              {columns.map((column) => (
                <td key={column.header} className="px-5 py-4 align-top">
                  {column.cell ? column.cell({ row: { original: item } }) : String(column.accessorKey ? item[column.accessorKey] ?? '' : '')}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
