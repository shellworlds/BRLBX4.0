"use client";

import useSWR from "swr";
import { Bar, BarChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { usePortalUser } from "@/components/PortalUserContext";
import { vendorApi } from "@/lib/api";

export default function VendorDashboardPage() {
  const me = usePortalUser();
  const vendorId = me?.rbac?.user?.vendor_id || "";

  const { data: vendor } = useSWR(vendorId ? ["vend", vendorId] : null, () =>
    vendorApi.getVendor(vendorId),
  );
  const { data: txs } = useSWR(vendorId ? ["txs", vendorId] : null, () =>
    vendorApi.listTransactions(vendorId, 200),
  );
  const { data: fins } = useSWR(vendorId ? ["fin", vendorId] : null, () =>
    vendorApi.listFinancing(vendorId),
  );

  const chartData = (txs?.items || []).slice(0, 14).map((t, i) => ({
    i: String(i + 1),
    amount: t.amount,
  }));

  const openFin = (fins?.items || []).filter((f) => (f.remaining_balance || 0) > 0);
  const fssai = vendor?.fssai_score ?? 0;
  const alertHygiene = fssai < 4;

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-ink-950">Dashboard</h1>
        <p className="mt-2 text-slate-600">Earnings, financing, and hygiene signals.</p>
      </div>
      {!vendorId && (
        <p className="rounded-lg border border-amber-200 bg-amber-50 p-4 text-sm text-amber-900">
          Attach vendor_id via auth-rbac sync to load live vendor APIs.
        </p>
      )}
      <section className="grid gap-4 md:grid-cols-3">
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">FSSAI hygiene score</p>
          <p className="mt-2 text-3xl font-bold text-ink-950">{fssai || "—"}</p>
          {alertHygiene && (
            <p className="mt-2 text-xs font-semibold text-red-600">
              Below threshold - inspect compliance tab.
            </p>
          )}
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">Open financing balance</p>
          <p className="mt-2 text-3xl font-bold text-ink-950">
            {openFin[0]?.remaining_balance?.toFixed?.(2) ?? "0.00"}
          </p>
        </div>
        <div className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
          <p className="text-sm text-slate-500">Recent transactions</p>
          <p className="mt-2 text-3xl font-bold text-ink-950">{txs?.items?.length ?? 0}</p>
        </div>
      </section>
      <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="text-lg font-semibold text-ink-950">Meal revenue (recent)</h2>
        <div className="mt-4 h-64">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="i" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="amount" fill="#15803d" />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </section>
    </div>
  );
}
