---
name: design-principles
description: "Function and module design in aya.is, in Go and the webclient's TypeScript: pure functions with thin wrappers, copies over in-place changes, small ports passed in, plain data, one clock read, one cancellation handle, collections by access pattern, no import side effects. Use when shaping signatures, types, modules, state or collections. Not for system layout (use architecture-guidelines)."
---

# Design Principles

How functions, modules, data and state are shaped. System structure lives in
architecture-guidelines; code-level quality rules in coding-practices.

## Always

- Logic lives in pure functions with explicit parameters, in their own file; a
  thin wrapper reads settings or the environment and performs the effect.
- Return new values instead of changing inputs; the webclient copies data
  instead of mutating it, and Go hands out clones of internal slices and maps.
- Accept small interfaces (ports) declared by the consumer and pass
  dependencies in as parameters, fields or props; no inheritance for reuse.
- Data is plain structs and objects plus functions; classes only for stateful
  services; no getters or setters with side effects.
- Read the clock once per operation and pass the value down; take the clock as
  a `now` function or `clockfx.Clock` so tests set it.
- One cancellation handle (`ctx`, `AbortSignal`) flows through the whole
  operation; never start a fresh one mid-chain.
- Pick maps and sets for lookups, slices and arrays for sequences.
- Importing a package or module does no I/O and reads no environment; no
  `init()`, no work at module scope.
- Use promises (`async` / `await`) over callbacks and template strings over
  concatenation in TypeScript.

## References

| File                                      | Read when                                                                  |
| ----------------------------------------- | -------------------------------------------------------------------------- |
| [principles.md](references/principles.md) | Designing a function, module, type or collection; reviewing state handling |
