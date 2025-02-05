import { useState, useEffect } from 'react';

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

const logoMap = new Map<string, string>();
logoMap.set("golang", "https://devicon-website.vercel.app/api/go/original.svg")
logoMap.set("Ruby", "https://devicon-website.vercel.app/api/ruby/original.svg")
logoMap.set("Python", "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/python/python-original.svg")
logoMap.set("JavaScript", "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/javascript/javascript-original.svg")
logoMap.set("c++", "")
logoMap.set("c", "")

type CodeSnippets = CodeSnippet[];

async function fetchSnippets(): Promise<CodeSnippets | ErrorResponse> {
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

  const data = (await resp.json()) as CodeSnippets;
  return data;
}

export default function Home() {
  const [snippets, setSnippets] = useState<CodeSnippets | ErrorResponse | null>(null);
  
  // Add your tags to filter here
  const tagsToCheck = ['JavaScript', 'Python']; // Example tags you want to check

  useEffect(() => {
    const loadSnippets = async () => {
      const fetchedSnippets = await fetchSnippets();
      setSnippets(fetchedSnippets);
    };

    loadSnippets();
  }, []);

  // Filter snippets by tags (only show snippets that contain one of the specified tags)
  const filteredSnippets = snippets && !('error' in snippets) 
    ? snippets.filter(snippet => snippet.tags.some(tag => tagsToCheck.includes(tag))) 
    : [];

  if ("error" in (snippets || {})) {
    return <div>Error: {(snippets as ErrorResponse).error} (Code: {(snippets as ErrorResponse).code})</div>
  }

  return (
    <>
      <div className="px-3 sm:w-2/3 mt-10 flex flex-col gap-5 w-full h-full">
        <div className="w-full h-fit lg:grid lg:gap-10 gap-5 justify-center grid-cols-3 flex flex-col">
          {filteredSnippets.length === 0 ? (
            <p className="text-white">No snippets match the selected tags.</p>
          ) : (
            filteredSnippets.map((snippet, idx) => (
              <div key={idx} className="w-full h-40 select-none hover:scale-105 duration-300 hover:cursor-pointer ease-in-out will-change-transform rounded-lg bg-black hover:bg-neutral-900 border gap-3 border-white border-opacity-15 p-3 flex flex-row">
                <div className="w-1/5 aspect-square rounded-md overflow-hidden"><img className="h-full w-full" src={logoMap.get(snippet.language)}></img></div>
                <div className="w-4/5 h-full gap-1 flex flex-col">
                  <div className="w-full pt-2 h-3/4 ">
                    <p className="w-full line-clamp-1 truncate text-ellipsis overflow-hidden whitespace-nowrap text-white font-bold text-2xl">{snippet.title}</p>
                    <p className="w-full line-clamp-2 text-white text-opacity-50 ">{snippet.description}</p>
                  </div>
                  <div className="w-full h-1/4 flex flex-row gap-4 items-center overflow-hidden">
                    {snippet.tags.map((tag, idx) => (
                      <p key={idx} className="text-white text-nowrap text-opacity-50 px-5 rounded-lg py-0.5 bg-neutral-800" >{tag}</p>
                    ))}
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </>
  );
}
