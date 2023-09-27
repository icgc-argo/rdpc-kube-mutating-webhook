#
# Copyright (c) 2020. Ontario Institute for Cancer Research
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/>.
#
ARG GO_VERSION=1.16
FROM golang:${GO_VERSION}-alpine AS build

COPY ./cmd /srv
COPY ./init-build.sh /srv
WORKDIR /srv

RUN apk add --no-cache git \
    && chmod +x ./init-build.sh \
	&& ./init-build.sh


env CGO_ENABLED=0 
env GOOS=linux 
RUN cd /srv/webhook-server \
        && go mod init \
	&& go build -ldflags="-s -w" -o webhook-server *.go


###############################################################

FROM golang:${GO_VERSION}-alpine

ENV APP_USER appuser
ENV APP_UID 9999
ENV APP_GID 9999
ENV APP_HOME /app

COPY --from=build /srv/webhook-server/webhook-server /tmp/webhook-server

RUN addgroup -S -g $APP_GID $APP_USER  \
	&& adduser -S -u $APP_UID -G $APP_USER $APP_USER  \
	&& mkdir -p $APP_HOME \
	&& mv /tmp/webhook-server $APP_HOME/webhook-server \
	&& chown -R $APP_UID:$APP_GID $APP_HOME

WORKDIR $APP_HOME

USER $APP_UID

CMD ["/app/webhook-server"]


