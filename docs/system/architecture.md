# Architecture

> Template. Replace all placeholders when adopting this standard.
> Purpose: describe the intended system architecture (architecture intent).
> Implementation truth lives in `app/`, `config/`, and `infra/`. If they disagree, report the discrepancy.

## Architectural Style

- Style: `[Project declares its own — e.g. modular monolith, layered, event-driven, microservices. Do not leave generic.]`
- Key patterns: [List — must match `architecture.patterns` intent in `.ai/project.yaml`.]
- Decision references: [Links to relevant `docs/decisions/*.md`.]

## System Overview

[Short paragraph: what the system is, its main responsibilities, and its runtime shape.]

```text
[Optional ASCII overview: major components and relationships.
Keep small — detailed graphs belong in docs/system/dependencies.md.]
```

## Components

Summarize here; detail belongs in [components.md](./components.md).

| Component | Responsibility | Location |
|-----------|----------------|----------|
| [Name] | [One line] | [Source path] |

## Communication Mechanisms

| From → To | Mechanism | Notes |
|-----------|-----------|-------|
| [Component → Component] | [e.g. function call / request-response / message / batch — project declares] | [Sync/async, contract location] |

## Data Flows

1. [Primary flow: trigger → processing → storage → output.]
2. [Secondary flow, if any.]

## Boundaries

- Module/component boundaries: [What belongs where, per `components.md`.]
- Trust boundaries: [Where authentication/authorization is enforced.]
- External boundaries: [What crosses outside the system.]

## External Integrations

| Integration | Direction | Contract | Failure Handling |
|-------------|-----------|----------|------------------|
| [External system] | [inbound/outbound] | [Protocol, schema location] | [Timeout, retry, fallback] |

## Security Boundaries

- [Boundary 1: enforcement point and policy.]
- [Boundary 2.]

## Scalability Considerations

- [Known bottleneck or scaling unit.]
- [Stateless vs. stateful notes.]

## Failure Modes

| Failure | Impact | Mitigation |
|---------|--------|------------|
| [e.g. Dependency unavailable] | [Effect] | [Fallback, degradation, alerting] |

---

**Rules:**

- Do not prescribe an architecture style in this template. The project must explicitly state which architecture it uses.
- Every non-trivial claim here should trace to code, configuration, or a decision record.
