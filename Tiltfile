allow_k8s_contexts('minikube')

k8s_yaml(['k8s/nats.yaml', 'k8s/nginx.yaml'])

local_resource(
    'build-frontend',
    'cd frontend && npm run build',
    deps=['frontend/src']
)

k8s_resource(
    'nginx-dep',
    resource_deps=['build-frontend'],
    port_forwards='8080:80'
)

def build_service(name):
    docker_build(
        'localhost:54886/{}-service'.format(name),
        context='backend',
        dockerfile='backend/services/{}/tilt.Dockerfile'.format(name),
        live_update=[
            sync('./backend/services/{}'.format(name), '/build'),
            sync('./backend/shared', '/build/../shared'),
            run('go build -o service ./cmd/main.go', trigger=[
                './backend/services/{}/**.go'.format(name),
                './backend/shared/**.go'
            ])
        ]
    )
    k8s_yaml('k8s/{}.yaml'.format(name))
    k8s_resource(
        '{}-service'.format(name),
        resource_deps=['nats']
    )

services = ['auth', 'user', 'course', 'billing', 'enrollment', 'notification', 'analytics', 'progress']
for service in services:
    build_service(service)

k8s_resource(
    'nats',
    port_forwards=['4222:4222', '8222:8222']
)