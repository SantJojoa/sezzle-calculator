import { useEffect } from 'react'
import { useCalculator } from '../hooks/useCalculator'
import { Display } from './Display'
import { Keypad } from './Keypad'

const INPUT_KEY = /^[0-9.+\-*/^()]$/

export function Calculator() {
    const { expression, result, error, loading, press, clear, deleteLast, submit } = useCalculator()

    useEffect(() => {
        function onKeyDown(e: KeyboardEvent) {
            // leave browser/OS shortcuts (Cmd+C, Ctrl+V, ...) alone
            if (e.ctrlKey || e.metaKey) return
            // never hijack typing in form fields
            if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return

            if (INPUT_KEY.test(e.key)) {
                e.preventDefault() // e.g. "/" would open quick find in some browsers
                press(e.key)
            } else if (e.key === 'Enter' || e.key === '=') {
                // Enter on a focused button must activate that button, not submit
                if (e.key === 'Enter' && e.target instanceof HTMLButtonElement) return
                e.preventDefault()
                if (!e.repeat) void submit()
            } else if (e.key === 'Backspace') {
                deleteLast()
            } else if (e.key === 'Escape') {
                clear()
            }
        }

        window.addEventListener('keydown', onKeyDown)
        return () => window.removeEventListener('keydown', onKeyDown)
    }, [press, deleteLast, clear, submit])

    return (
        <div className="flex w-full max-w-sm flex-col gap-4 rounded-3xl bg-surface p-4 shadow-xl shadow-black/30 sm:p-6">
            <Display expression={expression} result={result} error={error} loading={loading} />
            <Keypad
                onKey={press}
                onClear={clear}
                onDelete={deleteLast}
                onSubmit={() => void submit()}
                disabled={loading}
            />
        </div>
    )
}
