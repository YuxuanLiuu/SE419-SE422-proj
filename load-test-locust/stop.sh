kubectl delete -f locust-cm.yaml -f scripts-cm.yaml
kubectl delete -f master-deployment.yaml -f service.yaml -f slave-deployment.yaml
