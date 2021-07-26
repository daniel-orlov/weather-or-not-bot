FROM alpine:3.14

WORKDIR /service

COPY --chown=bot_operator weather-or-not-bot /service/.

USER bot_operator

ENTRYPOINT ["/service/weather-or-not-bot"]