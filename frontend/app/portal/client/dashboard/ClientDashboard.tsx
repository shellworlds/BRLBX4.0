"use client";

import useSWR from "swr";
import {
  Area,
  AreaChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import { usePortalUser } from "@/components/PortalUserContext";
import { energyApi, fetchPublicEnergySnapshot, type EnergySnapshot } from "@/lib/api";

export function ClientDashboard() {
  const me = usePortalUser();
  const clientId = me?.rbac?.user?.client_id || "";

  const { data: snap } = useSWR<EnergySnapshot>("/api/public/energy-snapshot", () =>
    fetchPublicEnergySnapshot(),
  );

  const to = new Date();
  const from = new Date();
  from.setDate(from.getDate() - 30);
  const fromS = from.toISOString().slice(0, 10);
  const toS = to.toISOString().slice(0, 10);

  const { data: reports } = useSWR(
    clientId ? ["reports", clientId, fromS, toS] : null,
    () => energyApi.clientReports(clientId, fromS, toS),
  );

  const chartData = [0, 1, 2, 3, 4, 5, 6].map((i) => ({
    day: "D" + String(i + 1),
    lcoe: 0.12 + i * 0.01,
    baseline: 0.18,
  }));

  return (
    <div className="space-y-8">
      {!clientId && (
        <p className="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900">
          No client_id on your RBAC profile yet. Sync the user via auth-rbac /api/v1/users/sync to
          load reports.
        </p>
      )}

      <section className="grid gap-4 md:grid-cols-3">
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">Fleet uptime (snapshot)</p>
          <p className="mt-2 text-3xl font-bold text-ink-950">{snap?.uptime_percent ?? "—"}%</p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">tCO2e avoided (program)</p>
          <p className="mt-2 text-3xl font-bold text-ink-950">{snap?.tco2e_avoided ?? "—"}</p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">Opex reduction vs baseline</p>
          <p className="mt-2 text-3xl font-bold text-ink-950">
            {snap?.opex_reduction_percent ?? "—"}%
          </p>
        </div>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-ink-950">LCOE reduction curve (illustrative)</h2>
        <p className="mt-1 text-sm text-slate-600">
          Stub series; wire to daily report payloads once client_id is bound.
        </p>
        <div className="mt-4 h-64">
          <ResponsiveContainer width="100%" height="100%">
            <AreaChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="day" />
              <YAxis />
              <Tooltip />
              <Area
                type="monotone"
                dataKey="baseline"
                stackId="1"
                stroke="#94a3b8"
                fill="#e2e8f0"
                name="Baseline"
              />
              <Area
                type="monotone"
                dataKey="lcoe"
                stackId="2"
                stroke="#15803d"
                fill="#bbf7d0"
                name="Hybrid stack"
              />
            </AreaChart>
          </ResponsiveContainer>
        </div>
      </section>

      <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-ink-950">Monthly ESG reports</h2>
        <p className="mt-1 text-sm text-slate-600">
          Rows from GET /api/v1/reports/client/:id
        </p>
        <ul className="mt-4 space-y-2 text-sm">
          {(reports?.items || []).slice(0, 8).map((row: unknown, idx: number) => {
            const s = JSON.stringify(row);
            return (
              <li key={idx} className="rounded border border-slate-100 bg-slate-50 px-3 py-2 font-mono">
                {s.slice(0, 240)}
                {s.length > 240 ? "…" : ""}
              </li>
            );
          })}
          {clientId && (!reports?.items || reports.items.length === 0) && (
            <li className="text-slate-500">No rows in range (aggregate jobs may not have run).</li>
          )}
        </ul>
      </section>
    </div>
  );
}
