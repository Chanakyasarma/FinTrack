export function formatCurrency(amount: number, currency = 'USD'): string {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
  }).format(amount)
}

export function formatDate(dateStr: string): string {
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  }).format(new Date(dateStr))
}

export function formatPercent(value: number): string {
  return `${value.toFixed(1)}%`
}

export const CATEGORY_COLORS: Record<string, string> = {
  food: '#f5a623',
  transport: '#7c6af7',
  shopping: '#f2605c',
  entertainment: '#2dd4a0',
  health: '#4fc3f7',
  salary: '#2dd4a0',
  other: '#9e9eb8',
}

export const CATEGORY_ICONS: Record<string, string> = {
  food: '🍜',
  transport: '🚗',
  shopping: '🛍️',
  entertainment: '🎬',
  health: '💊',
  salary: '💰',
  other: '📦',
}

export const ACCOUNT_TYPE_COLORS: Record<string, string> = {
  checking: '#7c6af7',
  savings: '#2dd4a0',
  credit: '#f2605c',
}
