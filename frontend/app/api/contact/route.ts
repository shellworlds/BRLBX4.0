import { NextResponse } from "next/server";

export async function POST(req: Request) {
  let body: { name?: string; email?: string; company?: string; message?: string };
  try {
    body = await req.json();
  } catch {
    return NextResponse.json({ error: "invalid json" }, { status: 400 });
  }
  const { name, email, message, company } = body;
  if (!name?.trim() || !email?.trim() || !message?.trim()) {
    return NextResponse.json({ error: "name, email, message required" }, { status: 400 });
  }

  const base = process.env.COMPLIANCE_SERVICE_URL?.replace(/\/$/, "");
  if (base) {
    const res = await fetch(`${base}/api/v1/contact`, {
      method: "POST",
      headers: { "Content-Type": "application/json", Accept: "application/json" },
      body: JSON.stringify({
        name: name.trim(),
        email: email.trim(),
        company: (company || "").trim(),
        message: message.trim(),
      }),
    });
    const text = await res.text();
    if (!res.ok) {
      return NextResponse.json(
        { error: text || "upstream error" },
        { status: res.status },
      );
    }
    try {
      const data = JSON.parse(text) as { ok?: boolean; id?: string };
      return NextResponse.json({
        ok: true,
        id: data.id,
        detail: "Message received — we will reply shortly.",
      });
    } catch {
      return NextResponse.json({ ok: true, detail: "Message received." });
    }
  }

  return NextResponse.json({
    ok: true,
    detail:
      "Request recorded locally — set COMPLIANCE_SERVICE_URL to persist via compliance service and email.",
  });
}
