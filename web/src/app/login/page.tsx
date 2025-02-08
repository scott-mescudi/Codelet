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

interface LoginResp {
    access_token: string
}

interface ErrorResponse {
  error: string;
  code: number;
}


async function LoginReq(loginData: LoginRequest): Promise<LoginResp | ErrorResponse> {
    try {
        const resp = await fetch("http://localhost:3021/api/v1/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(loginData),
        });

        if (!resp.ok) {
            try {
                const errorResponse = (await resp.json()) as ErrorResponse;
                return errorResponse;
            } catch {
                return { error: "Unexpected error format", code: resp.status };
            }
        }

        try {
            const loginResponse = (await resp.json()) as LoginResp;
            return loginResponse;
        } catch {
            return { error: "Invalid response format", code: 500 };
        }

    } catch (err) {
        return { error: "Failed to connect to server", code: 500 };
    }
}

export function ErrorBox({error}: ErrorBoxProps) {
    return (
        <>
            <div className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl">{error}</div>
        </>
    )
}

export default function Login() {
    const [email, setEmail] = useState<string>("")
    const [password, setPassword] = useState<string>("")
    const [loading, setLoading] = useState<boolean>(false)
    const [apiErr, setApiErr] = useState<string>("")
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


        try{
            const resp = await LoginReq(loginData)

            if ("error" in resp){
                console.error("Login failed:", resp.error);
                return
            } 

            localStorage.setItem("ACCESS_TOKEN", resp.access_token)
            setLogedin(true)
            router.push("/dashboard");
        } catch (err) {
            console.error("Unexpected error:", err);
        }

        setLoading(false)
    }

  return (
    <>
      <div className="w-full h-full mt-60 flex flex-col justify-center items-center">
        {loading && <div className="animate-spin"></div>}
        {!loading && apiErr === "" && !loggedIn && <LoginForm password={password} email={email} setEmail={setEmail} setPassword={setPassword} HandleData={HandleData} />}
        {apiErr !== "" && <ErrorBox error={apiErr} />}
      </div>
    </>
  )
}

