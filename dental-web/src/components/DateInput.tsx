import { useState, useEffect, type ChangeEvent } from 'react'

interface DateInputProps {
  value: string                   // yyyy-mm-dd or ''
  onChange: (v: string) => void   // emits yyyy-mm-dd or '' (while incomplete)
  required?: boolean
  className?: string
  placeholder?: string
}

function isoToDisplay(iso: string): string {
  if (!iso || iso.length < 10) return ''
  const [y, m, d] = iso.split('-')
  if (!y || !m || !d) return ''
  return `${d}/${m}/${y}`
}

function digitsToDisplay(digits: string): string {
  const d = digits.slice(0, 8)
  let r = ''
  if (d.length > 0) r += d.slice(0, 2)
  if (d.length > 2) r += '/' + d.slice(2, 4)
  if (d.length > 4) r += '/' + d.slice(4, 8)
  return r
}

function displayToISO(display: string): string {
  const digits = display.replace(/\D/g, '')
  if (digits.length < 8) return ''
  return `${digits.slice(4, 8)}-${digits.slice(2, 4)}-${digits.slice(0, 2)}`
}

export default function DateInput({
  value,
  onChange,
  required,
  className,
  placeholder = 'dd/mm/yyyy',
}: DateInputProps) {
  const [display, setDisplay] = useState(() => isoToDisplay(value))

  // Sync when parent resets (e.g. form clear) or loads edit data
  useEffect(() => {
    setDisplay(isoToDisplay(value))
  }, [value])

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    const digits = e.target.value.replace(/\D/g, '').slice(0, 8)
    const formatted = digitsToDisplay(digits)
    setDisplay(formatted)
    onChange(displayToISO(formatted))
  }

  return (
    <input
      className={className ?? 'form-input'}
      type="text"
      inputMode="numeric"
      value={display}
      onChange={handleChange}
      placeholder={placeholder}
      maxLength={10}
      required={required}
    />
  )
}
