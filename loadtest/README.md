# k6 load testing

Realistic, role-based load tests for the MahirLearning API. Two roles exist in
the backend — **admin** and **student** — so there are two journeys, run as
separate k6 scenarios and mixed ~10% admin / ~90% student to match real traffic.

Each virtual user (VU) **logs in once** via `POST /auth/public/login`, caches its
JWT, then loops through a realistic session with think-time between actions.

```
loadtest/
├── main.js                 # entry: profiles → scenarios, thresholds
├── journeys.js             # student() and admin() sessions
├── lib.js                  # login + token cache + auth helpers
├── data/
│   ├── users.example.json  # template
│   └── users.json          # you create this (gitignore it)
└── README.md
```

## ⚠️ Before you run

- **Never target production.** These scripts hit real endpoints hard and the
  student journey **writes** (marks progress, submits assignments). Point
  `BASE_URL` at local or a dedicated staging instance backed by a **throwaway
  DB**. Hitting prod would hammer Render/Neon/Upstash.
- **Raise the rate limit**, or you'll measure the limiter, not the app. The API
  has a global per-IP rate-limit middleware; all VUs share your one IP, so under
  load they'll get `429`s. Boot the test server with a very high limit:
  `RATE_LIMIT_REQUESTS=1000000 PORT=8080 go run ./cmd/server`.
- **A VU is not a user.** "50 VUs" is 50 concurrent execution loops, each pausing
  1–5s between requests — it does *not* equal 50 requests/sec. Request rate
  depends on the journey and think-time. Read the `http_reqs` rate and
  `iterations` in the summary for actual throughput.
- The Lambda code-runner (`/code/run`) is **excluded** on purpose — it costs
  money per call and is rate-limited to one run / 2s per user.

## 1. Install k6

```bash
brew install k6           # macOS
# or see https://grafana.com/docs/k6/latest/set-up/install-k6/
```

## 2. Create the accounts file

```bash
cp data/users.example.json data/users.json
# edit data/users.json — real seeded students (username) + admin(s) (email)
```

More distinct students = more realistic (VUs spread across rows instead of
contending on one). Seed a batch of test students in your dev DB first.

## 3. Run a profile

`TEST` selects the profile; `BASE_URL` selects the target.

```bash
BASE_URL=http://localhost:8080 TEST=smoke    k6 run main.js   # 1–2 VUs, verify script/API
BASE_URL=http://localhost:8080 TEST=baseline k6 run main.js   # ~10 VUs, normal performance
BASE_URL=http://localhost:8080 TEST=load     k6 run main.js   # ~50 VUs, expected classroom usage
BASE_URL=http://localhost:8080 TEST=heavy    k6 run main.js   # ~100 VUs, above expected
BASE_URL=http://localhost:8080 TEST=stress   k6 run main.js   # ramp to 200+, find breaking point
BASE_URL=http://localhost:8080 TEST=spike    k6 run main.js   # 0→100 sudden surge
BASE_URL=http://localhost:8080 TEST=soak     k6 run main.js   # ~40 VUs for 45m, resource leaks
```

| Profile   | VUs (total) | Duration | Purpose                  |
| --------- | ----------: | -------- | ------------------------ |
| smoke     |           2 | 30s      | Verify script/API        |
| baseline  |          10 | 2m       | Normal performance       |
| load      |          50 | 5m       | Expected classroom usage |
| heavy     |         100 | 5m       | Above expected usage     |
| stress    |    ramp 200 | 6m       | Find the breaking point  |
| spike     |    0→100    | ~2m      | Sudden traffic           |
| soak      |          40 | 45m      | Memory/resource leaks    |

## 4. Read the results

The run exits non-zero if a **threshold** in `main.js` is breached (tune them to
your SLOs):

- `http_req_failed` — error rate; target `< 1%`.
- `http_req_duration` p(95)/p(99) — latency budget (currently 800ms / 1500ms).
- `checks` — share of per-request assertions that passed (`> 99%`).

Requests are grouped by a `name` tag (e.g. `GET /dashboard/student/{courseId}/lessons`),
so the summary shows latency **per endpoint** — that's where you'll spot the slow
one. `login_duration` is tracked separately (bcrypt is intentionally heavy, and
happens once per VU, so it shouldn't dominate).

## Notes

- Journeys degrade gracefully on an empty DB: list endpoints return `200` with no
  items and the drill-down steps are skipped, so you still exercise every list
  endpoint. For deeper coverage, seed courses/lessons/assignments.
- To run one role only, temporarily comment out the other scenario in
  `main.js` (k6 has no built-in single-scenario flag).
