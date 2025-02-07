import { Sidebar } from "@/components/Sidebar";
import { Snippet } from "@/components/SmallSnippet";

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


type CodeSnippets  = CodeSnippet[]

async function fetchSnippets(): Promise<CodeSnippets | ErrorResponse> {
  const resp = await fetch("http://localhost:3021/api/v1/public/snippets?page=1&limit=30", {
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




export default async function Home() {
  const snippets = await fetchSnippets();

  if ("error" in snippets) {
    return <div>Error: {snippets.error} (Code: {snippets.code})</div>;
  }


  return (
    <>
      <div className="px-3 sm:w-2/3 mt-10 flex flex-row gap-5 w-full h-full">
        <Sidebar title="Components">
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
          <p className="text-white text-opacity-50 hover:text-opacity-100 duration-300 ease-in-out">hello</p>
        </Sidebar>

        <div className="w-10/12 h-fit  gap-5 justify-center  flex flex-col">
          {snippets.map((snippet, idx) => (
            <Snippet language={snippet.language} title={snippet.title} description={snippet.description} tags={snippet.tags} idx={idx} />
          ))}
        </div>
      </div>
    </>
  );
}


