export interface User {
  id: string
  email: string
  full_name: string
  created_at: string
}

export interface Account {
  id: string
  user_id: string
  name: string
  type: 'checking' | 'savings' | 'credit'
  balance: number
  currency: string
  created_at: string
}

export interface Transaction {
  id: string
  account_id: string
  user_id: string
  amount: number
  type: 'credit' | 'debit'
  category: TransactionCategory
  description: string
  merchant_name: string
  created_at: string
}

export type TransactionCategory =
  | 'food'
  | 'transport'
  | 'shopping'
  | 'entertainment'
  | 'health'
  | 'salary'
  | 'other'

export interface MonthlySummary {
  month: string
  total_credit: number
  total_debit: number
  net_flow: number
}

export interface Summary {
  total_balance: number
  monthly_income: number
  monthly_expenses: number
  savings_rate: number
  category_breakdown: Record<string, number>
  monthly_trend: MonthlySummary[]
}

export interface WSEvent {
  type: 'transaction.created' | 'balance.updated'
  payload: Transaction | { account_id: string; balance: number }
}

export interface ApiResponse<T> {
  data: T
  success: boolean
}
