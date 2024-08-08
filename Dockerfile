FROM ubuntu

RUN apt-get update

VOLUME ["/runtime"]
COPY ./disperser/bin/combined /bin/combined
WORKDIR /runtime
CMD /runtime/run.sh
