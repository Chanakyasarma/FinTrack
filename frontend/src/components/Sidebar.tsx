import { NavLink } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import {
  LayoutDashboard,
  ArrowLeftRight,
  CreditCard,
  TrendingUp,
  LogOut,
  Settings,
} from 'lucide-react'

const NAV_ITEMS = [
  { to: '/', icon: LayoutDashboard, label: 'Dashboard' },
  { to: '/transactions', icon: ArrowLeftRight, label: 'Transactions' },
  { to: '/accounts', icon: CreditCard, label: 'Accounts' },
  { to: '/analytics', icon: TrendingUp, label: 'Analytics' },
]

export default function Sidebar() {
  const { user, logout } = useAuth()

  return (
    <aside className="w-64 h-screen bg-surface-1 border-r border-surface-3 flex flex-col fixed left-0 top-0">
      {/* Logo */}
      <div className="flex items-center gap-3 px-6 py-6 border-b border-surface-3">
        <div className="w-8 h-8 bg-accent rounded-lg flex items-center justify-center">
          <TrendingUp className="w-4 h-4 text-white" />
        </div>
        <span className="text-lg font-bold text-white tracking-tight">FinTrack</span>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-4 space-y-1">
        {NAV_ITEMS.map(({ to, icon: Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            end={to === '/'}
            className={({ isActive }) =>
              `flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all duration-150 ${
                isActive
                  ? 'bg-accent/15 text-accent'
                  : 'text-gray-400 hover:text-white hover:bg-surface-2'
              }`
            }
          >
            <Icon className="w-4 h-4" />
            {label}
          </NavLink>
        ))}
      </nav>

      {/* Bottom user section */}
      <div className="px-3 py-4 border-t border-surface-3 space-y-1">
        <button className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-gray-400 hover:text-white hover:bg-surface-2 w-full transition-all duration-150">
          <Settings className="w-4 h-4" />
          Settings
        </button>
        <button
          onClick={logout}
          className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium text-gray-400 hover:text-danger hover:bg-danger/10 w-full transition-all duration-150"
        >
          <LogOut className="w-4 h-4" />
          Sign Out
        </button>

        {/* User avatar */}
        <div className="flex items-center gap-3 px-3 py-3 mt-2 bg-surface-2 rounded-xl">
          <div className="w-8 h-8 bg-accent/20 rounded-full flex items-center justify-center">
            <span className="text-xs font-bold text-accent">
              {user?.full_name?.charAt(0).toUpperCase() ?? 'U'}
            </span>
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-xs font-semibold text-white truncate">{user?.full_name}</p>
            <p className="text-xs text-gray-500 truncate">{user?.email}</p>
          </div>
        </div>
      </div>
    </aside>
  )
}
