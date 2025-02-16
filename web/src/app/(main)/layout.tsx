import type { Metadata } from "next";
import "../globals.css";
import { MainNavbar } from "@/components/MainNavbar";

export const metadata: Metadata = {
  title: "Codelet",
  description: "Codelet",
};

export default function RootLayout({children}: Readonly<{children: React.ReactNode;}>) {
  return (
		<>
			<MainNavbar />
			{children}
		</>
  )
}
