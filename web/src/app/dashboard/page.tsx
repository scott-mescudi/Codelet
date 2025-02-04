'use client'

import { useRouter } from "next/navigation";
import { useEffect, useState } from 'react';

export default function DashBoard() {
    const [loggedIn, setLogedin]= useState<boolean>(false)
    const router = useRouter()
    
    
   useEffect(() => {
        const key = localStorage.getItem("ACCESS_TOKEN")
        if (!key) {
            setLogedin(true)
            router.push("/login");
        }
    })

    return (
        <>
        {!loggedIn &&         
        <div className="bg-neutral-900 w-full h-full">
            {/* Your dashboard content here */}
        </div>}
        </>
    );
}
