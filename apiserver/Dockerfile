FROM scratch
COPY ca-certificates.crt /etc/ssl/certs/
ADD apiserver /
EXPOSE 80
ENTRYPOINT ["/apiserver"]

