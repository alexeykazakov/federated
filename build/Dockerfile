FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

LABEL maintainer "Devtools <devtools@redhat.com>"
LABEL author "Devtools <devtools@redhat.com>"

ENV PROXY_SERVICE=/usr/local/bin/proxy \
    USER_UID=1001 \
    USER_NAME=proxy \
    LANG=en_US.utf8

# install proxy binary
COPY build/_output/bin/proxy ${PROXY_SERVICE}

COPY build/bin /usr/local/bin
RUN  /usr/local/bin/user_setup

ENTRYPOINT ["/usr/local/bin/proxy"]

EXPOSE 6443

USER ${USER_UID}
