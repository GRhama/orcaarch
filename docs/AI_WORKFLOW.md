# AI-Assisted Development Workflow

OrcaArch was built using a structured AI-assisted development approach with Claude Code.

## Principles

**Scoped iterations**: each development step loaded only the files it needed — domain contracts, acceptance criteria, and the specific package under implementation. Not the full specification at once.

**Domain-first**: domain entities, value objects, and invariants were defined and tested before any infrastructure was touched. The domain package has zero external dependencies.

**Human review gates**: architecture decisions, trade-off choices, and scope boundaries were made by the developer — not delegated to the AI. The AI implemented within constraints the developer defined.

**Tests before integration**: domain business rules have unit tests before service orchestration or API wiring. Benchmarks validate O(n) scaling assumptions before moving on.

## What the AI accelerated

- Translating domain rules into Go structs and methods
- Writing test cases for reconciliation states and edge cases
- Implementing the reprocessing pattern (reversal + version increment)
- Wiring report mappers and CSV writers
- Maintaining consistency across packages (Money type, error handling conventions)

## What stayed human

- Choosing Clean Architecture as the structural constraint
- Defining the reconciliation states and their invariants
- Deciding on integer-scaled Money (no float) as a non-negotiable
- Scoping each iteration to avoid premature abstraction
- Reviewing and approving every change before continuing to the next step
- Designing the financial scenario domain model before any implementation

## Development cycle

```
define domain contract
→ write tests
→ implement domain
→ human review
→ implement service / adapter
→ human review
→ next iteration
```

Each iteration: small, independently testable, reviewed before continuing.
