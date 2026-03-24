"use client";

import { usePortalUser } from "@/components/PortalUserContext";

export default function ClientSettingsPage() {
  const me = usePortalUser();
  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Settings</h1>
      <p className="mt-2 text-slate-600">
        Profile, billing, and contract renewal hooks connect to vendor-ecosystem or a dedicated
        billing service in production.
      </p>
      <div className="mt-6 rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="font-semibold text-ink-950">Profile</h2>
        <pre className="mt-3 overflow-x-auto rounded-lg bg-slate-900 p-4 text-xs text-slate-100">
          {JSON.stringify(me, null, 2)}
        </pre>
        <h2 className="mt-8 font-semibold text-ink-950">Renewal</h2>
        <p className="mt-2 text-sm text-slate-600">
          Stub: POST to your billing orchestration with the next contract term. No charges are
          initiated from this UI in dev.
        </p>
        <button
          type="button"
          className="mt-4 rounded-xl bg-slate-200 px-4 py-2 text-sm font-medium text-slate-700"
          disabled
        >
          Request renewal (coming soon)
        </button>
      </div>
    </div>
  );
}
