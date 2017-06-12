FROM scratch

EXPOSE 80
WORKDIR /app
COPY html /app/html
COPY linkhide /app
ENTRYPOINT ["/app/linkhide"]
