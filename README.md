# RISKEN Common

![Build Status](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiZ09vd2t6d3hNNzl4aWp1Z2FTclBPMXBpb3pMSEFmR1dYclZHaVI0VVZCbEExb1poYkdZSHluZnVoS0V6VzJpY0IvaHVLb0ExV2NRaUcrOW95Z3E3TEhNPSIsIml2UGFyYW1ldGVyU3BlYyI6Ilc1aG0vemhscUlTVy9hVm8iLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=main)

`RISKEN` is a monitoring tool for your cloud platforms, web-site, source-code... 
`RISKEN Common` provides common packages, middleware components.

Please check [RISKEN Documentation](https://docs.security-hub.jp/).

## Installation

### Requirements

Common packages requires the following modules:

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)

### Packages

Common packages stored `pkg` directory. Sub-directories are created according to the purpose of the packages.

### Middleware

Middleware module stored `middleware` directory.
Build the containers on your machine with the following command:

```bash
$ cd middleware
$ make build
```

#### Running middlewares

Deploy the pre-built containers to the Kubernetes environment on your local machine.

- Follow the [documentation](https://docs.security-hub.jp/admin/infra_local/#risken) to download the Kubernetes manifest sample.
- Fix the Kubernetes object specs of the manifest file as follows and deploy it.

`k8s-sample/overlays/local/middleware.yaml`

| component | spec                                | before (public images)                          | after (pre-build images on your machine) |
| --------- | ----------------------------------- | ----------------------------------------------- | ---------------------------------------- |
| db        | spec.template.spec.containers.image | `public.ecr.aws/risken/middleware/db:latest`    | `middleware/db:latest`                   |
| queue     | spec.template.spec.containers.image | `public.ecr.aws/risken/middleware/queue:latest` | `middleware/queue:latest`                |

## Community

Info on reporting bugs, getting help, finding roadmaps,
and more can be found in the [RISKEN Community](https://github.com/ca-risken/community).

## License

[MIT](LICENSE).
