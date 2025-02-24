export function InfiniteScrollAnimationPage() {
  return (
    <div className="w-full flex border border-t-white/5 border-x-0 border-b-0 overflow-hidden relative ">

      <ul className="flex animate-infinite-scroll gap-10 z-0  py-4 text-white">
        {[...Logos, ...Logos].map((logo, idx) => {
          return (
            <li key={idx} className="flex items-center gap-2">
              <div className="h-14 aspect-square ">
                <img
                draggable="false"
                  className="h-full w-full flex items-center justify-center"
                  src={logo.href}
                  alt={logo.alt}
                />
              </div>
            </li>
          );
        })}
      </ul>
    </div>
  );
}


const Logos = [
  {
    alt: "go",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/go/go-original.svg",
  },
  {
    alt: "javascript",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/javascript/javascript-original.svg",
  },
  {
    alt: "python",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/python/python-original.svg",
  },
  {
    alt: "java",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/java/java-original.svg",
  },
  {
    alt: "c",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/c/c-original.svg",
  },
  {
    alt: "cpp",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/cplusplus/cplusplus-original.svg",
  },
  {
    alt: "csharp",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/csharp/csharp-original.svg",
  },
  {
    alt: "php",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/php/php-original.svg",
  },
  {
    alt: "ruby",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/ruby/ruby-original.svg",
  },
  {
    alt: "swift",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/swift/swift-original.svg",
  },
  {
    alt: "kotlin",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/kotlin/kotlin-original.svg",
  },
  {
    alt: "typescript",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/typescript/typescript-original.svg",
  },
  {
    alt: "rust",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/rust/rust-original.svg",
  },
  {
    alt: "dart",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/dart/dart-original.svg",
  },
  {
    alt: "perl",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/perl/perl-original.svg",
  },
  {
    alt: "scala",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/scala/scala-original.svg",
  },
  {
    alt: "lua",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/lua/lua-original.svg",
  },
  {
    alt: "haskell",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/haskell/haskell-original.svg",
  },
  {
    alt: "r",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/r/r-original.svg",
  },
  {
    alt: "elixir",
    href: "https://cdn.jsdelivr.net/gh/devicons/devicon@latest/icons/elixir/elixir-original.svg",
  },
];
