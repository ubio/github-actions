FROM alpine:3
RUN apk add --update --no-cache bash ca-certificates curl git jq openssh
RUN ["bin/sh", "-c", "mkdir -p ./src"]
RUN curl --retry 30 --retry-max-time 120 -s -S -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize/v3.9.3/kustomize_v3.9.3_linux_amd64.tar.gz | tar -xz -C /usr/bin 
RUN chmod +x /usr/bin/kustomize
COPY ["./src", "/src/"]
RUN chmod +x /src/entrypoint.sh
ENTRYPOINT ["/src/entrypoint.sh"]