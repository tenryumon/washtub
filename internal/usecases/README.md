# Backend Usecases

Folder `usecases` consist objects to do business process like Login, Create new program, etc.
Usecases will get data from handlers, validate it, call repositories to get data or do some action then return models for response.


## What it do:
---

Folder `usecases` main job is to:
- Get input from handlers.
- Do Validation to the input.
- Get Data from Repository or Do Action through Repository.
- Return output to handlers.

## Import Rules:
---
- `Usecases` **CAN** import `entities`, `models` and `interfaces` package.
- `Usecases` **MUST NOT** import any other package beside above import rules.
- `Usecases` **CAN** be imported by `cmd` package.
- `Usecases` **MUST NOT** be imported by other package.

