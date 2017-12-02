FROM kkarczmarczyk/node-yarn:8.0-wheezy

WORKDIR /build/app

ENV PATH=/build/node_modules/.bin:$PATH

ADD package.json /build/

RUN cd /build && ([[ -n ${https_proxy} ]] && yarn config set proxy ${https_proxy} -g || true) && yarn && chmod -R 777 /build

RUN mkdir /.config /.cache && chmod -R 777 /.config /.cache

ENTRYPOINT cd /build/app && rm -rf node_modules && mv /build/node_modules /build/app/ && yarn build
