import Link from "next/link";
import { PitchTable } from "@/components/PitchTable";
import { SiteFooter } from "@/components/SiteFooter";
import { SiteHeader } from "@/components/SiteHeader";
import { UptimeTicker } from "@/components/UptimeTicker";
import type { EnergySnapshot } from "@/lib/api";
import { serverServiceBase } from "@/lib/service-url";

async function loadSnapshot(): Promise<EnergySnapshot> {
  const base = serverServiceBase("energy");
  if (!base) {
    return {
      uptime_percent: 98.7,
      tco2e_avoided: 2840,
      opex_reduction_percent: 31,
      patent_pipeline_count: 12,
    };
  }
  try {
    const res = await fetch(`${base}/api/v1/public/snapshot`, {
      next: { revalidate: 30 },
    });
    if (!res.ok) {
      throw new Error("snapshot");
    }
    return res.json() as Promise<EnergySnapshot>;
  } catch {
    return {
      uptime_percent: 98.7,
      tco2e_avoided: 2840,
      opex_reduction_percent: 31,
      patent_pipeline_count: 12,
    };
  }
}

export default async function HomePage() {
  const snap = await loadSnapshot();

  return (
    <>
      <SiteHeader />
      <main>
        <section className="border-b border-slate-200 bg-gradient-to-b from-brand-50 to-slate-50">
          <div className="mx-auto max-w-6xl px-4 py-16 md:py-24">
            <p className="text-sm font-semibold uppercase tracking-widest text-brand-800">
              Institutional food-energy cloud
            </p>
            <h1 className="mt-3 max-w-3xl text-4xl font-bold tracking-tight text-ink-950 md:text-5xl">
              Powering resilient kitchens that feed cities without fragility.
            </h1>
            <p className="mt-6 max-w-2xl text-lg text-slate-700">
              Borel Sigma unifies tri-modal energy orchestration, vendor financing, and IoT
              compliance so public and private food programs scale with auditable ESG outcomes.
            </p>
            <div className="mt-8 flex flex-wrap gap-3">
              <Link
                href="/api/auth/login?returnTo=/portal"
                className="rounded-xl bg-brand-700 px-5 py-3 text-sm font-semibold text-white shadow hover:bg-brand-900"
              >
                Client / operator access
              </Link>
              <Link
                href="/api/auth/login?returnTo=/portal"
                className="rounded-xl border border-slate-300 bg-white px-5 py-3 text-sm font-semibold text-ink-900 hover:border-brand-700"
              >
                Vendor partner access
              </Link>
              <Link
                href="/contact"
                className="rounded-xl px-5 py-3 text-sm font-semibold text-brand-800 underline-offset-4 hover:underline"
              >
                Talk to sales
              </Link>
            </div>
          </div>
        </section>

        <section className="mx-auto max-w-6xl px-4 py-12">
          <h2 className="text-2xl font-bold text-ink-950">Key metrics</h2>
          <p className="mt-2 max-w-3xl text-slate-600">
            Narrative targets backed by platform telemetry; public snapshot refreshes from the
            energy-management service when reachable.
          </p>
          <div className="mt-8 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            {[
              {
                label: "Target uptime",
                value: `${snap.uptime_percent}%`,
                hint: "SLA-grade hybrid stacks",
              },
              {
                label: "tCO2e avoided (program)",
                value: String(snap.tco2e_avoided),
                hint: "Cohort-level ledger",
              },
              {
                label: "Opex reduction",
                value: `${snap.opex_reduction_percent}%`,
                hint: "Vs. diesel-first baseline",
              },
              {
                label: "Patent pipeline",
                value: String(snap.patent_pipeline_count ?? 12),
                hint: "Controls and thermal IP",
              },
            ].map((m) => (
              <div
                key={m.label}
                className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm"
              >
                <p className="text-sm text-slate-500">{m.label}</p>
                <p className="mt-2 text-3xl font-bold text-ink-950">{m.value}</p>
                <p className="mt-1 text-xs text-slate-600">{m.hint}</p>
              </div>
            ))}
          </div>
          <p className="mt-6 text-sm">
            <UptimeTicker />
          </p>
        </section>

        <section className="border-y border-slate-200 bg-white py-14">
          <div className="mx-auto max-w-6xl px-4">
            <h2 className="text-2xl font-bold text-ink-950">Investor diligence matrix</h2>
            <p className="mt-2 max-w-3xl text-slate-600">
              Twenty rows mapped to authoritative references. Each source opens in a new tab for IC
              review.
            </p>
            <div className="mt-8">
              <PitchTable />
            </div>
          </div>
        </section>

        <section className="mx-auto max-w-6xl px-4 py-14">
          <h2 className="text-2xl font-bold text-ink-950">Why institutions choose Borel Sigma</h2>
          <ul className="mt-6 grid gap-6 md:grid-cols-3">
            <li className="rounded-2xl border border-slate-200 bg-slate-50 p-6">
              <h3 className="font-semibold text-ink-950">Patent-backed control plane</h3>
              <p className="mt-2 text-sm text-slate-600">
                Deterministic handoffs across grid, battery, and solar with explainable ML overlays
                for load and demand forecasting.
              </p>
            </li>
            <li className="rounded-2xl border border-slate-200 bg-slate-50 p-6">
              <h3 className="font-semibold text-ink-950">Vendor retention economics</h3>
              <p className="mt-2 text-sm text-slate-600">
                Zero-interest advances with transparent repayment tied to meal throughput, surfaced
                in the vendor portal.
              </p>
            </li>
            <li className="rounded-2xl border border-slate-200 bg-slate-50 p-6">
              <h3 className="font-semibold text-ink-950">Compliance you can export</h3>
              <p className="mt-2 text-sm text-slate-600">
                IoT ingestion, alert acknowledgement, and monthly PDF-ready reports for clients and
                regulators.
              </p>
            </li>
          </ul>
        </section>
      </main>
      <SiteFooter />
    </>
  );
}
