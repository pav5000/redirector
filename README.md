# Redirector
Simple util to redirect tcp connections. Created it to replace "redir".

## Installation

Requirements:
- [Make](https://ru.wikipedia.org/wiki/Make)
- [Docker](https://www.docker.com/)

Steps:
1. Copy `config.example.yml` -> `config.yml` and configure there your desired port redirects and log level.
2. Build the docker image (just run `make build`).
3. Start the service (`make start`).
4. The service is running and will continue to run even after a reboot.

## Commands

- `make build` builds the docker image, you should run it after each update of the code from the repo.
- `make start` starts the service container.
- `make stop` stops and removes the service container.
- `make restart` re-creates the container, you should run it after each change in `config.yml`.
- `make logs` shows latest logs, here you can see if there are any errors (ctrl+c to exit).
- `make status` checks if the container is running and displays it's current CPU/memory usage (ctrl+c to exit).

#### Feel free to open an issue if you find a bug.
