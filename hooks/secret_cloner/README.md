# Secret Cloner 
## Overview
**Helm hook to manage a kubernetes secret lifecycle for application development**

Secret cloner 
1. Clones a secret to create a new one based on the label provided 
2. Updates the values of the secret with the given key values 

###Sample Cloner job 
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: sec
  annotations:
    "helm.sh/hook": pre-install,pre-upgrade
    "helm.sh/hook-weight": "1"
  namespace: secret-cloner
spec:
  backoffLimit: 0
  template:
    spec:
      containers:
        - env:
            - name: ACTION
              value: clone
            - name: NAMESPACE
              value: 'demo'
            - name: SECRETNAME
              value:  'demo-secret'
            - name: SECRETSUFFIX
              value: 'alice'
          image: 'docker/razorpay/devstack:sec'
          imagePullPolicy: IfNotPresent
          name: sec
      imagePullSecrets:
        - name: registry
      nodeSelector:
        node.kubernetes.io/worker-generic: ''
      restartPolicy: Never
```

On execution of the hook a new secret `demo-secret-alice` with the values of `demo-secret`

###Sample Updater Job 
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: sec-updater
  annotations:
    "helm.sh/hook": pre-install,pre-upgrade
    "helm.sh/hook-weight": "2"
  namespace: secret-cloner
spec:
  backoffLimit: 0
  ttlSecondsAfterFinished: 0
  template:
    spec:
      containers:
        - image: 'docker/razorpay/devstack:sec'
          imagePullPolicy: IfNotPresent
          name: sec
          volumeMounts:
          - name: config-volume
            mountPath: /src/config
      imagePullSecrets:
        - name: registry
      nodeSelector:
        node.kubernetes.io/worker-generic: ''
      volumes:
        - name: config-volume
          configMap:
            name: sec-updater
      restartPolicy: Never
---
apiVersion: v1
kind: ConfigMap
data:
  app.yaml: |
    updateEntries:
      s1:
        key: SAMPLE_KEY
        value: 'sample_value'
    action: update
    secretName: demo-secret-alice
    namespace: demo
metadata:
  labels:
    app: sec-updater
  name: sec-updater
  namespace: secret-cloner
```

On execution of the hook the secret `demo-secret-alice` would have the key value `SAMPLE_KEY:sample_value` added

