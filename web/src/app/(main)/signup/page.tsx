'use client'

import {RegisterForm} from '@/components/SignupForm'
import {useState} from 'react'

interface SignupRequest {
	username: string
	email: string
	password: string
	role: string
}

interface ErrorResponse {
	error: string
	code: number
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
			const error = (await resp.json()) as ErrorResponse
			console.error(error)
			return false
		}

		return true
	} catch (err) {
		console.error(err)
		return false
	}
}

async function Login(
	email: string,
	password: string
): Promise<boolean> {
	const data: LoginRequest = {
		email,
		password,
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
			const error = (await resp.json()) as ErrorResponse
			console.error(error)
			return false
		}

		const token = await resp.json() as LoginResponse
		localStorage.setItem("ACCESS_TOKEN", token.access_token)
		
		return true
	} catch (err) {
		console.error(err)
		return false
	}
}

export default function RegisterPage() {
	const [username, setUsername] = useState<string>('')
	const [email, setEmail] = useState<string>('')
	const [password, setPassword] = useState<string>('')

	const submit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		let success = await Signup(username, email, password)
		// add err handling logic here
		if (!success) console.log("Failed to create account")
		success = await Login(email, password)
	}

	return (
		<div className="flex min-h-screen items-center justify-center ">
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
		</div>
	)
}
