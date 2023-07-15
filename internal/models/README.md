# Backend Models

Folder `models` consist of input and output structs between `Handlers` and `Usecases`.
The main difference between `entities` and `models` are in the struct definition.
`Models` struct usually represent input/output parameter of an endpoint, which may contains entity data.

`entities` and `models` technically can be inside a package, but for now we separate it first.

## What it do:
---

Folder `models` main job is to:
- Store all input and output structs between `Handlers` and `Usecases`

## Import Rules:
---
- `Models` **CAN** import `entities` package.
- `Models` **MUST NOT** import any other package beside above import rules.
- `Models` **CAN** be imported by `Handlers` and `Usecases`.
- `Models` **MUST NOT** imported by other package beside above imported rules.

