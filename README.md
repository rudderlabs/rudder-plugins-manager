[![codecov](https://codecov.io/gh/rudderlabs/rudder-plugins-manager/branch/main/graph/badge.svg?token=ErUmduv9C8)](https://codecov.io/gh/rudderlabs/rudder-plugins-manager)
<p align="center">
  <a href="https://rudderstack.com/">
    <img src="https://user-images.githubusercontent.com/59817155/121357083-1c571300-c94f-11eb-8cc7-ce6df13855c9.png">
  </a>
</p>

<p align="center"><b>The Customer Data Platform for Developers</b></p>

<p align="center">
  <b>
    <a href="https://rudderstack.com">Website</a>
    ·
    <a href="">Documentation</a>
    ·
    <a href="https://rudderstack.com/join-rudderstack-slack-community">Community Slack</a>
  </b>
</p>

---

## Rudder Plugins Manager

Plugins, Workflows and Managers
## Overview

Provides managers, interfaces and types to write generic plugins and workflows.

## Features

* Generics [Plugin Interfaces](./plugins/types.go)
* Useful Plugins to get started
 * [Base Plugin](./plugins/base.go)
 * [Bloblang Plugin](./plugins/bloblang.go) (using [bloblang](https://www.benthos.dev/docs/guides/bloblang/about))
 * [Workflow Plugin](./plugins/workflow.go)
 * [Retryable Plugin](./plugins/retryable.go)
 * [Orchestrator Plugin](./plugins/orchestrator.go)
* Managers
  * [Plugin Manager](./plugins/manager.go)
  * [Workflow Manager](./plugins/manager.go)

## Getting started

```bash
go get github.com/rudderlabs/rudder-plugins-manager

```

## Examples
* [Base Plugin](./plugins/base_test.go)
* [Bloblang Plugin](./plugins/bloblang_test.go)
* [Workflow Plugin](./plugins/workflow_test.go)
 * [Sample Workflow](./test_data/workflows/sample.yaml)
* [Retryable Plugin](./plugins/retryable_test.go)
* [Orchestrator Plugin](./plugins/orchestrator_test.go)
* Managers
* [Plugin Manager](./plugins/manager_test.go)
* [Workflow Manager](./plugins/manager_test.go)

## Contribute

We would love to see you contribute to RudderStack. Get more information on how to contribute [**here**](CONTRIBUTING.md).

## License

The RudderStack \*\*software name\*\* is released under the [**MIT License**](https://opensource.org/licenses/MIT).
