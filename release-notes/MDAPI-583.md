* **\[MDAPI-583]\[GO]** Fixed the errors of a build without cgo (`CGO_ENABLED=0`, no C compiler found, the Go
  Playground): instead of undefined internal symbols, it fails with one error,
  `undefined: dxfeed_graal_go_api_requires_cgo_set_CGO_ENABLED_1_and_install_a_C_compiler`.
