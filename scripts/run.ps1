param (
    $command
)

if (-not $command)  {
    $command = "start"
}

$ProjectRoot = "${PSScriptRoot}/.."

$env:SURGEON_API_ENVIRONMENT="Development"
$env:SURGEON_API_PORT="8080"

switch ($command) {
    "start" {
        go run ${ProjectRoot}/cmd/surgeon-api-service
    }
    "openapi" {
        docker run --rm -ti -v ${ProjectRoot}:/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    }
    default {
        throw "Unknown command: $command"
    }
}