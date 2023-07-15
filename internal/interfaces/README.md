# Backend interfaces

Folder `interfaces` consist interface/abstraction of object returned by Initialization of `Repositories`, `Usecases` or `Handlers`
We return interface instead of pointer of object so it can be unit-tested without use the real initialization function.
An interface consist of public functions we want to expose to other packages.

## What it do:
---

Folder `interfaces` main job is to:
- Abstract object as interface or functions so we can unit-test it easier

## Import Rules:
---
- `Interfaces` **CAN** import `entities` and `models` package.
- `Interfaces` **MUST NOT** import any other package beside above import rules.
- `Interfaces` **MUST NOT** imported by `entities` and `models` package.

