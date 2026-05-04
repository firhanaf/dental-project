import { useEffect, useRef, useState, useCallback } from 'react'

const ACTIVITY_EVENTS = [
  'mousemove', 'mousedown', 'keydown', 'touchstart', 'scroll', 'click',
] as const

interface Options {
  idleMinutes?: number
  warningSeconds?: number
  onLogout: () => void
}

export function useIdleTimeout({
  idleMinutes = 30,
  warningSeconds = 60,
  onLogout,
}: Options) {
  const [showWarning, setShowWarning] = useState(false)
  const [countdown, setCountdown] = useState(warningSeconds)

  // Refs agar callback timer tidak stale
  const isWarningActiveRef = useRef(false)
  const onLogoutRef        = useRef(onLogout)
  const warningTimerRef    = useRef<ReturnType<typeof setTimeout> | null>(null)
  const logoutTimerRef     = useRef<ReturnType<typeof setTimeout> | null>(null)
  const countdownRef       = useRef<ReturnType<typeof setInterval> | null>(null)

  useEffect(() => { onLogoutRef.current = onLogout })

  const clearAll = useCallback(() => {
    if (warningTimerRef.current) clearTimeout(warningTimerRef.current)
    if (logoutTimerRef.current)  clearTimeout(logoutTimerRef.current)
    if (countdownRef.current)    clearInterval(countdownRef.current)
  }, [])

  const resetTimers = useCallback(() => {
    clearAll()
    isWarningActiveRef.current = false
    setShowWarning(false)
    setCountdown(warningSeconds)

    // Munculkan warning (idleMinutes * 60 - warningSeconds) detik setelah idle
    const warnAfterMs = (idleMinutes * 60 - warningSeconds) * 1000

    warningTimerRef.current = setTimeout(() => {
      isWarningActiveRef.current = true
      setShowWarning(true)

      let remaining = warningSeconds
      setCountdown(remaining)

      countdownRef.current = setInterval(() => {
        remaining -= 1
        setCountdown(remaining)
        if (remaining <= 0) clearInterval(countdownRef.current!)
      }, 1000)

      logoutTimerRef.current = setTimeout(() => {
        onLogoutRef.current()
      }, warningSeconds * 1000)
    }, warnAfterMs)
  }, [idleMinutes, warningSeconds, clearAll])

  // Jika user klik "Saya Masih Di Sini" di modal
  const stayLoggedIn = useCallback(() => {
    resetTimers()
  }, [resetTimers])

  useEffect(() => {
    resetTimers()

    const handleActivity = () => {
      // Jika warning sudah muncul, abaikan aktivitas — biarkan modal yang handle
      if (isWarningActiveRef.current) return
      resetTimers()
    }

    ACTIVITY_EVENTS.forEach((e) =>
      window.addEventListener(e, handleActivity, { passive: true })
    )
    return () => {
      clearAll()
      ACTIVITY_EVENTS.forEach((e) =>
        window.removeEventListener(e, handleActivity)
      )
    }
  }, [resetTimers, clearAll])

  return { showWarning, countdown, stayLoggedIn }
}
