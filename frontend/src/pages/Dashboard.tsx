import { useState, useCallback, useEffect } from 'react'
import { useAuth } from '@/context/AuthContext'
import { useAccounts, useTransactions, useSummary } from '@/hooks/useFinance'
import { useWebSocket } from '@/hooks/useWebSocket'
import { WSEvent, Transaction } from '@/types'
import StatCard from '@/components/StatCard'
import TransactionItem from '@/components/TransactionItem'
import AccountCard from '@/components/AccountCard'
import AddTransactionModal from '@/components/AddTransactionModal'
import AddAccountModal from '@/components/AddAccountModal'
import {
  AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  PieChart, Pie, Cell, Legend,
} from 'recharts'
import {
  DollarSign, TrendingUp, TrendingDown, Percent, Plus, Zap,
} from 'lucide-react'
import { formatCurrency, formatPercent, CATEGORY_COLORS } from '@/lib/format'
import toast from 'react-hot-toast'

const CHART_COLORS = ['#7c6af7', '#2dd4a0', '#f5a623', '#f2605c', '#4fc3f7', '#9e9eb8']

export default function Dashboard() {
  const { user, token } = useAuth()
  const { accounts, loading: acLoading, refetch: refetchAccounts, updateBalance } = useAccounts()
  const { transactions, loading: txLoading, refetch: refetchTx, prepend } = useTransactions(10)
  const { summary, loading: sumLoading, refetch: refetchSummary } = useSummary()

  const [newTxIds, setNewTxIds] = useState<Set<string>>(new Set())
  const [showAddTx, setShowAddTx] = useState(false)
  const [showAddAccount, setShowAddAccount] = useState(false)
  const [greeting, setGreeting] = useState(getGreeting())
  const [currentTime, setCurrentTime] = useState(formatTime())

  // Update greeting and time every minute
  useEffect(() => {
    const timer = setInterval(() => {
      setGreeting(getGreeting())
      setCurrentTime(formatTime())
    }, 60000)
    return () => clearInterval(timer)
  }, [])

  const handleWSEvent = useCallback((event: WSEvent) => {
    if (event.type === 'transaction.created') {
      const tx = event.payload as Transaction
      prepend(tx)
      setNewTxIds((prev) => new Set(prev).add(tx.id))
      setTimeout(() => {
        setNewTxIds((prev) => {
          const next = new Set(prev)
          next.delete(tx.id)
          return next
        })
      }, 3000)
      refetchSummary()
      toast.success(`New transaction: ${formatCurrency(tx.amount)}`, { icon: '⚡' })
    } else if (event.type === 'balance.updated') {
      const { account_id, balance } = event.payload as { account_id: string; balance: number }
      updateBalance(account_id, balance)
    }
  }, [prepend, refetchSummary, updateBalance])

  useWebSocket(token, handleWSEvent)

  const handleTxSuccess = () => {
    refetchTx()
    refetchAccounts()
    refetchSummary()
  }

  const pieData = summary
    ? Object.entries(summary.category_breakdown).map(([name, value]) => ({ name, value }))
    : []

  const trendData = summary?.monthly_trend ?? []

  return (
    <div className="space-y-8 animate-fade-up pt-4 md:pt-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <div>
          <div className="flex items-center gap-2 mb-1">
            <span className="text-2xl">{getGreetingEmoji()}</span>
            <h1 className="text-xl md:text-2xl font-bold text-white">
              Good {greeting}, {user?.full_name?.split(' ')[0]}!
            </h1>
          </div>
          <p className="text-sm text-gray-400">
            {currentTime} · Here's your financial overview
          </p>
        </div>
        <div className="flex items-center gap-3 flex-wrap">
          <div className="flex items-center gap-2 px-3 py-1.5 rounded-full bg-surface-2 border border-surface-3">
            <div className="live-dot" />
            <span className="text-xs text-gray-400 font-medium">Live</span>
            <Zap className="w-3 h-3 text-success" />
          </div>
          <button onClick={() => setShowAddAccount(true)} className="btn-ghost text-sm">
            <Plus className="w-4 h-4 inline mr-1" />
            Account
          </button>
          <button onClick={() => setShowAddTx(true)} className="btn-primary text-sm">
            <Plus className="w-4 h-4 inline mr-1" />
            Transaction
          </button>
        </div>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-4">
        <StatCard
          title="Total Balance"
          value={formatCurrency(summary?.total_balance ?? 0)}
          icon={<DollarSign className="w-5 h-5" />}
          accentColor="#7c6af7"
          loading={sumLoading}
        />
        <StatCard
          title="Monthly Income"
          value={formatCurrency(summary?.monthly_income ?? 0)}
          subtitle="This month"
          icon={<TrendingUp className="w-5 h-5" />}
          accentColor="#2dd4a0"
          loading={sumLoading}
        />
        <StatCard
          title="Monthly Expenses"
          value={formatCurrency(summary?.monthly_expenses ?? 0)}
          subtitle="This month"
          icon={<TrendingDown className="w-5 h-5" />}
          accentColor="#f2605c"
          loading={sumLoading}
        />
        <StatCard
          title="Savings Rate"
          value={formatPercent(summary?.savings_rate ?? 0)}
          subtitle="Income saved"
          icon={<Percent className="w-5 h-5" />}
          accentColor="#f5a623"
          loading={sumLoading}
        />
      </div>

      {/* Accounts row */}
      {accounts.length > 0 && (
        <div>
          <h2 className="text-sm font-semibold text-gray-400 uppercase tracking-wider mb-3">
            Accounts
          </h2>
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
            {acLoading
              ? Array.from({ length: 3 }).map((_, i) => (
                  <div key={i} className="h-32 bg-surface-1 rounded-2xl animate-pulse border border-surface-3" />
                ))
              : accounts.map((a) => <AccountCard key={a.id} account={a} />)}
          </div>
        </div>
      )}

      {/* Charts + Transactions */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Cash flow chart */}
        <div className="lg:col-span-2 card">
          <h2 className="text-sm font-semibold text-white mb-4">Cash Flow — Last 6 Months</h2>
          {trendData.length > 0 ? (
            <ResponsiveContainer width="100%" height={220}>
              <AreaChart data={trendData} margin={{ top: 0, right: 0, left: -20, bottom: 0 }}>
                <defs>
                  <linearGradient id="colorCredit" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#2dd4a0" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#2dd4a0" stopOpacity={0} />
                  </linearGradient>
                  <linearGradient id="colorDebit" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#7c6af7" stopOpacity={0.3} />
                    <stop offset="95%" stopColor="#7c6af7" stopOpacity={0} />
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#252534" />
                <XAxis dataKey="month" tick={{ fill: '#6b7280', fontSize: 11 }} axisLine={false} tickLine={false} />
                <YAxis tick={{ fill: '#6b7280', fontSize: 11 }} axisLine={false} tickLine={false} />
                <Tooltip
                  contentStyle={{ background: '#16161f', border: '1px solid #252534', borderRadius: '12px' }}
                  labelStyle={{ color: '#e2e2f0', fontSize: 12 }}
                  itemStyle={{ fontSize: 12 }}
                />
                <Area type="monotone" dataKey="total_credit" name="Income" stroke="#2dd4a0" fill="url(#colorCredit)" strokeWidth={2} dot={false} />
                <Area type="monotone" dataKey="total_debit" name="Expenses" stroke="#7c6af7" fill="url(#colorDebit)" strokeWidth={2} dot={false} />
              </AreaChart>
            </ResponsiveContainer>
          ) : (
            <EmptyChart />
          )}
        </div>

        {/* Spending breakdown */}
        <div className="card">
          <h2 className="text-sm font-semibold text-white mb-4">Spending Breakdown</h2>
          {pieData.length > 0 ? (
            <ResponsiveContainer width="100%" height={220}>
              <PieChart>
                <Pie
                  data={pieData}
                  cx="50%"
                  cy="50%"
                  innerRadius={55}
                  outerRadius={80}
                  paddingAngle={3}
                  dataKey="value"
                >
                  {pieData.map((entry, index) => (
                    <Cell
                      key={entry.name}
                      fill={CATEGORY_COLORS[entry.name] ?? CHART_COLORS[index % CHART_COLORS.length]}
                    />
                  ))}
                </Pie>
                <Tooltip
                  contentStyle={{ background: '#16161f', border: '1px solid #252534', borderRadius: '12px' }}
                  formatter={(v: number) => formatCurrency(v)}
                  itemStyle={{ fontSize: 12 }}
                />
                <Legend
                  iconType="circle"
                  iconSize={8}
                  formatter={(value) => (
                    <span style={{ fontSize: 11, color: '#9ca3af', textTransform: 'capitalize' }}>{value}</span>
                  )}
                />
              </PieChart>
            </ResponsiveContainer>
          ) : (
            <EmptyChart />
          )}
        </div>
      </div>

      {/* Recent transactions */}
      <div className="card">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-sm font-semibold text-white">Recent Transactions</h2>
          <span className="text-xs text-gray-500">Updates in real-time ⚡</span>
        </div>

        {txLoading ? (
          <div className="space-y-3">
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-16 bg-surface-2 rounded-xl animate-pulse" />
            ))}
          </div>
        ) : transactions.length === 0 ? (
          <p className="text-center text-gray-500 py-8 text-sm">
            No transactions yet. Add one to get started!
          </p>
        ) : (
          <div className="divide-y divide-surface-2">
            {transactions.map((tx) => (
              <TransactionItem key={tx.id} tx={tx} isNew={newTxIds.has(tx.id)} />
            ))}
          </div>
        )}
      </div>

      {/* Modals */}
      {showAddTx && (
        <AddTransactionModal
          accounts={accounts}
          onClose={() => setShowAddTx(false)}
          onSuccess={handleTxSuccess}
        />
      )}
      {showAddAccount && (
        <AddAccountModal
          onClose={() => setShowAddAccount(false)}
          onSuccess={refetchAccounts}
        />
      )}
    </div>
  )
}

function getGreeting() {
  const h = new Date().getHours()
  if (h >= 5 && h < 12) return 'morning'
  if (h >= 12 && h < 17) return 'afternoon'
  if (h >= 17 && h < 21) return 'evening'
  return 'night'
}

function getGreetingEmoji() {
  const h = new Date().getHours()
  if (h >= 5 && h < 12) return '🌅'
  if (h >= 12 && h < 17) return '☀️'
  if (h >= 17 && h < 21) return '🌆'
  return '🌙'
}

function formatTime() {
  return new Date().toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
  })
}

function EmptyChart() {
  return (
    <div className="h-[220px] flex items-center justify-center text-gray-600 text-sm">
      No data yet — add some transactions!
    </div>
  )
}
