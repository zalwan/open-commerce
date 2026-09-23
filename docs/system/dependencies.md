# Dependencies

> Template. Replace all placeholders when adopting this standard.
> Purpose: describe what depends on what, inside and outside the repository.
> The initial representation is Markdown. No graph-generation tooling is required.

## Component Dependency Graph

```text
[Project declares. Keep small and directional.
Example shape (replace entirely):

  [entry-point] --> [component-a] --> [data-store]
                    [component-a] --> [external-service]
                    [entry-point] --> [component-b]
]
```

## Runtime Dependencies

| Dependent | Depends On | Type | Notes |
|-----------|------------|------|-------|
| [Component/service] | [Library/service] | [internal / third-party / platform] | [Version constraint, purpose] |

## External Dependencies

| Dependency | Provider | Purpose | Contract | Fallback |
|------------|----------|---------|----------|----------|
| [Name] | [Who provides it] | [Why needed] | [API/schema location] | [Degradation or "none"] |

## Infrastructure Dependencies

| Dependent | Infrastructure | Notes |
|-----------|----------------|-------|
| [Service/component] | [Network, storage, hosting, pipeline dependency] | [Environment scope] |

---

**Guidance:**

- Document only dependencies that matter for building, running, or changing the system.
- If documented dependencies differ from actual dependencies, that is context drift — resolve or explicitly acknowledge it before changing dependent code.
