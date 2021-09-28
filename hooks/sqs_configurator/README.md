#SQS Configurator
## Overview
**Helm hook for SQS management**

SQS configurator 
1. Creates a SQS queue in the specified provider (AWS/Localstack)
2. Updates the kubernetes secret with the key specified if needed 

###Sample Configurator Job
```yaml
 apiVersion: batch/v1
 kind: Job
 metadata:
   name: sqs-configurator
   annotations:
     "helm.sh/hook": pre-install,pre-upgrade
     "helm.sh/hook-weight": "3"
   namespace: sqs-configurator
 spec:
   backoffLimit: 0
   ttlSecondsAfterFinished: 0
   template:
     spec:
       containers:
         - image: 'docker/razorpay/devstack:sqsc'
           imagePullPolicy: IfNotPresent
           name: irc
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
             name: sqs
       restartPolicy: Never
---
 apiVersion: v1
 kind: ConfigMap
 data:
   app.yaml: |
     queue:
       q1:
         name: demo_queue_alice
         secretKey: QUEUE_NAME
     updateSecret: true
     kubeSecret: demo-secret-alice
     namespace: demo
     provider: localstack
     enableEndpointPrefix: false
 metadata:
   name: sqs
   namespace: sqs-configurator
```
On execution the hook creates a queue with name `demo_queue_alice` in localstack and updates the secret `demo-secret-alice` with the queue name for key `QUEUE_NAME`
###Options 
*queue.\*.name* - The name of the queue 

*secretKey* - The key for the secret updation (Optional)

*provider* (AWS/Localstack) - The provider to create the SQS queue with 

*updateSecret* (true/false) - Enable / disable the updation of kubernetes secret 

*kubeSecret*  - The kubernetes secret name to be updated with the value (Optional)

*namespace* - The namespace of the secret (Optional)

*enableEndpointPrefix* (true/false) - Flag to update the secret with providers prefix URL (Optional)

Eg: 

i.  true updates the secret with value `https://localstack-services.dev.razorpay.com/000000000000/demo_queue_alice` 

ii. false updates the secret with value `demo_queue_alice`
 
*NOTE:* The optional fields are not required in case of update secret is false 