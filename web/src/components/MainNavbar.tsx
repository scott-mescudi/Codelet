"use client";

import Image from "next/image";
import Link from "next/link";
import logo from "../../public/logo.svg";
import { jwtDecode } from "jwt-decode";
import { useEffect, useState } from "react";
import GitHubIcon from "@mui/icons-material/GitHub";

function isTokenExpired(token: string): boolean {
  try {
    const decoded: { exp?: number } = jwtDecode(token);
    const now = Date.now() / 1000;

    return decoded.exp !== undefined ? decoded.exp < now : true;
  } catch (error) {
    console.error(error);
    return true;
  }
}

export function MainNavbar() {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);

  useEffect(() => {
    const token = localStorage.getItem("ACCESS_TOKEN");
    if (!token) return;

    if (!isTokenExpired(token)) {
      setIsLoggedIn(true);
    }
  }, []);

  return (
    <>
      <div className="w-full sm:w-3/4 h-12 flex flex-row overflow-hidden items-center px-5 rounded-xl sticky top-0  backdrop-blur-lg ">
        <Link href={"/"} className="flex flex-row h-full items-center">
          <div className="h-full aspect-square">
            <Image
              draggable={false}
              src={logo}
              className="h-full w-full "
              alt="codelet logo"
            />
          </div>
          <p className="text-2xl select-none ml-2 text-white font-bold">
            Codelet
          </p>
        </Link>
        <div className="w-full h-full flex flex-row gap-2 justify-end items-center">
          <Link
          aria-label="github"
            id="github"
            href={"https://github.com/scott-mescudi/Codelet"}
            target="_blank"
            className="h-full aspect-square flex items-center text-white justify-center"
          >
            <GitHubIcon />
          </Link>

          {isLoggedIn ? (
            <Link
              href={"/dashboard"}
              className="bg-white hover:bg-opacity-80 duration-300 ease-in-out  h-fit py-1 px-5 text-lg font-semibold rounded-lg"
            >
              Dashboard
            </Link>
          ) : (
            <Link
              href={"/signup"}
              className="bg-white hover:bg-opacity-80 duration-300 ease-in-out  h-fit py-1 px-5 text-lg font-semibold rounded-lg"
            >
              Signup
            </Link>
          )}
        </div>
      </div>
    </>
  );
}
