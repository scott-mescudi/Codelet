import Link from "next/link"
import { ReactNode } from "react"

interface DropdownItemProps {
    title: string
    subTitle: string
    iconBgColor?: string
    iconColor?:string 
    link?: string
    iconBgHoverColor?: string
    iconHoverColor?: string
    icon: ReactNode;
    onClick?: () => void

}

export function DropdownItem({title, link="#", onClick, subTitle, iconBgColor="bg-black", iconColor="text-white", iconBgHoverColor="group-hover:bg-white", iconHoverColor="group-hover:text-black", icon}:DropdownItemProps) {
    return (
        <>
            <div onClick={() => {if (onClick) onClick();}} className="h-12 group flex flex-row rounded-md w-52 gap-3">
                <div className={`h-full aspect-square flex items-center justify-center ${iconBgColor} rounded-sm ${iconColor} ${iconHoverColor} duration-300 ease-out border border-white border-opacity-15 ${iconBgHoverColor}`}>
                    {icon}
                </div>
                <Link href={link}  className="w-full h-full flex flex-col overflow-hidden">
                    <h1 className="text-lg text-left  truncate text-white font-semibold">{title}</h1>
                    <p  className="text-sm text-left text-white text-opacity-50 group-hover:text-opacity-100 duration-200 ease-in-out overflow-hidden">{subTitle}</p>
                </Link>
                
            </div>
        </>
    )
}