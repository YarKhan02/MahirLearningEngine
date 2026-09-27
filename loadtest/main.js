import { student, admin } from "./journeys.js";
import { warmup } from "./lib.js";

// Runs once before any VUs start: wake a cold/asleep instance so the one-time
// boot latency isn't charged to the load test.
export function setup() {
  warmup();
}

// Which test profile to run: TEST=smoke|baseline|load|heavy|stress|spike|soak
const TEST = (__ENV.TEST || "smoke").toLowerCase();

// Real traffic is ~90% students, ~10% admins. Split total VUs accordingly,
// always keeping at least one admin.
function split(total) {
  const adminVUs = Math.max(1, Math.round(total * 0.1));
  const studentVUs = Math.max(1, total - adminVUs);
  return { studentVUs, adminVUs };
}

// Ramp VUs up over ~20s, then hold at target for `duration`. The ramp staggers
// logins so they don't all fire at t=0 (a thundering login herd), which is both
// more realistic and lets the run reach steady state instead of stalling on auth.
function constant(total, duration) {
  const { studentVUs, adminVUs } = split(total);
  const ramp = "20s";
  return {
    student_flow: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [{ duration: ramp, target: studentVUs }, { duration, target: studentVUs }],
      gracefulRampDown: "10s",
      exec: "student",
    },
    admin_flow: {
      executor: "ramping-vus",
      startVUs: 0,
      stages: [{ duration: ramp, target: adminVUs }, { duration, target: adminVUs }],
      gracefulRampDown: "10s",
      exec: "admin",
    },
  };
}

function ramping(peak, stages) {
  const s = (frac) => stages.map((st) => ({ duration: st.duration, target: Math.max(0, Math.round(st.target * frac)) }));
  return {
    student_flow: { executor: "ramping-vus", startVUs: 0, stages: s(0.9), exec: "student" },
    admin_flow: { executor: "ramping-vus", startVUs: 0, stages: s(0.1), exec: "admin" },
  };
}

function scenariosFor(test) {
  switch (test) {
    case "smoke":
      return constant(2, "30s");
    case "baseline":
      return constant(10, "2m");
    case "load":
      return constant(50, "5m");
    case "heavy":
      return constant(100, "5m");
    case "soak":
      return constant(40, "45m");
    case "stress":
      return ramping(200, [
        { duration: "1m", target: 50 },
        { duration: "2m", target: 200 },
        { duration: "2m", target: 200 },
        { duration: "1m", target: 0 },
      ]);
    case "spike":
      return ramping(100, [
        { duration: "10s", target: 0 },
        { duration: "20s", target: 100 },
        { duration: "1m", target: 100 },
        { duration: "20s", target: 0 },
      ]);
    default:
      throw new Error(`unknown TEST="${test}" (use smoke|baseline|load|heavy|stress|spike|soak)`);
  }
}

export const options = {
  scenarios: scenariosFor(TEST),
  thresholds: {
    // Tune these to your SLOs. A run "fails" (non-zero exit) if any are breached.
    http_req_failed: ["rate<0.01"], // <1% of requests error
    http_req_duration: ["p(95)<800", "p(99)<1500"], // latency budget (ms)
    checks: ["rate>0.99"], // >99% of assertions pass
  },
};

// k6 requires the scenario `exec` targets to be exported from the entry script.
export { student, admin };
