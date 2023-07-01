if ( -not (Test-Path -Path "openapi-generator-cli.jar"))
{
    Invoke-WebRequest `
        -OutFile openapi-generator-cli.jar `
        https://repo1.maven.org/maven2/org/openapitools/openapi-generator-cli/6.6.0/openapi-generator-cli-6.6.0.jar
}

java -jar openapi-generator-cli.jar generate `
    -i "https://stoplight.io/api/v1/projects/spacetraders/spacetraders/nodes/reference/SpaceTraders.json" `
    -o "../stapi" `
    -g go `
    --git-repo-id "space-traders-go" `
    --git-user-id "bgreen" `
    --minimal-update `
    -p enumClassPrefix=true `
    -p packageName=stapi `
    -p isGoSubmodule=true
