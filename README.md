# Swish Analytics Test Consumer

The consumer ingests data from kafka published by the "swish publisher" and logs the data. That's it.

### Build and Deployment
Build and deployment is done automatically upon pushing changes to the `main` branch. However, the option for manual build and deployments is also there. 

### Under the hood
Configuration and deployment of the consumer is done using [Helm charts](https://helm.sh/docs).

### Local deployment
If deploying locally, a skaffold.yaml and env template is provided for convenience. Skaffold overlays [local.values.yaml](swish-test-consumer-ops/consumer-chart/charts/local.values.yaml) over the default [values.yaml](swish-test-consumer-ops/consumer-chart/values.yaml) so any configuration changes to the values file can go here.