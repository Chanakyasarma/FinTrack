import { useSummary } from '@/hooks/useFinance'
import {
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  RadarChart, PolarGrid, PolarAngleAxis, Radar,
  LineChart, Line, Legend,
} from 'recharts'
import { formatCurrency, CATEGORY_COLORS, CATEGORY_ICONS } from '@/lib/format'
import { TrendingUp, TrendingDown, BarChart2 } from 'lucide-react'

const tooltipStyle = {
  contentStyle: {
    background: '#1c1c28',
    border: '1px solid #7c6af7',
    borderRadius: '12px',
    padding: '10px 14px',
  },
  labelStyle: { color: '#ffffff', fontSize: 13, fontWeight: 600, marginBottom: 4 },
  itemStyle: { color: '#e2e2f0', fontSize: 12 },
  cursor: { fill: 'rgba(124,106,247,0.08)' },
}

export default function AnalyticsPage() {
  const { summary, loading } = useSummary()

  if (loading) {
    return (
      <div className="space-y-6 animate-fade-up pt-2 md:pt-6">
        <div className="h-8 w-48 bg-surface-2 rounded-lg animate-pulse" />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {Array.from({ length: 4 }).map((_, i) => (
            <div key={i} className="h-64 bg-surface-1 rounded-2xl animate-pulse border border-surface-3" />
          ))}
        </div>
      </div>
    )
  }

  const trendData = summary?.monthly_trend ?? []
  const categoryData = summary
    ? Object.entries(summary.category_breakdown)
        .map(([name, value]) => ({ name, value, icon: CATEGORY_ICONS[name] ?? '📦' }))
        .sort((a, b) => b.value - a.value)
    : []

  const radarData = categoryData.map((c) => ({
    subject: c.name,
    amount: c.value,
  }))

  const netFlowData = trendData.map((m) => ({
    ...m,
    net_flow_color: m.net_flow >= 0 ? '#2dd4a0' : '#f2605c',
  }))

  return (
    <div className="space-y-6 animate-fade-up pt-2 md:pt-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Analytics</h1>
        <p className="text-sm text-gray-400 mt-0.5">Detailed breakdown of your financial patterns</p>
      </div>

      {/* Income vs Expenses bar chart */}
      <div className="card">
        <div className="flex items-center gap-2 mb-5">
          <BarChart2 className="w-4 h-4 text-accent" />
          <h2 className="text-sm font-semibold text-white">Income vs Expenses — Monthly</h2>
        </div>
        {trendData.length > 0 ? (
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={trendData} margin={{ top: 0, right: 0, left: -20, bottom: 0 }} barGap={4}>
              <CartesianGrid strokeDasharray="3 3" stroke="#252534" vertical={false} />
              <XAxis dataKey="month" tick={{ fill: '#6b7280', fontSize: 11 }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fill: '#6b7280', fontSize: 11 }} axisLine={false} tickLine={false} />
              <Tooltip
                contentStyle={tooltipStyle.contentStyle}
                labelStyle={tooltipStyle.labelStyle}
                itemStyle={tooltipStyle.itemStyle}
                cursor={tooltipStyle.cursor}
                formatter={(v: number) => formatCurrency(v)}
              />
              <Legend
                iconType="circle"
                iconSize={8}
                formatter={(v) => <span style={{ fontSize: 11, color: '#9ca3af' }}>{v}</span>}
              />
              <Bar dataKey="total_credit" name="Income" fill="#2dd4a0" radius={[4, 4, 0, 0]} />
              <Bar dataKey="total_debit" name="Expenses" fill="#7c6af7" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        ) : (
          <EmptyState />
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {/* Net cash flow line */}
        <div className="card">
          <div className="flex items-center gap-2 mb-5">
            <TrendingUp className="w-4 h-4 text-success" />
            <h2 className="text-sm font-semibold text-white">Net Cash Flow</h2>
          </div>
          {netFlowData.length > 0 ? (
            <ResponsiveContainer width="100%" height={200}>
              <LineChart data={netFlowData} margin={{ top: 0, right: 0, left: -20, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#252534" />
                <XAxis dataKey="month" tick={{ fill: '#6b7280', fontSize: 10 }} axisLine={false} tickLine={false} />
                <YAxis tick={{ fill: '#6b7280', fontSize: 10 }} axisLine={false} tickLine={false} />
                <Tooltip
                  contentStyle={tooltipStyle.contentStyle}
                  labelStyle={tooltipStyle.labelStyle}
                  itemStyle={tooltipStyle.itemStyle}
                  formatter={(v: number) => [formatCurrency(v), 'Net Flow']}
                />
                <Line
                  type="monotone"
                  dataKey="net_flow"
                  stroke="#7c6af7"
                  strokeWidth={2.5}
                  dot={{ fill: '#7c6af7', r: 4 }}
                  activeDot={{ r: 6 }}
                />
              </LineChart>
            </ResponsiveContainer>
          ) : (
            <EmptyState />
          )}
        </div>

        {/* Spending radar */}
        <div className="card">
          <div className="flex items-center gap-2 mb-5">
            <TrendingDown className="w-4 h-4 text-danger" />
            <h2 className="text-sm font-semibold text-white">Spending by Category</h2>
          </div>
          {radarData.length > 0 ? (
            <ResponsiveContainer width="100%" height={200}>
              <RadarChart data={radarData}>
                <PolarGrid stroke="#252534" />
                <PolarAngleAxis
                  dataKey="subject"
                  tick={{ fill: '#9ca3af', fontSize: 10 }}
                />
                <Radar
                  name="Spending"
                  dataKey="amount"
                  stroke="#7c6af7"
                  fill="#7c6af7"
                  fillOpacity={0.2}
                  strokeWidth={2}
                />
                <Tooltip
                  contentStyle={tooltipStyle.contentStyle}
                  labelStyle={tooltipStyle.labelStyle}
                  itemStyle={tooltipStyle.itemStyle}
                  formatter={(v: number) => [formatCurrency(v), 'Spent']}
                />
              </RadarChart>
            </ResponsiveContainer>
          ) : (
            <EmptyState />
          )}
        </div>
      </div>

      {/* Category breakdown table */}
      {categoryData.length > 0 && (
        <div className="card">
          <h2 className="text-sm font-semibold text-white mb-4">Category Breakdown</h2>
          <div className="space-y-3">
            {categoryData.map((cat) => {
              const total = categoryData.reduce((s, c) => s + c.value, 0)
              const pct = total > 0 ? (cat.value / total) * 100 : 0
              const color = CATEGORY_COLORS[cat.name] ?? '#9e9eb8'

              return (
                <div key={cat.name} className="flex items-center gap-4">
                  <span className="text-lg w-6">{cat.icon}</span>
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm capitalize text-gray-300">{cat.name}</span>
                      <span className="text-sm font-semibold font-mono text-white">
                        {formatCurrency(cat.value)}
                      </span>
                    </div>
                    <div className="h-1.5 bg-surface-3 rounded-full overflow-hidden">
                      <div
                        className="h-full rounded-full transition-all duration-700"
                        style={{ width: `${pct}%`, background: color }}
                      />
                    </div>
                  </div>
                  <span className="text-xs text-gray-500 w-10 text-right">{pct.toFixed(0)}%</span>
                </div>
              )
            })}
          </div>
        </div>
      )}
    </div>
  )
}

function EmptyState() {
  return (
    <div className="h-[200px] flex items-center justify-center text-gray-600 text-sm">
      Add transactions to see analytics
    </div>
  )
}
