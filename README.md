# devstack
Devstack is Razorpay's Developer Experience Solution for cloud on laptop

## What is Devstack
At razorpay, we run all our workloads on kubernetes. Like any other mature organization, we have an involved CI/CD practice with extremely sophisticated pipelines.  While this works great for all production and pre-production workloads, we have been noticing over a period of time a bunch of development challenges. 
In essence, the goal is to **Simplify developer workflow and reduce the time taken to rollout features independently**. 
Devstack, offers a set of tools to help build and develop code on the individual developer's laptop, as if they are working
on a cloud environment. 

In a nutshell: **"Its a client based development tool for building cloud native applications on kubernetes"**

## Presentations and Videos
- [Devstack Presentation: Slides](https://static.sched.com/hosted_files/osselc21/50/Improving_Developer_Experience_Srinidhi_VV_09292021_v1.pdf)

## Design Goals
### High Level Goals
- Streamlined Dev Workflow: Provide a streamlined workflow and faster merges to `main` or `master` branches. 
- Consistent Environment: Provide a seamless consistent environment across dev, stage, pre-prod and production environments
- Faster Feedback: Reduce time to write and build containerized applications. Enable faster feedback loop on local development environment

### Design Choices
- Remove vendor lock-in (rely on OSS practically)
- Kubernetes native (At the moment, we don't have extensions for non K8s solutions). Our environment is kubernetes native
- Hassle Free onboarding - minimal changes to application and development lifecycle
- Cost Effective - We should eventually be able to bill developers and teams on usage patterns. 
- Slightly Opinionated - This isn't a PaaS offering at the momemt. And hence can be extended and deployed on any native kubernetes installation

### Features
- [x] Build, Test and Deploy applications from laptop directly into kubernetes using simple CLI tools
- [x] Ship code to remote container without tunneling: Sync files into container directly using File Sync(using [Devspace](https://github.com/loft-sh/devspace))
- [x] Ability to provide hot reloading of apps : sync directly into the containers without restart(e.g. [CompileDaemon](https://github.com/githubnemo/CompileDaemon) for statically typed languages)
- [x] Support out-of-the-box support using existing available [helm hooks](https://helm.sh/docs/topics/charts_hooks/#hooks-and-the-release-lifecycle) 
- [x] Support for extensible custom Helm Hooks that handles provisioning of AWS infrastructure components using [LocalStack](https://github.com/localstack/localstack)
- [x] Declarative-ly define service and service dependencies using [Helmfile](https://github.com/roboll/helmfile)
- [x] Traffic routing to right upstream. Uses traefik 2.0 [IngressRoute](https://doc.traefik.io/traefik/v2.0/providers/kubernetes-crd/)
- [x] Ability to selectively route traffic to different upstream. Done via [opentelemetry context/header propagation](https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/overview.md)
- [x] Ability to expose preview URL for all services

## Docs

### Pre-Requisites
As we mentioned earlier, our solution is slightly opinionated and requires the following stack:

- Cloud Provider : AWS(Note: Our custom infra helm hooks are all designed for AWS. However, it can be extended. See details below)
- Kubernetes: Requires verion 1.15+
- Traefik: 2.0+ to be deployed on the kubernetes cluster above. Please see the [official installation instructions](https://doc.traefik.io/traefik/v2.0/getting-started/install-traefik/)
- Helm: 3.0+
- LocalStack: To be deployed on kubernetes cluster above. Refer to [LocalStack](https://github.com/localstack/localstack#using-helm)

Other requirements (For hot-reload)
If you are using a loosely typed language like php / python etc, then you can safely skip this section. For static languages like golang, java, nodejs etc, please refer below:
- Golang: [CompileDaemon](https://github.com/githubnemo/CompileDaemon)
- NodeJs: [Nodemon](https://www.npmjs.com/package/nodemon)
- Java: [GradleDaemon](https://docs.gradle.org/current/userguide/gradle_daemon.html) or [MavenDaemon](https://github.com/mvndaemon/mvnd)

### Architecture
Please refer to the [Architecture Overview](Architecture.md#) for entire details on the devstack architecture.

### Examples
TBD

### Extension
We have provided a collection of custom helm hooks for AWS and kubernetes specific workloads. All of these are extensible. Please refer to the documentation of the helm hooks(TBD)

## Need Help

Please file an issue on this repo using the following labels: Clarification, Feature, Bug

## Contributing
Please refer to the [Contribution Guide](https://github.com/razorpay/devstack/blob/master/CONTRIBUTING.md)

## Roadmap
TBD
