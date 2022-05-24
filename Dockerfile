FROM alpine:3.16
COPY ./bin/linux/weigh /bin/weigh
RUN chmod 0700 /bin/weigh
RUN apk --update add ca-certificates tzdata
ENTRYPOINT ["/bin/weigh"]
