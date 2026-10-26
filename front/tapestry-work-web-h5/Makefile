test:
	echo empty

build:

	node -v
	npm install --no-fund
	echo "当前发布环境:" $(E) $(V)
	rm -rf ci-cd/build/dist/*
	npm run env-ci-cd $(E) $(V)
	npm run build
