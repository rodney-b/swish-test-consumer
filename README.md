# Swish Analytics Test Consumer

The consumer ingests data from kafka published by the "swish publisher" and logs the data. That's it.

### Build and Deployment
Build and deployment is done automatically upon pushing changes to the `main` branch. However, the option for manual build and deployments is also there. 

### Under the hood
Configuration and deployment of the consumer is done using [Helm charts](https://helm.sh/docs).