# 🐋+📦 dockerv

![build](https://github.com/theobori/dockerv/actions/workflows/build.yml/badge.svg)

A simple to use (KISS) CLI to backup Docker volumes.

## 📖 How to build and run ?

1. Install the dependencies
    - `go`
    - `make` (for tests)

2. Install the binary
   
```bash
go install github.com/theobori/dockerv@latest
```
3. Before starting to use it, make sure you have the need resources by doing:

```bash
dockerv init
```

## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## 🎉 Tasks

- [ ] Docker compose env var `COMPOSE_PROJECT_NAME`
- [ ] Target local driver only in docker compose files
- [x] tarball export
- [ ] zip export
- [x] Dynamic and scalable point identification
- [x] Command help
- [x] `import` without dest
- [x] Export output permissions
- [x] Support multi `--src`
