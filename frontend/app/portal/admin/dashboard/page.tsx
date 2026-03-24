"use client";

import useSWR from "swr";
import { fetchPublicEnergySnapshot } from "@/lib/api";

export default function AdminOverviewPage() {
  const { data: snap } = useSWR("adm-snap", () => fetchPublicEnergySnapshot());

  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-ink-950">Overview</h1>
      <p className="text-slate-600">
        Cross-tenant snapshot from public energy endpoint; extend with admin-only aggregates.
      </p>
      <div className="grid gap-4 md:grid-cols-3">
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">Kitchens (stub)</p>
          <p className="mt-2 text-3xl font-bold">128</p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">Vendors (stub)</p>
          <p className="mt-2 text-3xl font-bold">42</p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">Program uptime</p>
          <p className="mt-2 text-3xl font-bold">{snap?.uptime_percent ?? "—"}%</p>
        </div>
      </div>
    </div>
  );
}
