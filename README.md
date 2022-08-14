# git-pipeline-dashboard

[![Build Status](https://img.shields.io/github/workflow/status/jangraefen/git-pipeline-dashboard/Build/main?logo=GitHub)](https://github.com/jangraefen/git-pipeline-dashboard/actions?query=workflow:Build%20branch:main)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/jangraefen/git-pipeline-dashboard)](https://pkg.go.dev/mod/github.com/jangraefen/git-pipeline-dashboard)
[![Coverage](https://img.shields.io/codecov/c/github/jangraefen/git-pipeline-dashboard?logo=codecov)](https://codecov.io/gh/jangraefen/git-pipeline-dashboard)
[![Go Report Card](https://goreportcard.com/badge/github.com/jangraefen/git-pipeline-dashboard)](https://goreportcard.com/report/github.com/jangraefen/git-pipeline-dashboard)
[![Docker Pulls](https://img.shields.io/docker/pulls/jangraefen/git-pipeline-dashboard)](https://hub.docker.com/r/jangraefen/git-pipeline-dashboard)

A small dashboard application that provides an overview for git-based CI/CD pipelines like GitLab CI. While these tools are awesome, they usually not provide a good overview over multiple pipelines in multiple repositories, at least not for free.

A few alternatives do already exist, but none of them suited my personal needs, especially considering an enterpise context, so I decided to do my own implementation. Maybe it is of use for you aswell 😉.

## Configuration

Since the dashboard is designed to be served from a container, configuration is mainly done via environment variables.

```bash
PIPELINE_DASHBOARD_ENDPOINT=https://my-gitlab-host.com/api/v4
PIPELINE_DASHBOARD_GITLAB_TOKEN=XXX

PIPELINE_DASHBOARD_GROUPS=Test Group 1,Test Group 2

PIPELINE_DASHBOARD_TEST_GROUP_1_SOURCE=gitlab
PIPELINE_DASHBOARD_TEST_GROUP_1_NAMESPACES=some/group
PIPELINE_DASHBOARD_TEST_GROUP_1_USERS=SomeUsername

PIPELINE_DASHBOARD_TEST_GROUP_2_SOURCE=gitlab
PIPELINE_DASHBOARD_TEST_GROUP_2_REPOSITORIES=someRepo1,someRepo2
```

Settings these environment variables and starting the application, will lead to a dashboard with two collapsable groups with pipelines tiles inside of them. Selecting pipelines can be done by either scrapping a "namespace" or a user for all repositories or by manually selecting individual repositories.
