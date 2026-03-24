"use client";

import useSWR from "swr";
import { usePortalUser } from "@/components/PortalUserContext";
import { vendorApi } from "@/lib/api";

export default function VendorTransactionsPage() {
  const me = usePortalUser();
  const vendorId = me?.rbac?.user?.vendor_id || "";
  const { data } = useSWR(vendorId ? ["vtx", vendorId] : null, () =>
    vendorApi.listTransactions(vendorId, 500),
  );

  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Transactions</h1>
      <p className="mt-2 text-slate-600">Meal throughput tied to financing repayment.</p>
      {!vendorId && <p className="mt-4 text-amber-800">vendor_id required.</p>}
      <div className="mt-6 overflow-x-auto rounded-xl border border-slate-200 bg-white">
        <table className="min-w-full text-left text-sm">
          <thead className="bg-slate-100 text-xs uppercase text-slate-600">
            <tr>
              <th className="px-3 py-2">Kitchen</th>
              <th className="px-3 py-2">Amount</th>
              <th className="px-3 py-2">Meals</th>
            </tr>
          </thead>
          <tbody>
            {(data?.items || []).map((t, i) => (
              <tr key={i} className="border-t border-slate-100">
                <td className="px-3 py-2 font-mono text-xs">{t.kitchen_id}</td>
                <td className="px-3 py-2">{t.amount}</td>
                <td className="px-3 py-2">{t.meal_count ?? "—"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
