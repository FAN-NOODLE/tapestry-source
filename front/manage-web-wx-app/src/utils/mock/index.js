// mock入口文件
const baseUrl = 'https://crs-api-test.vchangyi.com/gw-assistant/app/v1/'

require('./member').mock(baseUrl)
require('./store').mock(baseUrl)
require('./withdrawal').mock(baseUrl)
