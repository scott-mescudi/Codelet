"use client";

import { RegisterForm } from "@/components/SignupForm";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Login, Signup } from "@/shared/api/UserApiReq";
import "./signup.css";
import { Loader } from "@/components/Loader";

export default function RegisterPage() {
  const [username, setUsername] = useState<string>("");
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [err, setErr] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();

  const submit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLoading(true);
    if (email === "" || username === "" || password === "") {
      setLoading(false);
      setErr("Please fill in all fields");
      return;
    }

    let code = await Signup(username, email, password);
    if (code !== 500 && code !== 200) {
      setLoading(false);
      setErr("A user with that email already exists");
      return;
    }

    if (code === 500) {
      setLoading(false);
      setErr("Server is down");
      return;
    }

    code = await Login(email, password);
    if (code === 200) {
      setLoading(false);
      router.push("/dashboard");
      return;
    }

    if (code === 500) {
      setLoading(false);
      setErr("Server is down");
      return;
    }

    setLoading(false);
  };

  useEffect(() => {
    if (err) {
      const t = setTimeout(() => {
        setErr("");
      }, 3000);
      return () => {
        clearTimeout(t);
      };
    }
  }, [err]);

  return (
    <div className="w-full h-[calc(100vh-5rem)] flex items-center justify-center">
      <div className="flex flex-col w-fit h-fit items-center justify-center ">
        {!loading && (
          <RegisterForm
            onSubmit={submit}
            href="/login"
            username={username}
            email={email}
            password={password}
            setEmail={setEmail}
            setPassword={setPassword}
            setUsername={setUsername}
          />
        )}

        {loading && <Loader />}

        <p
          className={`text-red-700 min-h-[24px] ${
            err === "" ? "opacity-0" : "wiggle"
          }`}
        >
          {err}
        </p>
      </div>
    </div>
  );
}
