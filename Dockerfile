FROM alpine:3.14

WORKDIR /service

RUN addgroup -S bot_operator && adduser -S bot_operator -G bot_operator && chown -R bot_operator:bot_operator /service

COPY --chown=bot_operator weather-or-not-bot /service/.

USER bot_operator

ENTRYPOINT ["/service/weather-or-not-bot"]