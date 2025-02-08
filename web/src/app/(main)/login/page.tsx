"use client";

import { LoginForm } from "@/components/LoginForm";
import { useState } from "react";1

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
			const error = (await resp.json()) as ErrorResponse
			console.error(error)
			return false
		}

		const token = (await resp.json()) as LoginResponse
		localStorage.setItem('ACCESS_TOKEN', token.access_token)

		return true
	} catch (err) {
		console.error(err)
		return false
	}
}

export default function LoginPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  	const submit = async (e: React.FormEvent<HTMLFormElement>) => {
		e.preventDefault()
		const success = await Login(email, password)
	}

  
  return (
    <div className="flex min-h-screen items-center justify-center ">
      <LoginForm onSubmit={submit} href="/signup" email={email} password={password} setEmail={setEmail} setPassword={setPassword} />
    </div>
  );
}
