import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { cleanup, render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { Calculator } from './Calculator'
import { evaluate } from '../api/calculatorClient'

vi.mock('../api/calculatorClient', async (importOriginal) => ({
    ...(await importOriginal<typeof import('../api/calculatorClient')>()),
    evaluate: vi.fn(),
}))

const evaluateMock = vi.mocked(evaluate)

describe('Calculator accessibility and keyboard handling', () => {
    beforeEach(() => {
        evaluateMock.mockReset()
    })
    afterEach(() => {
        cleanup() // vitest globals are off, so RTL does not auto-clean
        vi.restoreAllMocks()
    })

    it('announces the result in the status live region', async () => {
        evaluateMock.mockResolvedValue(14)
        const user = userEvent.setup()
        render(<Calculator />)

        expect(screen.getByRole('status')).toBeEmptyDOMElement() // mounted before any result

        await user.keyboard('2+3*4{Enter}')

        expect(evaluateMock).toHaveBeenCalledWith('2+3*4')
        expect(await screen.findByRole('status')).toHaveTextContent('14')
    })

    it('activates a focused key with Enter instead of submitting', async () => {
        const user = userEvent.setup()
        render(<Calculator />)

        screen.getByRole('button', { name: '7' }).focus()
        await user.keyboard('{Enter}')

        expect(screen.getByTestId('expression')).toHaveTextContent('7')
        expect(evaluateMock).not.toHaveBeenCalled()
    })

    it('still submits with Enter after clicking a key with the mouse', async () => {
        evaluateMock.mockResolvedValue(7)
        const user = userEvent.setup()
        render(<Calculator />)

        await user.click(screen.getByRole('button', { name: '7' }))
        await user.keyboard('{Enter}')

        expect(evaluateMock).toHaveBeenCalledWith('7')
    })

    it('labels the keypad group and the symbol keys', () => {
        render(<Calculator />)

        expect(screen.getByRole('group', { name: 'Calculator keypad' })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'multiply' })).toBeInTheDocument()
        expect(screen.getByRole('button', { name: 'square root' })).toBeInTheDocument()
    })

    it('shows errors in an alert', async () => {
        const { ApiError } = await import('../api/calculatorClient')
        evaluateMock.mockRejectedValue(new ApiError('DIVISION_BY_ZERO', 'cannot divide by zero', 422))
        const user = userEvent.setup()
        render(<Calculator />)

        await user.keyboard('1/0{Enter}')

        expect(await screen.findByRole('alert')).toHaveTextContent('cannot divide by zero')
    })

    it('builds the expression from clicked keys', async () => {
        const user = userEvent.setup()
        render(<Calculator />)

        for (const name of ['2', 'add', '3', 'multiply', '4']) {
            await user.click(screen.getByRole('button', { name }))
        }

        // display shows pretty symbols; the API still receives * (checked in the next test)
        expect(screen.getByTestId('expression')).toHaveTextContent('2+3×4')
    })

    it('shows the formatted result when "=" is clicked', async () => {
        evaluateMock.mockResolvedValue(0.1 + 0.2)
        const user = userEvent.setup()
        render(<Calculator />)

        for (const name of ['2', 'multiply', '3']) {
            await user.click(screen.getByRole('button', { name }))
        }
        await user.click(screen.getByRole('button', { name: 'equals' }))

        expect(evaluateMock).toHaveBeenCalledWith('2*3') // API character, not the symbol
        expect(await screen.findByRole('status')).toHaveTextContent('0.3') // not 0.30000000000000004
    })

    it('shows the backend ApiError message in the alert', async () => {
        const { ApiError } = await import('../api/calculatorClient')
        evaluateMock.mockRejectedValue(new ApiError('INVALID_EXPRESSION', 'invalid expression at position 2', 400))
        const user = userEvent.setup()
        render(<Calculator />)

        await user.click(screen.getByRole('button', { name: '2' }))
        await user.click(screen.getByRole('button', { name: 'equals' }))

        expect(await screen.findByRole('alert')).toHaveTextContent('invalid expression at position 2')
        expect(screen.getByRole('status')).toBeEmptyDOMElement()
    })

    it('shows a client-side validation error without calling the API', async () => {
        const user = userEvent.setup()
        render(<Calculator />)

        await user.click(screen.getByRole('button', { name: 'open parenthesis' }))
        await user.click(screen.getByRole('button', { name: '2' }))
        await user.click(screen.getByRole('button', { name: 'equals' }))

        expect(await screen.findByRole('alert')).toHaveTextContent('Missing closing parenthesis')
        expect(evaluateMock).not.toHaveBeenCalled()
    })

    it('continues from the result when an operator is pressed', async () => {
        evaluateMock.mockResolvedValue(14)
        const user = userEvent.setup()
        render(<Calculator />)

        await user.keyboard('2+3*4{Enter}')
        expect(await screen.findByRole('status')).toHaveTextContent('14')

        await user.click(screen.getByRole('button', { name: 'multiply' }))

        expect(screen.getByTestId('expression')).toHaveTextContent('14×')
        expect(screen.getByRole('status')).toBeEmptyDOMElement()
    })

    it('starts a new expression when a digit is pressed after a result', async () => {
        evaluateMock.mockResolvedValue(14)
        const user = userEvent.setup()
        render(<Calculator />)

        await user.keyboard('2+3*4{Enter}')
        expect(await screen.findByRole('status')).toHaveTextContent('14')

        await user.click(screen.getByRole('button', { name: '5' }))

        const expression = screen.getByTestId('expression')
        expect(expression).toHaveTextContent('5')
        expect(expression).not.toHaveTextContent('14')
    })

    it('accepts physical keyboard input: typing 2+3 and Enter', async () => {
        evaluateMock.mockResolvedValue(5)
        const user = userEvent.setup()
        render(<Calculator />)

        await user.keyboard('2+3')
        expect(screen.getByTestId('expression')).toHaveTextContent('2+3')

        await user.keyboard('{Enter}')

        expect(evaluateMock).toHaveBeenCalledExactlyOnceWith('2+3')
        expect(await screen.findByRole('status')).toHaveTextContent('5')
    })

    it('disables "=" while loading and re-enables it afterwards', async () => {
        let resolve!: (n: number) => void
        evaluateMock.mockReturnValue(new Promise<number>((r) => { resolve = r }))
        const user = userEvent.setup()
        render(<Calculator />)

        const equals = screen.getByRole('button', { name: 'equals' })
        await user.click(screen.getByRole('button', { name: '2' }))
        await user.click(equals)

        expect(equals).toBeDisabled()

        resolve(2)

        await waitFor(() => expect(equals).toBeEnabled())
        expect(screen.getByRole('status')).toHaveTextContent('2')
        expect(evaluateMock).toHaveBeenCalledTimes(1)
    })

    it('deletes the last character with Backspace', async () => {
        const user = userEvent.setup()
        render(<Calculator />)

        await user.keyboard('12{Backspace}')

        expect(screen.getByTestId('expression')).toHaveTextContent(/^Expression: 1$/)
    })

    it('clears the expression and result with Escape', async () => {
        evaluateMock.mockResolvedValue(5)
        const user = userEvent.setup()
        render(<Calculator />)

        await user.keyboard('2+3{Enter}')
        expect(await screen.findByRole('status')).toHaveTextContent('5')

        await user.keyboard('{Escape}')

        expect(screen.getByTestId('expression')).toBeEmptyDOMElement()
        expect(screen.getByRole('status')).toBeEmptyDOMElement()
    })

    it('clears everything when the C key is clicked', async () => {
        const user = userEvent.setup()
        render(<Calculator />)

        await user.click(screen.getByRole('button', { name: '7' }))
        await user.click(screen.getByRole('button', { name: 'clear' }))

        expect(screen.getByTestId('expression')).toBeEmptyDOMElement()
    })

    it('deletes the last character when the delete key is clicked', async () => {
        const user = userEvent.setup()
        render(<Calculator />)

        await user.click(screen.getByRole('button', { name: '7' }))
        await user.click(screen.getByRole('button', { name: '8' }))
        await user.click(screen.getByRole('button', { name: 'delete' }))

        expect(screen.getByTestId('expression')).toHaveTextContent(/^Expression: 7$/)
    })
})
