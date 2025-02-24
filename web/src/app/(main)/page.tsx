import { Safari } from "@/components/Browser";
import { CodeBox } from "@/components/CodeBlock";
import { InfiniteScrollAnimationPage } from "@/components/Marquee";
import Link from "next/link";

const code = `function HelloWorld () {
  console.log("Hello Codelet!!")
}`;

const code2 = `func HelloWorld() {
  fmt.Println("Hello Codelet!!")
}`;

const code3 = `void HelloWorld () {
  printf("Hello Codelet!!")
}`;

const LinesBackground = () => {
  return (
    <>
      <div className="w-full -z-10 absolute h-full">
        <div className="w-full h-full  flex flex-row">
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
          <div className="w-full h-full border border-y-0 b border-r-white/5 border-l-0"></div>
        </div>
      </div>
    </>
  );
};

export default function Home() {
  return (
    <>
      <div className="w-full flex flex-col items-center">
        <div className="h-fit overflow-hidden w-full lg:w-2/3 relative flex flex-col justify-center border border-t-0 border-r-0 border-white/5">
          <div className="h-full w-full mb-20 mt-10 flex flex-row items-center p-5">
            <div className="h-3/4 lg:w-1/2 flex  flex-col gap-2">
              <p className="w-full text-white font-bold text-5xl">
                All Your Code Snippets,{" "}
                <span className="text-green-700">Organized</span>
              </p>
              <p className="text-white/50 text-xl  w-full">
                Stop wasting time searching. Access all your code snippets
                instantly, organized in one place for easy use.
              </p>
              <Link
                href={"/signup"}
                className="bg-yellow-100 mt-5 w-fit px-4 py-2 rounded-lg text-2xl font-bold hover:bg-yellow-500 duration-300 ease-in-out will-change-transform"
              >
                Get started
              </Link>
            </div>
            <div className="h-3/4 relative w-1/2 px-2 hidden lg:flex flex-col gap-2">
              <CodeBox background="bg-neutral-950" code={code} />
              <CodeBox background="bg-orange-950" code={code2} />
              <CodeBox background="bg-blue-950" code={code3} />
            </div>
          </div>
          <LinesBackground />
        </div>

        <div className="h-fit w-full lg:w-2/3 relative flex flex-col items-center border border-t-0  border-white/5">
          <p className="w-fit px-4 text-center border border-y-0 border-x-white/5 py-2 text-white font-bold  text-2xl sm:text-4xl">
            Syntax highlighting for all programming languages
          </p>
          <InfiniteScrollAnimationPage />
        </div>

        <div className="h-fit w-full lg:w-2/3 pb-5 flex flex-col items-center border-b-0 border border-t-0  border-white/5">
          <p className="w-fit px-4 text-center border border-y-0 border-x-white/5 py-2 text-white font-bold text-2xl sm:text-4xl">
            Simple ui
          </p>
          <div className="h-fit w-full px-2 border flex items-center justify-center border-x-0 border-b-0 border-t-white/5 py-3">
            <Safari imageSrc="/image.png" url="codelet.dev" />
          </div>
        </div>
      </div>
    </>
  );
}
