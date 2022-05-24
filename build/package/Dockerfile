FROM alpine:3.16
COPY weigh /bin/weigh
RUN apk --update add ca-certificates
RUN apk add tzdata
ENTRYPOINT ['/bin/weigh']
