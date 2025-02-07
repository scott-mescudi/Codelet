"use client";

import { Sidebar } from "@/components/Sidebar";
import { Snippet } from "@/components/SmallSnippet";
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
    <div className="px-3  mt-10 min-h-screen  flex flex-row gap-5 sm:w-2/3 w-full">
      <div className="w-2/12 min-h-full max-h-full overflow-y-auto flex flex-col gap-3">
        {langs.map((lang, idx) => (
          <Sidebar key={idx} title={lang}>
            {snippets.map((snippet: CodeSnippet, idx:number) => (
              snippet.language === lang ? (
                <p className="text-white py-1 w-full border hover:border-opacity-100 border-l-white border-r-0 border-t-0 border-b-0 border-opacity-15 pl-5 text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden">{snippet.title}</p>
              ) : (null)
            ))}
          </Sidebar>
        ))}
      </div>

      <div className="w-10/12  min-h-full max-h-full">
        <div className="w-full gap-5 justify-center flex flex-col">
          {snippets.map((snippet, idx) => (
            <Snippet key={idx} language={snippet.language} title={snippet.title} description={snippet.description} tags={snippet.tags} idx={idx}/>
          ))}
        </div>
      </div>
    </div>
  );
}