import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 20,
  duration: "2m",
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<600"],
  },
};

const BASE = __ENV.VENDOR_BASE_URL || __ENV.BASE_URL || "http://127.0.0.1:8081";
const VENDOR = __ENV.VENDOR_ID || "00000000-0000-0000-0000-000000000002";
const KITCHEN = __ENV.KITCHEN_ID || "00000000-0000-0000-0000-000000000003";

export default function () {
  const url = `${BASE}/api/v1/transactions`;
  const body = JSON.stringify({
    vendor_id: VENDOR,
    kitchen_id: KITCHEN,
    amount: 10 + Math.random() * 5,
    meal_count: 20,
  });
  const res = http.post(url, body, {
    headers: { "Content-Type": "application/json" },
  });
  check(res, { "created or ok": (r) => r.status === 201 || r.status === 200 });
  sleep(0.3);
}
