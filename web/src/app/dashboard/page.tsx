"use client";

import { useRouter } from "next/navigation";
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

async function fetchSnippets(
  key: string
): Promise<CodeSnippets | ErrorResponse> {
  try {
    const resp = await fetch(
      "http://localhost:3021/api/v1/user/snippets?page=1&limit=100",
      {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: key,
        },
      }
    );

    if (!resp.ok) {
      const errorResponse = (await resp.json()) as ErrorResponse;
      return errorResponse;
    }

    return (await resp.json()) as CodeSnippets;
  } catch (error) {
    return { error: "Failed to fetch data", code: 500 };
  }
}

interface SnippetSectionProps {
  snippets: CodeSnippet[];
  langs: string[];
  inViewSnippet: CodeSnippet | null;
  setInViewSnippet: (snippet: CodeSnippet) => void;
}

export function SnippetSection({
  snippets,
  langs,
  inViewSnippet,
  setInViewSnippet,
}: SnippetSectionProps) {
  return (
    <>
      <div className="w-2/12 flex flex-col gap-3">
        {langs.map((lang, idx) => (
          <Sidebar key={idx} title={lang}>
            {snippets.map((snippet: CodeSnippet, idx: number) =>
              snippet.language === lang ? (
                <p
                  key={idx}
                  onClick={() => setInViewSnippet(snippet)}
                  className={`text-white py-1 w-full border ${
                    inViewSnippet?.title === snippet.title
                      ? "border-opacity-100"
                      : "hover:border-opacity-100 border-opacity-15 text-opacity-50 hover:text-opacity-100"
                  } border-l-white border-r-0 border-t-0 border-b-0 pl-5  duration-300 ease-in-out hover:cursor-pointer text-nowrap text-ellipsis overflow-hidden`}
                >
                  {snippet.title}
                </p>
              ) : null
            )}
          </Sidebar>
        ))}
      </div>

      <div className="w-10/12  min-h-full flex flex-col gap-1 text-white">
        <p className="w-full select-none pl-1 text-white text-left text-xl font-bold text-opacity-50 font-mono">
          {inViewSnippet?.language.toLowerCase()}
        </p>
        <p className="w-full select-none text-white text-left text-6xl font-bold ">
          {inViewSnippet?.title}
        </p>
        <div className="w-full select-none flex gap-5 mt-2 flex-row">
          {inViewSnippet?.tags.map((tag: string, idx: number) => (
            <p
              key={idx}
              className="text-white text-nowrap w-fit text-opacity-50 px-5 rounded-lg py-0.5 bg-neutral-800"
            >
              {tag}
            </p>
          ))}
        </div>
        <p className="w-full pl-1 text-white text-left mt-4 text-opacity-50">
          {inViewSnippet?.description}
        </p>
        <div className="w-full mt-10">
          <CodeBox code={inViewSnippet?.code ? inViewSnippet.code : ""} />
        </div>
      </div>
    </>
  );
}


export default function DashBoard() {
  const [loggedIn, setLogedin] = useState<boolean>(false);
  const router = useRouter();
  const [snippets, setSnippets] = useState<CodeSnippets | null>(null);
  const [langs, setLangs] = useState<string[]>([]);
  const [inViewSnippet, setInViewSnippet] = useState<CodeSnippet | null>(null);
  const [addSnippet, setAddSnippet] = useState<boolean>(false)

  useEffect(() => {
    const key = localStorage.getItem("ACCESS_TOKEN");
    if (!key) {
      setLogedin(true);
      router.push("/login");
    }

    async function getSnippets() {
      const data = await fetchSnippets(key ? key : "");
      if ("error" in data) {
        console.error(data.error);
        return;
      }

      setSnippets(data);
      setLangs([...new Set(data.map((snippet) => snippet.language))]);
    }
    getSnippets();
  }, [router]);

  useEffect(() => {
    if (!addSnippet) {
      setAddSnippet(false)
      return
    };

    return


    setAddSnippet(false)
  }, [addSnippet])


  return (
    <>
      {!loggedIn && (
        <div className="px-3 mt-10 relative flex flex-col gap-5 sm:w-2/3 w-full">
          {/* <div className="w-full  bg-black px-2 items-center  flex flex-rowh-20 rounded-xl">
            <button
              onClick={() => setAddSnippet((prev) => !prev)}
              className="bg-white ml-auto px-6 py-4 text-lg rounded-xl font-bold hover:bg-opacity-80 ease-in-out duration-300"
            >
              Add Snippet
            </button>
          </div> */}
          {!snippets ? (
            <div></div>
          ) : (
            <SnippetSection
              snippets={snippets || []}
              langs={langs}
              inViewSnippet={inViewSnippet}
              setInViewSnippet={setInViewSnippet}
            />
          )}
        </div>
      )}
{/* 
      {addSnippet && <div className="w-full h-fukk absolute bg-white"></div>} */}
    </>
  );
}
