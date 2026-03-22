import { describe, it, expect } from 'vitest'
import { formatCurrency, formatPercent, CATEGORY_COLORS, CATEGORY_ICONS } from '../lib/format'

describe('formatCurrency', () => {
  it('formats positive USD amounts', () => {
    expect(formatCurrency(1234.56)).toBe('$1,234.56')
  })

  it('formats zero', () => {
    expect(formatCurrency(0)).toBe('$0.00')
  })

  it('formats negative amounts (credit card balance)', () => {
    expect(formatCurrency(-500)).toBe('-$500.00')
  })

  it('respects currency argument', () => {
    const result = formatCurrency(100, 'EUR')
    expect(result).toContain('100')
    // EUR symbol or code depends on locale — just ensure it doesn't crash
  })

  it('rounds to 2 decimal places', () => {
    expect(formatCurrency(9.999)).toBe('$10.00')
  })
})

describe('formatPercent', () => {
  it('formats percentage with 1 decimal', () => {
    expect(formatPercent(42.567)).toBe('42.6%')
  })

  it('formats zero percent', () => {
    expect(formatPercent(0)).toBe('0.0%')
  })

  it('formats 100%', () => {
    expect(formatPercent(100)).toBe('100.0%')
  })
})

describe('CATEGORY_COLORS', () => {
  it('has a color for every expected category', () => {
    const expectedCategories = ['food', 'transport', 'shopping', 'entertainment', 'health', 'salary', 'other']
    for (const cat of expectedCategories) {
      expect(CATEGORY_COLORS[cat]).toBeDefined()
      expect(CATEGORY_COLORS[cat]).toMatch(/^#[0-9a-f]{6}$/i)
    }
  })
})

describe('CATEGORY_ICONS', () => {
  it('has an icon for every category', () => {
    const expectedCategories = ['food', 'transport', 'shopping', 'entertainment', 'health', 'salary', 'other']
    for (const cat of expectedCategories) {
      expect(CATEGORY_ICONS[cat]).toBeDefined()
      expect(typeof CATEGORY_ICONS[cat]).toBe('string')
    }
  })
})
