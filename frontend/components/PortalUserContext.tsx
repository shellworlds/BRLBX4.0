"use client";

import { createContext, useContext } from "react";
import useSWR from "swr";

export type MeResponse = {
  role: string;
  rbac: {
    user?: {
      role?: string;
      client_id?: string | null;
      vendor_id?: string | null;
    };
  } | null;
};

const Ctx = createContext<MeResponse | undefined>(undefined);

export function PortalUserProvider({ children }: { children: React.ReactNode }) {
  const { data } = useSWR<MeResponse>("/api/me");
  return <Ctx.Provider value={data}>{children}</Ctx.Provider>;
}

export function usePortalUser() {
  const v = useContext(Ctx);
  return v;
}
