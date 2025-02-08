import Link from "next/link";

interface FormProps {
  email: string;
  password: string;
  onSubmit?: (e: React.FormEvent<HTMLFormElement>) => void;
  setEmail: (value: string) => void;
  href?: string
  setPassword: (value: string) => void;
}

export function LoginForm({
  email,
  password,
  onSubmit,
  setEmail,
  href,
  setPassword,
}: FormProps) {
  return (
    <>
      <div className="mx-3 p-4">
        <form className="flex gap-3 flex-col" onSubmit={onSubmit}>
          <p className="text-3xl font-bold text-white w-full text-center">
            Log in to codelet
          </p>
          <input
            value={email}
            required
            onChange={(e) => setEmail(e.target.value)}
            className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl"
            type="email"
            placeholder="email"
          ></input>
          <input
            value={password}
            required
            onChange={(e) => setPassword(e.target.value)}
            className="text-white bg-black border transition-all duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl"
            type="password"
            placeholder="password"
          ></input>
          <button
            type="submit"
            className="bg-white px-6 py-4 text-lg rounded-xl hover:bg-opacity-80 ease-in-out duration-300"
          >
            Login
          </button>
          <Link
            href={href ? href : "#"}
            className="text-blue-600 w-full text-center hover:underline"
          >
            Don't have an account? Sign Up
          </Link>
        </form>
      </div>
    </>
  );
}
