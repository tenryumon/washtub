# Backend Repositories

Folder `repositories` consist objects to get or change data from a dependency.
A good repository is similar to a full-fledge services, like when we want to get staff info+roles, then there must be staff repositories

## What it do:
---

Folder `repositories` main job is to:
- Get data from dependency like database, redis or other service.
- Change value of data to dependency.
- Create abstraction so usecase doesn't need to know where and how the data is retrived or editted.

## Import Rules:
---
- `Repositories` **CAN** import `entities` package.
- `Repositories` **MUST NOT** import any other package beside above import rules.
- `Repositories` **CAN** be imported by `cmd` package.
- `Repositories` **MUST NOT** be imported by other package.

