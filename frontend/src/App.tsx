import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { Toaster } from 'react-hot-toast'
import { AuthProvider, useAuth } from '@/context/AuthContext'
import AppLayout from '@/layouts/AppLayout'
import AuthPage from '@/pages/AuthPage'
import Dashboard from '@/pages/Dashboard'
import TransactionsPage from '@/pages/TransactionsPage'
import AccountsPage from '@/pages/AccountsPage'
import AnalyticsPage from '@/pages/AnalyticsPage'

function PublicRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuth()
  return isAuthenticated ? <Navigate to="/" replace /> : <>{children}</>
}

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route
            path="/login"
            element={
              <PublicRoute>
                <AuthPage />
              </PublicRoute>
            }
          />
          <Route element={<AppLayout />}>
            <Route path="/" element={<Dashboard />} />
            <Route path="/transactions" element={<TransactionsPage />} />
            <Route path="/accounts" element={<AccountsPage />} />
            <Route path="/analytics" element={<AnalyticsPage />} />
          </Route>
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
      <Toaster
        position="top-right"
        toastOptions={{
          style: {
            background: '#16161f',
            color: '#e2e2f0',
            border: '1px solid #252534',
            borderRadius: '12px',
            fontSize: '13px',
            fontFamily: 'Sora, sans-serif',
          },
          success: { iconTheme: { primary: '#2dd4a0', secondary: '#16161f' } },
          error: { iconTheme: { primary: '#f2605c', secondary: '#16161f' } },
        }}
      />
    </AuthProvider>
  )
}
