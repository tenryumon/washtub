# Backend Core

Folder `core` consist of general functions needed for service operational.
The functions that created here must be able to be used in another golang repositories.
In the future, this folder will become it's own golang repository, but for now we put it here

## What it do:
---

Folder `core` main job is to:
- Simplify initialization of a dependency
- Use default value to make sure we prevent misconfiguration
- Unify the behaviour in every golang service

