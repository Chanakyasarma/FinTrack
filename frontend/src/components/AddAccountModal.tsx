import { useState, FormEvent } from 'react'
import { X } from 'lucide-react'
import api from '@/lib/api'
import toast from 'react-hot-toast'

interface AddAccountModalProps {
  onClose: () => void
  onSuccess: () => void
}

export default function AddAccountModal({ onClose, onSuccess }: AddAccountModalProps) {
  const [loading, setLoading] = useState(false)
  const [form, setForm] = useState({
    name: '',
    type: 'checking' as 'checking' | 'savings' | 'credit',
    currency: 'USD',
  })

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await api.post('/accounts', form)
      toast.success('Account created!')
      onSuccess()
      onClose()
    } catch (err: unknown) {
      const message =
        (err as { response?: { data?: { error?: string } } })?.response?.data?.error ??
        'Failed to create account'
      toast.error(message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={onClose} />
      <div className="relative w-full max-w-sm card animate-fade-up shadow-2xl shadow-black/50">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-bold text-white">New Account</h2>
          <button onClick={onClose} className="btn-ghost p-2">
            <X className="w-4 h-4" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Account Name</label>
            <input
              className="input"
              placeholder="e.g. Main Checking"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              required
            />
          </div>

          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Account Type</label>
            <div className="flex gap-2">
              {(['checking', 'savings', 'credit'] as const).map((t) => (
                <button
                  key={t}
                  type="button"
                  onClick={() => setForm({ ...form, type: t })}
                  className={`flex-1 py-2 rounded-xl text-sm font-semibold capitalize transition-all ${
                    form.type === t
                      ? 'bg-accent text-white'
                      : 'bg-surface-2 text-gray-400 hover:text-white'
                  }`}
                >
                  {t}
                </button>
              ))}
            </div>
          </div>

          <div>
            <label className="text-xs text-gray-400 font-medium mb-1.5 block">Currency</label>
            <select
              className="input"
              value={form.currency}
              onChange={(e) => setForm({ ...form, currency: e.target.value })}
            >
              <option value="USD">USD — US Dollar</option>
              <option value="EUR">EUR — Euro</option>
              <option value="GBP">GBP — British Pound</option>
              <option value="INR">INR — Indian Rupee</option>
            </select>
          </div>

          <button type="submit" disabled={loading} className="btn-primary w-full">
            {loading ? 'Creating…' : 'Create Account'}
          </button>
        </form>
      </div>
    </div>
  )
}
