'use client'

import {LoginForm} from '@/components/LoginForm'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'
import { Login } from '@/shared/api/UserApiReq'
import './login.css'

export default function LoginPage() {
	const [email, setEmail] = useState<string>('')
	const [password, setPassword] = useState<string>('')
	const [err, setErr] = useState<string>('')
	const router = useRouter()

	const submit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		if (email === '' || password === '') {
			setErr('Please fill in all fields')
			return
		}

		const success = await Login(email, password)
		if (success) {
			router.push('/dashboard')
			return
		}

		setErr('invalid Email or Password')
	}

	useEffect(() => {
		if (err) {
			const t = setTimeout(() => {
				setErr('')
			}, 3000)
			return () => {
				clearTimeout(t)
			}
		}
	}, [err])

	return (
		<>
			<div className="w-full h-full flex items-center justify-center">
				<div className="flex flex-col h-fit w-fit items-center justify-center ">
					<LoginForm
						onSubmit={submit}
						href="/signup"
						email={email}
						password={password}
						setEmail={setEmail}
						setPassword={setPassword}
					/>
					<p
						className={`text-red-700 min-h-[24px] ${
							err === '' ? 'opacity-0' : 'wiggle'
						}`}
					>
						{err}
					</p>
				</div>
			</div>
		</>
	)
}
