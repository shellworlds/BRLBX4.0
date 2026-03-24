"use client";

import useSWR from "swr";
import { usePortalUser } from "@/components/PortalUserContext";
import { energyApi } from "@/lib/api";

export default function ClientReportsPage() {
  const me = usePortalUser();
  const clientId = me?.rbac?.user?.client_id || "";
  const to = new Date();
  const from = new Date();
  from.setMonth(from.getMonth() - 3);
  const fromS = from.toISOString().slice(0, 10);
  const toS = to.toISOString().slice(0, 10);

  const { data, error, isLoading } = useSWR(
    clientId ? ["creports", clientId] : null,
    () => energyApi.clientReports(clientId, fromS, toS),
  );

  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Reports</h1>
      <p className="mt-2 text-slate-600">
        Download-ready JSON today; PDF viewer can wrap the same payload in a print stylesheet.
      </p>
      {!clientId && (
        <p className="mt-4 text-sm text-amber-800">Set client_id on your user to list reports.</p>
      )}
      {isLoading && <p className="mt-4 text-slate-500">Loading…</p>}
      {error && <p className="mt-4 text-red-600">Failed to load reports.</p>}
      <ul className="mt-6 space-y-3">
        {(data?.items || []).map((row: unknown, i: number) => (
          <li
            key={i}
            className="flex flex-col gap-2 rounded-xl border border-slate-200 bg-white p-4 shadow-sm md:flex-row md:items-center md:justify-between"
          >
            <span className="font-mono text-xs text-slate-600">Row {i + 1}</span>
            <a
              className="text-sm font-semibold text-brand-700 hover:underline"
              href={
                "data:application/json;charset=utf-8," +
                encodeURIComponent(JSON.stringify(row, null, 2))
              }
              download={"borel-report-" + String(i + 1) + ".json"}
            >
              Download JSON
            </a>
          </li>
        ))}
      </ul>
    </div>
  );
}
