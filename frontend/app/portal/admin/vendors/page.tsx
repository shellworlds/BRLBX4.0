export default function AdminVendorsPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Vendor management</h1>
      <p className="mt-2 text-slate-600">
        Approve financing, suspend accounts, and reconcile FSSAI scores. Wire to vendor-ecosystem
        admin routes when RBAC allows.
      </p>
      <ul className="mt-6 list-disc space-y-2 pl-6 text-sm text-slate-700">
        <li>Export vendor ledger for audit.</li>
        <li>Bulk notify vendors below hygiene threshold.</li>
      </ul>
    </div>
  );
}
