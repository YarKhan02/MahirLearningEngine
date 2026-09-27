import http from "k6/http";
import { check, sleep } from "k6";
import { SharedArray } from "k6/data";
import { Trend } from "k6/metrics";

// Base URL of the API under test. NEVER point this at production.
export const BASE_URL = (__ENV.BASE_URL || "http://localhost:8080").replace(/\/$/, "");

// Accounts are loaded once and shared across all VUs (memory-efficient).
// Copy data/users.example.json -> data/users.json and fill in real credentials.
export const students = new SharedArray("students", () => JSON.parse(open("./data/users.json")).students || []);
export const admins = new SharedArray("admins", () => JSON.parse(open("./data/users.json")).admins || []);

export const loginDuration = new Trend("login_duration", true);

// Per-VU auth state: a real user logs in ONCE, then reuses the token for the
// whole session. Keyed by __VU. `failed` guards against a login storm — if the
// server is down and login keeps failing, we don't re-attempt on every request
// (that's what produced hundreds of logins when the instance OOM-crashed).
const auth = {};

function state() {
  return auth[__VU] || (auth[__VU] = { token: null, failed: false });
}

// warmup wakes a sleeping/cold instance (Render free tier) BEFORE the load
// starts, so a one-time ~10s cold boot doesn't get charged to every VU's login.
// Call once from setup(). Polls /health until healthy or it gives up.
export function warmup() {
  const url = `${BASE_URL}/health`;
  for (let i = 1; i <= 24; i++) {
    const res = http.get(url, { tags: { name: "GET /health (warmup)" } });
    if (res.status === 200) {
      console.log(`warmup: instance healthy after ${i} attempt(s)`);
      return true;
    }
    console.warn(`warmup: /health -> ${res.status} (attempt ${i}/24), waiting…`);
    sleep(5);
  }
  console.error("warmup: instance never became healthy — check the deploy/logs before trusting results");
  return false;
}

// login retries a few times (with backoff) so a transient cold-start blip
// doesn't disable a VU for the whole run — but stops after MAX_LOGIN_TRIES so a
// genuinely down server doesn't cause a login storm.
const MAX_LOGIN_TRIES = 3;

function login(user) {
  for (let attempt = 1; attempt <= MAX_LOGIN_TRIES; attempt++) {
    const res = http.post(
      `${BASE_URL}/auth/public/login`,
      JSON.stringify({ identifier: user.identifier, password: user.password }),
      { headers: { "Content-Type": "application/json" }, tags: { name: "POST /auth/public/login" } },
    );
    loginDuration.add(res.timings.duration);
    if (res.status === 200) {
      try {
        const body = res.json();
        const token = body.token || body.access_token;
        if (token) return token;
      } catch (_) {
        /* fall through to retry */
      }
    }
    if (attempt === 1) console.warn(`login ${user.identifier} -> ${res.status} (${Math.round(res.timings.duration)}ms), retrying…`);
    sleep(1 + attempt);
  }
  return null;
}

function ensureToken(user) {
  const s = state();
  if (s.token) return s.token;
  if (s.failed) return null; // don't hammer login once it's clearly broken
  s.token = login(user);
  check(s.token, { "authenticated": (t) => !!t });
  if (!s.token) s.failed = true;
  return s.token;
}

function authOpts(token, name) {
  return {
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
    tags: { name: name },
  };
}

// get / post authenticate with the VU's cached token. On a 401 they re-login
// exactly once — but only if we HAD a working token (genuine expiry), never
// while login is already failing. `name` groups the request in the summary.
export function get(user, path, name) {
  let token = ensureToken(user);
  let res = http.get(`${BASE_URL}${path}`, authOpts(token, name));
  if (res.status === 401 && !state().failed) {
    state().token = null;
    token = ensureToken(user);
    res = http.get(`${BASE_URL}${path}`, authOpts(token, name));
  }
  check(res, { [`${name} -> 2xx`]: (r) => r.status >= 200 && r.status < 300 });
  return res;
}

export function post(user, path, payload, name) {
  let token = ensureToken(user);
  let res = http.post(`${BASE_URL}${path}`, JSON.stringify(payload), authOpts(token, name));
  if (res.status === 401 && !state().failed) {
    state().token = null;
    token = ensureToken(user);
    res = http.post(`${BASE_URL}${path}`, JSON.stringify(payload), authOpts(token, name));
  }
  check(res, { [`${name} -> 2xx`]: (r) => r.status >= 200 && r.status < 300 });
  return res;
}

export function pick(arr) {
  return arr[Math.floor(Math.random() * arr.length)];
}

// Human "think time" between actions so the load looks like real usage.
export function think(minSec, maxSec) {
  sleep(minSec + Math.random() * (maxSec - minSec));
}

export function safeJson(res) {
  try {
    return res.json();
  } catch (_) {
    return null;
  }
}

// Endpoints return either a bare array or a paginated { items: [...] } object.
export function asArray(x) {
  if (Array.isArray(x)) return x;
  if (x && Array.isArray(x.items)) return x.items;
  return [];
}
