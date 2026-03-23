import { useState } from 'react'
import { useAccounts } from '@/hooks/useFinance'
import AccountCard from '@/components/AccountCard'
import AddAccountModal from '@/components/AddAccountModal'
import AddTransactionModal from '@/components/AddTransactionModal'
import { useTransactions } from '@/hooks/useFinance'
import TransactionItem from '@/components/TransactionItem'
import { Account } from '@/types'
import { formatCurrency } from '@/lib/format'
import { Plus, Wallet } from 'lucide-react'

export default function AccountsPage() {
  const { accounts, loading, refetch } = useAccounts()
  const { transactions } = useTransactions(20)
  const [selected, setSelected] = useState<Account | null>(null)
  const [showAddAccount, setShowAddAccount] = useState(false)
  const [showAddTx, setShowAddTx] = useState(false)

  const filteredTx = selected
    ? transactions.filter((t) => t.account_id === selected.id)
    : transactions

  const totalBalance = accounts.reduce((sum, a) => sum + a.balance, 0)

  return (
    <div className="space-y-6 animate-fade-up pt-2 md:pt-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">Accounts</h1>
          <p className="text-sm text-gray-400 mt-0.5">
            Total across all accounts:{' '}
            <span className="text-white font-semibold font-mono">
              {formatCurrency(totalBalance)}
            </span>
          </p>
        </div>
        <div className="flex gap-3">
          <button onClick={() => setShowAddTx(true)} className="btn-ghost text-sm">
            <Plus className="w-4 h-4 inline mr-1" />
            Transaction
          </button>
          <button onClick={() => setShowAddAccount(true)} className="btn-primary text-sm">
            <Plus className="w-4 h-4 inline mr-1" />
            New Account
          </button>
        </div>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          {Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-36 bg-surface-1 rounded-2xl animate-pulse border border-surface-3" />
          ))}
        </div>
      ) : accounts.length === 0 ? (
        <div className="card text-center py-12">
          <Wallet className="w-10 h-10 text-gray-600 mx-auto mb-3" />
          <p className="text-gray-400 font-medium">No accounts yet</p>
          <p className="text-gray-600 text-sm mt-1">Create your first account to get started</p>
          <button onClick={() => setShowAddAccount(true)} className="btn-primary text-sm mt-4">
            <Plus className="w-4 h-4 inline mr-1" />
            Create Account
          </button>
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          {accounts.map((a) => (
            <AccountCard
              key={a.id}
              account={a}
              isSelected={selected?.id === a.id}
              onClick={() => setSelected(selected?.id === a.id ? null : a)}
            />
          ))}
        </div>
      )}

      {/* Transactions for selected account */}
      <div className="card">
        <h2 className="text-sm font-semibold text-white mb-4">
          {selected ? `Transactions — ${selected.name}` : 'All Recent Transactions'}
        </h2>
        {filteredTx.length === 0 ? (
          <p className="text-center text-gray-500 text-sm py-8">No transactions found.</p>
        ) : (
          <div className="divide-y divide-surface-2">
            {filteredTx.map((tx) => (
              <TransactionItem key={tx.id} tx={tx} />
            ))}
          </div>
        )}
      </div>

      {showAddAccount && (
        <AddAccountModal onClose={() => setShowAddAccount(false)} onSuccess={refetch} />
      )}
      {showAddTx && (
        <AddTransactionModal
          accounts={accounts}
          onClose={() => setShowAddTx(false)}
          onSuccess={refetch}
        />
      )}
    </div>
  )
}
