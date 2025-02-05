'use client'

import Image from "next/image";
import { ReactNode } from "react";
import { useState, useRef, useEffect, useCallback } from "react";

import CloseIcon from '@mui/icons-material/Close';
import SearchIcon from '@mui/icons-material/Search';
import { AccountBoxOutlined, ExitToApp, LoginOutlined, Home } from "@mui/icons-material";

import Link from "next/link";
import { useRouter } from "next/navigation"; 

interface SearchBarProps {
    inputValue: string
    setInputValue: (value: string) => void
}


interface DropdownItemProps {
    title: string
    subTitle: string
    iconBgColor?: string
    iconColor?:string 
    link?: string
    setClick: (value: boolean) => void
    iconBgHoverColor?: string
    iconHoverColor?: string
    icon: ReactNode;
    onClick?: () => void

}

function DropdownItem({title, link="#", onClick, setClick, subTitle, iconBgColor="bg-black", iconColor="text-white", iconBgHoverColor="group-hover:bg-white", iconHoverColor="group-hover:text-black", icon}:DropdownItemProps) {
    return (
        <>
            <div onClick={() => {if (onClick) onClick();}} className="h-12 group flex flex-row rounded-md w-52 gap-3">
                <div className={`h-full aspect-square flex items-center justify-center ${iconBgColor} rounded-sm ${iconColor} ${iconHoverColor} duration-300 ease-out border border-white border-opacity-15 ${iconBgHoverColor}`}>
                    {icon}
                </div>
                <Link href={link} onClick={() => {setClick(false)}}  className="w-full h-full flex flex-col overflow-hidden">
                    <h1 className="text-lg  truncate text-white font-semibold">{title}</h1>
                    <p  className="text-sm text-white text-opacity-50 group-hover:text-opacity-100 duration-200 ease-in-out overflow-hidden">{subTitle}</p>
                </Link>
            </div>
        </>
    )
}

function ProfilePicture() {
    const [click, setClick] = useState<boolean>(false)
    const [loggedIn, setLoggedIn] = useState<boolean>(false)

    const router = useRouter()

   useEffect(() => {
        const key = localStorage.getItem("ACCESS_TOKEN")
        if (key) {
            setLoggedIn(true)
        }
    })


    useEffect(()=>{
        const handleKeyPress = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                setClick(false);
            }
        };

        document.addEventListener("keydown", handleKeyPress);
        return () => {
            document.removeEventListener("keydown", handleKeyPress);
        }

    }, [])

    const logout = () => {
        localStorage.removeItem("ACCESS_TOKEN")
        router.push("/login")
        window.location.reload()
    }

    return (
        <>
            <div className="h-full aspect-square relative ml-2">
                <Image onClick={() => setClick((prev) => !prev)} draggable="false" className="object-cover hover:cursor-pointer rounded-full h-full w-full" src="/navbar/pfp_placeholder.png" width={500} height={500} alt="Pofile picture" />
                {click ? (
                    <div className="p-5 mt-5 z-10 origin-right grid rounded-lg right-0 absolute bg-black border border-opacity-15 border-white">
                        <div className="flex flex-col gap-2">
                            <DropdownItem setClick={setClick} link="/" icon={<Home fontSize="large" />} title="Home" subTitle="Explore public snippets" />
                            {!loggedIn && <DropdownItem setClick={setClick} link="/login" icon={ <LoginOutlined fontSize="large" />}  title="Login" subTitle="Secure Login Portal" /> }
                            {loggedIn && <DropdownItem setClick={setClick} link="/dashboard" icon={ <AccountBoxOutlined fontSize="large" />}  title="Profile" subTitle="Your Dashboard" />}
                            {loggedIn && <DropdownItem setClick={setClick} onClick={logout} icon={ <ExitToApp fontSize="large" />} iconHoverColor="group-hover:text-white" iconBgHoverColor="group-hover:bg-red-700" title="Logout" subTitle="Sign Out Securely" />}
                        </div>
                    </div>
                ) : null}
            </div>
        </>
    )
}

function SearchBar( {inputValue, setInputValue} : SearchBarProps) {
    const [isFocused, setIsFocused] = useState<boolean>(false)
    const searchRef = useRef<HTMLInputElement>(null)

    const handleChange = useCallback((event: React.ChangeEvent<HTMLInputElement>) => {
        setInputValue(event.target.value);
    }, [setInputValue]);

    useEffect(() => {
        const handleKeyPress = (event: KeyboardEvent) => {
            if (!isFocused && event.key === "/") {
                event.preventDefault();
                setIsFocused(true);
                searchRef.current?.focus();
            }
            if (event.key === "Escape") {
                setIsFocused(false);
                searchRef.current?.blur();
            }
        };

        const checkFocus = () => setIsFocused(document.activeElement === searchRef.current);

        document.addEventListener("keydown", handleKeyPress);
        document.addEventListener("focusin", checkFocus);
        document.addEventListener("focusout", checkFocus);

        return () => {
            document.removeEventListener("keydown", handleKeyPress);
            document.removeEventListener("focusin", checkFocus);
            document.removeEventListener("focusout", checkFocus);
        };
    }, [isFocused]);

    return (
        <>
            <div className="w-full relative">
                <div className={`h-full right-0  ${!isFocused ? "absolute" : "hidden"}`}>
                    <div className="h-full w-fit flex justify-center items-center ">
                        <div className="py-1 px-2 bg-neutral-900 rounded-xl">
                            <p className="text-white text-opacity-50 font-bold">/</p>
                        </div>
                    </div>
                </div>
                <input ref={searchRef} value={inputValue} onChange={handleChange}  type="text" className="h-full  outline-none w-full px-3 text-xl appearance-none bg-transparent text-white" placeholder="Search anything..."></input>
            </div>
            <button onClick={() => setInputValue("")}  className={`h-full aspect-square ${inputValue != "" ? "": "hidden"}`}>
                <CloseIcon fontSize="large" className="text-white opacity-50"  />
            </button>
        </>
    )

}

export function Navbar() {
    const [inputValue, setInputValue] = useState<string>("");

    return (
        <>
        <div className='sm:w-1/2 mx-3 mb-4 h-fit absolute z-50 bg-black border border-white border-opacity-15  rounded-3xl'>
            <div className={`w-full h-16 py-2 px-5  justify-center flex flex-row  `}>
                <div className="rounded-full h-full flex items-center justify-center aspect-square">
                    <SearchIcon fontSize="large" className="text-white opacity-50" />
                </div>

                <SearchBar inputValue={inputValue} setInputValue={setInputValue} />
                <ProfilePicture />
            </div>

            {inputValue && (
                <div id="items" className="max-h-96 flex">
                    <div className="p-5 overflow-auto scrollbar-none w-full space-y-2">
                        <div className="w-full h-16 hover:scale-[102%] animate-pulse rounded-xl ease-in-out duration-200 bg-neutral-800"></div>
                    </div>
                </div>
            )}

        </div>
        </>
    )
}