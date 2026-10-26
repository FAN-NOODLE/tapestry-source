VERSON=$1
npm run env work https://tapestry-dev.vchangyi.com/gw-work-scrm/
npm run build
cd deploy
docker build --build-arg HASH=work -t  test-work.b-${VERSON} .
docker tag test-work.b-${VERSON} harbor.vchangyi.com/oulaiya-public/crs_scrm-work-http:test-work.b-${VERSON}
docker login -u "prod-oulaiya" -p "<sY+{6^YLN3's>Ez" harbor.vchangyi.com
docker push harbor.vchangyi.com/oulaiya-public/crs_scrm-work-http:test-work.b-${VERSON}