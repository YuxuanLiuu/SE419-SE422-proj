kubectl create -f locust-cm.yaml -f scripts-cm.yaml
sleep 1
kubectl create -f master-deployment.yaml -f service.yaml -f slave-deployment.yaml
