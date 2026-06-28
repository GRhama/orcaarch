# Product Case

## Problem

Reconciliation pipelines — matching records from two sources, detecting discrepancies, generating accounting entries, handling exceptions, and reprocessing corrections — appear across industries with the same structural shape:

- Logistics: ERP shipment vs carrier invoice
- Finance: internal trade vs broker execution
- Payments: internal ledger vs bank statement
- Healthcare: insurance claim vs provider billing

Each domain has specific entities and rules. The underlying pattern is the same.

## Solution

OrcaArch is a scenario generator for reconciliation pipelines. Each scenario is a self-contained implementation:

- Deterministic mock data generator (fixed seed, reproducible at any volume)
- Reconciliation engine (domain-specific matching rules)
- Cost / position calculation (landed cost for logistics; avg cost + PnL for financial)
- Double-entry ledger (accounting entries per reconciliation status)
- Exception queue (unmatched records isolated, excluded from calculations)
- Idempotent reprocessing (reversal + recalculation + version increment)
- Reports (CSV + read-only HTTP API)

Shared patterns (Money, ledger, reprocessing) are kept separate per scenario — no forced abstraction before both implementations exist and duplication is visible and tested.

## What this demonstrates

**Product thinking**: multi-scenario architecture designed from the start. The financial scenario domain model, API shape, and report contracts are fully specified before implementation begins — API First in practice.

**Engineering discipline**: Clean Architecture enforced at the package level. Integer-scaled Money as a non-negotiable constraint. Idempotent reprocessing. O(n) processing. No premature optimization.

**Delivery**: testable via `npx create-orcaarch` — without cloning the source repo.

**Structured AI-assisted development**: scoped iterations, domain-first, human review gates at every step.

## Audience

- Engineering managers evaluating technical depth and product judgment
- Tech leads reviewing Go domain modeling and architecture decisions
- Interviewers assessing product + engineering balance
- Developers interested in structured AI-assisted development workflows
