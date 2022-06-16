# RISKEN OSINT

![Build Status](https://codebuild.ap-northeast-1.amazonaws.com/badges?uuid=eyJlbmNyeXB0ZWREYXRhIjoiNlFlQkVnOU94ek9aWlMxamhYc0dkcGk5aExqVWtoOVV3eWJhbWQ5YWVhNkljSk5RR0h1SFpSd2VTUzdSMk8yU1czNUVPTVh3M01CdSt5bzZ0RXNrNlc4PSIsIml2UGFyYW1ldGVyU3BlYyI6InA3QzlhcWRlUENkM3hRV1oiLCJtYXRlcmlhbFNldFNlcmlhbCI6MX0%3D&branch=master)

`RISKEN` is a monitoring tool for your cloud platforms, web-site, source-code... 
`RISKEN OSINT` uses *Open Source Inttelligence* to collect security threat informations.

Please check [RISKEN Documentation](https://docs.security-hub.jp/).

## Installation

### Requirements

This module requires the following modules:

- [Go](https://go.dev/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Protocol Buffer](https://grpc.io/docs/protoc-installation/)

### Install packages

This module is developed in the `Go language`, please run the following command after installing the `Go`.

```bash
$ make install
```

### Building

Build the containers on your machine with the following command

```bash
$ make build
```

### Running Apps

Deploy the pre-built containers to the Kubernetes environment on your local machine.

- Follow the [documentation](https://docs.security-hub.jp/admin/infra_local/#risken) to download the Kubernetes manifest sample.
- Fix the Kubernetes object specs of the manifest file as follows and deploy it.

`k8s-sample/overlays/local/osint.yaml`

| service   | spec                                | before (public images)                         | after (pre-build images on your machine) |
| --------- | ----------------------------------- | ---------------------------------------------- | ---------------------------------------- |
| subdomain | spec.template.spec.containers.image | `public.ecr.aws/risken/osint/subdomain:latest` | `osint/subdomain:latest`                 |
| website   | spec.template.spec.containers.image | `public.ecr.aws/risken/osint/website:latest`   | `osint/website:latest`                   |

## Community

Info on reporting bugs, getting help, finding roadmaps,
and more can be found in the [RISKEN Community](https://github.com/ca-risken/community).

## License

[MIT](LICENSE).
