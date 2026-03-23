import { Outlet, Navigate } from 'react-router-dom'
import { useAuth } from '@/context/AuthContext'
import Sidebar from '@/components/Sidebar'

export default function AppLayout() {
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      {/* Desktop: offset for sidebar. Mobile: offset for top bar + bottom nav */}
      <main className="flex-1 md:ml-64 pt-16 md:pt-0 pb-20 md:pb-0 px-4 md:px-8 py-4 md:py-8 max-w-[1400px] w-full">
        <Outlet />
      </main>
    </div>
  )
}
