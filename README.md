<p align="center">
  <img src="docs/assets/orcaarch-cover.png" alt="OrcaArch — Orchestrated Reconciliation & Control Architecture" width="900">
</p>

# OrcaArch

Reconciliation engine playground for domain-specific pipelines — logistics and financial.

Built in Go with Clean Architecture, integer-scaled Money, deterministic mock data, and a structured AI-assisted development workflow.

## How OrcaArch works

<p align="center">
  <img src="docs/assets/orcaarch-flow.png" alt="How OrcaArch works — from inputs to controls, reports, CSV and read-only API" width="1000">
</p>

## Try it

```bash
npx create-orcaarch
cd <project-name>
make run
```

## Source engine

```bash
go test ./...
go test ./internal/service/... -bench=. -benchmem
make run-api      # starts API on :8080
make export-csv   # writes reconciliation.csv, ledger.csv, risk.csv
```

## Scenarios

| Scenario | Status |
|---|---|
| Logistics / Supply Chain | ✅ Implemented |
| Financial / Banking | 🔵 In Design — see [docs/FINANCIAL_SCENARIO_DESIGN.md](docs/FINANCIAL_SCENARIO_DESIGN.md) |

## What this demonstrates

- Go engineering: Clean Architecture, no ORM, no framework
- Money without floating point: integer-scaled `int64`, no `float32`/`float64`
- Reconciliation logic: matching, discrepancy detection, exception queue
- Idempotent reprocessing: reversal + recalculation + version increment
- Double-entry ledger: accounting entries per reconciliation status
- Deterministic mock data: fixed-seed generator, reproducible at any volume
- Read-only API: stdlib `net/http` only
- CSV export: standard library only
- Benchmarks: O(n) processing validated at 50k–250k records
- Multi-scenario architecture: pluggable domains, shared patterns, no forced abstraction
- Structured AI-assisted development with human review gates

## Requirements

- Go 1.22+
- Node 16+ (for `npx create-orcaarch`)

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [API First](docs/API_FIRST.md)
- [AI Workflow](docs/AI_WORKFLOW.md)
- [Product Case](docs/PRODUCT_CASE.md)
- [Scenarios](docs/SCENARIOS.md)
- [Benchmarks](docs/BENCHMARKS.md)
- [Financial Scenario Design](docs/FINANCIAL_SCENARIO_DESIGN.md)
