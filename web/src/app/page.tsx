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

const logoMap: { [key: string]: string } = {
  go: "https://devicon-website.vercel.app/api/go/original.svg",
  ruby: "https://devicon-website.vercel.app/api/ruby/original.svg",
  python: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/python/python-original.svg",
  javascript: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/javascript/javascript-original.svg",
  swift: "https://devicon-website.vercel.app/api/swift/original.svg",
  php: "https://devicon-website.vercel.app/api/php/original.svg",
  java: "https://devicon-website.vercel.app/api/java/original.svg",
  c: "https://devicon-website.vercel.app/api/c/original.svg",
  "c#": "https://devicon-website.vercel.app/api/csharp/original.svg",
  rust: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/rust/rust-original.svg",
  zig: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/zig/zig-original.svg",
  docker: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/docker/docker-original.svg",
  json: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/json/json-original.svg",
  css: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/css3/css3-original.svg",
  html: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/html5/html5-original.svg",
  react: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/react/react-original.svg",
  typescript: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/typescript/typescript-original.svg",
};


type CodeSnippets  = CodeSnippet[]

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



export default async function Home() {
  const snippets = await fetchSnippets();

  if ("error" in snippets) {
    return <div>Error: {snippets.error} (Code: {snippets.code})</div>;
  }

  const getIcon = (name:string) => {
    const logo = logoMap[name.toLowerCase()]
    if (!logo) {
      return "/fallback.svg"
    }

    return logo
  }

  return (
    <>
      <div className="px-3 sm:w-2/3 mt-10 flex flex-col gap-5 w-full h-full">
        <div className="w-full h-fit lg:grid pb-10 lg:gap-10 2xl:grid-cols-3 grid-cols-2 gap-5 justify-center  flex flex-col">
          {snippets.map((snippet, idx) => (
            <div key={idx} className="w-full h-40 select-none hover:scale-105 duration-300 hover:cursor-pointer ease-in-out will-change-transform rounded-lg bg-black hover:bg-neutral-900 border gap-3 border-white border-opacity-15 p-3 flex flex-row">
              <div className="w-1/5 aspect-square rounded-md overflow-hidden">
                <img className="h-full w-full" src={getIcon(snippet.language)}></img>
              </div>
              <div className="w-4/5 h-full gap-1 flex flex-col">
                <div className="w-full pt-2 h-3/4 ">
                  <p className="w-full line-clamp-1 truncate text-ellipsis overflow-hidden whitespace-nowrap text-white font-bold text-2xl">{snippet.title}</p>
                  <p className="w-full line-clamp-2 text-white text-opacity-50 ">{snippet.description}</p>
                </div>
                <div className="w-full h-1/4 flex flex-row gap-4 items-center overflow-hidden">
                  {snippet.tags.slice(0, 2).map((tag, idx) => (
                    <p key={idx} className="text-white text-nowrap text-opacity-50 px-5 rounded-lg py-0.5 bg-neutral-800">{tag}</p>
                  ))}
                </div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </>
  );
}
