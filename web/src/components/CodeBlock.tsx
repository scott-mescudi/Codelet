'use client'

import { useEffect, useState } from "react"

interface codeBoxProps {
    code: string,
    fileName: string,
}

export function CodeBox({ code, fileName }: codeBoxProps) {
    const [copied, setCopied] = useState<boolean>(false)

    const handleClick = (text:string) =>{
        navigator.clipboard.writeText(text)
        setCopied(true)
    }

    useEffect(() => {
        const tt = setTimeout(() => {
            setCopied(false)
        }, 1000);

        return () => {clearTimeout(tt)}
    }, [copied])

    return (
        <>
        <div className="h-fit w-fit group flex flex-col">
            <div className="bg-neutral-800 relative w-full gap-2 flex flex-row h-10 rounded-t-xl">
                <div className=" w-fit flex px-2 py-2 flex-row gap-3 h-full">
                    <div className="size-6 rounded-full bg-red-700"></div>
                    <div className="size-6 rounded-full bg-green-700"></div>
                    <div className="size-6 rounded-full bg-yellow-500"></div>
                </div>
                <div className="w-fit bg-[#0b0e14] pr-3 rounded-t-lg mt-1 gap-1 overflow-hidden h-full flex flex-row justify-start items-center ">
                    <div className='h-full p-2'>
                        <img className='h-full' src="https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/react/react-original.svg" />
                    </div>
                    <p className="text-xl font-semibold text-white">{fileName}</p>
                </div>
                <div id="copy" className="absolute z-10 mt-12 mr-5 w-fit h-fit opacity-0 group-hover:opacity-100 duration-200 ease-in-out top-0 right-0 bg-black bg-opacity-50 rounded-lg text-[white]">
                    <button onClick={() => handleClick(code)} className="text-sm px-3 py-1">{copied ? "Copied!" : "Copy"}</button>
                </div>
            </div>
            <div className="sm:w-[50dvw] w-[80dvw] relative aspect-video rounded-b-xl scrollbar-track-[black] scrollbar-thumb-neutral-800 scrollbar-thin overflow-auto px-8  bg-[#0b0e14]">

                <pre className="text-white text-md">
                    {code}
                </pre>
            </div>
        </div>

        </>
    )
}