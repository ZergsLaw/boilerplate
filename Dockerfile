FROM scratch

COPY ./migration /migration

COPY ./bin/ /

CMD /boilerplate
