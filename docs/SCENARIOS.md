# Scenarios

## Status

| Scenario | Status |
|---|---|
| Logistics / Supply Chain | ✅ Implemented |
| Financial / Banking | 🔵 In Design — see [FINANCIAL_SCENARIO_DESIGN.md](FINANCIAL_SCENARIO_DESIGN.md) |

---

## Logistics / Supply Chain

**Implemented.** Full pipeline in Go.

Primary reconciliation key: `tracking_number`
Sources: ERP Shipment vs Carrier Invoice

### Reconciliation states

| State | Condition |
|---|---|
| `MATCHED` | Both sides present; freight within tolerance |
| `DISCREPANCY` | Both sides present; freight exceeds tolerance |
| `UNRECONCILED_ERP` | ERP shipment has no carrier invoice |
| `UNRECONCILED_CARRIER` | Carrier invoice has no ERP shipment → Quarantine |

### Calculations

- **Landed Cost**: freight + duties + insurance + customs fees
- **WAC**: weighted average cost per SKU, integer arithmetic, trade-date order
- **Ledger**: double-entry entries per reconciliation status
- **Reprocessing**: idempotent reversal + recalculation on corrected carrier invoice

### Outputs

```
GET /api/v1/reports/inventory   reconciliation rows
GET /api/v1/reports/ledger      ledger entries
GET /api/v1/reports/risk        quarantine entries
```

CSV export: `reconciliation.csv`, `ledger.csv`, `risk.csv`

---

## Financial / Banking

**In Design.** Domain contracts and API shape fully specified. No Go implementation yet.

Primary reconciliation key: `trade_id`
Sources: Internal Trade vs Broker Execution

### Reconciliation states

| State | Condition |
|---|---|
| `MATCHED` | Both sides present; quantity and price within tolerance |
| `PRICE_DISCREPANCY` | Both present; price differs beyond tolerance |
| `QUANTITY_DISCREPANCY` | Both present; quantity differs |
| `UNRECONCILED_INTERNAL` | Internal trade has no broker execution |
| `UNRECONCILED_BROKER` | Broker execution has no internal trade → ExceptionQueue |

### Calculations

- **Position**: net quantity per instrument from confirmed trades
- **Average Cost**: weighted acquisition cost per instrument
- **Realized PnL**: `(sell_price − avg_cost) × quantity_sold`
- **Unrealized PnL**: `(mark_price − avg_cost) × open_quantity`
- **Ledger**: double-entry per reconciliation status
- **Reprocessing**: reversal + version increment on corrected broker execution

### Outputs (designed)

```
GET /api/v1/financial/reports/reconciliation
GET /api/v1/financial/reports/position
GET /api/v1/financial/reports/pnl
GET /api/v1/financial/reports/ledger
GET /api/v1/financial/reports/risk
```

CSV export: all 5 reports

See [FINANCIAL_SCENARIO_DESIGN.md](FINANCIAL_SCENARIO_DESIGN.md) for full domain model, entity definitions, and invariants.

---

## Shared Patterns

Both scenarios implement the same structural patterns with domain-specific entities.

| Pattern | Logistics | Financial |
|---|---|---|
| `Money` value object | ✅ | ✅ designed |
| Reconciliation result | ✅ | ✅ designed |
| Ledger entry + reversal | ✅ | ✅ designed |
| Idempotent reprocessing | ✅ | ✅ designed |
| Report / export | ✅ | ✅ designed |
| Benchmark harness | ✅ | planned |

**Design rule**: no shared interface or abstraction is extracted until both scenarios are fully implemented and duplication is visible and tested.
