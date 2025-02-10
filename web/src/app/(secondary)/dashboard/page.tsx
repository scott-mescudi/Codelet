'use client'
import {Sidebar} from '@/components/Sidebar'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'
import {jwtDecode} from 'jwt-decode'
import {CodeBox} from '@/components/CodeBlock'
import MenuIcon from '@mui/icons-material/Menu';
import HouseIcon from '@mui/icons-material/House';

import logo from "../../../../public/logo.svg"
import Image from 'next/image'
import LogoutIcon from '@mui/icons-material/Logout';
import { DropdownItem } from '@/components/DropdownItem'

interface LoginResponse {
	access_token: string
}

interface CodeSnippet {
	id: number
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

interface UserContentProps {
	snippets: SmallSnippets
	categories: string[]
	setSnippetToGet: React.Dispatch<React.SetStateAction<number | undefined>>
	inViewSnippet: CodeSnippet | undefined
	setAddsnippet: React.Dispatch<React.SetStateAction<boolean>>
}

async function Logout(token:string) {
	if (token === "")  return;

	try {
		const resp = await fetch(
			`http://localhost:3021/api/v1/logout`,
			{
				method: 'POST',
				headers: {
					Authorization: token
				}
			}
		)

		if (!resp.ok) {
			const errResp = (await resp.json()) as ErrorResp
			console.log(errResp)
			return
		}
	} catch (err) {
		console.error(err)
		return
	}
}

export function UserContent({
	snippets,
	categories,
	setSnippetToGet,
	inViewSnippet,
	setAddsnippet
}: UserContentProps) {
	return (
		<>	
		<div className='h-[90vh] pb-10 relative lg:flex hidden overflow-auto scrollbar-none w-2/12'>
			<div id="sidebar" className="w-full flex flex-col gap-3">
				{snippets &&
					snippets.length > 0 &&
					categories.length > 0 &&
					categories.map((category: string) => (
						<Sidebar key={category} title={category}>
							{snippets
								.filter(
									snippet => snippet.language === category
								)
								.map(snippet => (
									<p
										onClick={() =>
											setSnippetToGet(snippet.id)
										}
										key={snippet.id}
										className={`text-white py-1 w-full border ${
											inViewSnippet?.id === snippet.id
												? 'border-opacity-100 text-opacity-100'
												: ' border-opacity-15 text-opacity-60 hover:text-opacity-100'
										} border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden`}
									>
										{snippet.title}
									</p>
								))}
						</Sidebar>
					))}
			</div>
		</div>
		
		<div className='h-[90vh] pb-10 w-10/12 overflow-auto scrollbar-none'>
			<div className="w-full flex flex-col">
				<p className="w-full line-clamp-1 h-20 select-none text-white  text-6xl font-bold">
					{inViewSnippet?.title}
				</p>
				<div className="w-full select-none flex gap-5 sm:mt-3 mt-5 flex-row">
					{inViewSnippet?.tags.map((tag: string, idx: number) => (
						<p
							key={idx}
							className="text-white text-nowrap w-fit text-opacity-60 px-5 rounded-lg py-0.5 bg-neutral-900"
						>
							{tag}
						</p>
					))}
				</div>
				<p className="w-full pl-1 text-white   mt-4 text-opacity-80">
					{inViewSnippet?.description}
				</p>
				<div className="w-full mt-10 ">
					<CodeBox
						background="bg-neutral-950"
						code={inViewSnippet?.code ? inViewSnippet.code : ''}
					/>
				</div>
			</div>
		</div>

		</>
	)
}

export default function DashboardPage() {
	const [loggedIn, setLoggedin] = useState<boolean>(false)
	const [snippets, setSnippets] = useState<SmallSnippets>([])
	const [categories, setCategories] = useState<string[]>([])
	const [snippetToGet, setSnippetToGet] = useState<number>()
	const [inViewSnippet, setInViewSnippet] = useState<CodeSnippet>()
	const [addSnippet, setAddsnippet] = useState<boolean>(false)
	const [dropdownOpen, setDropdownOpen] = useState <boolean>(false)
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
			} else {
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
			localStorage.setItem(
				'LastSnippet',
				snippetToGet ? String(snippetToGet) : ''
			)
		}

		getSnippet()
	}, [snippetToGet])

	const LogoutHandler = async () => {
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

		await Logout(token ? token : "")

		localStorage.removeItem("ACCESS_TOKEN")
		router.push('/login')
		setLoggedin(false)
	}

	return (
		<>
			{loggedIn && (
				<div className="flex  w-full flex-col gap-10 items-center">
					{/* navbar */}
					<div className="lg:w-2/3 h-16 mt-5  flex  items-center   rounded-xl">
						<div className='h-full aspect-square'>
							<Image draggable={false} src={logo} className='h-full w-full ' alt='codelet logo' />
						</div>
						<p className='text-3xl select-none ml-2 text-white font-bold'>Codelet</p>
						<button className='bg-white hover:bg-opacity-80 duration-300 ease-in-out ml-auto h-fit py-1 px-5 text-lg font-semibold rounded-lg' onClick={() => setAddsnippet(true)}>new snippet</button>
						<button onClick={() => setDropdownOpen(prev => !prev)} className='p-1 rounded-md ml-3 relative text-white'> 
							<MenuIcon fontSize='large' />
							{dropdownOpen &&
							<div  className='absolute mt-4 z-50 bg-black right-0 mr-2'>
								<div className='w-fit h-fit p-4 flex flex-col gap-2  border border-white border-opacity-15'>
									<DropdownItem link='/' title='Home' subTitle='Back to home' icon={<HouseIcon fontSize='large' />} />
									<DropdownItem onClick={LogoutHandler} title='Logout' subTitle='Secure Logout portal' icon={<LogoutIcon fontSize='large' />} />
								</div>
							</div>
							}
						</button>
					</div>
					<div id="user-content" className="lg:w-2/3 h-full gap-5 flex flex-row justify-center">
						<UserContent snippets={snippets} categories={categories} setSnippetToGet={setSnippetToGet} inViewSnippet={inViewSnippet} setAddsnippet={setAddsnippet}/>
					</div>
					{addSnippet && (
						<div id="parent" onClick={() => setAddsnippet(false)} className="fixed h-screen w-screen backdrop-blur-lg">
							<div className="w-full h-full flex justify-center items-center">
								<div onClick={e => e.stopPropagation()} className="sm:w-2/3 w-11/12 mt-auto overflow-auto scrollbar-hidden rounded-t-xl h-3/4 bg-neutral-900">
									
								</div>
							</div>
						</div>
					)}
				</div>
			)}
		</>
	)
}
