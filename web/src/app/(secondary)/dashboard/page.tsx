'use client'
import {Sidebar} from '@/components/Sidebar'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'

interface LoginResponse {
	access_token: string
}

async function getRefreshtoken(): Promise<boolean> {
	try {
		const resp = await fetch('http://localhost:3021/api/v1/refresh', {
			method: 'GET',
			credentials: 'include'
		})

		if (!resp.ok) {
			console.log(await resp.json())
			return false
		}

		const token = (await resp.json()) as LoginResponse
		localStorage.setItem('ACCESS_TOKEN', token.access_token)
		return true
	} catch (err) {
		return false
	}
}

export default function DashboardPage() {
	const [loggedIn, setLoggedin] = useState<boolean>(false)
	const router = useRouter()

	useEffect(() => {
		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			router.push('/login')
	
		} else {
			setLoggedin(true)
		}
	}, [])

	return (
		<>
			{loggedIn && (
				<div className="flex w-full flex-col items-center">
					{/* navbar */}
					<div
						id="user-content"
						className="w-2/3 h-full flex flex-row"
					>
						<div className="w-1/12 flex flex-col gap-3">
							{/* place holder, eventally map over snippets */}
							<Sidebar title="Golang">
								<p className="text-white py-1 w-full border hover:border-opacity-100 border-opacity-15 text-opacity-50 hover:text-opacity-100 border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden">
									main
								</p>
								<p className="text-white py-1 w-full border hover:border-opacity-100 border-opacity-15 text-opacity-50 hover:text-opacity-100 border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden">
									main
								</p>
							</Sidebar>
							<Sidebar title="Rust">
								<p className="text-white py-1 w-full border hover:border-opacity-100 border-opacity-15 text-opacity-50 hover:text-opacity-100 border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden">
									main
								</p>
								<p className="text-white py-1 w-full border hover:border-opacity-100 border-opacity-15 text-opacity-50 hover:text-opacity-100 border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden">
									main
								</p>
							</Sidebar>
						</div>
						<div className="w-11/12 bg-neutral-900 ">
							{/* categorie */}
							{/* title in big */}
							{/* tags */}
							{/* desription in 50  opacity */}
							{/* space */}
							{/* codeblock */}
							{/* extra shit maybe */}
						</div>
					</div>
				</div>
			)}
		</>
	)
}
