apiVersion: v1
kind: ServiceAccount
metadata:
  name: notebook-auth-policy-checker
  namespace: daaas
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: notebook-auth-policy-checker
  namespace: daaas
  labels:
    apps.kubernetes.io/name: notebook-auth-policy-checker
spec:
  selector:
    matchLabels:
      apps.kubernetes.io/name: notebook-auth-policy-checker
  template:
    metadata:
      labels:
        apps.kubernetes.io/name: notebook-auth-policy-checker
      annotations:
        sidecar.istio.io/inject: 'false'
    spec:
      serviceAccountName: notebook-auth-policy-checker
      imagePullSecrets:
        - name: k8scc01covidacr-registry-connection
      containers:
      - name: notebook-auth-policy-checker
        imagePullPolicy: IfNotPresent
        image: k8scc01covidacr.azurecr.io/prob-notebook-controller:${IMAGE_SHA}
        resources:
          limits:
            memory: "256Mi"
            cpu: "500m"
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: notebook-auth-policy-checker
rules:
- apiGroups:
    - 'kubeflow.org'
  resources:
    - 'notebooks'
  verbs:
    - get
    - list
    - watch
- apiGroups:
    - security.istio.io
  resources:
    - 'authorizationpolicies'
  verbs:
    - get
    - list
    - watch
    - create
    - update
    - delete
- apiGroups:
    - ''
  resources:
    - 'events'
  verbs:
    - create
    - patch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: notebook-auth-policy-checker
subjects:
- kind: ServiceAccount
  name: notebook-auth-policy-checker
  namespace: daaas
roleRef:
  kind: ClusterRole
  name: notebook-auth-policy-checker
  apiGroup: rbac.authorization.k8s.io
