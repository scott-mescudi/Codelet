'use client'
import {Sidebar} from '@/components/Sidebar'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'
import {jwtDecode} from 'jwt-decode'
import {CodeBox} from '@/components/CodeBlock'

interface LoginResponse {
	access_token: string
}

interface CodeSnippet {
	language: string
	title: string
	code: string
	favorite: boolean
	private: boolean
	tags: string[]
	description: string
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

interface SmallSnippet {
	id: number
	language: string
	title: string
	favorite: boolean
}

interface ErrorResp {
	error: string
	code: number
}

type SmallSnippets = SmallSnippet[]

function isTokenExpired(token: string): boolean {
	try {
		const decoded: {exp?: number} = jwtDecode(token)
		const now = Date.now() / 1000

		return decoded.exp !== undefined ? decoded.exp < now : true
	} catch (error) {
		return true
	}
}

async function getSmallSnippets(
	token: string
): Promise<SmallSnippets | undefined> {
	if (token === '') {
		return undefined
	}

	try {
		const resp = await fetch(
			'http://localhost:3021/api/v1/user/small/snippets',
			{
				method: 'GET',
				headers: {
					Authorization: token
				}
			}
		)

		if (!resp.ok) {
			return undefined
		}

		const snippets = (await resp.json()) as SmallSnippets
		if (snippets.length <= 0) {
			return undefined
		}

		return snippets
	} catch (err) {
		console.error(err)
		return undefined
	}
}

async function GetSnippetByID(
	snippetID: number,
	token: string
): Promise<CodeSnippet | undefined> {
	if (!snippetID || snippetID < 0) {
		return undefined
	}

	if (token === '') {
		return undefined
	}

	try {
		const resp = await fetch(
			`http://localhost:3021/api/v1/user/snippets/${snippetID}`,
			{
				method: 'GET',
				headers: {
					Authorization: token
				}
			}
		)

		if (!resp.ok) {
			const errResp = (await resp.json()) as ErrorResp
			console.log(errResp)
			return undefined
		}

		const snippets = (await resp.json()) as CodeSnippet
		return snippets
	} catch (err) {
		console.error(err)
		return undefined
	}
}

export default function DashboardPage() {
	const [loggedIn, setLoggedin] = useState<boolean>(false)
	const [snippets, setSnippets] = useState<SmallSnippets>([])
	const [categories, setCategories] = useState<string[]>([])
	const [snippetToGet, setSnippetToGet] = useState<number>()
	const [inViewSnippet, setInViewSnippet] = useState<CodeSnippet>()
	const router = useRouter()

	useEffect(() => {
		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			router.push('/login')
		} else {
			setLoggedin(true)
		}

		// hankle refresh logi chere
		if (isTokenExpired(token ? token : '')) {
			console.log('Session expired')
			router.push('/login')
			return
		}

		const get = async () => {
			const snippets = await getSmallSnippets(token ? token : '')
			if (!snippets) {
				return
			}

			setSnippets(snippets)

			const LastSnippet = localStorage.getItem('LastSnippet')
			if (!LastSnippet) {
				setSnippetToGet(snippets[0].id)
			}else{
				setSnippetToGet(Number(LastSnippet))
			}

			const uniqueLanguages = [
				...new Set(snippets.map(snippet => snippet.language))
			]

			setCategories(uniqueLanguages)
		}

		get()
	}, [])

	useEffect(() => {
		if (!snippetToGet) {
			return
		}

		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			router.push('/login')
		} else {
			setLoggedin(true)
		}

		if (isTokenExpired(token ? token : '')) {
			console.log('Session expired')
			router.push('/login')
			return
		}

		const getSnippet = async () => {
			const snippet = await GetSnippetByID(
				snippetToGet,
				token ? token : ''
			)
			if (snippet === undefined) {
				console.error('fialed to get snippet')
				return
			}

			setInViewSnippet(snippet)
			localStorage.setItem('LastSnippet', snippetToGet ? String(snippetToGet) : '')
		}

		getSnippet()
	}, [snippetToGet])

	return (
		<>
			{loggedIn && (
				<div className="flex w-full flex-col items-center">
					{/* navbar */}
					<div
						id="user-content"
						className="w-2/3 h-full flex flex-row"
					>
						<div id='sidebar' className="w-2/12 flex flex-col gap-3">
							{snippets &&
								snippets.length > 0 &&
								categories.length > 0 &&
								categories.map((category: string) => (
									<Sidebar key={category} title={category}>
										{snippets
											.filter(
												snippet =>
													snippet.language ===
													category
											)
											.map(snippet => (
												<p
													onClick={() =>
														setSnippetToGet(
															snippet.id
														)
													}
													key={snippet.id}
													className="text-white py-1 w-full border hover:border-opacity-100 border-opacity-15 text-opacity-50 hover:text-opacity-100 border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden"
												>
													{snippet.title}
												</p>
											))}
									</Sidebar>
								))}
						</div>
						<div className="w-10/12 ml-3 flex flex-col">
							<p className="w-full select-none text-white text-left text-6xl font-bold ">
								{inViewSnippet?.title}
							</p>
							<div className="w-full select-none flex gap-5 mt-2 flex-row">
								{inViewSnippet?.tags.map(
									(tag: string, idx: number) => (
										<p
											key={idx}
											className="text-white text-nowrap w-fit text-opacity-50 px-5 rounded-lg py-0.5 bg-neutral-800"
										>
											{tag}
										</p>
									)
								)}
							</div>
							<p className="w-full pl-1 text-white text-left mt-4 text-opacity-50">
								{inViewSnippet?.description}
							</p>
							<div className="w-full mt-10">
								<CodeBox
									code={
										inViewSnippet?.code
											? inViewSnippet.code
											: ''
									}
								/>
							</div>
						</div>
					</div>
				</div>
			)}
		</>
	)
}

//   <div className="w-10/12  min-h-full flex flex-col gap-1 text-white">
// 	<p className="w-full select-none pl-1 text-white text-left text-xl font-bold text-opacity-50 font-mono">{inViewSnippet?.language.toLowerCase()}</p>
// 	<p className="w-full select-none text-white text-left text-6xl font-bold ">{inViewSnippet?.title}</p>
// 	<div className="w-full select-none flex gap-5 mt-2 flex-row">
// 	  {inViewSnippet?.tags.map((tag:string, idx:number) => (
// 		<p key={idx} className="text-white text-nowrap w-fit text-opacity-50 px-5 rounded-lg py-0.5 bg-neutral-800">{tag}</p>
// 	  ))}
// 	</div>
// 	<p className="w-full pl-1 text-white text-left mt-4 text-opacity-50">{inViewSnippet?.description}</p>
// 	<div className="w-full mt-10">
// 	  <CodeBox code={inViewSnippet?.code ? inViewSnippet.code : ""} />
// 	</div>
//   </div>
