interface KeypadProps {
    onKey: (key: string) => void
    onClear: () => void
    onDelete: () => void
    onSubmit: () => void
    /** Disables "=" (e.g. while a request is in flight) to avoid double submits. */
    disabled?: boolean
}

type Action = 'key' | 'clear' | 'delete' | 'submit'

interface KeyDef {
    label: string // what the user sees
    action: Action
    value?: string // what the API/expression receives (for action "key")
    ariaLabel?: string // needed when the label is a symbol screen readers can't read
    className?: string // grid span / colour overrides
}

const DIGIT = 'bg-key text-white'
const OPERATOR = 'bg-brand text-white'
const CLEAR = 'bg-danger text-white'

const KEYS: KeyDef[] = [
    { label: 'C', action: 'clear', ariaLabel: 'clear', className: CLEAR },
    { label: '(', action: 'key', value: '(', ariaLabel: 'open parenthesis', className: OPERATOR },
    { label: ')', action: 'key', value: ')', ariaLabel: 'close parenthesis', className: OPERATOR },
    { label: '⌫', action: 'delete', ariaLabel: 'delete', className: OPERATOR },

    { label: '√', action: 'key', value: 'sqrt(', ariaLabel: 'square root', className: OPERATOR },
    { label: '^', action: 'key', value: '^', ariaLabel: 'power', className: OPERATOR },
    { label: '÷', action: 'key', value: '/', ariaLabel: 'divide', className: OPERATOR },
    { label: '×', action: 'key', value: '*', ariaLabel: 'multiply', className: OPERATOR },

    { label: '7', action: 'key', value: '7', className: DIGIT },
    { label: '8', action: 'key', value: '8', className: DIGIT },
    { label: '9', action: 'key', value: '9', className: DIGIT },
    { label: '−', action: 'key', value: '-', ariaLabel: 'subtract', className: OPERATOR },

    { label: '4', action: 'key', value: '4', className: DIGIT },
    { label: '5', action: 'key', value: '5', className: DIGIT },
    { label: '6', action: 'key', value: '6', className: DIGIT },
    { label: '+', action: 'key', value: '+', ariaLabel: 'add', className: OPERATOR },

    { label: '1', action: 'key', value: '1', className: DIGIT },
    { label: '2', action: 'key', value: '2', className: DIGIT },
    { label: '3', action: 'key', value: '3', className: DIGIT },
    // "=" spans the last two rows so the grid has no empty cell
    {
        label: '=',
        action: 'submit',
        ariaLabel: 'equals',
        className: 'row-span-2 bg-accent text-brand-dark',
    },

    { label: '0', action: 'key', value: '0', className: `col-span-2 ${DIGIT}` },
    { label: '.', action: 'key', value: '.', ariaLabel: 'decimal point', className: DIGIT },
]

export function Keypad({ onKey, onClear, onDelete, onSubmit, disabled = false }: KeypadProps) {
    function handle(k: KeyDef) {
        switch (k.action) {
            case 'clear':
                return onClear()
            case 'delete':
                return onDelete()
            case 'submit':
                return onSubmit()
            case 'key':
                return onKey(k.value ?? k.label)
        }
    }

    return (
        <div role="group" aria-label="Calculator keypad" className="grid grid-cols-4 gap-2 md:gap-3 short:flex-1 short:gap-1.5">
            {KEYS.map((k) => (
                <button
                    key={k.label}
                    type="button"
                    aria-label={k.ariaLabel}
                    disabled={k.action === 'submit' ? disabled : false}
                    onClick={() => handle(k)}
                    // keep focus off the button on mouse click, so a following Enter
                    // still submits (see the Enter-on-button rule in Calculator)
                    onMouseDown={(e) => e.preventDefault()}
                    className={`min-h-12 rounded-xl p-3 text-xl font-semibold md:min-h-16 md:text-2xl short:min-h-9 short:p-1 short:text-base transition hover:brightness-110 active:scale-95 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-accent disabled:cursor-not-allowed disabled:opacity-50 ${k.className ?? ''}`}
                >
                    {k.label}
                </button>
            ))}
        </div>
    )
}
