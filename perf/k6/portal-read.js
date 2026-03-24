import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  vus: 30,
  duration: "2m",
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<400"],
  },
};

const BASE = __ENV.BASE_URL || "http://127.0.0.1:8080";

export default function () {
  const res = http.get(`${BASE}/api/v1/public/snapshot`);
  check(res, { "snapshot ok": (r) => r.status === 200 });
  sleep(0.5);
}
