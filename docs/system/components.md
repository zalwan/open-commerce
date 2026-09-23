# Components

> Template. Replace all placeholders when adopting this standard.
> Purpose: component-level map of responsibilities, interfaces, and ownership.
> One section per significant component. Trivial helpers do not need entries.

## `[Component Name]`

- **Responsibility:** [What it does and what it explicitly does not do.]
- **Source location:** `[e.g. app/... path]`
- **Dependencies:** [Other components, libraries, or services it requires.]
- **Consumers:** [Who calls or depends on it.]
- **Interfaces:** [Public functions, APIs, events, or schemas — with paths to contracts.]
- **Data stores:** [What it reads/writes, if any — or "none".]
- **Operational notes:** [Scaling, background work, failure behavior, or "none".]

---

## `[Component Name]`

- **Responsibility:** [...]
- **Source location:** [...]
- **Dependencies:** [...]
- **Consumers:** [...]
- **Interfaces:** [...]
- **Data stores:** [...]
- **Operational notes:** [...]

---

**Guidance:**

- Copy the block above for each significant component.
- Delete this guidance and all placeholder blocks when populating.
- Keep entries short. Link to code and contracts instead of duplicating them.
- If code introduces a component absent from this map, that is context drift — update this file or report the discrepancy.
