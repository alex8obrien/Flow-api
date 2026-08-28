# Architecture

- `cmd`: Contains folders with `main` functions for each application.
  - `api`: Main application and web server.
  - `cli`: Companion CLI application to interact with the web server.
- `docs`: Contains all documentation for the project
- `hooks`: Git hooks
- `internal`: Contains the majority of the application code
  - `controller`: Echo v5 HTTP handler.
  - `db`: PostgreSQL pool initialisation.
  - `domain`: Common models and interfaces for core application processes.
  - `environment`: Loads `.env` files and validates variables are initialised.
  - `logger`: Application logger
  - `repository`: User & mock repository implementations.
  - `server`: Web server configuration.
  - `service`: Main business logic processing code.
  - `telemetry`: Definitions and interfaces for logs, traces and metrics.
    - `otel`: Open Telemetry implementation.
  - `test`: Shared mock implementations for unit tests.
  - `util`: Shared project utilities.
