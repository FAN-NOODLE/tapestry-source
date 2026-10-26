const Mock = require('../mock.js')

const mock = baseUrl => {
  Mock.mock(`${baseUrl}member/detail`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'success',
      'request_id': '2bd70c934aa44882018e3c661696165c17537',
      'data': {
        'code': 200,
        'message': 'success',
        'request_id': '2bd70c934aa44882018e3c661696165c17537',
        'id': 80467573631680,
        'worker_status': 1,
        'name': '贺树鹏',
        'mobile': '15014002531',
        'worker_no': 'DY20200330775',
        'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330775.352.png',
        'avatar': 'https://wework.qpic.cn/bizmail/gY2LjWt7XrRTEcaoV06WOCyK8oDO3EEcEGfGibNxAy7SibiaocFicv04Wg/0',
        'mem_uid': 'C304C02C7F0000017A24589056B4D30B',
        'is_boss': 1,
        'is_send_sms': 1,
        'sms_mobile': 15014002535
      }
    }
  })
}

module.exports = {
  mock
}

