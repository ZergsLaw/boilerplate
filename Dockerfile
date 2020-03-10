FROM alpine

COPY ./migration /migration

COPY ./bin/ /

CMD /boilerplate
