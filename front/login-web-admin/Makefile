test:
	echo empty

build:

	node -v
	npm install --no-fund
	echo "乾坤版本:" 
	npm view qiankun version
	echo "当前发布环境:" $(E)
	rm -rf ci-cd/build/dist/*
	npm run env-ci-cd $(E)
	npm run build
