import { Transaction } from '@/types'
import { formatCurrency, formatDate, CATEGORY_COLORS, CATEGORY_ICONS } from '@/lib/format'
import { ArrowDownLeft, ArrowUpRight } from 'lucide-react'

interface TransactionItemProps {
  tx: Transaction
  isNew?: boolean
}

export default function TransactionItem({ tx, isNew = false }: TransactionItemProps) {
  const isCredit = tx.type === 'credit'
  const color = CATEGORY_COLORS[tx.category] ?? '#9e9eb8'
  const icon = CATEGORY_ICONS[tx.category] ?? '📦'

  return (
    <div
      className={`flex items-center gap-4 p-4 rounded-xl transition-all duration-300 hover:bg-surface-2 ${
        isNew ? 'animate-fade-up bg-accent/5 border border-accent/20' : ''
      }`}
    >
      {/* Category icon */}
      <div
        className="w-10 h-10 rounded-xl flex items-center justify-center text-lg flex-shrink-0"
        style={{ background: `${color}15` }}
      >
        {icon}
      </div>

      {/* Details */}
      <div className="flex-1 min-w-0">
        <p className="text-sm font-semibold text-white truncate">
          {tx.merchant_name || tx.description || 'Transaction'}
        </p>
        <div className="flex items-center gap-2 mt-0.5">
          <span
            className="text-xs px-1.5 py-0.5 rounded-md font-medium capitalize"
            style={{ color, background: `${color}15` }}
          >
            {tx.category}
          </span>
          <span className="text-xs text-gray-500">{formatDate(tx.created_at)}</span>
        </div>
      </div>

      {/* Amount */}
      <div className="flex items-center gap-2 flex-shrink-0">
        <span
          className={`text-sm font-bold font-mono ${
            isCredit ? 'text-success' : 'text-white'
          }`}
        >
          {isCredit ? '+' : '-'}
          {formatCurrency(tx.amount)}
        </span>
        <div
          className={`w-6 h-6 rounded-full flex items-center justify-center ${
            isCredit ? 'bg-success/10' : 'bg-surface-3'
          }`}
        >
          {isCredit ? (
            <ArrowDownLeft className="w-3 h-3 text-success" />
          ) : (
            <ArrowUpRight className="w-3 h-3 text-gray-400" />
          )}
        </div>
      </div>
    </div>
  )
}
