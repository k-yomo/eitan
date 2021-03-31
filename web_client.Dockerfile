FROM node:14.8.0-alpine
WORKDIR /usr/src/app

COPY ./src/web_client /usr/src/app

CMD npm run dev