'use client'
import Image from "next/image";
import { useState } from "react";

import CloseIcon from '@mui/icons-material/Close';
import SearchIcon from '@mui/icons-material/Search';

export function Navbar() {
  const [inputValue, setInputValue] = useState("");

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(event.target.value);
  };

    return (
        <>
            <div className="w-1/2 h-16 py-2 px-5 bg-neutral-950 gap-4 justify-center flex flex-row rounded-full">
                <div className="rounded-full h-full flex items-center justify-center aspect-square">
                    <SearchIcon fontSize="large" className="text-white opacity-50" />
                </div>
                <div className="w-full relative">
                    <div className={`h-full aspect-square right-0 ${inputValue != "" ? "absolute": "hidden"}`}>
                        <button onClick={() => setInputValue("")} className="h-full w-full">
                            <CloseIcon fontSize="large" className="text-white opacity-50"  />
                        </button>
                    </div>
                    <input value={inputValue} onChange={handleChange}  type="text" className="h-full outline-none w-full px-3 text-xl appearance-none bg-transparent text-white" placeholder="Search anything..."></input>
                </div>
                <div className=" h-full aspect-square">
                    <Image className="object-cover  rounded-full h-full w-full" src="/pfp_placeholder.png" width={500} height={500} alt="Pofile picture" />
                </div>
            </div>
        </>
    )
}