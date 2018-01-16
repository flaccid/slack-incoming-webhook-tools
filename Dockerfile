FROM centurylink/ca-certs

COPY bin/siwp /usr/local/bin/siwp

WORKDIR /usr/local/bin

CMD ["siwp"]
