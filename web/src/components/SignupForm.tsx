import Link from "next/link";

interface RegisterFormProps {
  username: string;
  email: string;
  password: string;
  href?: string
  setUsername: React.Dispatch<React.SetStateAction<string>>;
  setEmail: React.Dispatch<React.SetStateAction<string>>;
  setPassword: React.Dispatch<React.SetStateAction<string>>;
  onSubmit?: (e: React.FormEvent<HTMLFormElement>) => void;
}

export function RegisterForm({
  username,
  email,
  password,
  setUsername,
  setEmail,
  setPassword,
  onSubmit,
  href
}: RegisterFormProps) {
  return (
    <>
      <div className="mx-3 p-4">
        <form className="flex gap-3 flex-col" onSubmit={onSubmit}>
          <p className="text-3xl font-bold text-white w-full text-center">
            Sign Up to codelet
          </p>
          <input
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl"
            type="text"
            placeholder="Username"
          ></input>
          <input
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="text-white bg-black border transition-all  duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl"
            type="email"
            placeholder="email"
          ></input>
          <input
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="text-white bg-black border transition-all duration-300 ease-in-out border-white border-opacity-15 px-6 py-4 text-lg rounded-xl"
            type="password"
            placeholder="password"
          ></input>
          <button
            type="submit"
            className="bg-white px-6 py-4 text-lg rounded-xl hover:bg-opacity-80 ease-in-out duration-300"
          >
            Sign Up
          </button>
          <Link
            href={href ? href : "#"}
            className="text-blue-600 w-full text-center hover:underline"
          >
            Already have a account? Login
          </Link>
        </form>
      </div>
    </>
  );
}
