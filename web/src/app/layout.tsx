import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Codelet - Your Code Snippet Manager",
  description:
    "Codelet is the ultimate code snippet manager for developers. Organize, store, and access your snippets easily.",
  icons: {
    icon: "/logo.svg",
  },
  keywords:
    "code snippets, developer tools, code storage, programming, code sharing",
  openGraph: {
    title: "Codelet - Your Code Snippet Manager",
    description:
      "Organize, store, and access your code snippets with ease. A must-have tool for developers.",
    url: "https://codelet-mu.vercel.app/",
    siteName: "Codelet",
    images: [
      {
        url: "https://codelet-mu.vercel.app/seo.png",
        width: 1200,
        height: 630,
        alt: "Codelet - Your Code Snippet Manager",
      },
    ],
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "Codelet - Your Code Snippet Manager",
    description: "Organize, store, and access your code snippets easily.",
    images: ["https://codelet-mu.vercel.app/seo.png"],
  },
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" className="overflow-y-scroll">
      <head>
        <meta name="robots" content="index, follow" />
      </head>
      <body className="overflow-hidden flex flex-col items-center">
        {children}
      </body>
    </html>
  );
}
