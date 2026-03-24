export default function AdminHealthPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">System health</h1>
      <p className="mt-2 text-slate-600">
        Observe KEDA ScaledObjects, EMQX, and service /metrics via Prometheus/Grafana. This page is a
        checklist for operators.
      </p>
      <ul className="mt-6 space-y-3 text-sm text-slate-700">
        <li>ArgoCD application sync status for dev/staging/prod.</li>
        <li>Energy / vendor / IoT / auth / ML healthz probes green.</li>
        <li>KEDA polling intervals within SLO.</li>
      </ul>
    </div>
  );
}
