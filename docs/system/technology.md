# Technology

> Template. Replace all placeholders when adopting this standard.
> Purpose: authoritative inventory of the project's technology choices.
> Keep this file current. Technology versions should be explicitly documented when relevant.

## Languages

| Language | Version | Purpose | Notes |
|----------|---------|---------|-------|
| [Name] | [Version or "unspecified"] | [Where used] | [Optional] |

## Frameworks and Libraries

| Name | Version | Purpose | Notes |
|------|---------|---------|-------|
| [Name] | [Version] | [Where used] | [Optional] |

## Runtimes

| Runtime | Version | Purpose | Notes |
|---------|---------|---------|-------|
| [Name] | [Version] | [Where executed] | [Optional] |

## Data Stores

| Store | Type | Purpose | Notes |
|-------|------|---------|-------|
| [Name/role] | [e.g. relational / document / key-value / file — project declares its own] | [What it stores] | [Version, hosting, backup notes] |

## Infrastructure and Hosting

| Area | Choice | Notes |
|------|--------|-------|
| Hosting | [Project declares] | [Environments] |
| Networking | [Project declares] | [Optional] |
| Storage | [Project declares] | [Optional] |

## CI/CD

| Stage | Tool / Mechanism | Notes |
|-------|------------------|-------|
| Build | [Project declares] | |
| Test | [Project declares] | |
| Deploy | [Project declares] | |

## Observability

| Area | Tool / Mechanism | Notes |
|------|------------------|-------|
| Logs | [Project declares] | |
| Metrics | [Project declares] | |
| Traces | [Project declares or "none"] | |

## External Services

| Service | Purpose | Integration | Notes |
|---------|---------|-------------|-------|
| [Name] | [Why needed] | [Protocol/API] | [Auth, SLA, fallback] |

---

**Guidance:** list only what the project actually uses. Do not list candidate or illustrative technologies here. Illustrative stacks belong in `examples/`, not in this file.
