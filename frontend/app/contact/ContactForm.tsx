"use client";

import { useState } from "react";
import { toast } from "sonner";

export function ContactForm() {
  const [pending, setPending] = useState(false);

  async function onSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const fd = new FormData(e.currentTarget);
    const payload = {
      name: String(fd.get("name") || ""),
      email: String(fd.get("email") || ""),
      company: String(fd.get("company") || ""),
      message: String(fd.get("message") || ""),
    };
    setPending(true);
    try {
      const res = await fetch("/api/contact", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      const data = await res.json();
      if (!res.ok) {
        throw new Error(data.error || "failed");
      }
      toast.success(data.detail || "Sent");
      e.currentTarget.reset();
    } catch {
      toast.error("Could not send — try sales@borelsigma.com directly");
    } finally {
      setPending(false);
    }
  }

  return (
    <form onSubmit={onSubmit} className="mt-8 grid max-w-lg gap-4">
      <label className="grid gap-1 text-sm">
        <span className="font-medium text-slate-700">Name</span>
        <input
          name="name"
          required
          className="rounded-lg border border-slate-300 px-3 py-2"
        />
      </label>
      <label className="grid gap-1 text-sm">
        <span className="font-medium text-slate-700">Work email</span>
        <input
          name="email"
          type="email"
          required
          className="rounded-lg border border-slate-300 px-3 py-2"
        />
      </label>
      <label className="grid gap-1 text-sm">
        <span className="font-medium text-slate-700">Organization</span>
        <input name="company" className="rounded-lg border border-slate-300 px-3 py-2" />
      </label>
      <label className="grid gap-1 text-sm">
        <span className="font-medium text-slate-700">How can we help?</span>
        <textarea
          name="message"
          required
          rows={5}
          className="rounded-lg border border-slate-300 px-3 py-2"
        />
      </label>
      <button
        type="submit"
        disabled={pending}
        className="rounded-xl bg-brand-700 px-4 py-3 text-sm font-semibold text-white disabled:opacity-60"
      >
        {pending ? "Sending…" : "Send to sales@borelsigma.com"}
      </button>
    </form>
  );
}
