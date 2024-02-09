build-prod:
	yarn
	yarn nx affected -t build --base=origin/master --head=HEAD --parallel --max-parallel=3
build-prod-core:
	yarn
	yarn nx build noxite-core
build-prod-web-adapters:
	yarn
	yarn nx build noxite-web-adapters
build-prod-app:
	yarn
	yarn nx build noxite-webapp
test:
	yarn
	yarn nx affected -t test --base=origin/master --parallel --max-parallel=3