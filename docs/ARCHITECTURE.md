# Architecture

## Design Principles

- **Clean Architecture**: domain has zero infrastructure dependencies
- **Money without float**: all monetary values as `int64` cents via `domain.MustMoney` — no `float32` or `float64` anywhere
- **Idempotent reprocessing**: reversal + re-reconcile + version increment
- **Read-only API**: no writes, no auth, no database
- **Deterministic mock data**: fixed seed, reproducible output at any volume

## Layer Map

```
domain/     business rules, entities, value objects
            imports: nothing internal

mock/       deterministic data generator
            imports: domain

service/    pipeline orchestration
            imports: domain, mock

report/     DTOs and row mappers
            imports: domain

export/     CSV writers
            imports: report

api/        HTTP handlers (read-only)
            imports: service, report

cmd/api/    entrypoint — start HTTP server
cmd/export/ entrypoint — write CSV files
```

The domain package is the invariant boundary. It never imports from service, api, export, mock, or report.

## Logistics Pipeline

![OrcaArch processing flow](assets/orcaarch-flow.png)

```
mock.Generate(n)
→ Reconcile      ERP shipment vs carrier invoice, matched by tracking_number
→ Landed Cost    freight + duties + insurance + customs fees
→ WAC            weighted average cost per SKU, integer arithmetic, trade-date order
→ Ledger         double-entry entries per reconciliation status
→ Quarantine     carrier invoices without matching ERP shipment
→ Reprocess      reversal + recalculation when a carrier invoice is corrected
→ Reports        reconciliation, ledger, risk — CSV + read-only API
```

## Reconciliation States

| State | Condition |
|---|---|
| `MATCHED` | tracking_number in both; freight within tolerance |
| `DISCREPANCY` | tracking_number in both; freight exceeds tolerance |
| `UNRECONCILED_ERP` | ERP shipment has no carrier invoice |
| `UNRECONCILED_CARRIER` | Carrier invoice has no ERP shipment → Quarantine |

## Money

```go
type Money struct {
    Amount   int64
    Currency string
}
```

No `float32`. No `float64`. All arithmetic in integer-scaled cents.

## Reprocessing

```
Input: original reconciliation result + original ledger entries + corrected invoice + ERP record
If corrected invoice == original: Changed=false, no new entries, no version increment
If changed: generate reversals → re-reconcile → Version+1 → new ledger entries
Same input always produces same output.
```

## API

Stdlib `net/http` only. No framework, no middleware, no auth.

```
GET /api/v1/reports/inventory  → reconciliation rows (JSON)
GET /api/v1/reports/ledger     → ledger entries (JSON)
GET /api/v1/reports/risk       → quarantine entries (JSON)
```

All monetary values in responses: `{ "Amount": int64, "Currency": string }`. No float.
