'use client'

import {LoginForm} from '@/components/LoginForm'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'
import "./login.css"


interface LoginRequest {
	email: string
	password: string
}

interface LoginResponse {
	access_token: string
}

async function Login(email: string, password: string): Promise<boolean> {
	const data: LoginRequest = {
		email,
		password
	}

	try {
		const resp = await fetch('http://localhost:3021/api/v1/login', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		})

		if (!resp.ok) {
			return false
		}

		const token = (await resp.json()) as LoginResponse
		localStorage.setItem('ACCESS_TOKEN', token.access_token)

		return true
	} catch (err) {
		return false
	}
}

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
		setErr("invalid Email or Password")
	}

	useEffect(() => {
		if (err) {
			const t = setTimeout(() => {
				setErr("")
			}, 3000);
			return () => {clearTimeout(t)}
		}
	}, [err])

	return (
		<div className="flex flex-col min-h-screen items-center justify-center ">
			<LoginForm
				onSubmit={submit}
				href="/signup"
				email={email}
				password={password}
				setEmail={setEmail}
				setPassword={setPassword}
			/>
			<p className={`text-red-700 min-h-[24px] ${err === '' ? 'opacity-0' : 'wiggle'}`}>{err}</p>
		</div>
	)
}
