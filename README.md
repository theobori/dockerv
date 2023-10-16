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
3. Then you can start using it

```bash
dockerv -h
```

## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## 🎉 Tasks

- [x] tarball export
- [ ] zip export
- [x] Dynamic and scalable point identification
- [x] Command help
- [x] `import` without dest
- [x] Export output permissions
- [x] Support multi `--src`
