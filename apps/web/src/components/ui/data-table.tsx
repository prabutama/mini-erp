import type { ReactNode } from 'react'

export type ColumnDef<T> = {
  header: string
  accessorKey?: keyof T
  cell?: (ctx: { row: { original: T } }) => ReactNode
}

export function DataTable<T>({ data, columns, empty }: { data: T[]; columns: ColumnDef<T>[]; empty: string }) {
  if (data.length === 0) {
    return <div className="rounded-lg border border-dashed border-hairline bg-surface-soft p-8 text-center text-body">{empty}</div>
  }

  return (
    <div className="overflow-x-auto rounded-lg border border-hairline">
      <table className="min-w-full divide-y divide-hairline text-sm">
        <thead className="bg-surface-soft text-left text-muted">
          <tr>
            {columns.map((column) => (
              <th key={column.header} className="px-4 py-3 font-medium">{column.header}</th>
            ))}
          </tr>
        </thead>
        <tbody className="divide-y divide-hairline bg-canvas text-ink">
          {data.map((item, index) => (
            <tr key={index}>
              {columns.map((column) => (
                <td key={column.header} className="px-4 py-3 align-top">
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
