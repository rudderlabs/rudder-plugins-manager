# Changelog

## [0.8.0](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.7.0...v0.8.0) (2023-10-13)


### Features

* add support expr plugins ([31bf9aa](https://github.com/rudderlabs/rudder-plugins-manager/commit/31bf9aa31bc85c342858f6fbc29f944c2007915f))

## [0.7.0](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.6.1...v0.7.0) (2023-10-12)


### Features

* add has plugin method to manager ([25ff059](https://github.com/rudderlabs/rudder-plugins-manager/commit/25ff0599df51ccfe14d739871d8d0707bdcd8068))

## [0.6.1](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.6.0...v0.6.1) (2023-10-12)


### Miscellaneous

* downgrade benthos/v4 ([6265e5e](https://github.com/rudderlabs/rudder-plugins-manager/commit/6265e5ef39444d262eb89ff918a5c60cb09eab16))

## [0.6.0](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.5.1...v0.6.0) (2023-10-11)


### Features

* add workflow execution status ([c0ce788](https://github.com/rudderlabs/rudder-plugins-manager/commit/c0ce788a211bc8c80cc8eb0fadd2cc80c5142537))


### Miscellaneous

* clone to use go-clone ([86eb849](https://github.com/rudderlabs/rudder-plugins-manager/commit/86eb84947edac5950083d4d4619048fc8057c2f3))
* workflow status ([ec402a4](https://github.com/rudderlabs/rudder-plugins-manager/commit/ec402a4c23bc3ff3476a07529c1b579dc18a456f))

## [0.5.1](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.5.0...v0.5.1) (2023-03-24)


### Miscellaneous

* workflow to return partial completion output ([8228147](https://github.com/rudderlabs/rudder-plugins-manager/commit/82281471bf9a21b45d3a9ef7c301bab0d6c482d1))

## [0.5.0](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.4.0...v0.5.0) (2023-03-23)


### Features

* add support for reply of workflows ([f28f9ff](https://github.com/rudderlabs/rudder-plugins-manager/commit/f28f9ffb7bd4372b753f5ce55ac3cf7199dd3021))
* add support for step outputs ([6fa3a6c](https://github.com/rudderlabs/rudder-plugins-manager/commit/6fa3a6cfb4cd2295225abbefe1f3e171f2f56bcc))


### Miscellaneous

* add execution manager interface ([e55f01e](https://github.com/rudderlabs/rudder-plugins-manager/commit/e55f01efad9e8899593ec28ebbbcb66a64cec0d8))
* message ([2b8e27c](https://github.com/rudderlabs/rudder-plugins-manager/commit/2b8e27cdd32cc594d7c339437343a60ee946766c))
* update readme ([688b54d](https://github.com/rudderlabs/rudder-plugins-manager/commit/688b54d11032d3dbf929544edb8b6ccce364666f))
* update readme with examples ([1374c2b](https://github.com/rudderlabs/rudder-plugins-manager/commit/1374c2b3697ea0e320c632b43905d55537e618f7))

## [0.4.0](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.3.0...v0.4.0) (2023-03-16)


### Features

* add load workflow file util function ([56a8810](https://github.com/rudderlabs/rudder-plugins-manager/commit/56a8810378156842deb9298ba9e4985b8a2b4cec))


### Bug Fixes

* lint issue in message.go ([479802c](https://github.com/rudderlabs/rudder-plugins-manager/commit/479802c9c1acbb7a54644b23163cf1c9e5db76ee))


### Miscellaneous

* add log message ([23e47b4](https://github.com/rudderlabs/rudder-plugins-manager/commit/23e47b49965cd81e9f6a466ba4cb837aac31e3bb))
* bloblang plugin ([0a77cb7](https://github.com/rudderlabs/rudder-plugins-manager/commit/0a77cb730f4b3f32c0eaaeb32af502f3d013b9e8))
* manager ([9c6274e](https://github.com/rudderlabs/rudder-plugins-manager/commit/9c6274ef63b9bd638503d079defa0951e5d6f47d))
* message cloning ([c778b6e](https://github.com/rudderlabs/rudder-plugins-manager/commit/c778b6e9414fa740f07ffc3242762e1d6dde901f))
* message to map function ([84b3d18](https://github.com/rudderlabs/rudder-plugins-manager/commit/84b3d1876528cb0b5176a13cdc5dfaca08aa94d6))
* message to map function ([bba8347](https://github.com/rudderlabs/rudder-plugins-manager/commit/bba834792fa12246e37ed24c42f561c5046ce429))
* split into individual files ([24de94b](https://github.com/rudderlabs/rudder-plugins-manager/commit/24de94b7317ba0b36cb96fd6961843fb2aa78273))

## [0.3.0](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.2.2...v0.3.0) (2023-03-15)


### Features

* add retryable executors ([6f8b520](https://github.com/rudderlabs/rudder-plugins-manager/commit/6f8b520ae13183cd9ad2ddd056c1b2a199f3fabb))
* use generics for base manager ([d61df0b](https://github.com/rudderlabs/rudder-plugins-manager/commit/d61df0b10607f2cd35c3ff153b4b6230d30f5e33))


### Miscellaneous

* message clone test case ([cfa5586](https://github.com/rudderlabs/rudder-plugins-manager/commit/cfa558641738e2ec5da3f56c754e78c2dee502ba))

## [0.2.2](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.2.0...v0.2.2) (2023-03-14)


### Features

* add get name to pipeline interface ([3a8090a](https://github.com/rudderlabs/rudder-plugins-manager/commit/3a8090a33dc6ac6e49d26f7aa809e0ce1941c578))
* add github workflows ([3d1d0a0](https://github.com/rudderlabs/rudder-plugins-manager/commit/3d1d0a0d0ef44e801b017148200b38c667356b32))
* add json transform command ([93d69a5](https://github.com/rudderlabs/rudder-plugins-manager/commit/93d69a5ed1a216e578d4d2d630f0f5cefe6730a7))
* add pipeline types ([fd58d9c](https://github.com/rudderlabs/rudder-plugins-manager/commit/fd58d9c6bbbb2846cbbfc658e7497e8bdbd7b6c5))
* add plugin execute function ([e6749c0](https://github.com/rudderlabs/rudder-plugins-manager/commit/e6749c0d8292ec41992f22cfebbc57a8f5b5b825))
* add plugin manager and types ([0fe21d9](https://github.com/rudderlabs/rudder-plugins-manager/commit/0fe21d99579603a188ae75702de00637eb79aee5))
* add plugin types for easy using of this library ([010c45e](https://github.com/rudderlabs/rudder-plugins-manager/commit/010c45ec9e90f3c6f6096248c15381f72578069b))
* add workflow engine ([8a8eeae](https://github.com/rudderlabs/rudder-plugins-manager/commit/8a8eeae851fbbc515eaf7f8fa064c3bafff7a1ae))


### Bug Fixes

* github workflows ([121fe8e](https://github.com/rudderlabs/rudder-plugins-manager/commit/121fe8e9bccb5b465f33f3fda4d37f46f7cfcc4d))


### release

* 0.2.1 ([79a63ba](https://github.com/rudderlabs/rudder-plugins-manager/commit/79a63badae268e61ae7256fcdcd7c77f30f9fe64))


### Miscellaneous

* add coverage badge ([48bad6d](https://github.com/rudderlabs/rudder-plugins-manager/commit/48bad6d89ce0db9894af3b4c4913c87f28f96390))
* add coverage badge ([3d433c5](https://github.com/rudderlabs/rudder-plugins-manager/commit/3d433c58e3803d29d59513d2b40e10c4a9d2ff94))
* add submit function pipeline manager ([ed82e65](https://github.com/rudderlabs/rudder-plugins-manager/commit/ed82e65c1bcc24cc111ef234c873ce5224e2233d))
* plugins and add tests ([cb95596](https://github.com/rudderlabs/rudder-plugins-manager/commit/cb95596d92970d6a153845754c60eab71cc8d767))
* release 0.1.0 ([d5ac936](https://github.com/rudderlabs/rudder-plugins-manager/commit/d5ac9365e308857ca3421ee7758b1bdd3eb03b25))
* release 0.2.0 ([693fdb4](https://github.com/rudderlabs/rudder-plugins-manager/commit/693fdb41168a2b241fa2fbbbb178c67dda331fc7))
* release 0.2.2 ([b57e559](https://github.com/rudderlabs/rudder-plugins-manager/commit/b57e559ec8947f80bfd5a8abc798243f31bebe75))
* remove pipeline types ([a492c6b](https://github.com/rudderlabs/rudder-plugins-manager/commit/a492c6b0246f2c81992a025a38c6e709a36d3735))
* update make file and github workflows ([76f48b5](https://github.com/rudderlabs/rudder-plugins-manager/commit/76f48b545daf74052e9da6e5dc0460f56463b81e))
* update readme ([86936fb](https://github.com/rudderlabs/rudder-plugins-manager/commit/86936fb01efc48f521e7aedb076f71873f589503))
* update readme ([d6e7743](https://github.com/rudderlabs/rudder-plugins-manager/commit/d6e7743a3e50e683e2fa31e5d7b0cac427f2329d))
* upgrade github actions ([046f92c](https://github.com/rudderlabs/rudder-plugins-manager/commit/046f92cfc52f15e29dbd9c1d141efc460b43500a))

## [0.2.0](https://github.com/rudderlabs/rudder-plugins-manager/compare/v0.1.0...v0.2.0) (2023-03-07)


### Features

* add plugin execute function ([e6749c0](https://github.com/rudderlabs/rudder-plugins-manager/commit/e6749c0d8292ec41992f22cfebbc57a8f5b5b825))
* add plugin types for easy using of this library ([010c45e](https://github.com/rudderlabs/rudder-plugins-manager/commit/010c45ec9e90f3c6f6096248c15381f72578069b))


### Miscellaneous

* plugins and add tests ([cb95596](https://github.com/rudderlabs/rudder-plugins-manager/commit/cb95596d92970d6a153845754c60eab71cc8d767))
* update make file and github workflows ([76f48b5](https://github.com/rudderlabs/rudder-plugins-manager/commit/76f48b545daf74052e9da6e5dc0460f56463b81e))
* upgrade github actions ([046f92c](https://github.com/rudderlabs/rudder-plugins-manager/commit/046f92cfc52f15e29dbd9c1d141efc460b43500a))

## 0.1.0 (2023-02-23)

### Features

* add github workflows ([3d1d0a0](https://github.com/rudderlabs/rudder-plugins-manager/commit/3d1d0a0d0ef44e801b017148200b38c667356b32))
* add json transform command ([93d69a5](https://github.com/rudderlabs/rudder-plugins-manager/commit/93d69a5ed1a216e578d4d2d630f0f5cefe6730a7))
* add plugin manager and types ([0fe21d9](https://github.com/rudderlabs/rudder-plugins-manager/commit/0fe21d99579603a188ae75702de00637eb79aee5))


### Bug Fixes

* github workflows ([121fe8e](https://github.com/rudderlabs/rudder-plugins-manager/commit/121fe8e9bccb5b465f33f3fda4d37f46f7cfcc4d))
