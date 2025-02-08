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

async function Signup(
  username: string,
	email: string,
	password: string
) {
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
      const error = await resp.json() as ErrorResponse
      console.error(error)
      return
    }

    console.log(resp.json())
	} catch (err) {console.error(err)}

}

export default function RegisterPage() {
	const [username, setUsername] = useState<string>('')
	const [email, setEmail] = useState<string>('')
	const [password, setPassword] = useState<string>('')


	return (
		<div className="flex min-h-screen items-center justify-center ">
			<RegisterForm
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
