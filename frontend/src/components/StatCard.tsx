import { ReactNode } from 'react'
import { TrendingUp, TrendingDown } from 'lucide-react'

interface StatCardProps {
  title: string
  value: string
  subtitle?: string
  trend?: number
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
          className="w-8 h-8 md:w-10 md:h-10 rounded-xl flex items-center justify-center flex-shrink-0"
          style={{ background: `${accentColor}20` }}
        >
          <span style={{ color: accentColor }}>{icon}</span>
        </div>

        {trend !== undefined && (
          <div
            className={`flex items-center gap-1 text-xs font-semibold px-1.5 py-0.5 rounded-full ${
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

      <div className="min-w-0">
        <p className="text-xs md:text-sm text-gray-400 font-medium truncate">{title}</p>
        {loading ? (
          <div className="h-6 w-24 bg-surface-3 rounded-lg animate-pulse mt-1" />
        ) : (
          <p className="text-base md:text-2xl font-bold text-white mt-1 font-mono truncate leading-tight">
            {value}
          </p>
        )}
        {subtitle && (
          <p className="text-xs text-gray-500 mt-0.5 truncate">{subtitle}</p>
        )}
      </div>
    </div>
  )
}
