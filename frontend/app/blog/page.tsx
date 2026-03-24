import { SiteFooter } from "@/components/SiteFooter";
import { SiteHeader } from "@/components/SiteHeader";

export default function BlogPage() {
  return (
    <>
      <SiteHeader />
      <main className="mx-auto max-w-3xl px-4 py-14">
        <h1 className="text-3xl font-bold text-ink-950">Insights</h1>
        <p className="mt-4 text-slate-700">
          Field notes, regulatory updates, and technical deep dives will appear here. Subscribe via{" "}
          <a className="text-brand-700 hover:underline" href="mailto:sales@borelsigma.com">
            sales@borelsigma.com
          </a>{" "}
          to be notified.
        </p>
      </main>
      <SiteFooter />
    </>
  );
}
