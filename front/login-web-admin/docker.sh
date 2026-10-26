VERSON=$1
npm run env login https://tapestry-dev.vchangyi.com/gw-scrm/ prod
npm run build
cd deploy
docker build --build-arg HASH=login -t  test-login.b-${VERSON} .
docker tag test-login.b-${VERSON} harbor.vchangyi.com/oulaiya-public/crs_admin-login-http:test-login.b-${VERSON}
docker login -u "prod-oulaiya" -p "<sY+{6^YLN3's>Ez" harbor.vchangyi.com
docker push harbor.vchangyi.com/oulaiya-public/crs_admin-login-http:test-login.b-${VERSON}