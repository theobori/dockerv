# 🐋+📦 dockerv

A simple to use (KISS) CLI to backup Docker volumes.

## 📖 How to build and run ?

1. Install the dependencies
    - `go`
    - `make` (for tests)

2. Install the binary
   
```bash
go install github.com/theobori/dockerv
```

## 🤝 Contribute

If you want to help the project, you can follow the guidelines in [CONTRIBUTING.md](./CONTRIBUTING.md).

## ⭐ Use cases

#### Export Docker volumes

```bash
dockerv export \
    --src "relative_or_absolute_path/docker-compose.yml" \
    --dest "volumes.tar.gz
```

#### Import Docker volumes

```bash
dockerv import \
    --src "volumes.tar.gz" \
    --dest "relative_or_absolute_path/docker-compose.yml"
```

#### List recursively Docker volumes in docker compose files

```bash
dockerv list \
    --src "./" \
    --state
```

## 🎉 Tasks

- [x] tarball export
- [ ] zip export
- [x] Dynamic and scalable point identification
- [ ] Documentation 80%
- [ ] Custom volume destination for single source volume packed
