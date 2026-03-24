import { ClientDashboard } from "./ClientDashboard";

export default function ClientDashboardPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Dashboard</h1>
      <p className="mt-2 text-slate-600">
        Energy telemetry, ESG rollups, and cost curves for your contracted kitchens.
      </p>
      <ClientDashboard />
    </div>
  );
}
