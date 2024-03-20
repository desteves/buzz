# Buzz

chatGPT-powered Golang app that is gated with Google Auth (OAuth). 

## Prereqs

- Configured Google Project w/ OAuth and matching callback URL
- Env vars: `GOOGLE_OAUTH_CLIENT_ID` and `GOOGLE_OAUTH_CLIENT_SECRET` with correct values

## Run

```bash
# sans secrets management
go run main.go

# with pulumi esc + 1password integration ✨🔐✨
esc run buzz-dev-environment go run main.go
```

Note: the later uses [Pulumi ESC](https://www.pulumi.com/product/esc/) for secrets management.

## Bundle

```bash
TAG="nullstring/buzz:dev"
docker build . -t $TAG
docker push $TAG

####################################################
# sans secrets management #
docker login
####################################################
# OR with pulumi esc + 1password integration ✨🔐✨ #
esc run buzz-dev-environment  -- bash -c 'echo "$PAT" | docker login -u $U --password-stdin'
####################################################

docker build . -t $TAG
docker push $TAG
```

Note: the later uses [Pulumi ESC](https://www.pulumi.com/product/esc/) for secrets management.

## Credits

- OAuth piece of code originates from this [blog](https://www.kungfudev.com/blog/2018/07/10/oauth2-example-with-go).
- <TBD>