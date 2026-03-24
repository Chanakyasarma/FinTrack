import { useState, useEffect, useCallback } from 'react'
import { Account, Transaction, Summary } from '@/types'
import api from '@/lib/api'

export function useAccounts() {
  const [accounts, setAccounts] = useState<Account[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetch = useCallback(async () => {
    try {
      setLoading(true)
      const { data } = await api.get('/accounts')
      setAccounts(data.data ?? [])
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to load accounts')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetch() }, [fetch])

  const updateBalance = useCallback((accountId: string, balance: number) => {
    setAccounts((prev) =>
      prev.map((a) => (a.id === accountId ? { ...a, balance } : a))
    )
  }, [])

  return { accounts, loading, error, refetch: fetch, updateBalance }
}

export function useTransactions(limit = 20) {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetch = useCallback(async () => {
    try {
      setLoading(true)
      const { data } = await api.get(`/transactions?limit=${limit}`)
      setTransactions(data.data ?? [])
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to load transactions')
    } finally {
      setLoading(false)
    }
  }, [limit])

  useEffect(() => { fetch() }, [fetch])

  const prepend = useCallback((tx: Transaction) => {
    setTransactions((prev) => [tx, ...prev.slice(0, limit - 1)])
  }, [limit])

  const updateOne = useCallback((updated: Transaction) => {
    setTransactions((prev) => prev.map((t) => t.id === updated.id ? updated : t))
  }, [])

  const removeOne = useCallback((id: string) => {
    setTransactions((prev) => prev.filter((t) => t.id !== id))
  }, [])

  return { transactions, loading, error, refetch: fetch, prepend, updateOne, removeOne }
}

export function useSummary() {
  const [summary, setSummary] = useState<Summary | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetch = useCallback(async () => {
    try {
      setLoading(true)
      const { data } = await api.get('/transactions/summary')
      setSummary(data.data)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : 'Failed to load summary')
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetch() }, [fetch])

  return { summary, loading, error, refetch: fetch }
}
