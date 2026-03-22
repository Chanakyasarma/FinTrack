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
      <main className="flex-1 ml-64 p-8 max-w-[1400px]">
        <Outlet />
      </main>
    </div>
  )
}
