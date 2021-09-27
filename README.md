# devstack
Devstack is Razorpay's Developer Experience Solution for cloud of laptop Solution

## What is Devstack
At razorpay, we run all our workloads on kubernetes. As with any mature organization, we have a fairly laid out CI and CD pipelines. While this works great for all production and pre-production workloads, we have been noticing over a period of time a bunch of issues. In essnece, the goal is to *Simplify developer workflow and reduce time to rollout features independently*. Devstack, offers a set of tools to help build and develop code on the individual developer's laptop, as if they are working
on a cloud environment. Its a client based development tool for building cloud native applications on kubernetes

## Presentations and Videos
TBD

## Design Goals
### High Level Goals
- Streamlined Dev Workflow: Provide a streamlined workflow and faster merges to `main` or `master` branches. 
- Consisten Environment: Provide a seamless consistent environment across dev, stage, pre-prod and production environments
- Faster Feedback: Reduce time to write and build containerized applications. Enable faster feedback loop on local development environment

### Design Choices
- Remove vendor lock-in (rely on OSS practically)
- Kubernetes native (At the moment, we don't have extensions for non K8s solutions). Our environment is kubernetes native
- Hassle Free onboarding - minimal changes to application and development lifecycle
- Cost Effective - We should eventually be able to bill developers and teams on usage patterns. 
- Slightly Opinionated - This isn't a PaaS offering at the momemt. And hence can be extended and deployed on any native kubernetes installation

### Features
- [ ] Build, Test and Deploy applications from laptop directly into kubernetes using simple CLI tools
- [ ] Ship code to remote container without tunneling: Sync files into container directly using File Sync(using [Devspace](https://github.com/loft-sh/devspace))
- [ ] Ability to provide hot reloading of apps : sync directly into the containers without restart
- [ ] Support for extensible custom Helm Hooks that handles provisioning of AWS infrastructure components using [LocalStack](https://github.com/localstack/localstack)
- [ ] Declaratively define service and service dependencies using [Helmfile](https://github.com/roboll/helmfile)
- [ ] Traffic routing to right upstream (Using traefik 2.0 with custom ingress route)
- [ ] Ability to selectively route traffic to different upstream(using header propogation)
- [ ] Ability to expose preview URL for all services

## Docs

## Need Help

Please file an issue on this repo using the following labels: Clarification, Feature, Bug

## Contributing
Please refer to the [Contribution Guide](https://github.com/razorpay/devstack/blob/master/CONTRIBUTING.md)

## Roadmap
TBD
