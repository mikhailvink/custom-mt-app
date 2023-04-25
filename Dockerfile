FROM ubuntu:20.04

LABEL description="Image of mau/crowdin-custom-mt-app" \
      maintainer="Andrey Vasin <andrey.vasin@jetbrains.com>" \
      source="https://jetbrains.team/p/mau/repositories/crowdin-custom-mt-app"

EXPOSE 8080

RUN groupadd -r crowdin_grazie_mt && useradd --uid 1000 --no-log-init -r -g crowdin_grazie_mt crowdin_grazie_mt
RUN apt-get update && \
    apt-get install -y ca-certificates && \
    rm -rfv /var/cache/apt

COPY static /static/
COPY crowdin_grazie_mt /crowdin_grazie_mt
RUN chmod +x /crowdin_grazie_mt

USER 1000

CMD ["/crowdin_grazie_mt"]
