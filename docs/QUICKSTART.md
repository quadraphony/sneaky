# Quickstart

This quickstart shows how to validate a config and run the CLI for basic verification.

1. Validate a config file:

   ```bash
   go run ./cmd/sneaky/validate <path/to/config.json>
   ```

2. Build the CLI:

   ```bash
   make build
   ```

3. (Dev) Run using Docker image:

   ```bash
   make docker-build
   docker run --rm -v $(pwd):/work -w /work quadraphony/sneaky:dev validate examples/singbox.example.json
   ```

4. Use the CLI `start`, `status`, `stop` commands once implemented by the runtime manager.
