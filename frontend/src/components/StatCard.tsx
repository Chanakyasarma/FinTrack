import { ReactNode } from 'react'
import { TrendingUp, TrendingDown } from 'lucide-react'

interface StatCardProps {
  title: string
  value: string
  subtitle?: string
  trend?: number // positive = up, negative = down
  icon: ReactNode
  accentColor?: string
  loading?: boolean
}

export default function StatCard({
  title,
  value,
  subtitle,
  trend,
  icon,
  accentColor = '#7c6af7',
  loading = false,
}: StatCardProps) {
  return (
    <div className="stat-card">
      <div className="flex items-start justify-between">
        <div
          className="w-10 h-10 rounded-xl flex items-center justify-center"
          style={{ background: `${accentColor}20` }}
        >
          <span style={{ color: accentColor }}>{icon}</span>
        </div>

        {trend !== undefined && (
          <div
            className={`flex items-center gap-1 text-xs font-semibold px-2 py-1 rounded-full ${
              trend >= 0
                ? 'text-success bg-success/10'
                : 'text-danger bg-danger/10'
            }`}
          >
            {trend >= 0 ? (
              <TrendingUp className="w-3 h-3" />
            ) : (
              <TrendingDown className="w-3 h-3" />
            )}
            {Math.abs(trend).toFixed(1)}%
          </div>
        )}
      </div>

      <div>
        <p className="text-sm text-gray-400 font-medium">{title}</p>
        {loading ? (
          <div className="h-7 w-32 bg-surface-3 rounded-lg animate-pulse mt-1" />
        ) : (
          <p className="text-2xl font-bold text-white mt-1 font-mono">{value}</p>
        )}
        {subtitle && (
          <p className="text-xs text-gray-500 mt-0.5">{subtitle}</p>
        )}
      </div>
    </div>
  )
}
