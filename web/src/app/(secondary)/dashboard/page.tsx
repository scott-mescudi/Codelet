'use client'
import {Sidebar} from '@/components/Sidebar'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'
import {CodeBox} from '@/components/CodeBlock'
import MenuIcon from '@mui/icons-material/Menu'
import HouseIcon from '@mui/icons-material/House'

import logo from '@public/logo.svg'
import Image from 'next/image'
import LogoutIcon from '@mui/icons-material/Logout'
import {DropdownItem} from '@/components/DropdownItem'
import {AppRouterInstance} from 'next/dist/shared/lib/app-router-context.shared-runtime'
import {Delete} from '@mui/icons-material'
import EditIcon from '@mui/icons-material/Edit'
import Link from 'next/link'
import {Logout, GetUsername} from '@/shared/api/UserApiReq'
import {
	GetSnippetByID,
	getSmallSnippets,
	DeleteReq,
	addSnippetReq,
	updateSnippetReq
} from '@/shared/api/SnippetApiReq'
import {IsTokenExpired} from '@/shared/helper/jwtAuth'

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

interface SnippetFormProps {
	setAddsnippet: (p: boolean) => void
	router: AppRouterInstance
}

interface SmallSnippet {
	id: number
	language: string
	title: string
	favorite: boolean
}

type SmallSnippets = SmallSnippet[]

interface UpdateSnippetFromProps {
	setUpdateSnippet: (p: boolean) => void
	snippet: CodeSnippet
	router: AppRouterInstance
}

export function UpdateSnippetForm({
	setUpdateSnippet,
	router,
	snippet
}: UpdateSnippetFromProps) {
	const [language, setLanguage] = useState<string>(snippet.language)
	const [title, setTitle] = useState<string>(snippet.title)
	const [tags, setTags] = useState<string>(snippet.tags.join(','))
	const [description, setDescription] = useState<string>(snippet.description)
	const [code, setCode] = useState<string>(snippet.code)

	const HandleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()

		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			router.push('/login')
		}

		if (IsTokenExpired(token ? token : '')) {
			console.log('Session expired')
			router.push('/login')
			return
		}

		const tokens = tags
			.split(',')
			.map(str => str.trim())
			.filter(str => str.length > 0)
		const resp = await updateSnippetReq(
			snippet.id,
			token ? token : '',
			language,
			title,
			tokens,
			description,
			code
		)
		if (!resp) {
			console.error('Failed to add snippet')
			return
		}

		setUpdateSnippet(false)
		window.location.reload()
	}

	return (
		<>
			<h2 className="text-white text-3xl font-semibold text-center mb-6">
				Update a New Snippet
			</h2>
			<form onSubmit={HandleSubmit} className="flex flex-col gap-5">
				<input
					value={language}
					maxLength={49}
					onChange={e => setLanguage(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Language"
					required
				/>
				<input
					value={title}
					maxLength={250}
					onChange={e => setTitle(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Title"
					required
				/>
				<input
					value={tags}
					onChange={e => setTags(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Tags (comma separated)"
				/>
				<input
					value={description}
					onChange={e => setDescription(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Description"
				/>
				<textarea
					maxLength={3000}
					value={code}
					onChange={e => setCode(e.target.value)}
					placeholder="Paste code here..."
					className="p-5 bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 scrollbar-thin min-h-96 max-h-96 rounded-lg text-white whitespace-pre aspect-video"
					required
				/>

				<div className="flex w-full flex-row gap-2">
					<button
						type="button"
						onClick={() => setUpdateSnippet(false)}
						className="mt-4 w-full bg-red-600 hover:bg-red-700 text-white font-medium py-3 rounded-lg transition"
					>
						Cancel
					</button>
					<button
						type="submit"
						className="mt-4 w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 rounded-lg transition"
					>
						Update
					</button>
				</div>
			</form>
		</>
	)
}

export function SnippetForm({setAddsnippet, router}: SnippetFormProps) {
	const [language, setLanguage] = useState<string>('')
	const [title, setTitle] = useState<string>('')
	const [tags, setTags] = useState<string>('')
	const [description, setDescription] = useState<string>('')
	const [code, setCode] = useState<string>('')

	const HandleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()

		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			router.push('/login')
		}

		if (IsTokenExpired(token ? token : '')) {
			console.log('Session expired')
			router.push('/login')
			return
		}

		const tokens = tags
			.split(',')
			.map(str => str.trim())
			.filter(str => str.length > 0)
		const resp = await addSnippetReq(
			token ? token : '',
			language,
			title,
			tokens,
			description,
			code
		)
		if (!resp) {
			console.error('Failed to add snippet')
			return
		}

		setAddsnippet(false)
	}

	return (
		<>
			<h2 className="text-white text-3xl font-semibold text-center mb-6">
				Add a New Snippet
			</h2>
			<form onSubmit={HandleSubmit} className="flex flex-col gap-5">
				<input
					value={language}
					maxLength={49}
					onChange={e => setLanguage(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Language"
					required
				/>
				<input
					value={title}
					maxLength={250}
					onChange={e => setTitle(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Title"
					required
				/>
				<input
					value={tags}
					onChange={e => setTags(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Tags (comma separated)"
				/>
				<input
					value={description}
					onChange={e => setDescription(e.target.value)}
					className="bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 text-white p-5 rounded-lg"
					type="text"
					placeholder="Description"
				/>
				<textarea
					maxLength={3000}
					value={code}
					onChange={e => setCode(e.target.value)}
					placeholder="Paste code here..."
					className="p-5 bg-neutral-900 focus:ring-blue-500 outline-none focus:ring-2 scrollbar-thin min-h-96 max-h-96 rounded-lg text-white whitespace-pre aspect-video"
					required
				/>

				<div className="flex w-full flex-row gap-2">
					<button
						type="button"
						onClick={() => setAddsnippet(false)}
						className="mt-4 w-full bg-red-600 hover:bg-red-700 text-white font-medium py-3 rounded-lg transition"
					>
						Cancel
					</button>
					<button
						type="submit"
						className="mt-4 w-full bg-blue-600 hover:bg-blue-700 text-white font-medium py-3 rounded-lg transition"
					>
						Add
					</button>
				</div>
			</form>
		</>
	)
}

interface DeleteProps {
	snippets: SmallSnippets
	setDeleteSnippet: React.Dispatch<React.SetStateAction<boolean>>
	deleteSnippet: boolean
	setSnippetToGet: React.Dispatch<React.SetStateAction<number | undefined>>
	id: number
}

export function DeleteButton({
	id,
	snippets,
	setDeleteSnippet,
	setSnippetToGet
}: DeleteProps) {
	if (id < 0) {
		console.error('invalid id in delete button', id)
		return
	}

	const [isSure, setIsSure] = useState<boolean>(false)

	const HandleDelete = async () => {
		setDeleteSnippet(true)
		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			console.error('Session expired')
		}

		if (IsTokenExpired(token ? token : '')) {
			console.log('Session expired')
			console.error('Session expired')
			return
		}

		const idx = snippets.findIndex(snippet => snippet.id === id)

		const ok = await DeleteReq(token || '', id)
		if (!ok) {
			console.error('Failed to delete snippet')
		} else {
			const newSnippets = [...snippets]
			newSnippets.splice(idx, 1)

			if (snippets.length == 1) {
				window.location.reload()
			}

			if (newSnippets.length > 0) {
				const nextSnippet =
					idx === newSnippets.length
						? newSnippets[idx - 1]
						: newSnippets[idx]
				setSnippetToGet(nextSnippet.id)
			} else {
				setSnippetToGet(undefined)
			}

			localStorage.removeItem('LastSnippet')
		}

		setDeleteSnippet(false)
	}

	const dl = async () => {
		HandleDelete()
	}

	const handleClick = async () => {
		setIsSure(true)
	}

	return (
		<>
			<button
				onClick={handleClick}
				className="text-red-700 ml-auto h-fit w-fit"
			>
				<Delete fontSize="medium" />
			</button>

			{isSure && (
				<div className="h-screen w-screen top-0 z-50 left-0 fixed backdrop-blur-lg">
					<div className="h-full w-full flex items-center justify-center">
						<div className="bg-neutral-950 border border-white border-opacity-15 h-fit w-96 p-5 gap-1 rounded-lg flex flex-col items-center">
							<p className="text-white text-2xl font-bold  antialiased">
								Are you sure?
							</p>
							<p className="text-white text-opacity-50 w-full text-center">
								Do you really want to delete this snippet
							</p>
							<div className="w-full flex flex-row mt-3 gap-2 px-2">
								<button
									onClick={() => {
										setIsSure(false)
										setDeleteSnippet(false)
									}}
									className="bg-neutral-900 w-full text-white  py-2 rounded-lg duration-300 ease-out will-change-contents active:scale-95 hover:bg-neutral-800"
								>
									Cancel
								</button>
								<button
									onClick={() => {
										setIsSure(false)
										dl()
									}}
									className="bg-red-700 w-full text-white  py-2 rounded-lg duration-300 ease-out will-change-contents active:scale-95 hover:bg-red-600"
								>
									Delete
								</button>
							</div>
						</div>
					</div>
				</div>
			)}
		</>
	)
}

interface UserContentProps {
	snippets: SmallSnippets
	categories: string[]
	setSnippetToGet: React.Dispatch<React.SetStateAction<number | undefined>>
	inViewSnippet: CodeSnippet | undefined
	setAddsnippet: React.Dispatch<React.SetStateAction<boolean>>
	deleteSnippet: boolean
	setDeleteSnippet: React.Dispatch<React.SetStateAction<boolean>>
	setUpdateSnippet: React.Dispatch<React.SetStateAction<boolean>>
}

export function UserContent({
	snippets,
	categories,
	setSnippetToGet,
	inViewSnippet,
	deleteSnippet,
	setUpdateSnippet,
	setDeleteSnippet
}: UserContentProps) {
	return (
		<>
			<div className="h-full relative lg:flex hidden  overflow-hidden pb-10  w-2/12">
				<div
					id="sidebar"
					className="w-full h-full overflow-auto flex flex-col scrollbar-none gap-3"
				>
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

			<div className="h-full pb-10 w-full overflow-auto scrollbar-none">
				<div className="w-full flex flex-col">
					<div className="w-full flex flex-row items-center">
						<p className="w-11/12 line-clamp-1 h-20 select-none text-white  text-6xl font-bold">
							{inViewSnippet?.title}
						</p>
						{inViewSnippet?.id && (
							<>
								<button
									onClick={() => setUpdateSnippet(true)}
									className="text-white ml-auto h-fit w-fit"
								>
									<EditIcon fontSize="medium" />
								</button>

								<DeleteButton
									deleteSnippet={deleteSnippet}
									setDeleteSnippet={setDeleteSnippet}
									snippets={snippets}
									setSnippetToGet={setSnippetToGet}
									id={inViewSnippet.id}
								/>
							</>
						)}
					</div>

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
					{inViewSnippet && inViewSnippet?.code != '' && (
						<div className="w-full mt-10 ">
							<CodeBox
								background="bg-neutral-950"
								code={
									inViewSnippet?.code
										? inViewSnippet.code
										: ''
								}
							/>
						</div>
					)}
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
	const [deleteSnippet, setDeleteSnippet] = useState<boolean>(false)
	const [dropdownOpen, setDropdownOpen] = useState<boolean>(false)
	const [updateSnippet, setUpdateSnippet] = useState<boolean>(false)
	const [username, setUsername] = useState<string>('')
	const router = useRouter()

	useEffect(() => {
		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			router.push('/login')
		} else {
			setLoggedin(true)
		}

		// hankle refresh logi chere
		if (IsTokenExpired(token ? token : '')) {
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
	}, [addSnippet, deleteSnippet, updateSnippet])

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

		if (IsTokenExpired(token ? token : '')) {
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

	useEffect(() => {
		const req = async () => {
			const token = localStorage.getItem('ACCESS_TOKEN')
			if (!token) {
				router.push('/login')
			} else {
				setLoggedin(true)
			}

			if (IsTokenExpired(token ? token : '')) {
				console.log('Session expired')
				router.push('/login')
				return
			}

			const username = await GetUsername(token ? token : '')
			if (username === undefined) {
				return
			}

			setUsername(username)
		}

		req()
	}, [])

	const LogoutHandler = async () => {
		const token = localStorage.getItem('ACCESS_TOKEN')
		if (!token) {
			router.push('/login')
		} else {
			setLoggedin(true)
		}

		if (IsTokenExpired(token ? token : '')) {
			console.log('Session expired')
			router.push('/login')
			return
		}

		await Logout(token ? token : '')

		localStorage.removeItem('ACCESS_TOKEN')
		router.push('/login')
		setLoggedin(false)
	}

	return (
    <>
      {loggedIn && (
        <div className="flex h-full w-full flex-col gap-10 items-center">
          <div className="w-full px-10 border border-white border-opacity-15 border-l-0 border-t-0 border-r-0 h-fit py-3  flex  items-center ">
            <div className="w-fit h-full flex flex-row gap-5 items-center">
              <Link href={"/"} className="flex flex-row h-full items-center">
                <div className="h-full aspect-square">
                  <Image
                    draggable={false}
                    src={logo}
                    className="h-full w-full "
                    alt="codelet logo"
                  />
                </div>
                <p className="text-2xl select-none ml-2 text-white font-bold">
                  Codelet
                </p>
              </Link>

              <div className="h-7  w-0.5 rotate-12 bg-white bg-opacity-25" />
              <p className="text-white text-lg text-opacity-50 font-semibold tracking-wide antialiased">
                {username}'s CodeSnippets
              </p>
            </div>

            <button
              className="bg-white hover:bg-opacity-80 duration-300 ease-in-out ml-auto h-fit py-1 px-5 text-lg font-semibold rounded-lg"
              onClick={() => setAddsnippet(true)}
            >
              new snippet
            </button>
            <div className="w-fit h-full relative">
              <button
                onClick={() => setDropdownOpen((prev) => !prev)}
                className="p-1 rounded-md ml-3 relative text-white"
              >
                <MenuIcon fontSize="large" />
              </button>
              {dropdownOpen && (
                <div className="absolute mt-4 z-0 bg-black right-0 mr-2">
                  <div className="w-fit h-fit p-4 flex flex-col gap-2  border border-white border-opacity-15">
                    <DropdownItem
                      link="/"
                      title="Home"
                      subTitle="Back to home"
                      icon={<HouseIcon fontSize="large" />}
                    />
                    <DropdownItem
                      onClick={LogoutHandler}
                      title="Logout"
                      subTitle="Secure Logout portal"
                      icon={<LogoutIcon fontSize="large" />}
                    />
                  </div>
                </div>
              )}
            </div>
          </div>

          {snippets && snippets.length > 0 && (
            <div
              id="user-content"
              className="lg:w-2/3 h-full gap-5 flex flex-row justify-center"
            >
              <UserContent
                setUpdateSnippet={setUpdateSnippet}
                deleteSnippet={deleteSnippet}
                setDeleteSnippet={setDeleteSnippet}
                snippets={snippets}
                categories={categories}
                setSnippetToGet={setSnippetToGet}
                inViewSnippet={inViewSnippet}
                setAddsnippet={setAddsnippet}
              />
            </div>
          )}

          {snippets.length <= 0 && (
            <div className="p-3 bg-black border border-white border-opacity-15 rounded-lg">
              <p className="text-white font-bold">No Snippets added yet</p>
            </div>
          )}

          {addSnippet && (
            <div
              id="parent"
              className="fixed h-screen w-screen backdrop-blur-lg"
            >
              <div className="w-full h-full flex justify-center items-center">
                <div
                  onClick={(e) => {
                    e.stopPropagation();
                    setDropdownOpen(false);
                  }}
                  className=" overflow-auto scrollbar-hidden rounded-xl  p-5  bg-neutral-950"
                >
                  <SnippetForm setAddsnippet={setAddsnippet} router={router} />
                </div>
              </div>
            </div>
          )}

          {updateSnippet && (
            <div
              id="parent"
              className="fixed h-screen w-screen backdrop-blur-lg"
            >
              <div className="w-full h-full flex justify-center items-center">
                <div
                  onClick={(e) => {
                    e.stopPropagation();
                    setDropdownOpen(false);
                  }}
                  className=" overflow-auto scrollbar-hidden rounded-xl  p-5  bg-neutral-950"
                >
                  <UpdateSnippetForm
                    snippet={
                      inViewSnippet ? inViewSnippet : ({} as CodeSnippet)
                    }
                    setUpdateSnippet={setUpdateSnippet}
                    router={router}
                  />
                </div>
              </div>
            </div>
          )}
        </div>
      )}
    </>
  );
}
