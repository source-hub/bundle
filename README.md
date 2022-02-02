# bundle
Cluster things and share them to the internet

# Setting up using docker compose
`docker-compose up --build -d`

# Setting up using kubernetes
- Install [kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
- Install [helm](helm.sh) - package manager for k8s
- Create a k8s cluster using: `kind create cluster --name bundle-cluster --config cluster-config/2-workers.yaml`
- Deploy the app using helm:
  `helm install bundle -n bundle-ns --create-namespace helm`
- Deploy the ingress controller:
  `kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml`
- Get the ip of the control plane: `k get nodes -o wide`
- add the ip to /etc/hosts