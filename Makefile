.PHONY: build
build: test
	( cd src && CGO_ENABLED=0 GOOS=linux go build -mod=vendor -o ../crowdin_grazie_mt ./main/main.go )

.PHONY: test
test:
	( cd src && go test ./... )

.PHONY: docker-image
docker-image: image_tag_required
	docker-compose build crowdin_grazie_mt
	docker-compose push crowdin_grazie_mt

.PHONY: deploy-prod
deploy-prod: image_tag_required k8s_token_required
	( cd deploy && k8s-handle deploy -s production --sync-mode --strict)

.PHONY: destroy-prod
destroy-prod: image_tag_required k8s_token_required
	( cd deploy && k8s-handle destroy -s production --sync-mode --strict)

.PHONY: all-direct
all-prod: build docker-image deploy-prod

# protect envvar
env-%:
	@ if [ "${${*}}" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

image_tag_required: env-IMAGE_TAG
k8s_token_required: env-K8S_TOKEN
