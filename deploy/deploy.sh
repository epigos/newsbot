#!/bin/sh

echo "Deploying version $VERSION"

# connect kubectl cluster
gcloud container clusters get-credentials cluster-epigos --zone europe-west2-a --project ${PROJECT_ID}
# create namespace with version
kubectl create namespace ${VERSION}
# add kube config map
sed 's/VERSION.*/VERSION="'${VERSION}'"/g' env/${ENV}.env > env/${ENV}.tmp
eval "kubectl delete configmap env-config -n ${VERSION}"
eval "kubectl create configmap env-config -n ${VERSION} $(cat env/${ENV}.tmp | sed -e 's/^/--from-literal=/' | tr "\n" ' ')"
eval "rm env/${ENV}.tmp"
# add web deployment
kubectl apply -f deploy/kube/web.yaml,deploy/kube/worker.yaml -n ${VERSION}
# # set image
kubectl set image deployment web -n ${VERSION} web=${IMAGE_PATH}
kubectl set image deployment worker -n ${VERSION} worker=${IMAGE_PATH}
# # apply patch
PATCH="{\"spec\":{\"template\":{\"metadata\":{\"annotations\":{\"date\":\"$(date +'%s')\"}}}}}"
kubectl patch deployment web -n ${VERSION} -p ${PATCH}
kubectl patch deployment worker -n ${VERSION} -p ${PATCH}
# create datastore index
gcloud datastore create-indexes models/index.yaml