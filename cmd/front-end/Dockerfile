FROM node:17-buster as node
COPY ./front/public /entry-task/public
COPY ./front/src /entry-task/src
COPY ./front/package.json /entry-task/
COPY ./front/package-lock.json /entry-task/
COPY ./front/.browserslistrc /entry-task/
COPY ./front/.editorconfig /entry-task/
COPY ./front/babel.config.js /entry-task/
COPY ./front/jest.config.js /entry-task/
RUN npm install -g npm@8.8.0
RUN cd /entry-task && npm install && npm run build

FROM nginx:1.21.6 as nginx
WORKDIR /var/
COPY --from=node /entry-task/dist /entry-task
COPY ./cmd/front-end/nginx.conf /etc/nginx/nginx.conf
