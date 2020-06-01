FROM alpine

COPY ./migrate /migrate

COPY ./bin/ /

CMD /boilerplate
