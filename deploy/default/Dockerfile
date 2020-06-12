FROM alpine:latest

ARG BUILD_DATE
ARG VCS_URL
ARG VCS_REF

LABEL maintainer="devs@joinhoney.com" \
      description="Go Gateway" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-url=$VCS_URL \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.schema-version="1.0"

RUN addgroup -S honey && adduser -h /opt/app -G honey -S honey
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

USER honey

ADD gateway/schema /opt/app/schema
WORKDIR /opt/app

COPY build/gateway_docker ./gateway

ENTRYPOINT ["./gateway"]
CMD ["start"]