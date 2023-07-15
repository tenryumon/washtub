# Backend CMD

Folder `cmd` will consist of `main()` function of a binary. It always start from this folder to initialize all dependencies needed for a service to run.

## What it do::
---

Folder `cmd` main job is to:
- Read configuration file 
- Initialize dependencies needed by the service
- Initialize Repositories, Usecases and Handlers needed by the service
- Initialize HTTP Endpoint (For http service)
- Start the service

