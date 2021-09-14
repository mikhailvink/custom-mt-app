FROM golang:1.17-alpine as build

LABEL description="Image of mau/crowdin-custom-mt-app" \
      maintainer="Semyon Atamas <semyon.atamas@jetbrains.com>" \
      source="https://jetbrains.team/p/mau/repositories/crowdin-custom-mt-app/commits?tab=changes"

COPY src /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o /crowdin_grazie_mt ./main/main.go

FROM alpine:latest
RUN addgroup  -S app_group && adduser -S app_user -G app_group -u 1001
USER 1001

EXPOSE 8080
COPY static /static/
COPY --from=build /crowdin_grazie_mt /crowdin_grazie_mt
ENTRYPOINT ["/crowdin_grazie_mt"]
