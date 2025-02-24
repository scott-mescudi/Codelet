'use client'

import {RegisterForm} from '@/components/SignupForm'
import {useRouter} from 'next/navigation'
import {useEffect, useState} from 'react'
import {Login, Signup} from '@/shared/api/UserApiReq'
import './signup.css'

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
    <div className="w-full h-[calc(100vh-5rem)] flex items-center justify-center">
      <div className="flex flex-col w-fit h-fit items-center justify-center ">
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
            err === "" ? "opacity-0" : "wiggle"
          }`}
        >
          {err}
        </p>
      </div>
    </div>
  );
}
