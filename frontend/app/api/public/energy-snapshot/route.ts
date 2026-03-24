import { NextResponse } from "next/server";
import type { EnergySnapshot } from "@/lib/api";
import { serverServiceBase } from "@/lib/service-url";

export const revalidate = 15;

export async function GET() {
  const base = serverServiceBase("energy");
  if (!base) {
    return NextResponse.json<EnergySnapshot>({
      uptime_percent: 98.7,
      tco2e_avoided: 2840,
      opex_reduction_percent: 31,
      as_of: new Date().toISOString(),
    });
  }
  try {
    const res = await fetch(`${base}/api/v1/public/snapshot`, {
      next: { revalidate: 15 },
    });
    if (!res.ok) {
      throw new Error(String(res.status));
    }
    const data = (await res.json()) as EnergySnapshot;
    return NextResponse.json(data);
  } catch {
    return NextResponse.json<EnergySnapshot>({
      uptime_percent: 98.7,
      tco2e_avoided: 2840,
      opex_reduction_percent: 31,
      as_of: new Date().toISOString(),
    });
  }
}
