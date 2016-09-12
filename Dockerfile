FROM alpine

ADD cloud-config-server /cloud-config-server

ENV PORT 3000
ENV WORKDIR .

CMD /cloud-config-server --port $PORT --work-dir $WORKDIR
