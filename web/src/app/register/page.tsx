'use client'

import Link from "next/link"
import React, { useState } from "react"
import axios from "axios";
import { useRouter } from "next/navigation";

interface SignupRequest {
    username: string
    email: string
    password: string
    role: string
}

interface LoginRequest {
  email: string;
  password: string;
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

interface RegisterFormProps {
    username: string;
    email: string;
    password: string;
    setUsername: React.Dispatch<React.SetStateAction<string>>;
    setEmail: React.Dispatch<React.SetStateAction<string>>;
    setPassword: React.Dispatch<React.SetStateAction<string>>;
    handleSubmit: (e: React.FormEvent<HTMLFormElement>) => void;
}

export function RegisterForm({ username, email, password, setUsername, setEmail, setPassword, handleSubmit }: RegisterFormProps) {
    return (
        <>
            <div className="mx-3 p-4">
                <form className="flex gap-3 flex-col" onSubmit={handleSubmit}>
                    <p className="text-3xl font-bold text-white w-full text-center">Sign Up to codelet</p>
                    <input value={username} onChange={(e) => setUsername(e.target.value)} className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="text" placeholder="Username"></input>
                    <input value={email} onChange={(e) => setEmail(e.target.value)} className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="email" placeholder="email"></input>
                    <input value={password} onChange={(e) => setPassword(e.target.value)} className="text-white bg-black border transition-all duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="password" placeholder="password"></input>
                    <button type="submit" className="bg-white px-6 py-4 text-lg rounded-xl hover:bg-opacity-80 ease-in-out duration-300">Sign Up</button>
                    <Link href="/login" className="text-blue-600 w-full text-center hover:underline">Already have a account? Login</Link>
                </form>
            </div>
        </>
    )
}

export default function Register() {
    const [username, setUsername] = useState<string>("")
    const [email, setEmail] = useState<string>("")
    const [password, setPassword] = useState<string>("")
    const [apiErr, setApiErr] = useState<string>("")
    const [loading, setLoading] = useState<boolean>(false)
    const router = useRouter();

    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault()

        const RegisterData: SignupRequest = {
            username, email, password, role: "user"
        }

        setLoading(true)

        try {
            const resp = await axios.post("http://localhost:3021/api/v1/register", RegisterData, {
            headers: {
                "Content-Type": "application/json",
            }})

            if (resp.status === 400 || resp.status === 422) {
                throw new Error ("Failed to register")
            }

        }catch (err){
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

        const LoginRequest: LoginRequest = {
            email, password
        }

        try {

            const resp = await axios.post("http://localhost:3021/api/v1/login", LoginRequest, {
            headers: {
                "Content-Type": "application/json",
            }})

            if (resp.status === 400 || resp.status === 422) {
                throw new Error ("Invalid Email or Password")
            }

            if (resp.status === 429) {
                throw new Error ("Too many login attempts, please try again after 30 seconds")
            }

            localStorage.setItem("ACCESS_TOKEN", resp.data)
            router.push("/dashboard");

        }catch (err){
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
            <RegisterForm username={username} email={email} password={password} setUsername={setUsername} setEmail={setEmail} setPassword={setPassword} handleSubmit={handleSubmit}
            />
        </div>
        </>
    )
}
