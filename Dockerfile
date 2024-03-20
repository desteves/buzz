# syntax=docker/dockerfile:1

# Build Code
FROM golang AS build-stage
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /buzz

# Unit Test
# FROM build-stage AS run-test-stage
# RUN  go test --cover ./... 

# Create Release 
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /buzz /buzz
ENTRYPOINT ["/buzz"]