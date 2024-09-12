FROM --platform=$BUILDPLATFORM golang:1.22.1 as base

COPY . /ChargePi/src
WORKDIR /ChargePi/src

ARG READER_TYPE
ENV READER_TYPE=$READER_TYPE

ARG DEVICE
ENV DEVICE=$DEVICE

# Install dependencies
RUN chmod +x /ChargePi/src/scripts/install-dependencies.sh
RUN /ChargePi/src/scripts/install-dependencies.sh $READER_TYPE 0

ENV GOOS=$BUILDOS
ENV GOARCH=$BUILDARCH
ENV CGO=1

RUN go mod download
RUN go mod verify

# Statically link the binary, remove debugging tools and set the firmware version at compile time using ldflags
RUN go build -ldflags="-extldflags '-static' -s -w -X 'internal/chargepoint/charge_point.FirmwareVersion=v1.0.0'" --tags=$DEVICE -o ChargePi .

FROM node:19 as build-ui

COPY ./ui /ui
WORKDIR /ui

RUN npm install
RUN npm run build

# Test the client
FROM --platform=$BUILDPLATFORM vektra/mockery:v2.9.4 as test

COPY --from=base /ChargePi/src /ChargePi/src
WORKDIR /ChargePi/src

# Install dependencies
RUN chmod +x /ChargePi/src/scripts/install-dependencies.sh
RUN /ChargePi/src/scripts/install-dependencies.sh $READER_TYPE 0

RUN cp -r configs/ test/  \
    && cd test/  \
    && chmod +x create-test-certs.sh \
    && ./create-test-certs.sh

# Gen mocks
RUN mockery .

# Run tests
CMD ["go", "test","-v","./..."]

FROM --platform=$BUILDPLATFORM alpine as chargepi

WORKDIR /etc/ChargePi

ARG READER_TYPE
ENV READER_TYPE=$READER_TYPE

COPY --from=base /ChargePi/src/scripts /etc/ChargePi/scripts

# Install dependencies
RUN chmod +x ./scripts/install-dependencies.sh
RUN ./scripts/install-dependencies.sh $READER_TYPE 0

# Copy the binary and the configuration files
COPY --from=base /ChargePi/src/configs /etc/ChargePi/configs
COPY --from=base /ChargePi/src/ChargePi /usr/bin/ChargePi

# Copy the UI static files
COPY --from=build-ui /ui/build/ /usr/bin/ChargePi/ui/build

EXPOSE 8080
HEALTHCHECK --interval=5s --timeout=3s CMD curl --fail http://localhost:8080/healthz || exit 1

ENTRYPOINT ["ChargePi"]
CMD ["-c", "/etc/ChargePi/configs/config.yaml", "run"]