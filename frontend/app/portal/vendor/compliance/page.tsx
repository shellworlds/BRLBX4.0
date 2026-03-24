"use client";

import useSWR from "swr";
import { Line, LineChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { iotApi } from "@/lib/api";

export default function VendorCompliancePage() {
  const { data, mutate } = useSWR("alerts", () => iotApi.alerts());

  const chartData = [20, 21, 22, 23, 24, 25].map((t, i) => ({
    t,
    temp: 4 + Math.sin(i) * 0.4,
    hum: 55 + i,
  }));

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-ink-950">Compliance</h1>
        <p className="mt-2 text-slate-600">IoT alerts from ingestion service.</p>
      </div>
      <section className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
        <h2 className="text-lg font-semibold">Sensor history (illustrative)</h2>
        <div className="mt-4 h-56">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="t" />
              <YAxis />
              <Tooltip />
              <Line type="monotone" dataKey="temp" stroke="#0ea5e9" dot={false} name="Temp C" />
              <Line type="monotone" dataKey="hum" stroke="#64748b" dot={false} name="RH%" />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold">Active alerts</h2>
        <ul className="mt-3 space-y-2">
          {(data?.items || []).map((a) => (
            <li
              key={a.id}
              className="flex items-center justify-between rounded-lg border border-slate-200 bg-white px-4 py-3 text-sm"
            >
              <span>{a.message || a.id}</span>
              <button
                type="button"
                className="text-brand-700 hover:underline"
                onClick={async () => {
                  await iotApi.ackAlert(a.id);
                  void mutate();
                }}
              >
                Ack
              </button>
            </li>
          ))}
        </ul>
      </section>
    </div>
  );
}
