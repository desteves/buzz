# 🐝 Buzz 🤖

**☁️ Deployed to GCP | 🚀 Managed by Pulumi | 🔐 Secured with 1Password**

![Buzz Logo](./app/static/buzz.jpg)

## What is Buzz?

Buzz is a Gemini-powered, Google Auth (OAuth)-gated Golang web application. The application takes in a string as input and returns the [NATO](https://en.wikipedia.org/wiki/NATO_phonetic_alphabet) spelling of the input. Examples,

```plain
pulumi -> Papa, Uniform, Lima, Uniform, Mike, India
cool -> Charlie, Oscar, Oscar, Lima
```

A running version of the Buzz application can be found at [buzz.atxyall.com](https://buzz.atxyall.com/) **However** the OAuth has been configured to **ONLY** work with the author's email 😈.

## Run the app locally

- See the [App README](./app/README.md)

## Run the app in GCP

- See the [Infra README](./infra/README.md)

## Reference Material

- [1Password Blog post](https://blog.1password.com/1password-pulumi-developer-secrets-guide/)
- [Pulumi Workshop](https://github.com/pulumi/workshops/tree/main/1password-pulumi-esc)
