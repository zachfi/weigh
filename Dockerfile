FROM alpine:3.15 as certs
COPY ./weigh /bin/weigh
RUN chmod 0700 /bin/weigh
RUN mkdir /var/weigh
RUN apk --update add ca-certificates
RUN apk add libc6-compat
RUN apk add tzdata
ENTRYPOINT ["/bin/weigh"]
