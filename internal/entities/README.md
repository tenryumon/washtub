# Backend Entities

Folder `entities` consist of input and output structs between `Repositories` and `Usecases`.
The main difference between `entities` and `models` are in the struct definition.
`Entities` struct key usually represent a table or combination of tables that determine an object.
If an entity have a fixed rules, like `isActive()` or `hasRole(role)`, you can create a function inside this folder instead of put it everywhere in usecases.

## What it do:
---

Folder `entities` main job is to:
- Store all input and output structs between `Repositories` and `Usecases`
- Store fixed rule function of the object

## Import Rules:
---
- `Entities` **MUST NOT** import any package from `backend`. 
- `Entities` **CAN** be imported by other packages in `backend` (beside `core` package)

