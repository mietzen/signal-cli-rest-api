# How to Debug

## Build debugable container

Build with debug flags `-gcflags='all=-N -l'` (this will disable compiler optimizations and inlining) and add dvl Debugger to the image:

```shell
docker build -t signal-cli-rest-api:$(git rev-parse --abbrev-ref HEAD) --build-arg BUILD_VERSION_ARG=0.0.1 --build-arg DEBUG_BUILD=true . --no-cache

docker build -t signal-cli-rest-api:debug --build-arg IMAGE=signal-cli-rest-api:$(git rev-parse --abbrev-ref HEAD) -f ./Dockerfile-debug . --no-cache
```

Be sure to add:
```yaml
    ports:
      - 2345:2345
```
to your `docker-compose.yml` to open the debugger port. 

**Beware signal-cli-rest-api only starts if you attach the debugger!!!**



## Debug with vscode

To attach to the debugger using vscode use the following `launch.json`

```yaml
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Attach to signal-cli-rest-api",
            "type": "go",
            "request": "attach",
            "mode": "remote",
            "cwd": "${workspaceFolder}/src",
            "remotePath": "/tmp/signal-cli-rest-api-src",
            "port": 2345,
            "host": "127.0.0.1",
            "showLog": true,
            "logOutput": "rpc"
        }
    ]
}
```

## Dev-Cycle
- Make changes
- Build the debug image
- Start your stack with `compose up`
- Attach the Debugger and test your changes
- Repeat
