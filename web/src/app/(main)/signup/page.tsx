'use client'

import {RegisterForm} from '@/components/SignupForm'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'
import './signup.css'

interface SignupRequest {
	username: string
	email: string
	password: string
	role: string
}

interface LoginRequest {
	email: string
	password: string
}

interface LoginResponse {
	access_token: string
}

async function Signup(
	username: string,
	email: string,
	password: string
): Promise<boolean> {
	const data: SignupRequest = {
		username,
		email,
		password,
		role: 'user'
	}

	try {
		const resp = await fetch('http://localhost:3021/api/v1/register', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		})

		if (!resp.ok) {
			return false
		}

		return true
	} catch (err) {
		return false
	}
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



export default function RegisterPage() {
	const [username, setUsername] = useState<string>('')
	const [email, setEmail] = useState<string>('')
	const [password, setPassword] = useState<string>('')
	const [err, setErr] = useState<string>('')
	const router = useRouter()

	const submit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		if (email === '' || username === '' || password === '') {
			setErr('Please fill in all fields')
			return
		}

		let success = await Signup(username, email, password)
		if (!success) {
			setErr('A user with that email already exists')
			return
		}

		success = await Login(email, password)
		if (success) {
			router.push('/dashboard')
			return
		}
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
		<div className="flex flex-col min-h-screen items-center justify-center ">
			<RegisterForm
				onSubmit={submit}
				href="/login"
				username={username}
				email={email}
				password={password}
				setEmail={setEmail}
				setPassword={setPassword}
				setUsername={setUsername}
			/>
			<p
				className={`text-red-700 min-h-[24px] ${
					err === '' ? 'opacity-0' : 'wiggle'
				}`}
			>
				{err}
			</p>
		</div>
	)
}
