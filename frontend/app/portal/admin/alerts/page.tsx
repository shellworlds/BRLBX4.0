"use client";

import useSWR from "swr";
import { iotApi } from "@/lib/api";

export default function AdminAlertsPage() {
  const { data, mutate } = useSWR("adm-alerts-full", () => iotApi.alerts());

  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Alert management</h1>
      <p className="mt-2 text-slate-600">Unresolved IoT alerts across kitchens.</p>
      <ul className="mt-6 space-y-2">
        {(data?.items || []).map((a) => (
          <li
            key={a.id}
            className="flex items-center justify-between rounded-xl border border-slate-200 bg-white px-4 py-3 text-sm shadow-sm"
          >
            <span>{a.message || String(a.id)}</span>
            <button
              type="button"
              className="text-brand-700 hover:underline"
              onClick={async () => {
                await iotApi.ackAlert(a.id);
                void mutate();
              }}
            >
              Acknowledge
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
}
