---
PLAN: "feat: Sensitive capability so a handler can mark its value unfit to render or log"
EXECUTOR: jules
REVIEWER: none
---

> This plan is dispatched via the CodeJob workflow. See skill: agents-workflow.

# Plan — `Sensitive` capability interface

## Part of a multi-repo wave

This is one piece of `KEYRING_DOTENV_MASTER_PLAN.md` (orchestrator:
`github.com/tinywasm/app-releases`, `docs/KEYRING_DOTENV_MASTER_PLAN.md`). This
piece has **no dependency on any other piece of that wave** — dispatch it
immediately, in parallel with the `tinywasm/kvdb` and `tinywasm/keyring`
plans. `tinywasm/devtui` and `tinywasm/wizard` each depend on this publishing
first.

## Why

A future handler (a secret/password field, first consumer: a
`tinywasm/wizard.Step` that asks for a credential) must never have its typed
value shown on screen or written into a log line — `tinywasm/devtui` streams
its transcript over SSE (`GET /logs`, also read by `tinywasm/app`'s MCP tool
`app_get_logs`), so anything logged in plaintext is not just on-screen, it is
exported.

Today nothing in this package can express "do not render or log my value" —
verified: `grep -rln "mask\|password\|Secret\|Sensitive" *.go` in this repo
returns nothing.

## The change

File: **[`interfaces.go`](interfaces.go)**. Add a new optional capability
interface, next to the existing handler interfaces (`HandlerDisplay`,
`HandlerEdit`, `HandlerExecution`, `HandlerInteractive`, `HandlerSelection`):

```go
// Sensitive is an optional capability a handler implements alongside
// HandlerEdit, HandlerInteractive, or HandlerSelection to mark its current
// Value() unfit to render or log in plaintext — a credential being typed,
// for example.
//
// Consumers detect it the same way the ecosystem already detects Cancelable
// or TabAware: a type assertion on the handler passed to AddHandler, never a
// change to the core Handler* interfaces (that would force every existing
// handler to implement a method it has no use for).
//
//	if s, ok := handler.(Sensitive); ok && s.Sensitive() {
//		// render/log a mask instead of the real Value()
//	}
//
// Sensitive() is called on every render/log — for a handler whose
// sensitivity can change per field (e.g. a multi-step wizard where only some
// steps collect a secret), return the CURRENT step's state, not a fixed
// value decided once.
type Sensitive interface {
	Sensitive() bool
}
```

This is additive only — no existing interface's method set changes, no
existing handler needs a new method to keep compiling.

## Tests

File: **`sensitive_test.go`** (new, next to `interfaces.go`).

| Test | Asserts |
|---|---|
| `TestSensitiveIsOptionalCapability` | `editSample` (existing fixture, defined alongside `displaySample`/`executionSample` in this package's test files) does **not** satisfy `Sensitive` — `_, ok := any(editSample{}).(Sensitive); ok` is `false` |
| `TestSensitiveDetected` | a fixture type implementing `HandlerEdit` **and** `Sensitive` (`Sensitive() bool { return true }`) is detected by the same type assertion — `ok` is `true` and the returned value's `Sensitive()` is `true` |

## Acceptance criteria

- [ ] `go build ./...` and `go vet ./...` clean.
- [ ] `go test ./...` green, including the two new tests.
- [ ] `grep -n "type Sensitive interface" interfaces.go` → exactly one match.
- [ ] No existing interface in `interfaces.go` had a method added or removed —
      `git diff interfaces.go` shows only new lines, no changed ones.

## Out of scope

Do not touch `tinywasm/devtui` or `tinywasm/wizard` here — consuming this
capability (masking the on-screen render, masking the auto-logged
confirmation) is their own plan in this same wave, dispatched once this one
publishes.
