"use client";

import { UserProvider } from "@auth0/nextjs-auth0/client";
import { SWRConfig } from "swr";
import { Toaster } from "sonner";

async function jsonFetcher(url: string) {
  const res = await fetch(url, { credentials: "include" });
  if (!res.ok) {
    throw new Error(`HTTP ${res.status}`);
  }
  return res.json();
}

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <UserProvider>
      <SWRConfig value={{ fetcher: jsonFetcher }}>
        {children}
        <Toaster richColors position="top-right" />
      </SWRConfig>
    </UserProvider>
  );
}
