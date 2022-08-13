FROM debian:11.4-slim

COPY credentials.tpl assume_role.sh /irsa-emu/

RUN apt-get update \
    && \
    # Install less, ca-certificates, curl, gettext-base, jq, tini, and unzip
    apt-get install -y --no-install-recommends \
      less=551-2 \
      openssl=1.1.1n-0+deb11u3 \
      ca-certificates=20210119 \
      libbrotli1=1.0.9-2+b2 \
      libsasl2-modules-db=2.1.27+dfsg-2.1+deb11u1 \
      libsasl2-2=2.1.27+dfsg-2.1+deb11u1 \
      libldap-2.4-2=2.4.57+dfsg-3+deb11u1 \
      libnghttp2-14=1.43.0-1 \
      libpsl5=0.21.0-1.2 \
      librtmp1=2.4+20151223.gitfa8646d.1-2+b2 \
      libssh2-1=1.9.0-2 \
      libcurl4=7.74.0-1.3+deb11u2 \
      curl=7.74.0-1.3+deb11u2 \
      libonig5=6.9.6-1.1 \
      libjq1=1.6-2.1 \
      gettext-base=0.21-4 \
      jq=1.6-2.1 \
      tini=0.19.0-1 \
      unzip=6.0-26+deb11u1 \
    && \
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64-2.0.30.zip" -o "awscliv2.zip" \
    && \
    unzip awscliv2.zip \
    && \
    ./aws/install \
    && \
    rm -rf awscliv2.zip aws \
    && \
    apt-get purge -y curl ca-certificates unzip \
    && \
    apt-get -y clean \
    && \
    rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/usr/bin/tini", "--", "/bin/bash", "/irsa-emu/assume_role.sh"]
