import Link from "next/link";

const links = [
  { href: "/portal/client/dashboard", label: "Dashboard" },
  { href: "/portal/client/reports", label: "Reports" },
  { href: "/portal/client/kitchens", label: "Kitchens" },
  { href: "/portal/client/settings", label: "Settings" },
];

export default function ClientPortalLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="mx-auto max-w-6xl px-4 py-6">
      <p className="text-sm font-semibold uppercase tracking-wide text-brand-800">Client portal</p>
      <nav className="mt-4 flex flex-wrap gap-2 border-b border-slate-200 pb-3">
        {links.map((l) => (
          <Link
            key={l.href}
            href={l.href}
            className="rounded-lg px-3 py-1.5 text-sm text-slate-700 hover:bg-brand-50 hover:text-brand-900"
          >
            {l.label}
          </Link>
        ))}
      </nav>
      <div className="mt-6">{children}</div>
    </div>
  );
}
