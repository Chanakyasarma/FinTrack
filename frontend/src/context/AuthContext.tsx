import { createContext, useContext, useState, useCallback, ReactNode } from 'react'
import { User } from '@/types'
import api from '@/lib/api'

interface AuthContextType {
  user: User | null
  token: string | null
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, fullName: string) => Promise<void>
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | null>(null)

function loadUser(): User | null {
  try {
    const raw = localStorage.getItem('fintrack_user')
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(loadUser)
  const [token, setToken] = useState<string | null>(
    () => localStorage.getItem('fintrack_token')
  )

  const persist = (token: string, user: User) => {
    localStorage.setItem('fintrack_token', token)
    localStorage.setItem('fintrack_user', JSON.stringify(user))
    setToken(token)
    setUser(user)
  }

  const login = useCallback(async (email: string, password: string) => {
    const { data } = await api.post('/auth/login', { email, password })
    persist(data.data.token, data.data.user)
  }, [])

  const register = useCallback(async (email: string, password: string, fullName: string) => {
    const { data } = await api.post('/auth/register', {
      email,
      password,
      full_name: fullName,
    })
    persist(data.data.token, data.data.user)
  }, [])

  const logout = useCallback(() => {
    localStorage.removeItem('fintrack_token')
    localStorage.removeItem('fintrack_user')
    setToken(null)
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider
      value={{ user, token, login, register, logout, isAuthenticated: !!token }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error('useAuth must be used within AuthProvider')
  return ctx
}
