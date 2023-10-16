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

## ⭐ Use cases

Export Docker volumes

```bash
dockerv export \
    --src "relative_or_absolute_path/docker-compose.yml" \
    --src "relative_or_absolute_path2/" \
    --dest "volumes.tar.gz
```

Import Docker volumes

```bash
dockerv import \
    --src "volumes.tar.gz" \
    --src "volumes2.tar.gz" \
    --dest "relative_or_absolute_path/docker-compose.yml"
```

List recursively Docker volumes in docker compose files or tarball

```bash
dockerv list \
    --src "./" \
    --state
```

Copy Docker volume content to another 

```bash
dockerv copy \
    --src "volume_src" \
    --dest "volume_dest" \
    --force
```

Copy then remove the Docker volume 

```bash
dockerv move \
    --src "volume_src" \
    --dest "volume_dest" \
    --force
```

## 🎉 Tasks

- [x] tarball export
- [ ] zip export
- [x] Dynamic and scalable point identification
- [ ] Command help
- [ ] Custom volume destination for single source volume packed
- [ ] `import` without dest
- [x] Export output permissions
- [x] Support multi `--src`
