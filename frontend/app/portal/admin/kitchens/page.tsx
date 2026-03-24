export default function AdminKitchensPage() {
  return (
    <div>
      <h1 className="text-2xl font-bold text-ink-950">Kitchen management</h1>
      <p className="mt-2 text-slate-600">
        Create kitchens via energy-management POST /api/v1/kitchens with admin credentials; UI form
        can post through the BFF proxy using an admin JWT.
      </p>
    </div>
  );
}
