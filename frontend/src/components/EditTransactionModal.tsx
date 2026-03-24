import { useState, FormEvent } from 'react'
import { X } from 'lucide-react'
import { Transaction, TransactionCategory } from '@/types'
import api from '@/lib/api'
import toast from 'react-hot-toast'

interface EditTransactionModalProps {
  transaction: Transaction
  onClose: () => void
  onSuccess: (updated: Transaction) => void
}

const CATEGORIES: TransactionCategory[] = [
  'food', 'transport', 'shopping', 'entertainment', 'health', 'salary', 'other',
]

export default function EditTransactionModal({
  transaction,
  onClose,
  onSuccess,
}: EditTransactionModalProps) {
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({
    amount: String(transaction.amount),
    type: transaction.type,
    category: transaction.category,
    description: transaction.description,
    merchant_name: transaction.merchant_name,
  })

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      const { data } = await api.put(`/transactions/${transaction.id}`, {
        ...form,
        amount: parseFloat(form.amount),
      })
      toast.success('Transaction updated!')
      onSuccess(data.data)
      onClose()
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
        'Failed to update transaction'
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-md card animate-fade-up shadow-2xl shadow-black/50">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-bold text-white">Edit Transaction</h2>
          <button onClick={onClose} className="btn-ghost p-2">
            <X className="w-4 h-4" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Type toggle */}
          <div className="flex bg-surface-2 rounded-xl p-1">
            {(['debit', 'credit'] as const).map((t) => (
              <button
                key={t}
                type="button"
                onClick={() => setForm({ ...form, type: t })}
                className={`flex-1 py-2 rounded-lg text-sm font-semibold transition-all duration-200 ${
                  form.type === t
                    ? t === 'credit'
                      ? 'bg-success text-surface shadow-md'
                      : 'bg-danger text-white shadow-md'
                    : 'text-gray-400 hover:text-white'
                }`}
              >
                {t === 'credit' ? '+ Income' : '− Expense'}
              </button>
            ))}
          </div>

          {/* Amount */}
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Amount</label>
            <input
              className="input font-mono"
              type="number"
              min="0.01"
              step="0.01"
              value={form.amount}
              onChange={(e) => setForm({ ...form, amount: e.target.value })}
              required
            />
          </div>

          {/* Category */}
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Category</label>
            <div className="flex flex-wrap gap-2">
              {CATEGORIES.map((cat) => (
                <button
                  key={cat}
                  type="button"
                  onClick={() => setForm({ ...form, category: cat })}
                  className={`text-xs px-3 py-1.5 rounded-full font-medium capitalize transition-all ${
                    form.category === cat
                      ? 'bg-accent text-white'
                      : 'bg-surface-2 text-gray-400 hover:text-white'
                  }`}
                >
                  {cat}
                </button>
              ))}
            </div>
          </div>

          {/* Merchant */}
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Merchant / Source</label>
            <input
              className="input"
              type="text"
              placeholder="e.g. Amazon, Starbucks"
              value={form.merchant_name}
              onChange={(e) => setForm({ ...form, merchant_name: e.target.value })}
            />
          </div>

          {/* Description */}
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Description</label>
            <input
              className="input"
              type="text"
              placeholder="What was this for?"
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
            />
          </div>

          <div className="flex gap-3 pt-1">
            <button
              type="button"
              onClick={onClose}
              className="btn-ghost flex-1"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="btn-primary flex-1"
            >
              {loading ? 'Saving…' : 'Save Changes'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
