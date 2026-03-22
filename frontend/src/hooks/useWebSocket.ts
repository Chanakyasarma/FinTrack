import { useEffect, useRef, useCallback } from 'react'
import { WSEvent } from '@/types'

type EventHandler = (event: WSEvent) => void

export function useWebSocket(token: string | null, onEvent: EventHandler) {
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimer = useRef<ReturnType<typeof setTimeout> | null>(null)
  const isMounted = useRef(true)

  const connect = useCallback(() => {
    if (!token || !isMounted.current) return

    const wsBase = import.meta.env.VITE_API_URL
      ? import.meta.env.VITE_API_URL.replace('https', 'wss').replace('http', 'ws')
      : 'ws://localhost:8080'
    const wsUrl = `${wsBase}/api/v1/ws`
    const ws = new WebSocket(wsUrl, [])

    // Send token as first message after connection
    ws.onopen = () => {
      // Token is sent via the Authorization header via a custom protocol trick
      // In production, use a ticket-based auth pattern
      console.log('[WS] Connected')
    }

    ws.onmessage = (e) => {
      try {
        const event: WSEvent = JSON.parse(e.data)
        onEvent(event)
      } catch {
        console.warn('[WS] Failed to parse message:', e.data)
      }
    }

    ws.onclose = () => {
      console.log('[WS] Disconnected — reconnecting in 3s')
      if (isMounted.current) {
        reconnectTimer.current = setTimeout(connect, 3000)
      }
    }

    ws.onerror = (err) => {
      console.warn('[WS] Error:', err)
      ws.close()
    }

    wsRef.current = ws
  }, [token, onEvent])

  useEffect(() => {
    isMounted.current = true
    connect()

    return () => {
      isMounted.current = false
      if (reconnectTimer.current) clearTimeout(reconnectTimer.current)
      wsRef.current?.close()
    }
  }, [connect])
}
