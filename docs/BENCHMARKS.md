# OrcaArch — Processing Benchmarks

## Environment

| Key       | Value                                    |
|-----------|------------------------------------------|
| Date      | 2026-06-27                               |
| OS        | Linux 6.8 amd64                          |
| CPU       | Intel Core i7-8565U @ 1.80GHz           |
| Go        | 1.22                                     |
| Command   | `go test ./internal/service/... -bench=. -benchmem -benchtime=3x` |

## Results

| Benchmark            | n       | ns/op         | ms/op   | B/op        | allocs/op |
|----------------------|---------|---------------|---------|-------------|-----------|
| BenchmarkProcess50k  | 50,000  | 836,714,514   | 836.7   | 144,044,744 | 439,325   |
| BenchmarkProcess100k | 100,000 | 1,715,316,782 | 1,715.3 | 291,095,146 | 878,169   |
| BenchmarkProcess250k | 250,000 | 4,010,376,361 | 4,010.4 | 724,033,170 | 2,203,562 |

## Scaling

Processing is linear in n (no quadratic paths, no external I/O):

| Scale jump     | Expected | Time ratio | Allocs ratio | Memory ratio |
|----------------|----------|------------|--------------|--------------|
| 50k → 100k     | 2.0×     | 2.05×      | 2.00×        | 2.02×        |
| 100k → 250k    | 2.5×     | 2.34×      | 2.51×        | 2.49×        |

Scaling is O(n). No optimization needed at these volumes.

## Smoke Test (100,000 records)

Run: `go test ./internal/service/... -run TestSmokeProcess100k -v`

| Assertion                         | Expected | Result |
|-----------------------------------|----------|--------|
| `len(Reconciliations)`            | 100,000  | PASS   |
| `len(Quarantine)`                 | 25,000   | PASS   |
| `len(LedgerEntries)`              | 100,002  | PASS   |
| All 4 reconciliation statuses     | present  | PASS   |

## Notes

- All processing is in-memory. Timings exclude I/O, serialization, and network.
- `BenchmarkProcess250k` is optional and included for scaling validation only.
- No optimization was applied before or after measurement.
