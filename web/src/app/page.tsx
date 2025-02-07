"use client";

import { CodeBox } from "@/components/CodeBlock";
import { Sidebar } from "@/components/Sidebar";

import { useEffect, useState } from "react";

interface CodeSnippet {
  language: string;
  title: string;
  code: string;
  favorite: boolean;
  private: boolean;
  tags: string[];
  description: string;
}

interface ErrorResponse {
  error: string;
  code: number;
}

type CodeSnippets = CodeSnippet[];

async function fetchSnippets(): Promise<CodeSnippets | ErrorResponse> {
  try {
    const resp = await fetch("http://localhost:3021/api/v1/public/snippets?page=1&limit=100", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!resp.ok) {
      const errorResponse = (await resp.json()) as ErrorResponse;
      return errorResponse;
    }

    return (await resp.json()) as CodeSnippets;
  } catch (error) {
    return { error: "Failed to fetch data", code: 500 };
  }
}

export default function Home() {
  const [snippets, setSnippets] = useState<CodeSnippets | null>(null);
  const [langs, setLangs] = useState<string[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [inViewSnippet, setInViewSnippet] = useState<CodeSnippet | null>(null);

  useEffect(() => {
    if (snippets && snippets.length > 0) {
      setInViewSnippet(snippets[0]);
    }
  }, [snippets]);

  useEffect(() => {
    async function getSnippets() {
      const data = await fetchSnippets();
      if ("error" in data) {
        setError(data.error);
      } else {
        setSnippets(data);
        setLangs([...new Set(data.map((snippet) => snippet.language))]); 
      }
    }
    getSnippets();
  }, []);

  if (error) {
    return <div>Error: {error}</div>;
  }

  if (!snippets) {
    return <div>Loading...</div>;
  }

  return (
    <div className="px-3 mt-10  flex flex-row gap-5 sm:w-2/3 w-full">
      <div className="w-2/12 flex flex-col gap-3">
        {langs.map((lang, idx) => (
          <Sidebar key={idx} title={lang}>
            {snippets.map((snippet: CodeSnippet, idx:number) => (
              snippet.language === lang ? (
                <p key={idx} onClick={() => setInViewSnippet(snippet)} className={`text-white py-1 w-full border ${inViewSnippet?.title === snippet.title ? ("border-opacity-100") : ("hover:border-opacity-100 border-opacity-15 text-opacity-50 hover:text-opacity-100")} border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden`}>{snippet.title}</p>
              ) : (null)
            ))}
          </Sidebar>
        ))}
      </div>

      <div className="w-10/12  min-h-full flex flex-col gap-1 text-white">
        <p className="w-full select-none pl-1 text-white text-left text-xl font-bold text-opacity-50 font-mono">{inViewSnippet?.language.toLowerCase()}</p>
        <p className="w-full select-none text-white text-left text-6xl font-bold ">{inViewSnippet?.title}</p>
        <div className="w-full select-none flex gap-5 mt-2 flex-row">
          {inViewSnippet?.tags.map((tag:string, idx:number) => (
            <p key={idx} className="text-white text-nowrap w-fit text-opacity-50 px-5 rounded-lg py-0.5 bg-neutral-800">{tag}</p>
          ))}
        </div>
        <p className="w-full pl-1 text-white text-left mt-4 text-opacity-50">{inViewSnippet?.description}</p>
        <div className="w-full mt-10">   
          <CodeBox code={inViewSnippet?.code ? inViewSnippet.code : ""} />
        </div>
      </div>
    </div>
  );
}