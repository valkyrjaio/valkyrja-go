<p align="center"><a href="https://valkyrja.io" target="_blank">
    <img src="https://raw.githubusercontent.com/valkyrjaio/art/refs/heads/master/long-banner/orange/go.png" width="100%">
</a></p>

# Valkyrja

[Valkyrja][Valkyrja url] is a Go framework for web and console applications.

Valkyrja (pronounced "Valk-ear-ya") is the Old Norse spelling for Valkyrie, a
mythical creature that would guide warriors to Valhalla (the afterlife and a
better place) after death. In a similar sense, the Valkyrja framework guides
your application to be in a better state. Fast, light, and robust, Valkyrja
does the heavy lifting so you can focus on your application.

<p>
    <a href="https://pkg.go.dev/github.com/valkyrjaio/valkyrja-go/v26"><img src="https://pkg.go.dev/badge/github.com/valkyrjaio/valkyrja-go/v26.svg" alt="Go Reference"></a>
    <a href="https://github.com/valkyrjaio/valkyrja-go/releases"><img src="https://img.shields.io/github/v/release/valkyrjaio/valkyrja-go" alt="Latest Version"></a>
    <a href="https://github.com/valkyrjaio/valkyrja-go/blob/26.x/go.mod"><img src="https://img.shields.io/badge/Go-1.26-orange" alt="Go Version"></a>
    <a href="https://github.com/valkyrjaio/valkyrja-go/blob/26.x/LICENSE.md"><img src="https://img.shields.io/github/license/valkyrjaio/valkyrja-go.svg" alt="License"></a>
    <a href="https://github.com/valkyrjaio/valkyrja-go/actions/workflows/ci.yml?query=branch%3A26.x"><img src="https://github.com/valkyrjaio/valkyrja-go/actions/workflows/ci.yml/badge.svg?branch=26.x" alt="CI Status"></a>
    <a href="https://coveralls.io/github/valkyrjaio/valkyrja-go?branch=26.x"><img src="https://coveralls.io/repos/github/valkyrjaio/valkyrja-go/badge.svg?branch=26.x" alt="Coverage Status"></a>
    <a href="https://sonarcloud.io/summary/new_code?id=valkyrjaio_valkyrja-go"><img src="https://sonarcloud.io/api/project_badges/measure?project=valkyrjaio_valkyrja-go&metric=sqale_rating" alt="Maintainability Rating"></a>
</p>

Port Status
-----------

The Go port is in progress, and this repository holds the scaffolding for it.
PHP is the reference implementation, and each component lands here in the order
the architecture sets: container, dispatch, event, application, CLI, HTTP, then
bin.

The section below describes the framework that every port implements. A
component is not in this repository until its own package exists. Read
[`PORTS.md`][ports url] for the state of each port.

What's Included
---------------

- **HTTP and CLI kernels** — unified application architecture serving both web
  requests and command-line invocations
- **Dependency injection container** — deferred bindings, contextual
  resolution, and generated data for fast resolution at runtime
- **Routing** — route registration with middleware pipelines for both HTTP and
  CLI
- **Event dispatcher** — decoupled event handling with typed listeners
- **Type system** — primitive wrappers, identifiers, models, and collections
- **Errors as values** — a failure returns an `error`, and each error type
  carries the `*Error` suffix rather than the `*Exception` name the other ports
  use

Installation
------------

### Start a New Application

The fastest way to start a new Valkyrja application is with the starter
application or the Sindri build tool:

- Use the [`valkyrja-starter-app-go`][starter url] GitHub template ("Use this
  template" button on the repository page)
- Or use [Sindri][sindri url] to scaffold a new project

### Add to an Existing Project

Add the framework as a dependency:

```bash
go get github.com/valkyrjaio/valkyrja-go/v26
```

The `/v26` suffix is part of the module path. Go's semantic import versioning
encodes a major version above 1 in the path, and the path tracks the annual
major version.

Documentation
-------------

Documentation is baked into the repository so you can browse it offline. Each
component carries its own `README.md`, and that document arrives with the
component. The areas are `http`, `cli`, `container`, and `event`.

Until then, read [`valkyrja-php`][php url]. It is the reference implementation,
and every port mirrors its structure and its naming.

Ecosystem
---------

Valkyrja is the core framework. Surrounding it is an ecosystem of related
projects in the Valkyrjaio organization:

- [**Sindri**][sindri url] — build tool and application creator
- [**Starter (App)**][starter url] — starter application for new projects
- [**golangci-lint**][lint url] — the shared lint configuration and the
  copyright header tool that every Go repository runs

See the [Valkyrjaio organization page][org url] for the complete listing.

Versioning and Release Process
------------------------------

Valkyrja follows [semantic versioning][semantic versioning url] with a major
release every year, and support for each major version for 2 years from the
date of release.

### Supported Versions

Bug fixes are provided until 3 months after the next major release. Security
fixes are provided for 2 years after the initial release.

| Version | Go   | Release | Bug Fixes Until | Security Fixes Until |
| :------ | :--- | :------ | :-------------- | :------------------- |
| 26      | 1.26 | Q3 2026 | Q2 2027         | Q1 2028              |
| 27      | 1.27 | Q1 2027 | Q2 2028         | Q1 2029              |
| 28      | 1.28 | Q1 2028 | Q2 2029         | Q1 2030              |

Contributing
------------

Valkyrja is an open-source, community-driven project. Thank you for your
interest in helping develop, maintain, and release it.

See [`CONTRIBUTING.md`][contributing url] for the submission process and
[`VOCABULARY.md`][vocabulary url] for the terminology used across Valkyrja.

Run the full gate before you open a pull request:

```bash
make ci
```

Security Issues
---------------

If you discover a security vulnerability within Valkyrja, please follow our
[disclosure procedure][security vulnerabilities url].

License
-------

Valkyrja is open-source software licensed under the
[MIT license][MIT license url]. See [`LICENSE.md`](./LICENSE.md).

[Valkyrja url]: https://valkyrja.io
[org url]: https://github.com/valkyrjaio
[sindri url]: https://github.com/valkyrjaio/sindri-go
[starter url]: https://github.com/valkyrjaio/valkyrja-starter-app-go
[lint url]: https://github.com/valkyrjaio/ci-golangcilint-go
[ports url]: https://github.com/valkyrjaio/architecture/blob/master/PORTS.md
[php url]: https://github.com/valkyrjaio/valkyrja-php
[contributing url]: https://github.com/valkyrjaio/.github/blob/26.x/CONTRIBUTING.md
[vocabulary url]: https://github.com/valkyrjaio/.github/blob/26.x/VOCABULARY.md
[security vulnerabilities url]: https://github.com/valkyrjaio/.github/blob/26.x/SECURITY.md
[semantic versioning url]: https://semver.org/
[MIT license url]: https://opensource.org/licenses/MIT
