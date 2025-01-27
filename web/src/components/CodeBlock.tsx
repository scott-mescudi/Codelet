'use client'

import { useEffect, useState } from "react"
import { Prism as SyntaxHighlighter } from 'react-syntax-highlighter';
import { oneDark } from 'react-syntax-highlighter/dist/esm/styles/prism'; 

interface codeBoxProps {
    code: string,
    fileName: string,
    extension: string
    style: string,
}

export function CodeBox({ code, fileName, extension, style }: codeBoxProps) {
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

    const customStyle = {
        lineHeight: '1.25rem',
        fontSize: '0.875rem',
        borderRadius: '0',
        margin: 'none',
        height: '100%',
        background: 'transparent',
    };


    return (
        <>
        <div className={`w-full h-full group flex flex-col ${style}`}>

            <div className="sm:w-full relative w-full aspect-video rounded-b-xl scrollbar-track-[black] scrollbar-thumb-neutral-800 scrollbar-thin overflow-auto  bg-[#0b0e14]">
                <div id="copy" className="absolute z-10 mt-3 mr-5 w-fit h-fit opacity-0 group-hover:opacity-100 duration-200 ease-in-out top-0 right-0 bg-black bg-opacity-50 rounded-lg text-[white]">
                    <button onClick={() => handleClick(code)} className="text-sm px-3 py-1">{copied ? "Copied!" : "Copy"}</button>
                </div>
                <SyntaxHighlighter language="go" style={oneDark} customStyle={customStyle} >
                    {code}
                </SyntaxHighlighter>
            </div>
        </div>

        </>
    )
}