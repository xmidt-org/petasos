# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]
- Updated spec file and rpkg version macro to be able to choose when the 'v' is included in the version. [#75](https://github.com/xmidt-org/petasos/pull/75)
- Patched failing Docker image, removed deprecated Maintainer information, fixed linter issues and enabled linte. [#107](https://github.com/xmidt-org/petasos/pull/107)

## [v0.1.7]
### Added
- use configured scheme filter. [#71](https://github.com/xmidt-org/petasos/pull/71)

## [v0.1.6]
### Changed
- Update mentions of the default branch from 'master' to 'main'. [#54](https://github.com/xmidt-org/petasos/pull/54)
- Update buildtime format in Makefile to match RPM spec file. [#61](https://github.com/xmidt-org/petasos/pull/61)
- Migrate to github actions, normalize analysis tools, Dockerfiles and Makefiles. [#62](https://github.com/xmidt-org/petasos/pull/62)
- Add optional OpenTelemetry tracing feature. [#67](https://github.com/xmidt-org/petasos/pull/67) thanks to @utsavbatra5

## [v0.1.5]
### Added
- adding docker automation [#48](https://github.com/xmidt-org/petasos/pull/48)

### Changed
- switch dependency tooling from glide to go modules
- updated release pipeline to use travis [#47](https://github.com/xmidt-org/petasos/pull/47)
- register for specific OS signals [#51](https://github.com/xmidt-org/petasos/pull/51)

### Fixed
- dependency updates including webpa-common's which fixes the SD metric label value for a service [#51](https://github.com/xmidt-org/petasos/pull/51)

## [v0.1.4]
fixed build upload

## [v0.1.2]
Switching to new build process

## [0.1.1]
### Added
- Initial creation

[Unreleased]: https://github.com/Comcast/petasos/compare/v0.1.7...HEAD
[v0.1.7]: https://github.com/Comcast/petasos/compare/v0.1.6...v0.1.7
[v0.1.6]: https://github.com/Comcast/petasos/compare/v0.1.5...v0.1.6
[v0.1.5]: https://github.com/Comcast/petasos/compare/v0.1.4...v0.1.5
[v0.1.4]: https://github.com/Comcast/petasos/compare/v0.1.2...v0.1.4
[v0.1.2]: https://github.com/Comcast/petasos/compare/0.1.1...v0.1.2
[0.1.1]: https://github.com/Comcast/petasos/compare/0.0.0...0.1.1
