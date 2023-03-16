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

Sometimes we want to write custom code specific to some user event and we should write it as a plugin so that the core logic is clean. This library helps write such plugins and workflows. It also helps to manage them using the Managers.

## Features

* [Plugin Interfaces](./plugins/types.go)
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
* Install `go get github.com/rudderlabs/rudder-plugins-manager`
* Create a plugin
```go
plugin := plugins.NewBasePlugin("no-op",
		plugins.ExecuteFunc(
			func(ctx context.Context, data *plugins.Message) (*plugins.Message, error) {
				return data, nil
			},
		),
	)
// Add to Manager
pluginManager := plugins.NewBasePluginManager()
pluginManager.Add(plugin)
// Execute Plugin
pluginManager.Execute(context.Background(), "no-op", plugins.NewMessage("some data"))
```
* Create a workflow
```go
workflow := lo.Must(plugins.NewBaseWorkflowPlugin(pluginManager, plugins.WorkflowConfig{
	Name: "test",
	Steps: []plugins.StepConfig{
		{
			Name:     "blobl",
			Bloblang: `root.test = "test"`,
		},
    {
			Name:     "plugin",
			Plugin: "no-op",
		},
	},
}))

// Add to Manager
workflowManager := plugins.NewBaseWorkflowManager()
workflowManager.Add(workflow)
// Execute Plugin
workflowManager.Execute(context.Background(), "test", plugins.NewMessage("some data"))
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

The RudderStack Plugins Manager is released under the [**MIT License**](https://opensource.org/licenses/MIT).
