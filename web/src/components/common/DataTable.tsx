import React from 'react';

interface Column<T> {
  key: keyof T | string;
  label: string;
  render?: (item: T, index: number) => React.ReactNode;
}

interface DataTableProps<T> {
  data: T[];
  columns: Column<T>[];
  keyExtractor: (item: T, index: number) => string | number;
}

export function DataTable<T>({ data, columns, keyExtractor }: DataTableProps<T>) {
  return (
    <table className="table">
      <thead>
        <tr>
          {columns.map(col => (
            <th key={String(col.key)}>{col.label}</th>
          ))}
        </tr>
      </thead>
      <tbody>
        {data.map((item, index) => (
          <tr key={keyExtractor(item, index)}>
            {columns.map(col => (
              <td key={String(col.key)}>
                {col.render ? col.render(item, index) : String(item[col.key as keyof T] || '')}
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
