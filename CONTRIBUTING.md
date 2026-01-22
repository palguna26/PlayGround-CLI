# Contributing to PlayGround CLI

Thank you for your interest in contributing to PlayGround CLI! We welcome contributions from the community to help make AI-assisted development safer and more accessible.

## Getting Started

1.  **Fork the repository** on GitHub.
2.  **Clone your fork** locally.
    ```bash
    git clone https://github.com/yourusername/PlayGround-CLI.git
    cd PlayGround-CLI
    ```
3.  **Install Go** (version 1.21 or higher).

## Development Workflow

1.  **Create a branch** for your feature or fix.
    ```bash
    git checkout -b feature/my-new-feature
    ```
2.  **Make your changes**.
3.  **Run tests** to ensure everything is working.
    ```bash
    go test ./...
    ```
4.  **Build the CLI** to verify compilation.
    ```bash
    go build -o pg ./cmd/pg
    ```
5.  **Commit your changes** with clear, descriptive messages.
6.  **Push to your fork** and submit a **Pull Request**.

## Project Structure

- `cmd/pg`: specific main entry point.
- `internal/`: Private application code.
  - `agent/`: Core agent logic, chat loop, tools.
  - `llm/`: LLM provider integrations (OpenAI, Gemini).
  - `cli/`: Cobra command definitions (`agent`, `setup`, etc.).
  - `workspace/`: Git and Snapshot workspace abstractions.
  - `session/`: Session state management.

## Code Style

- Follow standard Go conventions (use `gofmt`).
- Ensure code is well-documented.
- Add unit tests for new functionality where possible.

## Reporting Issues

If you find a bug or have a feature request, please open an issue on GitHub. Provide as much detail as possible, including steps to reproduce the issue.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
