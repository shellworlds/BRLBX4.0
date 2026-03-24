import Link from "next/link";

export function SiteHeader() {
  return (
    <header className="border-b border-slate-200 bg-white/90 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center justify-between gap-4 px-4 py-4">
        <Link href="/" className="text-lg font-semibold text-ink-950">
          Borel Sigma
        </Link>
        <nav className="flex flex-wrap items-center gap-4 text-sm text-slate-700">
          <Link href="/about" className="hover:text-brand-700">
            About
          </Link>
          <Link href="/contact" className="hover:text-brand-700">
            Contact
          </Link>
          <Link href="/blog" className="hover:text-brand-700">
            Insights
          </Link>
          <Link
            href="/portal"
            className="rounded-lg bg-brand-700 px-3 py-1.5 font-medium text-white hover:bg-brand-900"
          >
            Portals
          </Link>
          <a
            href="/api/auth/login?returnTo=/portal"
            className="rounded-lg border border-slate-300 px-3 py-1.5 font-medium hover:border-brand-700 hover:text-brand-700"
          >
            Log in
          </a>
        </nav>
      </div>
    </header>
  );
}
