import Link from "next/link";
import { SiteFooter } from "@/components/SiteFooter";
import { SiteHeader } from "@/components/SiteHeader";

export default function AboutPage() {
  return (
    <>
      <SiteHeader />
      <main className="mx-auto max-w-3xl px-4 py-14">
        <h1 className="text-3xl font-bold text-ink-950">Mission</h1>
        <p className="mt-4 text-slate-700">
          Borel Sigma exists to make institutional food infrastructure boringly reliable: energy
          that stays on, capital that reaches vendors, and compliance signals that travel with every
          meal.
        </p>
        <h2 className="mt-10 text-xl font-semibold text-ink-950">Technology</h2>
        <p className="mt-3 text-slate-700">
          Our stack couples Go microservices (energy, vendor, IoT, auth/RBAC, ML predictor) with a
          Next.js experience layer. Auth0 issues role-aware tokens; Kubernetes and ArgoCD carry
          configurations from commit to cluster.
        </p>
        <h2 className="mt-10 text-xl font-semibold text-ink-950">Team</h2>
        <p className="mt-3 text-slate-700">
          Operators, power systems engineers, and product builders working across India and global
          corridors. We hire for field empathy and measurable outcomes, not slide decks.
        </p>
        <p className="mt-10">
          <Link href="/contact" className="text-brand-700 hover:underline">
            Partner with us
          </Link>
        </p>
      </main>
      <SiteFooter />
    </>
  );
}
