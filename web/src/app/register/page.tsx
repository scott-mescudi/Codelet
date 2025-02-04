'use client'

import Link from "next/link"


export default function Register() {
    return (
        <>
         <div className="w-full h-full flex flex-col justify-center items-center">
            <div className="mx-3 p-4">
                <form className="flex gap-3 flex-col">
                    <p className="text-3xl font-bold text-white w-full text-center">Sign Up to codelet</p>
                    <input className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="text" placeholder="Username"></input>
                    <input className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="email" placeholder="email"></input>
                    <input className="text-white bg-black border transition-all duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl" type="password" placeholder="password"></input>
                    <button type="submit" className="bg-white px-6 py-4 text-lg rounded-xl hover:bg-opacity-80 ease-in-out duration-300">Sign Up</button>
                    <Link href="/login" className="text-blue-600 w-full text-center hover:underline">Already have a account? Login</Link>
                </form>
            </div>
        </div>
        </>

    
    )
}
