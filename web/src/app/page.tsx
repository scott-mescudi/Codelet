import { CodeBox } from "@/components/CodeBlock";
import Image from "next/image";

const code = `
<code className="text-white text-md">
    {code}
</code>
`

export default function Home() {
  return (
    <>
    <CodeBox fileName="main" extension=".jsx" code={code} />
    </>
  );
}
