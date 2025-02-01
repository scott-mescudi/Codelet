'use client'
import Image from "next/image";
import { useState, useRef, useEffect, useCallback } from "react";

import CloseIcon from '@mui/icons-material/Close';
import SearchIcon from '@mui/icons-material/Search';

interface SearchBarProps {
    inputValue: string
    setInputValue: (value: string) => void
}


function ProfilePicture() {
    return (
        <>
            <div className=" h-full ml-2 aspect-square">
                <Image draggable="false" className="object-cover  rounded-full h-full w-full" src="/navbar/pfp_placeholder.png" width={500} height={500} alt="Pofile picture" />
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
        <div className='sm:w-1/2 mx-3  h-fit bg-neutral-950  rounded-3xl'>
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