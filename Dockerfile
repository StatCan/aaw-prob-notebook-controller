# Build with the golang image
FROM golang:1.16-alpine AS build

# Add git
RUN apk add git

# Set workdir
WORKDIR /work

# Add dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build

# Generate final image
FROM scratch
COPY --from=build /work/prob-notebook-controller /prob-notebook-controller
ENTRYPOINT [ "/prob-notebook-controller" ]
