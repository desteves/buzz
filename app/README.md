# Buzz Application

## Prereqs

- Google OAuth creds avaialble
- Gemini API Key available

## Run localy

### Sans secrets manager 🙈🙊🙊

> [!WARNING]
> You're storing sensitve data available to _everything_ running in your terminal session!

```bash
export GEMINI_API_KEY=123
export GOOGLE_OAUTH_CLIENT_ID=xyz
export GOOGLE_OAUTH_CLIENT_SECRET=987abc
go run .          
```

### With [Pulumi ESC](https://www.pulumi.com/docs/esc/), a secrets manager 🔐😎✅

- Store your secrets in a new ESC Environment

  ```bash
  ESC_ENV=buzz-app-env
  esc login
  esc env init $ESC_ENV
  esc env set $ESC_ENV --secret  environmentVariables.GEMINI_API_KEY 123abc
  esc env set $ESC_ENV  environmentVariables.GOOGLE_OAUTH_CLIENT_ID 123abc
  esc env set $ESC_ENV --secret  environmentVariables.GOOGLE_OAUTH_CLIENT_SECRET 123abc
  ```

- Run the Buzz app

  ```bash
  esc run buzz-app-env go run .
  ```

### With ✨🔐 1Password-stored secrets, accessed via [Pulumi ESC](https://www.pulumi.com/docs/esc/) 🚀🦾😎✅

- Store your enviornment variables in a 1Password Vault
- Create a [1Password service account](https://developer.1password.com/docs/service-accounts/) with read access to your vault
- Configure a Pulumi ESC Environment to [reference the 1Password-stored secrets](https://www.pulumi.com/docs/esc/integrations/dynamic-secrets/1password-secrets/):

  ```bash
  ESC_ENV=buzz-app-1p-env
  esc login
  esc env init $ESC_ENV
  esc env edit $ESC_ENV
  ```

- Paste the yaml contents below then save the changes

  > [!IMPORTANT]
  > Update the [secret `ref` syntax](https://developer.1password.com/docs/cli/secret-reference-syntax/) placeholders to match **your** 1Password Vault and items configuration
  > Update the `serviceAccountToken` value

  ```yaml
  values:
  1password:
    secrets:
      fn::open::1password-secrets:
        login:
          serviceAccountToken:
            fn::secret: ABC123
        get:
          google_oauth_client_id:
            ref: "op://dev-vault/google-oauth/username"
          google_oauth_client_secret:
            ref: "op://dev-vault/google-oauth/credential"
          gemini:
            ref: "op://dev-vault/google-gemini/credential"
  environmentVariables:
    GOOGLE_OAUTH_CLIENT_ID: ${1password.secrets.google_oauth_client_id}
    GOOGLE_OAUTH_CLIENT_SECRET: ${1password.secrets.google_oauth_client_secret}
    GEMINI_API_KEY: ${1password.secrets.gemini}
  ```

- Run the Buzz app

  ```bash
  esc run buzz-app-1p-env go run .
  ```
