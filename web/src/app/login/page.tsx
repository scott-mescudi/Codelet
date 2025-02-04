'use client'

import Link from "next/link";
import React, { useEffect, useState } from "react";
import axios from "axios";
import { useRouter } from "next/navigation"; 

interface LoginRequest {
  email: string;
  password: string;
}

interface FormProps {
    email: string
    password: string
    HandleData: (e: React.FormEvent<HTMLFormElement>) => void
    setEmail: (value: string) => void
    setPassword: (value: string) => void
}


export function LoginForm({email, password, HandleData, setEmail, setPassword}:FormProps) {
    return (
        <>
            <div className="mx-3 p-4">
                <form className="flex gap-3 flex-col" onSubmit={HandleData}>
                    <p className="text-3xl font-bold text-white w-full text-center">Log in to codelet</p>
                    <input value={email} required onChange={(e) => setEmail(e.target.value)} className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="email" placeholder="email"></input>
                    <input value={password} required onChange={(e) => setPassword(e.target.value)} className="text-white bg-black border transition-all duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="password" placeholder="password"></input>
                    <button type="submit" className="bg-white px-6 py-4 text-lg rounded-xl hover:bg-opacity-80 ease-in-out duration-300">Login</button>
                    <Link href="/register" className="text-blue-600 w-full text-center hover:underline">Don't have an account? Sign Up</Link>
                </form>
            </div>
        </>
    )
}

interface ErrorBoxProps {
    error: any
}
export function ErrorBox({error}: ErrorBoxProps) {
    return (
        <>
            <div className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl">{error}</div>
        </>
    )
}

export default function Login() {
    const [email, setEmail] = useState<string>("sdsd@gmail.com")
    const [password, setPassword] = useState<string>("abcd")
    const [loading, setLoading] = useState<boolean>(false)
    const [apiErr, setApiErr] = useState("")
    const [loggedIn, setLogedin]= useState<boolean>(false)

    const router = useRouter();

    useEffect(() => {
        const key = localStorage.getItem("ACCESS_TOKEN")
        if (key) {
            setLogedin(true)
            router.push("/dashboard");
        }
    })

    const HandleData = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()

    const loginData: LoginRequest = {
        email,
        password,
        };

        try {
            setLoading(true)
            const resp = await axios.post("http://localhost:3021/api/v1/login", loginData, {
            headers: {
                "Content-Type": "application/json",
            }
            })

            if (resp.status === 400 || resp.status === 422) {
                throw new Error ("Failed to login")
            }

            if (resp.status === 429) {
                throw new Error ("Too many login attempts, please try again after 30 seconds")
            }

            localStorage.setItem("ACCESS_TOKEN", resp.data)
            setLogedin(true)
            router.push("/dashboard");
        } catch (err: any) {
            if (axios.isAxiosError(err)) {
                if (err.response) {
                    if (err.response.status === 429) {
                        setApiErr("Too many login attempts, please try again after 30 seconds")
                    } else {
                        setApiErr(err.response.data?.message || "An unexpected error occurred")
                    }
                } else {
                    setApiErr("An unexpected error occurred")
                }
            } else { 
                setApiErr("An unknown error occurred")
            }
        }

        setLoading(false)
    }

  return (
    <>
      <div className="w-full h-full flex flex-col justify-center items-center">
        {loading && <div className="animate-spin"></div>}
        {!loading && apiErr === "" && !loggedIn && <LoginForm password={password} email={email} setEmail={setEmail} setPassword={setPassword} HandleData={HandleData} />}
        {apiErr !== "" && <ErrorBox error={apiErr} />}
      </div>
    </>
  )
}

