import { useState, FormEvent } from 'react'
import { X } from 'lucide-react'
import { Account, TransactionCategory } from '@/types'
import api from '@/lib/api'
import toast from 'react-hot-toast'

interface AddTransactionModalProps {
  accounts: Account[]
  onClose: () => void
  onSuccess: () => void
}

const CATEGORIES: TransactionCategory[] = [
  'food', 'transport', 'shopping', 'entertainment', 'health', 'salary', 'other',
]

export default function AddTransactionModal({
  accounts,
  onClose,
  onSuccess,
}: AddTransactionModalProps) {
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({
    account_id: accounts[0]?.id ?? '',
    amount: '',
    type: 'debit' as 'credit' | 'debit',
    category: 'other' as TransactionCategory,
    description: '',
    merchant_name: '',
  })

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    if (!form.account_id) {
      toast.error('Please select an account')
      return
    }

    setLoading(true)
    try {
      await api.post('/transactions', {
        ...form,
        amount: parseFloat(form.amount),
      })
      toast.success('Transaction added!')
      onSuccess()
      onClose()
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
        'Failed to add transaction'
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div
        className="absolute inset-0 bg-black/60 backdrop-blur-sm"
        onClick={onClose}
      />
      <div className="relative w-full max-w-md card animate-fade-up shadow-2xl shadow-black/50">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-bold text-white">New Transaction</h2>
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
                className={`flex-1 py-2 rounded-lg text-sm font-semibold transition-all duration-200 capitalize ${
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

          {/* Account */}
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Account</label>
            <select
              className="input"
              value={form.account_id}
              onChange={(e) => setForm({ ...form, account_id: e.target.value })}
              required
            >
              {accounts.map((a) => (
                <option key={a.id} value={a.id}>
                  {a.name}
                </option>
              ))}
            </select>
          </div>

          {/* Amount */}
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Amount (USD)</label>
            <input
              className="input font-mono"
              type="number"
              min="0.01"
              step="0.01"
              placeholder="0.00"
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
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">
              Merchant / Source
            </label>
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
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">
              Description (optional)
            </label>
            <input
              className="input"
              type="text"
              placeholder="What was this for?"
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
            />
          </div>

          <button type="submit" disabled={loading} className="btn-primary w-full">
            {loading ? 'Adding…' : 'Add Transaction'}
          </button>
        </form>
      </div>
    </div>
  )
}
