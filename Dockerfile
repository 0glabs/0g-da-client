FROM ubuntu

RUN apt-get update
RUN apt-get install -y ca-certificates

VOLUME ["/runtime"]
COPY ./disperser/bin/combined /bin/combined
WORKDIR /runtime
CMD /runtime/run.sh
