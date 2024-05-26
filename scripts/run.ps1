param (
    $command
)

if (-not $command)  {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:SURGEON_API_ENVIRONMENT="Development"
$env:SURGEON_API_PORT="8080"
$env:SURGEON_API_MONGODB_USERNAME="root"
$env:SURGEON_API_MONGODB_PASSWORD="neUhaDnes"

function mongo {
    docker compose --file ${ProjectRoot}/deployments/docker-compose/compose.yaml $args
}

switch ($command) {
    "start" {
        try {
            mongo up --detach
            go run ${ProjectRoot}/cmd/surgeon-api-service
        } finally {
            mongo down
        }
    }
    "mongo" {
        mongo up
    }
    "openapi" {
        docker run --rm -ti -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    "docker" {
        docker build -t xnemect/surgeon-wl-webapi:local-build -f ${ProjectRoot}/build/docker/Dockerfile .
    }
    default {
        throw "Unknown command: $command"
    }
}