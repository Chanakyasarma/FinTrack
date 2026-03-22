import { Account } from '@/types'
import { formatCurrency } from '@/lib/format'
import { ACCOUNT_TYPE_COLORS } from '@/lib/format'
import { CreditCard, PiggyBank, Wallet } from 'lucide-react'

interface AccountCardProps {
  account: Account
  isSelected?: boolean
  onClick?: () => void
}

const TYPE_ICONS = {
  checking: Wallet,
  savings: PiggyBank,
  credit: CreditCard,
}

export default function AccountCard({ account, isSelected, onClick }: AccountCardProps) {
  const color = ACCOUNT_TYPE_COLORS[account.type] ?? '#7c6af7'
  const Icon = TYPE_ICONS[account.type] ?? Wallet

  return (
    <button
      onClick={onClick}
      className={`w-full text-left p-5 rounded-2xl border transition-all duration-200 ${
        isSelected
          ? 'border-accent bg-accent/10 shadow-lg shadow-accent/10'
          : 'border-surface-3 bg-surface-1 hover:border-surface-3 hover:bg-surface-2'
      }`}
    >
      <div className="flex items-start justify-between mb-4">
        <div
          className="w-9 h-9 rounded-lg flex items-center justify-center"
          style={{ background: `${color}20` }}
        >
          <Icon className="w-4 h-4" style={{ color }} />
        </div>
        <span
          className="text-xs font-semibold px-2 py-1 rounded-full capitalize"
          style={{ color, background: `${color}15` }}
        >
          {account.type}
        </span>
      </div>

      <p className="text-xs text-gray-500 font-medium mb-1">{account.name}</p>
      <p className="text-xl font-bold text-white font-mono">
        {formatCurrency(account.balance, account.currency)}
      </p>
      <p className="text-xs text-gray-600 mt-1">{account.currency}</p>
    </button>
  )
}
