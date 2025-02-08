"use client";

import { LoginForm } from "@/components/LoginForm";
import { useState } from "react";1


export default function LoginPage() {
  const [email, setEmail] = useState<string>("");
  const [password, setPassword] = useState<string>("");

  
  return (
    <div className="flex min-h-screen items-center justify-center ">
      <LoginForm href="/signup" email={email} password={password} setEmail={setEmail} setPassword={setPassword} />
    </div>
  );
}
