FROM alpine:3.2
ADD test-api /test-api
ENTRYPOINT [ "/test-api" ]
