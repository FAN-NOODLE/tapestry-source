const assert = require('assert')
const {
  DCL_DELAY_MESSAGE,
  PARTIAL_MESSAGE,
  normalizeInventoryFailure,
  normalizeInventoryPayload
} = require('./inventory-status')

const partial = normalizeInventoryPayload({
  list: [{ upc_code: 'UPC-1', dcl_on_hand: 3 }],
  inventory_status: 'partial',
  failed_upcs: ['UPC-2'],
  dcl_status: 'success'
})
assert.deepStrictEqual(partial.failedUPCs, ['UPC-2'])
assert.deepStrictEqual(partial.warnings, [PARTIAL_MESSAGE])
assert.strictEqual(partial.list[0].dcl_on_hand, 3)

const dclFailed = normalizeInventoryPayload({
  list: [{ upc_code: 'UPC-1', dcl_on_hand: 0 }],
  inventory_status: 'success',
  failed_upcs: [],
  dcl_status: 'failed'
})
assert.deepStrictEqual(dclFailed.warnings, [DCL_DELAY_MESSAGE])
assert.strictEqual(dclFailed.list[0].dcl_on_hand, '--')

const both = normalizeInventoryPayload({
  list: [{ upc_code: 'UPC-1', dcl_on_hand: 9 }],
  inventory_status: 'partial',
  failed_upcs: ['UPC-2'],
  dcl_status: 'failed'
})
assert.deepStrictEqual(both.warnings, [PARTIAL_MESSAGE, DCL_DELAY_MESSAGE])
assert.strictEqual(both.list[0].dcl_on_hand, '--')

const partialWithoutData = normalizeInventoryPayload({
  list: [],
  inventory_status: 'partial',
  failed_upcs: ['UPC-FAIL'],
  dcl_status: 'success'
})
assert.deepStrictEqual(partialWithoutData.warnings, [])
assert.strictEqual(partialWithoutData.emptyMessage, '库存信息暂时无法获取，请稍后再试。\n给您带来不便，敬请谅解。')

const dclFailedWithoutData = normalizeInventoryPayload({
  list: [],
  inventory_status: 'success',
  failed_upcs: [],
  dcl_status: 'failed'
})
assert.deepStrictEqual(dclFailedWithoutData.warnings, [])
assert.strictEqual(dclFailedWithoutData.emptyMessage, '库存信息暂时无法获取，请稍后再试。\n给您带来不便，敬请谅解。')

const unavailable = normalizeInventoryFailure()
assert.deepStrictEqual(unavailable.list, [])
assert.deepStrictEqual(unavailable.failedUPCs, [])
assert.deepStrictEqual(unavailable.warnings, [])
assert.strictEqual(unavailable.emptyMessage, '库存信息暂时无法获取，请稍后再试。\n给您带来不便，敬请谅解。')

console.log('inventory-status tests passed')
