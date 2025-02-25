"use client";

import { useEffect, useState } from "react";

export function Loader() {
  const chars = ["|", "/", "-", "\\"];
  const [index, setIndex] = useState(0);

  useEffect(() => {
    const interval = setInterval(() => {
      setIndex((prevIndex) => (prevIndex + 1) % chars.length);
    }, 100);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="text-4xl text-white select-none transition-all duration-300 ease-in-out opacity-50">
      {chars[index]}
    </div>
  );
}
