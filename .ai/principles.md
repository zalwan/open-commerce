# AI Principles

> Standard: AI-Native Repository Standard v1.0
> Scope: Defines required behavior for any agent or human modifying this repository.
> This file contains behavior instructions only. It does NOT contain system-specific technical documentation.

System knowledge lives in `docs/`. Implementation lives in `app/`, `infra/`, and `config/`. This file defines how to act, not what the system is.

## 1. Understand Before Modifying

- Do not modify code based solely on filenames or search results.
- Read the relevant implementation, its interfaces, and its documented responsibility before changing it.
- Confirm the applicable scope: root `AGENTS.md`, applicable local `AGENTS.md`, `.ai/context.md`, relevant `docs/system/*`, and relevant decision records.
- If the intent of the existing code is unclear, stop and resolve the uncertainty before modifying it.

## 2. Prefer Existing Patterns

- Reuse established project patterns for structure, naming, error handling, testing, and configuration.
- When a pattern exists in `docs/development/conventions.md` or in adjacent code, follow it.
- Introduce a new pattern only when no existing pattern fits, and document the deviation in the change summary.

## 3. Minimize Change Surface

- Prefer the smallest correct change that satisfies the request.
- Avoid unrelated refactoring, reformatting, or scope expansion.
- If an unrelated improvement is noticed, report it separately; do not bundle it silently.

## 4. Preserve Boundaries

- Respect documented architecture and component boundaries in `docs/system/architecture.md` and `docs/system/components.md`.
- Do not bypass interfaces, layering rules, or ownership boundaries.
- Cross-boundary changes require explicit impact analysis per `.ai/workflow.md`.

## 5. Explicit Uncertainty

- Never invent missing architecture or technology decisions.
- When required context is absent or contradictory, state what is missing and what was assumed, if anything.
- Use the discrepancy format:

```text
Architecture discrepancy detected.

Documentation states:
<documented behavior>

Implementation indicates:
<actual behavior>

This should be resolved or explicitly acknowledged before
making a change that depends on this information.
```

- Human intent remains authoritative. Do not override an explicit user instruction with an inference from documentation.

## 6. Verify

- A code modification is not complete until relevant validation has been performed.
- Run the applicable checks: tests, static analysis, formatting, type checking, infrastructure validation, security checks.
- Report what was run and the result. An unvalidated change is an incomplete change.

## 7. Keep Knowledge Current

- System changes that affect architecture, dependencies, constraints, operations, or behavior must update the relevant documentation in `docs/`.
- Update only the documents affected by the change; do not rewrite unrelated documentation.
- When an architectural decision changes, create a new decision record. Do not silently rewrite history (see `docs/decisions/README.md`).
