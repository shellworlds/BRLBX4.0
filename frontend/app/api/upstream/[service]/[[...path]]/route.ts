import { NextRequest, NextResponse } from "next/server";
import { getAccessToken } from "@auth0/nextjs-auth0";

export const dynamic = "force-dynamic";

const allowed = new Set(["energy", "vendor", "iot", "auth", "ml"]);

function baseFor(service: string): string | undefined {
  const map: Record<string, string | undefined> = {
    energy: process.env.ENERGY_SERVICE_URL,
    vendor: process.env.VENDOR_SERVICE_URL,
    iot: process.env.IOT_SERVICE_URL,
    auth: process.env.AUTH_SERVICE_URL,
    ml: process.env.ML_SERVICE_URL,
  };
  const b = map[service];
  return b?.replace(/\/$/, "");
}

async function proxy(
  req: NextRequest,
  service: string,
  pathParts: string[] | undefined,
) {
  if (!allowed.has(service)) {
    return NextResponse.json({ error: "unknown service" }, { status: 400 });
  }
  let base = baseFor(service);
  if (!base && process.env.NODE_ENV === "development") {
    const origin =
      process.env.AUTH0_BASE_URL?.replace(/\/$/, "") || "http://localhost:3000";
    base = `${origin}/dev-proxy/${service}`;
  }
  if (!base) {
    return NextResponse.json(
      { error: "service URL not configured" },
      { status: 502 },
    );
  }

  const tail = pathParts?.length ? pathParts.join("/") : "";
  const path = tail ? `/${tail}` : "";
  const target = `${base}${path}${req.nextUrl.search}`;

  let accessToken: string | undefined;
  const publicPath =
    path.includes("/public/") || path === "/healthz" || path === "/metrics";
  if (!publicPath) {
    try {
      const t = await getAccessToken();
      accessToken = t?.accessToken;
    } catch {
      accessToken = undefined;
    }
    if (!accessToken) {
      return NextResponse.json({ error: "unauthorized" }, { status: 401 });
    }
  }

  const headers = new Headers();
  const accept = req.headers.get("accept");
  if (accept) {
    headers.set("Accept", accept);
  }
  const ct = req.headers.get("content-type");
  if (ct) {
    headers.set("Content-Type", ct);
  }
  if (accessToken) {
    headers.set("Authorization", `Bearer ${accessToken}`);
  }

  const init: RequestInit = {
    method: req.method,
    headers,
    cache: "no-store",
  };
  if (req.method !== "GET" && req.method !== "HEAD") {
    init.body = await req.text();
  }

  const res = await fetch(target, init);
  const body = await res.arrayBuffer();
  const out = new NextResponse(body, { status: res.status });
  const outCt = res.headers.get("content-type");
  if (outCt) {
    out.headers.set("content-type", outCt);
  }
  return out;
}

type RouteCtx = { params: { service: string; path?: string[] } };

export async function GET(req: NextRequest, ctx: RouteCtx) {
  return proxy(req, ctx.params.service, ctx.params.path);
}

export async function POST(req: NextRequest, ctx: RouteCtx) {
  return proxy(req, ctx.params.service, ctx.params.path);
}

export async function PUT(req: NextRequest, ctx: RouteCtx) {
  return proxy(req, ctx.params.service, ctx.params.path);
}

export async function PATCH(req: NextRequest, ctx: RouteCtx) {
  return proxy(req, ctx.params.service, ctx.params.path);
}

export async function DELETE(req: NextRequest, ctx: RouteCtx) {
  return proxy(req, ctx.params.service, ctx.params.path);
}
