import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 50,
  duration: "2m",
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"],
  },
};

const BASE = __ENV.BASE_URL || "http://127.0.0.1:8080";
const TOKEN = __ENV.INGEST_BEARER_TOKEN || "";
const kitchen = __ENV.KITCHEN_ID || "00000000-0000-0000-0000-000000000001";

export default function () {
  const url = `${BASE}/api/v1/kitchens/${kitchen}/readings`;
  const payload = JSON.stringify({
    grid_power: 2.1 + Math.random(),
    battery_power: 0.5,
    solar_power: 1.2,
    lpg_status: "standby",
    uptime_percent: 99.5,
  });
  const headers = { "Content-Type": "application/json" };
  if (TOKEN) headers.Authorization = `Bearer ${TOKEN}`;
  const res = http.post(url, payload, { headers });
  check(res, { "2xx": (r) => r.status >= 200 && r.status < 300 });
  sleep(0.25);
}
