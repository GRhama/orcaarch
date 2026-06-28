<p align="center">
  <img src="docs/assets/orcaarch-cover.png" alt="OrcaArch — Orchestrated Reconciliation & Control Architecture" width="900">
</p>

# OrcaArch

[English](#orcaarch) · [Português](#orcaarch-pt-br)

---

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

---

# OrcaArch PT-BR

OrcaArch é um motor de reconciliação construído em Go, organizado em pipelines por domínio. O projeto demonstra que um Product Manager com visão técnica consegue arquitetar e entregar uma solução multi-cenário — com suporte de IA estruturado e revisão humana em cada etapa.

## O que é

Motor de reconciliação configurável por domínio. Cada cenário tem suas próprias regras de negócio, mas compartilha os mesmos padrões arquiteturais: Clean Architecture, Money em inteiro escalado, dados mock determinísticos, API somente-leitura e exportação CSV.

## Por que existe

Provar que a combinação de visão de produto + workflow de IA bem estruturado + Go clean entrega software técnico de qualidade sem precisar de uma equipe de engenharia. O fluxo completo — do insumo até o CSV e a API — está acima neste README.

## Testar com npx

```bash
npx create-orcaarch
cd <nome-do-projeto>
make run
```

Gera um scaffold funcional com o cenário escolhido, pronto para rodar localmente.

## Rodar API local

```bash
make run-api
# API disponível em http://localhost:8080
```

Endpoints disponíveis (somente leitura, sem autenticação):

- `GET /inventory` — inventário com status de reconciliação
- `GET /reconciliation` — resumo de discrepâncias
- `GET /risk` — itens com risco calculado

## Exportar CSV

```bash
make export-csv
# Gera: reconciliation.csv, ledger.csv, risk.csv
```

## Cenários

| Cenário | Status |
|---|---|
| Logistics / Supply Chain | ✅ Implementado |
| Financial / Banking | 🔵 Em Design — ver [docs/FINANCIAL_SCENARIO_DESIGN.md](docs/FINANCIAL_SCENARIO_DESIGN.md) |

## Desenvolvimento assistido por IA

O projeto foi desenvolvido com workflow de IA estruturado: cada tarefa tem contexto mínimo definido, gates de revisão humana e rastreabilidade de decisão. Nenhum artefato foi gerado e aceito cegamente — toda decisão arquitetural foi revisada e documentada.

## API First

Toda a lógica de negócio é acessível via API HTTP antes de qualquer integração externa. Isso permite validar o comportamento do motor de reconciliação de forma isolada, sem depender de banco de dados ou serviços externos.

## Documentação

- [Arquitetura](docs/ARCHITECTURE.md)
- [API First](docs/API_FIRST.md)
- [Workflow com IA](docs/AI_WORKFLOW.md)
- [Caso de Produto](docs/PRODUCT_CASE.md)
- [Cenários](docs/SCENARIOS.md)
- [Benchmarks](docs/BENCHMARKS.md)
- [Design do Cenário Financeiro](docs/FINANCIAL_SCENARIO_DESIGN.md)
