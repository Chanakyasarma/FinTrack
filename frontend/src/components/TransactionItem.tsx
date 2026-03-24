import { useState } from 'react'
import { Transaction } from '@/types'
import { formatCurrency, formatDate, CATEGORY_COLORS, CATEGORY_ICONS } from '@/lib/format'
import { ArrowDownLeft, ArrowUpRight, Pencil, Trash2 } from 'lucide-react'
import api from '@/lib/api'
import toast from 'react-hot-toast'
import EditTransactionModal from './EditTransactionModal'

interface TransactionItemProps {
  tx: Transaction
  isNew?: boolean
  onUpdate?: (updated: Transaction) => void
  onDelete?: (id: string) => void
}

export default function TransactionItem({
  tx,
  isNew = false,
  onUpdate,
  onDelete,
}: TransactionItemProps) {
  const [showEdit, setShowEdit] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const isCredit = tx.type === 'credit'
  const color = CATEGORY_COLORS[tx.category] ?? '#9e9eb8'
  const icon = CATEGORY_ICONS[tx.category] ?? '📦'

  const handleDelete = async () => {
    if (!confirm('Delete this transaction? This will reverse the balance.')) return
    setDeleting(true)
    try {
      await api.delete(`/transactions/${tx.id}`)
      toast.success('Transaction deleted')
      onDelete?.(tx.id)
    } catch {
      toast.error('Failed to delete transaction')
    } finally {
      setDeleting(false)
    }
  }

  return (
    <>
      <div
        className={`group flex items-center gap-4 p-4 rounded-xl transition-all duration-300 hover:bg-surface-2 ${
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

        {/* Amount + Actions */}
        <div className="flex items-center gap-2 flex-shrink-0">
          {/* Edit/Delete — visible on hover */}
          {(onUpdate || onDelete) && (
            <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity duration-150">
              {onUpdate && (
                <button
                  onClick={() => setShowEdit(true)}
                  className="w-7 h-7 rounded-lg flex items-center justify-center text-gray-500 hover:text-accent hover:bg-accent/10 transition-all"
                  title="Edit transaction"
                >
                  <Pencil className="w-3.5 h-3.5" />
                </button>
              )}
              {onDelete && (
                <button
                  onClick={handleDelete}
                  disabled={deleting}
                  className="w-7 h-7 rounded-lg flex items-center justify-center text-gray-500 hover:text-danger hover:bg-danger/10 transition-all disabled:opacity-50"
                  title="Delete transaction"
                >
                  <Trash2 className="w-3.5 h-3.5" />
                </button>
              )}
            </div>
          )}

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

      {showEdit && (
        <EditTransactionModal
          transaction={tx}
          onClose={() => setShowEdit(false)}
          onSuccess={(updated) => {
            onUpdate?.(updated)
            setShowEdit(false)
          }}
        />
      )}
    </>
  )
}
