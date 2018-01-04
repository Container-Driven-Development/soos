
FROM node:9.2.0

WORKDIR /build/app

ENV PATH=/build/node_modules/.bin:$PATH

ADD package.json /build/

RUN npm install && chmod -R 777 /build

RUN mkdir /.config /.cache && chmod -R 777 /.config /.cache

ENTRYPOINT [ "npm" ]

CMD ["start"]
    