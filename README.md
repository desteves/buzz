# Buzz

A chatGPT-powered Golang app that is gated with Google Auth (OAuth). 

This app also works with [Pulumi ESC](https://www.pulumi.com/product/esc/) for secrets management. [Pulumi ESC](https://www.pulumi.com/product/esc/) integrates with 1Password and many others.

## Prereqs

- Configured Google Project w/ OAuth and matching callback URL
- Env vars: `GOOGLE_OAUTH_CLIENT_ID` and `GOOGLE_OAUTH_CLIENT_SECRET` with correct values

## Run

```bash
####################################################
# sans secrets management #
go run main.go
####################################################
# OR with pulumi esc + 1password integration ✨🔐✨ #
esc run buzz-dev-environment go run main.go
####################################################
```

## Bundle

```bash
####################################################
# sans secrets management #
docker login
####################################################
# OR with pulumi esc + 1password integration ✨🔐✨ #
esc run buzz-dev-environment  -- bash -c 'echo "$PAT" | docker login -u $U --password-stdin'
####################################################

TAG="nullstring/buzz:dev"
docker build . -t $TAG
docker push $TAG
```

## Credits

- OAuth piece of code originates from this [blog](https://www.kungfudev.com/blog/2018/07/10/oauth2-example-with-go).
