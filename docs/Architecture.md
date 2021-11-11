# Architecture

The following section explains some of the high level architecture components of devstack.

## Solution Overview
At a high level the devstack solution works like this. Imagine we have a `Base Application(Or Master replicas)` of some app. Every developer / feature needs to essentially
have a copy of this base application with their own custom features. So traditional development process on a kubernetes cluster would be composed of the following:
- Create a new PR (on github / bitbucket etc)
- Build the docker image (using some CI tool like github actions)
- Push the docker image to a remote docker registry (dockerhub, ecs, gcs, harbor etc)
- Deploy the docker image to kubernetes (jenkins, spinnakar, argo, helm etc)

In addition, there might be other operational challenges on cloud providers / kubernetes, like the following:
- Create namespace, ingress, svc, load balancers etc for separate copies of the application
- Additionally in a cloud environment like AWS, if the app requires resources cloud specific resources like SQS, SNS etc, provision these through IaaS tools like terraform
- Clone kubernetes secrets from the base app if needed
- Provision other app resources like database, kafka, cache etc

Service Dependencies:
- In addtion to all the above, imagine additional use case of developers working through features that span across multiple services. In essence, the entire process
above needs to be repeated for every single service in the dependency chain. 

The solution tries to simplify this entire process by automating a clone of the base application using the following tools:
- [HelmFile]((https://github.com/roboll/helmfile)) - Allows to declaratively define service dependencies and also provides a CLI for applying these services on a kubernetes cluster
- [helm hooks](https://helm.sh/docs/topics/charts_hooks/#hooks-and-the-release-lifecycle) provide ways for injesting custom code as part a helm apply. We have provided a good number of ever growing custom helm hooks for all sorts
of infrastructure provisioning like SQS, SNS, Cloning Secrets, Creating Kafka Queues, Databases etc and a wide variety of use cases.
- [LocalStack](https://github.com/localstack/localstack) for reducing the burden of provisioning IaaS model of cloud resources. Localstack is native for AWS but similar such
solutions can be extended for other cloud providers as well.
- Preview URL: Over and above all of these, we also create a preview URL. So, image the same `Base Application` that is exposed on a URL called example.com. We have 2 features feat1 and feat2, with 
2 different devstack implementations, the solution (using helm hooks) can provision different urls called feat1.example.com and feat2.example.com
- Service Routing: The service routing is configured using Traefik 2.0's [IngressRoute](https://doc.traefik.io/traefik/v2.0/providers/kubernetes-crd/). In essence, as 
part of the feature deployment, the helm hooks configure a preview url alongside a feature label(e.g. `feature1`) in this case. So, if an application receives a header called `feature1` or uses an endpoint like feature1.example.com, 
  the ingress route sends the traffic to the feature 1 deployment/service. Else, its routed to the `Base Application`'s endpoint(in this case example.com). Lets explain with the following table:
  
  | Request                                              | Target Service |
  |------------------------------------------------------|----------------|
  | curl http://feature1.example.com                     | Feature 1      |
  | curl -H 'uber-ctx-key: feature1' http://example.com  | Feature 1      |
  | curl -H 'uber-ctx-key: feature2'  http://example.com | Feature 2      |
  | curl  http://example.com | Base Application      |

An interesting feature to note here is lets say example.com depends on example2.com as a downstream service. So, if `feature1` in this case requires development on both example.com and example2.com, then
just passing the header with `feature1` will automatically put the request flow as Feature 1(example.com) -> Feature 1(example2.com). We can see how this chain of dependencies
can be extended. At any point, whenever a targer service(say example5.com) doesn't have a requisite feature label deployed or enabled, it is automatically routed to the base version of the target service

With the above, lets see how the high level architecture for this would look like:

<p align="center">
  <img src="https://github.com/razorpay/devstack/blob/master/images/Solution-Overview.png?raw=true" alt="Solution Overview"/>
</p>

The below visualization shows how routing would look for a class of sample applications with different feature endpoints:
<p align="center">
<img src="https://github.com/razorpay/devstack/blob/master/images/Multi-App-Architecture.png?raw=true" alt="Multi Application Network Architecture" />
</p>


With all the above, overall devstack solution, appears like the following:
<p align="center">
<img src="https://github.com/razorpay/devstack/blob/master/images/Network-Architecture.png?raw=true" alt="DevStack Network Architecture" />
</p>

## Helmfile Workflow
The following diagram explains the helmfile worklow along with the custom helm hooks:
<p align="center">
<img src="https://github.com/razorpay/devstack/blob/master/images/Helm-Hooks-Workflow.png?raw=true" alt="Helmfile Workflow" />
</p>

