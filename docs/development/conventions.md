# Engineering Conventions

> Template. Replace all placeholders when adopting this standard.
> Purpose: project-specific rules that keep contributions consistent.
> Project-specific conventions must be defined after the technology stack is selected. Until then, the generic rules below apply.

## Generic Rules (apply to every adopter)

1. Follow the workflow in `.ai/workflow.md` and the principles in `.ai/principles.md`.
2. Prefer the smallest correct change; no unrelated refactoring.
3. Respect component boundaries in `docs/system/components.md`.
4. Respect constraints in `docs/system/constraints.md`.
5. Validate before declaring a change complete; report what was run.
6. Update `docs/` only where the change affects documented knowledge.

## Project Conventions (define after stack selection)

### Code Style

- [Formatter, linter, naming rules — with commands.]

### Branching and Commits

- [Branch naming, commit message format, review requirements.]

### Testing

- [Where tests live (`tests/`), what coverage is expected, how to run them.]

### Configuration and Secrets

- [How configuration is layered; secret handling rules — never commit secrets.]

### Documentation

- [When to update `docs/system/`, when to write an ADR, changelog expectations.]

---

**Guidance:** conventions that restate the generic rules add no value. Define concrete, checkable expectations (commands, paths, formats) specific to this project's stack.
