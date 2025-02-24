'use client'

import Image from "next/image";
import Link from "next/link";
import logo from '../../public/logo.svg'
import {jwtDecode} from 'jwt-decode'
import { useEffect, useState } from "react";


function isTokenExpired(token: string): boolean {
  try {
    const decoded: {exp?: number} = jwtDecode(token)
    const now = Date.now() / 1000

    return decoded.exp !== undefined ? decoded.exp < now : true
  } catch (error) {
    return true
  }
}

export function MainNavbar() {

  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false)

  useEffect(() => {
    const token = localStorage.getItem('ACCESS_TOKEN')
    if (!token) return

    if (!isTokenExpired(token)) {
      setIsLoggedIn(true)
    }
  }, [])


    return (
		<>
			<div className="w-2/3 h-12 flex flex-row overflow-hidden items-center px-5 rounded-xl sticky top-0  backdrop-blur-lg ">
				<Link href={'/'} className="flex flex-row h-full items-center">
					<div className="h-full aspect-square">
						<Image
							draggable={false}
							src={logo}
							className="h-full w-full "
							alt="codelet logo"
						/>
					</div>
				</Link>

				{isLoggedIn ? (
					<Link
						href={'/dashboard'}
						className="bg-white hover:bg-opacity-80 duration-300 ease-in-out ml-auto h-fit py-1 px-5 text-lg font-semibold rounded-lg"
					>
						Dashboard
					</Link>
				) : (
					<Link
						href={'/signup'}
						className="bg-white hover:bg-opacity-80 duration-300 ease-in-out ml-auto h-fit py-1 px-5 text-lg font-semibold rounded-lg"
					>
						Signup
					</Link>
				)}
			</div>
		</>
	)
}