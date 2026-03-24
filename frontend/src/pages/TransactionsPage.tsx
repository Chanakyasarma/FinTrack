import { useState } from 'react'
import { useTransactions, useAccounts } from '@/hooks/useFinance'
import TransactionItem from '@/components/TransactionItem'
import AddTransactionModal from '@/components/AddTransactionModal'
import { Plus, Search, Filter } from 'lucide-react'
import { Transaction, TransactionCategory } from '@/types'

const CATEGORIES: Array<TransactionCategory | 'all'> = [
  'all', 'food', 'transport', 'shopping', 'entertainment', 'health', 'salary', 'other',
]

export default function TransactionsPage() {
  const { transactions, loading, refetch, updateOne, removeOne } = useTransactions(50)
  const { accounts, refetch: refetchAccounts } = useAccounts()
  const [showModal, setShowModal] = useState(false)
  const [search, setSearch] = useState('')
  const [categoryFilter, setCategoryFilter] = useState<TransactionCategory | 'all'>('all')
  const [typeFilter, setTypeFilter] = useState<'all' | 'credit' | 'debit'>('all')

  const filtered = transactions.filter((tx) => {
    const matchSearch =
      search === '' ||
      tx.merchant_name.toLowerCase().includes(search.toLowerCase()) ||
      tx.description.toLowerCase().includes(search.toLowerCase())
    const matchCategory = categoryFilter === 'all' || tx.category === categoryFilter
    const matchType = typeFilter === 'all' || tx.type === typeFilter
    return matchSearch && matchCategory && matchType
  })

  const handleUpdate = (updated: Transaction) => {
    updateOne(updated)
    refetchAccounts()
  }

  const handleDelete = (id: string) => {
    removeOne(id)
    refetchAccounts()
  }

  return (
    <div className="space-y-6 animate-fade-up pt-2 md:pt-6">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold text-white">Transactions</h1>
          <p className="text-sm text-gray-400 mt-0.5">{filtered.length} transactions</p>
        </div>
        <button onClick={() => setShowModal(true)} className="btn-primary text-sm">
          <Plus className="w-4 h-4 inline mr-1" />
          Add Transaction
        </button>
      </div>

      {/* Filters */}
      <div className="card space-y-4">
        <div className="relative">
          <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-500" />
          <input
            className="input pl-10"
            placeholder="Search by merchant or description…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        <div className="flex items-center gap-3 flex-wrap">
          <Filter className="w-4 h-4 text-gray-500 flex-shrink-0" />
          <div className="flex gap-2 flex-wrap">
            {(['all', 'credit', 'debit'] as const).map((t) => (
              <button
                key={t}
                onClick={() => setTypeFilter(t)}
                className={`text-xs px-3 py-1.5 rounded-full font-semibold capitalize transition-all ${
                  typeFilter === t ? 'bg-accent text-white' : 'bg-surface-2 text-gray-400 hover:text-white'
                }`}
              >
                {t === 'all' ? 'All Types' : t === 'credit' ? '+ Income' : '− Expenses'}
              </button>
            ))}
          </div>

          <div className="w-px h-4 bg-surface-3 hidden sm:block" />

          <div className="flex gap-2 flex-wrap">
            {CATEGORIES.map((cat) => (
              <button
                key={cat}
                onClick={() => setCategoryFilter(cat)}
                className={`text-xs px-3 py-1.5 rounded-full font-medium capitalize transition-all ${
                  categoryFilter === cat
                    ? 'bg-accent text-white'
                    : 'bg-surface-2 text-gray-400 hover:text-white'
                }`}
              >
                {cat === 'all' ? 'All Categories' : cat}
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* List */}
      <div className="card">
        {loading ? (
          <div className="space-y-3">
            {Array.from({ length: 8 }).map((_, i) => (
              <div key={i} className="h-16 bg-surface-2 rounded-xl animate-pulse" />
            ))}
          </div>
        ) : filtered.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-500 text-sm">No transactions match your filters.</p>
          </div>
        ) : (
          <div className="divide-y divide-surface-2">
            {filtered.map((tx) => (
              <TransactionItem
                key={tx.id}
                tx={tx}
                onUpdate={handleUpdate}
                onDelete={handleDelete}
              />
            ))}
          </div>
        )}
      </div>

      {showModal && (
        <AddTransactionModal
          accounts={accounts}
          onClose={() => setShowModal(false)}
          onSuccess={() => { refetch(); refetchAccounts() }}
        />
      )}
    </div>
  )
}
