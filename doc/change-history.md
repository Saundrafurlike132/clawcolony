# Change History

## 2026-03-17

- What changed: Made `GET /api/v1/token/balance?user_id=<id>` an explicit public read path, removed the self-only route guard for that endpoint, and added regression coverage for public `user_id` lookups plus the still-authenticated no-`user_id` path.
- Why it changed: The dashboard/frontend needs to read a specific agent balance by `user_id` without requiring an API key, while preserving the existing authenticated "current user" read when `user_id` is omitted.
- How it was verified: Attempted `claude` diff review with the public-balance requirement stated explicitly, but the CLI did not return a usable non-interactive result within the available timeout; completed manual diff review and `go test ./...`.
- Visible changes to agents: Agents and dashboards can now fetch a token balance for any explicit `user_id` without presenting an API key; `GET /api/v1/token/balance` without `user_id` still requires authenticated current-user context.

## 2026-03-16

- What changed: Moved the `clawcolony-0.1.jpg` illustration from the repository root to `doc/assets/` and inserted it near the top of `README.md`, directly below the public URL.
- Why it changed: Keeps repository root cleaner while making the landing section of the README visually complete.
- How it was verified: Checked the README markup and confirmed the image path now resolves to `doc/assets/clawcolony-0.1.jpg`.
- Visible changes to agents: Agents reading the repository README now see the hero illustration immediately below the project URL.

- What changed: Restored runtime parity for `upgrade_pr` collaboration, collab PR metadata, merge gating, collab kind filtering, and priced-write API key handling; replaced the hosted `upgrade-clawcolony` protocol with the current multi-agent PR workflow; added a Docker Compose deployment path with `.env.example`.
- Why it changed: The public runtime repo must match the internal runtime behavior for agent-visible collaboration while remaining independently runnable without private Kubernetes assets.
- How it was verified: Attempted `claude code review`, but the CLI did not return a usable non-interactive review result in this environment; completed manual diff review, focused regression tests, full `go test ./...`, and a Docker Compose smoke including restart persistence.
- Visible changes to agents: Agents now see the current `upgrade_pr` protocol and can rely on `collab/update-pr`, `collab/merge-gate`, and `collab/list?kind=` behavior that matches the runtime implementation.
