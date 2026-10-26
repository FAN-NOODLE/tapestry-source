const UNAVAILABLE_MESSAGE = '库存信息暂时无法获取，请稍后再试。\n给您带来不便，敬请谅解。'
const PARTIAL_MESSAGE = '部分库存信息可能不完整，请稍后刷新查看。'
const DCL_DELAY_MESSAGE = '大仓库存信息可能存在延迟，请稍后重试或人工核对库存。'

function normalizeInventoryPayload(payload) {
  const source = payload || {}
  const list = Array.isArray(source)
    ? source
    : (Array.isArray(source.list) ? source.list : [])
  const failedUPCs = Array.isArray(source.failed_upcs) ? source.failed_upcs : []
  const warnings = []
  const hasInventoryError = source.inventory_status === 'partial' ||
    failedUPCs.length > 0 || source.dcl_status === 'failed'

  if (list.length === 0 && hasInventoryError) {
    return {
      list: [],
      failedUPCs,
      warnings: [],
      emptyMessage: UNAVAILABLE_MESSAGE
    }
  }

  if (source.inventory_status === 'partial' || failedUPCs.length > 0) {
    warnings.push(PARTIAL_MESSAGE)
  }
  if (source.dcl_status === 'failed') {
    warnings.push(DCL_DELAY_MESSAGE)
  }

  return {
    list: source.dcl_status === 'failed'
      ? list.map(item => ({ ...item, dcl_on_hand: '--' }))
      : list,
    failedUPCs,
    warnings,
    emptyMessage: ''
  }
}

function normalizeInventoryFailure() {
  return {
    list: [],
    failedUPCs: [],
    warnings: [],
    emptyMessage: UNAVAILABLE_MESSAGE
  }
}

module.exports = {
  UNAVAILABLE_MESSAGE,
  PARTIAL_MESSAGE,
  DCL_DELAY_MESSAGE,
  normalizeInventoryFailure,
  normalizeInventoryPayload
}
