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
        <div className='sm:w-1/2 mx-3  h-fit bg-neutral-950 rounded-3xl'>
            <div className={`w-full h-16 py-2 px-5  justify-center flex flex-row  `}>
                <div className="rounded-full h-full flex items-center justify-center aspect-square">
                    <SearchIcon fontSize="large" className="text-white opacity-50" />
                </div>
                <div className="w-full relative">
                    <div className={`h-full aspect-square right-0 ${inputValue != "" ? "absolute": "hidden"}`}>
                        <button onClick={() => setInputValue("")} className="h-full w-full">
                            <CloseIcon fontSize="large" className="text-white opacity-50"  />
                        </button>
                    </div>
                    <input value={inputValue} onChange={handleChange}  type="text" className="h-full  outline-none w-full px-3 text-xl appearance-none bg-transparent text-white" placeholder="Search anything..."></input>
                </div>
                <div className=" h-full ml-2 aspect-square">
                    <Image className="object-cover  rounded-full h-full w-full" src="/pfp_placeholder.png" width={500} height={500} alt="Pofile picture" />
                </div>
            </div>
            <div id="items" className={`${inputValue != "" ? "flex": "hidden"} max-h-96  `}>
                <div className="p-5 overflow-auto scrollbar scrollbar-none w-full space-y-2">
                    <div className="w-full h-16 hover:scale-[102%] rounded-xl ease-in-out duration-200 bg-neutral-800"></div>

                </div>
                
                
            </div>
        </div>
            

            {/* <div className={`w-1/2 ${inputValue != "" ? "scale-y-100": "scale-y-0"} origin-top  h-96 duration-300 ease-in-out bg-neutral-950 rounded-b-3xl`}></div> */}
        </>
    )
}