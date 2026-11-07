import { Calculator } from './components/Calculator'

export default function App() {
    return (
        <main className="flex min-h-dvh flex-col items-center justify-center gap-6 bg-brand-dark px-4 py-8 sm:py-12 short:gap-2 short:py-3">
            <h1 className="text-2xl font-semibold text-white short:sr-only">Calculator</h1>
            <Calculator />
        </main>
    )
}
