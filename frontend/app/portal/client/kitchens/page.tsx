"use client";

import useSWR from "swr";
import {
  Line,
  LineChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { useMemo, useState } from "react";
import { usePortalUser } from "@/components/PortalUserContext";
import { energyApi, type Kitchen } from "@/lib/api";

export default function ClientKitchensPage() {
  const me = usePortalUser();
  const vendorId = me?.rbac?.user?.vendor_id || "";
  const [selected, setSelected] = useState<Kitchen | null>(null);

  const { data } = useSWR(vendorId ? ["kitchens", vendorId] : null, () =>
    energyApi.kitchensByVendor(vendorId),
  );

  const to = useMemo(() => new Date(), []);
  const from = useMemo(() => {
    const d = new Date();
    d.setHours(d.getHours() - 24);
    return d;
  }, []);

  const { data: readings } = useSWR(
    selected ? ["readings", selected.id] : null,
    () =>
      energyApi.kitchenReadings(selected!.id, from.toISOString(), to.toISOString()),
  );

  const chartRows = (readings?.items || []).map((r: unknown, i: number) => {
    const row = r as { grid_power?: number; solar_power?: number };
    return {
      t: i,
      grid: Number(row.grid_power) || 0,
      solar: Number(row.solar_power) || 0,
    };
  });

  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Kitchens</h1>
      <p className="mt-2 text-slate-600">
        Kitchens are listed per vendor UUID via the energy-management service.
      </p>
      {!vendorId && (
        <p className="mt-4 text-sm text-amber-800">
          No vendor_id on profile. Attach vendor_id in auth-rbac to load kitchens.
        </p>
      )}
      <div className="mt-6 grid gap-4 md:grid-cols-2">
        <ul className="space-y-2">
          {(data?.items || []).map((k) => (
            <li key={k.id}>
              <button
                type="button"
                onClick={() => setSelected(k)}
                className={`w-full rounded-xl border px-4 py-3 text-left text-sm transition ${
                  selected?.id === k.id
                    ? "border-brand-700 bg-brand-50"
                    : "border-slate-200 bg-white hover:border-brand-400"
                }`}
              >
                <span className="font-semibold">{k.name}</span>
                <span className="mt-1 block text-xs text-slate-500">{k.location}</span>
              </button>
            </li>
          ))}
        </ul>
        <div className="rounded-2xl border border-slate-200 bg-white p-4 shadow-sm">
          <h2 className="text-sm font-semibold text-ink-950">24h power mix</h2>
          {selected ? (
            <div className="mt-4 h-56">
              <ResponsiveContainer width="100%" height="100%">
                <LineChart data={chartRows}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="t" />
                  <YAxis />
                  <Tooltip />
                  <Line type="monotone" dataKey="grid" stroke="#64748b" dot={false} />
                  <Line type="monotone" dataKey="solar" stroke="#16a34a" dot={false} />
                </LineChart>
              </ResponsiveContainer>
            </div>
          ) : (
            <p className="mt-4 text-sm text-slate-500">Select a kitchen for time series.</p>
          )}
        </div>
      </div>
    </div>
  );
}
