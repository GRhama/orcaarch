# Financial / Banking Scenario Design

> **Status: Design Only — not yet implemented.**
> Domain contracts, entities, reconciliation states, PnL model, and API shape are fully specified here.
> No Go implementation exists for this scenario yet.
> Implementation begins in the next iteration.

## Status

Design only. Not yet implemented.

## Purpose

Define the conceptual contracts, domain entities, reconciliation states, ledger model, position model, PnL model, and output boundaries for the OrcaArch financial/banking scenario.

This document is the design base for EPIC-004 (TASK-018 to TASK-027). It does not implement any of those tasks. Each future task will have its own scope and context files.

---

## Epic Scope — TASK-018 to TASK-027

Candidate tasks for EPIC-004 as defined in `docs/ROADMAP.md`. This document defines the shared design foundation. Each task is independently scoped and will declare its own minimum context.

| Task | Name | Output |
|------|------|--------|
| TASK-018 | Financial Scenario Design | This document + updated SCENARIO_MATRIX.md |
| TASK-019 | Financial Domain Models | Go structs: InternalTrade, BrokerExecution, Position, AverageCost, RealizedPnL, UnrealizedPnL, FinancialLedgerEntry, ExceptionQueue |
| TASK-020 | Trade Reconciliation Engine | Reconcile() for financial trades; 5 reconciliation states |
| TASK-021 | Position and Average Cost | CalculatePosition(), CalculateAverageCost() |
| TASK-022 | PnL Calculation | CalculateRealizedPnL(), CalculateUnrealizedPnL() |
| TASK-023 | Financial Ledger Rules | GenerateFinancialLedgerEntries() |
| TASK-024 | Financial Mock Generator | Generate(n) for financial scenario |
| TASK-025 | Financial Processing Service | Full financial pipeline in service layer |
| TASK-026 | Financial Reports | CSV + read-only API for financial reports |
| TASK-027 | Installer Financial Scenario Option | Unlock financial scenario in create-orcaarch CLI |

**TASK-018 does not implement any task above.** It defines what each should produce.

---

## Domain Entities

### InternalTrade

Internal record of a trade order/execution.

```text
Fields:
  trade_id        string     — primary reconciliation key
  instrument      string     — asset ticker or ISIN
  direction       string     — BUY | SELL
  quantity        int64      — in base units (no float)
  price           Money      — price per unit in USD
  notional        Money      — quantity × price (derived)
  trade_date      time.Time
  status          string     — PENDING | CONFIRMED | CANCELLED
```

### BrokerExecution

Broker or exchange confirmation of a trade.

```text
Fields:
  trade_id        string     — matches InternalTrade.trade_id
  instrument      string
  direction       string     — BUY | SELL
  executed_qty    int64
  executed_price  Money
  execution_date  time.Time
  broker_ref      string     — broker-side reference
```

### ExchangeStatement

Daily or settlement statement from exchange or custodian.

```text
Fields:
  statement_date  time.Time
  instrument      string
  net_position    int64
  settlement_pnl  Money
```

### Settlement

Settlement record for a trade.

```text
Fields:
  trade_id        string
  settled_qty     int64
  settled_amount  Money
  settlement_date time.Time
  status          string     — PENDING | SETTLED | FAILED
```

### Position

Net holding per instrument, derived from confirmed trades.

```text
Fields:
  instrument      string
  net_quantity    int64
  as_of           time.Time
```

### AverageCost

Weighted average acquisition cost per instrument.

```text
Fields:
  instrument      string
  avg_cost        Money      — cost per unit, no float
  total_quantity  int64
  as_of           time.Time
```

### RealizedPnL

PnL from closed positions (partial or full sell).

```text
Fields:
  trade_id        string
  instrument      string
  realized_pnl    Money
  quantity_sold   int64
  avg_cost_used   Money
  sell_price      Money
```

### UnrealizedPnL

PnL from open positions, mark-to-market.

```text
Fields:
  instrument      string
  open_quantity   int64
  avg_cost        Money
  mark_price      Money      — from mock static price per instrument
  unrealized_pnl  Money
```

### FinancialLedgerEntry

Double-entry ledger record for financial operations.

```text
Fields:
  trade_id        string
  entry_date      time.Time
  account         string
  debit           Money
  credit          Money
  version         int        — incremented on reprocessing
```

Invariant: sum(debit) == sum(credit) per trade_id.

### ExceptionQueue

Broker execution without a matching internal trade. Analogous to logistics `Quarantine`.

```text
Fields:
  broker_ref      string
  execution       BrokerExecution
  reason          string     — ReasonMissingInternalTrade
```

Exception entries do not update Position or AverageCost.

---

## Reconciliation

### Primary Key

`trade_id` — must match between `InternalTrade` and `BrokerExecution`.

### States

| State | Condition |
|-------|-----------|
| `MATCHED` | trade_id in both; quantity and price within tolerance |
| `PRICE_DISCREPANCY` | trade_id in both; quantity matches; price differs beyond tolerance |
| `QUANTITY_DISCREPANCY` | trade_id in both; price matches; quantity differs |
| `UNRECONCILED_INTERNAL` | InternalTrade has no BrokerExecution |
| `UNRECONCILED_BROKER` | BrokerExecution has no InternalTrade → ExceptionQueue |

### Tolerance

Basis points, same mechanism as logistics reconciliation. No hardcoded threshold.

---

## Mock Data States

Generator must produce all 5 states deterministically (seed-based, same pattern as logistics mock):

| State | Target fraction | Notes |
|-------|----------------|-------|
| MATCHED | ~n/5 | identical qty + price |
| PRICE_DISCREPANCY | ~n/5 | broker price = internal × 1.01 (beyond tolerance) |
| QUANTITY_DISCREPANCY | ~n/5 | broker qty = internal qty − 10 |
| UNRECONCILED_INTERNAL | ~n/5 | trade_id only in InternalTrade list |
| UNRECONCILED_BROKER | ~n/5 | trade_id only in BrokerExecution list |

All monetary values via `domain.MustMoney`. No float32 or float64.

---

## Position & Average Cost Model

- Position = net quantity per instrument, derived from MATCHED and PRICE_DISCREPANCY trades.
- AverageCost = weighted average acquisition cost, computed in `trade_date` ascending order.
- Same pattern as logistics WAC; formula:

```text
new_avg_cost = (old_qty × old_avg + new_qty × new_price) / (old_qty + new_qty)
```

- Integer arithmetic only. Truncation rules documented in TASK-021.
- UNRECONCILED_INTERNAL: provisional position entry (analogous to `ProvisionalLandedCostInput`).
- UNRECONCILED_BROKER / ExceptionQueue: no position or average cost update.

---

## PnL Model

### RealizedPnL

Triggered when a SELL trade is reconciled:

```text
realized_pnl = (sell_price − avg_cost) × quantity_sold
```

### UnrealizedPnL

Triggered for open BUY positions:

```text
unrealized_pnl = (mark_price − avg_cost) × open_quantity
```

`mark_price` sourced from static mock price per instrument. No real market data.

All calculations in integer-scaled Money. No float.

---

## Ledger Model

Double-entry. Invariant: debit == credit per trade_id.

| Reconciliation State | Debit | Credit |
|----------------------|-------|--------|
| MATCHED | Position (asset) | Cash / Settlement Payable |
| PRICE_DISCREPANCY | Position (at broker price) | Cash + PnL Suspense |
| QUANTITY_DISCREPANCY | Position (at matched qty) | Cash + Quantity Suspense |
| UNRECONCILED_INTERNAL | Suspense (provisional) | Provisions |
| UNRECONCILED_BROKER | — exception queue only, no ledger — | — |

Reversal = swap debit↔credit on all original entries (same pattern as logistics).

---

## Reprocessing

Same pattern as logistics reprocessing engine:

- Input: original reconciliation result + original ledger entries + corrected BrokerExecution + InternalTrade
- If corrected execution == original: `Changed=false`, no new entries, no version increment
- If changed: generate reversals → re-reconcile → `Version+1` → new ledger entries
- Position/AverageCost recalculation is caller's responsibility (not inside Reprocess)
- Must be idempotent

---

## Reports

| Report | Rows | Key Fields |
|--------|------|------------|
| ReconciliationReport | per trade_id | trade_id, instrument, status, difference_bp |
| PositionReport | per instrument | instrument, net_quantity, avg_cost |
| PnLReport | per instrument | realized_pnl, unrealized_pnl, mark_price |
| LedgerReport | per entry | trade_id, account, debit, credit, version |
| RiskReport | exceptions + discrepancies | trade_id, state, reason |

---

## API / CSV Outputs

Read-only, same pattern as logistics:

```text
GET /api/v1/financial/reports/reconciliation
GET /api/v1/financial/reports/position
GET /api/v1/financial/reports/pnl
GET /api/v1/financial/reports/ledger
GET /api/v1/financial/reports/risk
```

CSV export for all 5 reports. Money exported as `amount int64 + currency string`. No float in output.

---

## Reuse Boundaries

**Reusable (shared concept, separate implementation for now):**
- `Money` value object
- Reconciliation result pattern
- Ledger entry + reversal pattern
- Reprocessing pattern
- Report/export pattern
- Benchmark harness

**Domain-specific (financial only, no forced abstraction):**
- `InternalTrade`, `BrokerExecution`
- `Position`, `AverageCost`
- `RealizedPnL`, `UnrealizedPnL`
- Financial ledger accounts (Trade, Settlement, PnL Suspense)

Do not extract shared interfaces before both scenarios are fully implemented and duplication is visible.

---

## Non-Goals

- No real market data or broker API integrations
- No FX conversion (all USD in MVP)
- No regulatory reporting (MiFID, EMIR, etc.)
- No real customer or trade data
- No `float32` / `float64` for any monetary or quantity value
- No custody system simulation
- No authentication
- No database
- No optimization before measuring (see EPIC-002)
- No implementation in TASK-018 — design-only
