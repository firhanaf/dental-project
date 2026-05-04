import { Button } from './ui'

interface Props {
  countdown: number
  onStay: () => void
  onLogout: () => void
}

export default function IdleWarningModal({ countdown, onStay, onLogout }: Props) {
  const isUrgent = countdown <= 10

  return (
    <div className="modal-overlay" style={{ zIndex: 9999 }}>
      <div className="modal" style={{ maxWidth: 340 }}>
        <div className="modal-body" style={{ padding: '32px 28px', textAlign: 'center' }}>

          <div className="w-12 h-12 rounded-full flex items-center justify-center text-2xl mx-auto mb-4"
            style={{ background: 'var(--teal-l)' }}>
            🔒
          </div>

          <h3 className="text-[15px] font-semibold mb-2" style={{ color: 'var(--text)' }}>
            Sesi Hampir Berakhir
          </h3>
          <p className="text-[13px] mb-5" style={{ color: 'var(--text2)' }}>
            Tidak ada aktivitas terdeteksi.<br />
            Anda akan otomatis keluar dalam:
          </p>

          <div className="mb-6">
            <span
              className="text-[60px] font-bold leading-none"
              style={{
                fontFamily: 'DM Mono, monospace',
                color: isUrgent ? 'var(--danger-t)' : 'var(--teal)',
                transition: 'color 0.5s',
              }}
            >
              {countdown}
            </span>
            <p className="text-[12px] mt-1" style={{ color: 'var(--text3)' }}>detik</p>
          </div>

          <div className="flex flex-col gap-2">
            <Button
              variant="primary"
              onClick={onStay}
              className="w-full justify-center"
            >
              Saya Masih Di Sini
            </Button>
            <Button
              variant="secondary"
              onClick={onLogout}
              className="w-full justify-center"
            >
              Keluar Sekarang
            </Button>
          </div>
        </div>
      </div>
    </div>
  )
}
