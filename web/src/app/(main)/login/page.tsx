"use client";

import { LoginForm } from "@/components/LoginForm";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { Login } from "@/shared/api/UserApiReq";
import "./login.css";
import { Loader } from "@/components/Loader";

export default function LoginPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");
  const [err, setErr] = useState<string>("");
  const [loading, setLoading] = useState<boolean>(false);
  const router = useRouter();

  const submit = async (e: React.FormEvent<HTMLFormElement>) => {
    setLoading(true);
    e.preventDefault();
    if (email === "" || password === "") {
      setErr("Please fill in all fields");
      return;
    }

    const code = await Login(email, password);
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
    setErr("invalid Email or Password");
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
    <>
      <div className="w-full h-[calc(100vh-5rem)] flex items-center justify-center">
        <div className="flex flex-col h-fit w-fit items-center justify-center ">
          {!loading && (
            <LoginForm
              onSubmit={submit}
              href="/signup"
              email={email}
              password={password}
              setEmail={setEmail}
              setPassword={setPassword}
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
    </>
  );
}
