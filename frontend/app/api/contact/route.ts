import { NextResponse } from "next/server";

export async function POST(req: Request) {
  let body: { name?: string; email?: string; company?: string; message?: string };
  try {
    body = await req.json();
  } catch {
    return NextResponse.json({ error: "invalid json" }, { status: 400 });
  }
  const { name, email, message } = body;
  if (!name?.trim() || !email?.trim() || !message?.trim()) {
    return NextResponse.json({ error: "name, email, message required" }, { status: 400 });
  }
  // Production: connect to SendGrid / SES and deliver to sales@borelsigma.com
  return NextResponse.json({
    ok: true,
    detail:
      "Request recorded for sales@borelsigma.com — connect this route to your mail provider in production.",
  });
}
