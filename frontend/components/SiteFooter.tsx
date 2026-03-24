import Link from "next/link";

export function SiteFooter() {
  return (
    <footer className="border-t border-slate-200 bg-white">
      <div className="mx-auto flex max-w-6xl flex-col gap-4 px-4 py-10 text-sm text-slate-600 md:flex-row md:items-center md:justify-between">
        <p>© {new Date().getFullYear()} Borel Sigma. Institutional food–energy infrastructure.</p>
        <div className="flex gap-4">
          <Link href="/contact" className="hover:text-brand-700">
            sales@borelsigma.com
          </Link>
          <a href="/api/auth/logout" className="hover:text-brand-700">
            Sign out
          </a>
        </div>
      </div>
    </footer>
  );
}
