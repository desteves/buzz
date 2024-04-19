# Buzz

A Gemini-powered Golang app that is gated with Google Auth (OAuth). 

This app also works with [Pulumi ESC](https://www.pulumi.com/product/esc/) for secrets management. [Pulumi ESC](https://www.pulumi.com/product/esc/) integrates with 1Password and many others. To configure Pulumi ESC with 1Password define an environment following the syntax shown [here](https://www.pulumi.com/docs/esc/providers/1password-secrets/). Example for the `buzz-dev-environment` Environment:

```yaml
values:
  1password:
    secrets:
      fn::open::1password-secrets:
        login:
          serviceAccountToken:
            fn::secret:
              ciphertext: ZXN...= // not shown
        get:
          google_oauth_client_id:
            ref: "op://dev-environment/buzz/username"
          google_oauth_client_secret:
            ref: "op://dev-environment/buzz/credential"
          docker_pat:
            ref: "op://dev-environment/dockerhub/password"
          docker_usr:
            ref: "op://dev-environment/dockerhub/username"
  environmentVariables:
    GOOGLE_OAUTH_CLIENT_ID: ${1password.secrets.google_oauth_client_id}
    GOOGLE_OAUTH_CLIENT_SECRET: ${1password.secrets.google_oauth_client_secret}
    DOCKER_PAT: ${1password.secrets.docker_pat}
    DOCKER_USR: ${1password.secrets.docker_usr}

```


## Prereqs

- Configured Google Project w/ OAuth and matching callback URL
- Env vars: `GOOGLE_OAUTH_CLIENT_ID` and `GOOGLE_OAUTH_CLIENT_SECRET` with correct values
- `DOCKER_DEFAULT_PLATFORM=linux/amd64`

## Run

### From source

```bash
$ cd app
####################################################
# sans secrets management, set envrs locally
$ go run main.go
####################################################
# OR with pulumi esc + 1password integration ✨🔐✨ #
# I have defined the environment buzz-dev-environment 
$ esc run buzz-dev-environment go run main.go
####################################################
```

### From docker image

```bash
####################################################
# sans secrets management, set envrs locally
$ docker run --platform linux/amd64  -p 8000:8000  -e GOOGLE_OAUTH_CLIENT_ID=my_id -e GOOGLE_OAUTH_CLIENT_SECRET=my_secret_value nullstring/buzz:dev
####################################################
# OR with pulumi esc + 1password integration ✨🔐✨ #
# I have defined the environment buzz-dev-environment 
$ esc run buzz-dev-environment  -- bash -c 'docker run --platform linux/amd64  -p 8000:8000  -e GOOGLE_OAUTH_CLIENT_ID=$GOOGLE_OAUTH_CLIENT_ID -e GOOGLE_OAUTH_CLIENT_SECRET=$GOOGLE_OAUTH_CLIENT_SECRET nullstring/buzz:dev'
####################################################
```

## Build image

```bash
$ cd app

####################################################
# sans secrets management, provide creds locally   #
$ docker login
####################################################
# OR with pulumi esc + 1password integration ✨🔐✨ #
# I have defined the environment buzz-dev-environment 
$ esc run buzz-dev-environment  -- bash -c 'echo "$DOCKER_PAT" | docker login -u $DOCKER_USR --password-stdin'
####################################################

$ TAG="nullstring/buzz:dev"
$ docker build . -t $TAG
$ docker push $TAG
```

## Deploy to GCP (Work in progress)

```bash
$ cd infra
$ pulumi up

```

## Credits

- OAuth piece of code originates from this [example](https://www.kungfudev.com/blog/2018/07/10/oauth2-example-with-go).
