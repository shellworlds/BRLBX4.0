"use client";

import useSWR from "swr";
import type { EnergySnapshot } from "@/lib/api";

export function UptimeTicker() {
  const { data, error, isLoading } = useSWR<EnergySnapshot>(
    "/api/public/energy-snapshot",
    {
      refreshInterval: 30_000,
    },
  );

  if (isLoading) {
    return <span className="text-slate-500">Loading live metrics…</span>;
  }
  if (error || !data) {
    return <span className="text-amber-700">Using cached narrative metrics</span>;
  }

  return (
    <span className="font-mono text-brand-800">
      Live snapshot: {data.uptime_percent}% uptime · {data.tco2e_avoided} tCO2e avoided ·{" "}
      {data.opex_reduction_percent}% opex reduction (vs. baseline narrative)
    </span>
  );
}
