export default function Home() {
    return (
        <main className="space-y-4">
            <p className="text-gray-700">Pilih menu:</p>
            <ul className="list-disc pl-5">
                <li><a className="underline" href="/classes">Classes</a></li>
                <li><a className="underline" href="/tasks">Tasks</a></li>
            </ul>
        </main>
    );
}