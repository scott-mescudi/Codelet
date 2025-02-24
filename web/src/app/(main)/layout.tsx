import {MainNavbar} from '@/components/MainNavbar'
import type {Metadata} from 'next'
import '../globals.css'

export const metadata: Metadata = {
	title: 'Codelet',
	description: 'Codelet'
}

export default function RootLayout({
	children
}: Readonly<{children: React.ReactNode}>) {
	return (
		<>
			<div className='w-full flex py-2 justify-center border border-t-0 border-r-0 border-l-0 border-b-white/5'>
				<MainNavbar />
			</div>
			{children}
		</>
	)
}
