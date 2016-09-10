FROM registry.dataos.io/datafoundry/golang:1.6.2
ENV DATAFOUNDRY_HOST_ADDR dev.dataos.io:8443
ENV DATAFOUNDRY_ADMIN_USER wangmeng5
ENV DATAFOUNDRY_ADMIN_PASS AsiainfoLDPwangmeng5
ENV NAMESPACE team12
ENV SERVICE_PORT 8080
EXPOSE 8080
ENV SERVICE_SOURCE_URL github.com/asiainfoLDP/datafoundry_slack
WORKDIR $GOPATH/src/$SERVICE_SOURCE_URL
ADD . .
RUN go build
CMD ["sh", "-c", "./datafoundry_slack"]
