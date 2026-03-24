import Link from "next/link";
import { PortalUserProvider } from "@/components/PortalUserContext";

export default function PortalLayout({ children }: { children: React.ReactNode }) {
  return (
    <PortalUserProvider>
      <div className="min-h-screen bg-slate-50">
        <div className="flex items-center justify-between border-b border-slate-200 bg-white px-4 py-3">
          <Link href="/" className="font-semibold text-ink-950">
            Borel Sigma
          </Link>
          <div className="flex gap-4 text-sm">
            <Link href="/portal" className="text-slate-600 hover:text-brand-700">
              Portal home
            </Link>
            <a href="/api/auth/logout" className="text-slate-600 hover:text-brand-700">
              Sign out
            </a>
          </div>
        </div>
        {children}
      </div>
    </PortalUserProvider>
  );
}
