import { toast } from "sonner";

export type ServiceName =
  | "energy"
  | "vendor"
  | "iot"
  | "auth"
  | "ml"
  | "payments"
  | "compliance";

/** Browser-side path to the authenticated BFF proxy. */
export function upstreamPath(service: ServiceName, apiPath: string): string {
  const p = apiPath.startsWith("/") ? apiPath : `/${apiPath}`;
  return `/api/upstream/${service}${p}`;
}

export class ApiError extends Error {
  status: number;
  body: string;

  constructor(status: number, body: string) {
    super(`API ${status}: ${body}`);
    this.status = status;
    this.body = body;
  }
}

export async function fetchUpstream(
  service: ServiceName,
  apiPath: string,
  init?: RequestInit,
): Promise<Response> {
  const url = upstreamPath(service, apiPath);
  const res = await fetch(url, {
    ...init,
    credentials: "include",
    headers: {
      Accept: "application/json",
      ...(init?.headers || {}),
    },
  });
  if (!res.ok) {
    const text = await res.text();
    toast.error(`Request failed (${res.status})`);
    throw new ApiError(res.status, text);
  }
  return res;
}

export async function fetchJSON<T>(
  service: ServiceName,
  apiPath: string,
  init?: RequestInit,
): Promise<T> {
  const res = await fetchUpstream(service, apiPath, init);
  return res.json() as Promise<T>;
}

/** Typed helpers aligned with backend Swagger routes. */
export async function fetchPublicEnergySnapshot(): Promise<EnergySnapshot> {
  const res = await fetch("/api/public/energy-snapshot", { cache: "no-store" });
  if (!res.ok) {
    throw new ApiError(res.status, await res.text());
  }
  return res.json() as Promise<EnergySnapshot>;
}

export const energyApi = {
  kitchenMetrics: (id: string) =>
    fetchJSON<unknown>("energy", `/api/v1/kitchens/${id}/metrics`),
  kitchenReadings: (id: string, from: string, to: string) =>
    fetchJSON<{ items: unknown[] }>(
      "energy",
      `/api/v1/kitchens/${id}/readings?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`,
    ),
  clientReports: (clientId: string, from: string, to: string) =>
    fetchJSON<{ items: unknown[] }>(
      "energy",
      `/api/v1/reports/client/${clientId}?from=${from}&to=${to}`,
    ),
  ghgReport: (clientId: string, from: string, to: string, region?: string) => {
    const r = region ? `&region=${encodeURIComponent(region)}` : "";
    return fetchJSON<unknown>(
      "energy",
      `/api/v1/reports/ghg?client_id=${encodeURIComponent(clientId)}&from=${from}&to=${to}${r}`,
    );
  },
  kitchensByVendor: (vendorId: string) =>
    fetchJSON<{ items: Kitchen[] }>(
      "energy",
      `/api/v1/kitchens/vendor/${vendorId}`,
    ),
};

export interface PayoutRow {
  id: string;
  vendor_id: string;
  amount: number;
  status: string;
  stripe_transfer_id?: string;
  failure_reason?: string;
  created_at?: string;
}

export const vendorApi = {
  listTransactions: (vendorId: string, limit = 100) =>
    fetchJSON<{ items: Transaction[] }>(
      "vendor",
      `/api/v1/vendors/${vendorId}/transactions?limit=${limit}`,
    ),
  requestFinancing: (vendorId: string, amount: number) =>
    fetchUpstream("vendor", `/api/v1/vendors/${vendorId}/financing`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ amount }),
    }).then((r) => r.json()),
  listFinancing: (vendorId: string) =>
    fetchJSON<{ items: Financing[] }>(
      "vendor",
      `/api/v1/vendors/${vendorId}/financing`,
    ),
  getVendor: (vendorId: string) =>
    fetchJSON<VendorRecord>("vendor", `/api/v1/vendors/${vendorId}`),
  myWallet: () => fetchJSON<unknown>("vendor", "/api/v1/vendors/me/wallet"),
  requestWalletWithdraw: (amount: number) =>
    fetchUpstream("vendor", "/api/v1/vendors/me/wallet/withdraw", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ amount }),
    }).then((r) => r.json()),
  listPayouts: () =>
    fetchJSON<{ items: PayoutRow[] }>("vendor", "/api/v1/vendors/me/payouts"),
  connectOnboarding: () =>
    fetchUpstream("vendor", "/api/v1/vendors/me/connect/onboarding", {
      method: "POST",
    }).then((r) => r.json() as Promise<{ url?: string }>),
};

export const iotApi = {
  alerts: () => fetchJSON<{ items: AlertRow[] }>("iot", "/api/v1/alerts"),
  ackAlert: (id: string) =>
    fetchUpstream("iot", "/api/v1/alerts/ack", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ id }),
    }).then((r) => r.json()),
};

export const paymentsApi = {
  mySubscription: () =>
    fetchJSON<unknown>("payments", "/api/v1/subscriptions/me"),
};

export const complianceApi = {
  consentStatus: () =>
    fetchJSON<unknown>("compliance", "/api/v1/consent/status"),
};

export const mlApi = {
  predictEnergy: (kitchenId: string, hoursAhead: number) =>
    fetchUpstream("ml", "/api/v1/predict/energy", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        kitchen_id: kitchenId,
        hours_ahead: hoursAhead,
      }),
    }).then((r) => r.json()),
};

export interface EnergySnapshot {
  uptime_percent: number;
  tco2e_avoided: number;
  opex_reduction_percent: number;
  meals_served_daily_stub?: number;
  patent_pipeline_count?: number;
  as_of?: string;
}

export interface Kitchen {
  id: string;
  name: string;
  location: string;
  vendor_id: string;
  capacity_kw: number;
  region?: string;
}

export interface Transaction {
  id?: string;
  vendor_id: string;
  kitchen_id: string;
  amount: number;
  meal_count?: number;
  created_at?: string;
}

export interface Financing {
  id?: string;
  vendor_id: string;
  amount: number;
  status: string;
  remaining_balance?: number;
  repayment_schedule?: string;
}

export interface VendorRecord {
  id: string;
  name: string;
  fssai_score: number;
  location: string;
  contact: string;
}

export interface AlertRow {
  id: string;
  kitchen_id?: string;
  severity?: string;
  message?: string;
  created_at?: string;
  acknowledged?: boolean;
}
