import { afterEach, describe, expect, it } from 'vitest'
import { cleanup, render, screen } from '@testing-library/react'
import App from './App'

describe('App', () => {
    afterEach(cleanup)

    it('renders the page heading and the calculator keypad', () => {
        render(<App />)

        expect(screen.getByRole('heading', { name: 'Calculator' })).toBeInTheDocument()
        expect(screen.getByRole('group', { name: 'Calculator keypad' })).toBeInTheDocument()
    })
})
