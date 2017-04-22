FROM scratch
COPY ca-certificates.crt /etc/ssl/certs/
ADD apiserver apiserver
EXPOSE 80 443
ENTRYPOINT ["/apiserver"]

