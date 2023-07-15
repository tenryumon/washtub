# Backend Handlers

Folder `handlers` consist objects as act as trigger.
In HTTP Server, handler usually means an endpoint. In Scheduler Binary, handlers usually means the trigger. In Queue Binary, handler means the message consumer.
Handler will get the input from external party (http form, http body, queue message) and push the data to usecase.
When get result from usecase, handler will either write response, requeue or finish the queue message, or finish the cron.


## What it do:
---

Folder `handlers` main job is to:
- Get input from external parties (http service).
- Get queue message (queue service).
- Get trigger from scheduler (scheduler service).
- Call usecase function.
- Write response (http service).
- Requeue or Finish the message (queue service).

## Import Rules:
---
- `Handlers` **CAN** import `entities`, `models` and `interfaces` package.
- `Handlers` **MUST NOT** import any other package beside above import rules.
- `Handlers` **CAN** be imported by `cmd` package.
- `Handlers` **MUST NOT** be imported by other package.

