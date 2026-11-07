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

const DIGIT = 'bg-gray-700 hover:bg-gray-600'
const OPERATOR = 'bg-amber-600 hover:bg-amber-500'
const UTILITY = 'bg-gray-500 hover:bg-gray-400'

const KEYS: KeyDef[] = [
    { label: 'C', action: 'clear', ariaLabel: 'clear', className: UTILITY },
    { label: '(', action: 'key', value: '(', ariaLabel: 'open parenthesis', className: UTILITY },
    { label: ')', action: 'key', value: ')', ariaLabel: 'close parenthesis', className: UTILITY },
    { label: '⌫', action: 'delete', ariaLabel: 'delete', className: UTILITY },

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
        className: 'row-span-2 bg-emerald-600 hover:bg-emerald-500 disabled:opacity-50 disabled:hover:bg-emerald-600',
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
        <div className="grid grid-cols-4 gap-2">
            {KEYS.map((k) => (
                <button
                    key={k.label}
                    type="button"
                    aria-label={k.ariaLabel}
                    disabled={k.action === 'submit' ? disabled : false}
                    onClick={() => handle(k)}
                    className={`rounded-lg p-4 text-xl font-medium text-white transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-white ${k.className ?? ''}`}
                >
                    {k.label}
                </button>
            ))}
        </div>
    )
}
