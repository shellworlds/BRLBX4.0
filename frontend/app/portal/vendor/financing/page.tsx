"use client";

import { useState } from "react";
import useSWR from "swr";
import { toast } from "sonner";
import { usePortalUser } from "@/components/PortalUserContext";
import { vendorApi } from "@/lib/api";

export default function VendorFinancingPage() {
  const me = usePortalUser();
  const vendorId = me?.rbac?.user?.vendor_id || "";
  const [amount, setAmount] = useState("50000");
  const { data, mutate } = useSWR(vendorId ? ["vf", vendorId] : null, () =>
    vendorApi.listFinancing(vendorId),
  );

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!vendorId) {
      return;
    }
    try {
      const res = await vendorApi.requestFinancing(vendorId, Number(amount));
      toast.success("Financing request submitted");
      void mutate();
      console.info(res);
    } catch {
      toast.error("Request failed");
    }
  }

  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Financing</h1>
      <p className="mt-2 text-slate-600">Zero-interest advance workflow (vendor-ecosystem API).</p>
      <form onSubmit={submit} className="mt-6 max-w-md space-y-3 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <label className="grid gap-1 text-sm">
          <span>Amount (INR)</span>
          <input
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="rounded border border-slate-300 px-3 py-2"
          />
        </label>
        <button
          type="submit"
          className="rounded-xl bg-brand-700 px-4 py-2 text-sm font-semibold text-white"
          disabled={!vendorId}
        >
          Request advance
        </button>
      </form>
      <h2 className="mt-10 text-lg font-semibold">Repayment schedule</h2>
      <ul className="mt-3 space-y-2 text-sm">
        {(data?.items || []).map((f, i) => (
          <li key={i} className="rounded border border-slate-100 bg-slate-50 px-3 py-2">
            Status: {f.status} · Balance: {f.remaining_balance ?? "—"} · {f.repayment_schedule}
          </li>
        ))}
      </ul>
    </div>
  );
}
