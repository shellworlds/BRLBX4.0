import { SiteFooter } from "@/components/SiteFooter";
import { SiteHeader } from "@/components/SiteHeader";
import { ContactForm } from "./ContactForm";

export default function ContactPage() {
  return (
    <>
      <SiteHeader />
      <main className="mx-auto max-w-3xl px-4 py-14">
        <h1 className="text-3xl font-bold text-ink-950">Contact</h1>
        <p className="mt-4 text-slate-700">
          Reach the Borel Sigma sales desk for pilots, diligence rooms, and partnership structuring.
          Direct email:{" "}
          <a className="text-brand-700 hover:underline" href="mailto:sales@borelsigma.com">
            sales@borelsigma.com
          </a>
          .
        </p>
        <ContactForm />
      </main>
      <SiteFooter />
    </>
  );
}
