// ----------------------------------------
// 售后状态
// ----------------------------------------
// 待后台审核
const _WAIT_CHECK = 21
// 后台通过，待OPEN审核
const _WAIT_CHECK_A = 22
// OPEN 审核通过，待用户寄回商品
const _WAIT_CHECK_B = 23
// 用户已寄出商品
const _WAIT_CHECK_C = 24
// 人工中
const _WAIT_REFUND = 25
// 同意退款
const _REFUND_SUCCESS = 26

export {
  _WAIT_CHECK,
  _WAIT_CHECK_A,
  _WAIT_CHECK_B,
  _WAIT_CHECK_C,
  _WAIT_REFUND,
  _REFUND_SUCCESS
}
