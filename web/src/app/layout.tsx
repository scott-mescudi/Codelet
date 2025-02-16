import type {Metadata} from 'next'
import './globals.css'

export const metadata: Metadata = {
	title: 'Codelet',
	icons: {
		icon: '/logo.svg'
	},
	description: 'Codelet'
}

export default function RootLayout({children}: Readonly<{children: React.ReactNode}>) {
	return (
		<html lang="en" className="overflow-y-scroll">
			<body className="overflow-hidden flex flex-col items-center">{children}</body>
		</html>
	)
}