# Ingress Route Configurator 

## Overview
**Helm hook to configure the routes of ingressroute that enables header based routing**

The hook adds the rules based on the configuration provided

### Sample Hook Job 
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: irc
  annotations:
    "helm.sh/hook": post-install,post-upgrade
  namespace: ir-configurator
spec:
  backoffLimit: 0
  ttlSecondsAfterFinished: 0
  template:
    spec:
      containers:
        - env:
            - name: ACTION
              value: 'update'
            - name: HEADERVALUE
              value: 'alice'
            - name: INGRESSROUTENAME
              value: 'webapp'
            - name: INGRESSURL
              value: 'demo.dev.com'
            - name: NAMESPACE
              value: 'demo'
            - name: SERVICENAME
              value: 'alice-app'
            - name: SERVICEPORT
              value: '80'
          image: 'docker/razorpay:irc'
          imagePullPolicy: Always
          name: irc
      imagePullSecrets:
        - name: registry
      nodeSelector:
        node.kubernetes.io/worker-generic: ''
      restartPolicy: Never
```
###Sample Ingressroute Rule
```yaml 
kind: IngressRoute
apiVersion: traefik.containo.us/v1alpha1
metadata:
  name: demo
  namespace: demo
spec:
  entryPoints:
    - http
  routes:
    - kind: Rule
      match: Host(`demo.dev.com`)
      services:
        - name: demo-app
          port: 80
    - kind: Rule
      match: Host(`demo.dev.com`) && Headers(`uberctx-dev-serve-user`,`alice`)
      services:  
        - name: alice-app
          port: 80
```

###Options 
The configurator supports Update and Deletion of rules to the specified ingress route 

