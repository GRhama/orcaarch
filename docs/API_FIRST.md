# API First

Contracts are designed before implementation. The financial API shape is fully specified; the Go implementation follows.

## Logistics (implemented)

```
GET /api/v1/reports/inventory
GET /api/v1/reports/ledger
GET /api/v1/reports/risk
```

All responses: JSON array. Monetary values: `{ "Amount": int64, "Currency": string }`. No float.

```bash
make run-api

curl localhost:8080/api/v1/reports/inventory
curl localhost:8080/api/v1/reports/ledger
curl localhost:8080/api/v1/reports/risk
```

## Financial / Banking (designed, not yet implemented)

```
GET /api/v1/financial/reports/reconciliation
GET /api/v1/financial/reports/position
GET /api/v1/financial/reports/pnl
GET /api/v1/financial/reports/ledger
GET /api/v1/financial/reports/risk
```

Full API shape, report schemas, and response contracts defined in [docs/FINANCIAL_SCENARIO_DESIGN.md](FINANCIAL_SCENARIO_DESIGN.md).

## Design Rules

- Read-only: no POST, PUT, DELETE
- No database
- No authentication
- No framework — stdlib `net/http` only
- Money: `int64` amount + `string` currency in every response, no float
- Contracts designed before implementation — the financial endpoints above are specified and stable before any Go code is written
