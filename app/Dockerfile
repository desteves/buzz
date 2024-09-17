# syntax=docker/dockerfile:1

# Build Code
FROM --platform=linux/amd64 golang AS build-stage
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o buzz

# Unit Test
# FROM build-stage AS run-test-stage
# RUN  go test --cover ./... 

# Create Release 
FROM --platform=linux/amd64 gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /app
COPY --from=build-stage /app/buzz buzz
COPY --from=build-stage /app/static static

#  Environment variables placeholders to indicate app needs
ENV GOOGLE_OAUTH_CLIENT_ID=abc123 \
    GOOGLE_OAUTH_CLIENT_SECRET=xyz789 \
    GEMINI_API_KEY=gemini123  \
    REDIR=http://localhost:8000/app

ENTRYPOINT ["/app/buzz"]