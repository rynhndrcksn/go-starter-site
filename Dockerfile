FROM scratch
WORKDIR /app
COPY ./bin/linux_amd64_web main
EXPOSE 4000
CMD ["/app/main"]
