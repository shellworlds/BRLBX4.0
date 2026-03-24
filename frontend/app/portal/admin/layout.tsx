import Link from "next/link";

const links = [
  { href: "/portal/admin/dashboard", label: "Overview" },
  { href: "/portal/admin/vendors", label: "Vendors" },
  { href: "/portal/admin/kitchens", label: "Kitchens" },
  { href: "/portal/admin/alerts", label: "Alerts" },
  { href: "/portal/admin/health", label: "System health" },
];

export default function AdminPortalLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="mx-auto max-w-6xl px-4 py-6">
      <p className="text-sm font-semibold uppercase tracking-wide text-amber-800">Admin portal</p>
      <nav className="mt-4 flex flex-wrap gap-2 border-b border-slate-200 pb-3">
        {links.map((l) => (
          <Link
            key={l.href}
            href={l.href}
            className="rounded-lg px-3 py-1.5 text-sm text-slate-700 hover:bg-amber-50 hover:text-amber-950"
          >
            {l.label}
          </Link>
        ))}
      </nav>
      <div className="mt-6">{children}</div>
    </div>
  );
}
