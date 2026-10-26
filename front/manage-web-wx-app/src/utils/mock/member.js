const Mock = require('../mock.js')

console.log(Mock)
const mock = baseUrl => {
//   Mock.mock(`${baseUrl}member/detail`, {
//     statusCode: 200,
//     data: {
//       code: 200,
//       message: 'success',
//       request_id: '050444fbb4ecc0c07a3ee84d3efdc4e412363',
//       data: {
//         'code': 200,
//         'message': 'success',
//         'request_id': 'e098b79deab07a348c7775ca2572946a59317',
//         'data': {
//           'id': 80467573631680,
//           'worker_status': 1,
//           'name': '贺树鹏',
//           is_boss: 1,
//           'mobile': '15014002531',
//           'worker_no': 'DY20200330775',
//           'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330775.352.png',
//           'avatar': 'http://wework.qpic.cn/bizmail/gY2LjWt7XrRTEcaoV06WOCyK8oDO3EEcEGfGibNxAy7SibiaocFicv04Wg/0',
//           'mem_uid': 'C304C02C7F0000017A24589056B4D30B'
//         }
//       }
//     }
//   })
//
//   Mock.mock(`${baseUrl}member/list`, {
//     statusCode: 200,
//     data: {
//       code: 200,
//       message: 'success',
//       request_id: '050444fbb4ecc0c07a3ee84d3efdc4e412363',
//       data: {
//         'code': 200,
//         'message': 'success',
//         'request_id': 'e098b79deab07a348c7775ca2572946a59317',
//         'list': [
//           {
//             'id': 80467573631680,
//             'worker_no': 'DY20200330775',
//             'name': '贺树鹏',
//             'mobile': '15014002531',
//             'avatar': 'http://wework.qpic.cn/bizmail/gY2LjWt7XrRTEcaoV06WOCyK8oDO3EEcEGfGibNxAy7SibiaocFicv04Wg/0',
//             'is_boss': 2,
//             'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330775.352.png',
//             'source': '企微店员',
//             'worker_status': 1,
//             'created_at': 1585562807,
//             'nickname': '贺树鹏'
//           },
//           {
//             'id': 80467442362624,
//             'worker_no': 'DY20200330474',
//             'name': '郭天航',
//             'mobile': '18966532738',
//             'avatar': 'https://wework.qpic.cn/wwhead/duc2TvpEgSQO4BpE0WZSZ3kUp2eCRuPj5EIL52F7mMriaHQ4oUHs2rTCBmrXE9ibdzVlwjFkOksrs/0',
//             'is_boss': 2,
//             'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330474.352.png',
//             'source': '企微店员',
//             'worker_status': 1,
//             'created_at': 1585562743,
//             'nickname': '郭天航'
//           },
//           {
//             'id': 80466827254016,
//             'worker_no': 'DY20200330274',
//             'name': '齐文豪',
//             'mobile': '13279488135',
//             'avatar': 'https://wework.qpic.cn/bizmail/6l9AkZNfEu1Yr75whOhdArt2I7t7RLWVuUbd3EO1iaOdUEOYMRwdI2Q/0',
//             'is_boss': 2,
//             'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330274.178.png',
//             'source': '企微店员',
//             'worker_status': 1,
//             'created_at': 1585562442,
//             'nickname': '齐文豪'
//           },
//           {
//             'id': 80466727821568,
//             'worker_no': 'DY20200330783',
//             'name': '李云翔',
//             'mobile': '18082232357',
//             'avatar': 'http://wework.qpic.cn/bizmail/lLxfrD9dvz65r3FwVvYlWia190kWdreN6KE8TWQk5C1ib6jXsV6ibwMiaA/0',
//             'is_boss': 2,
//             'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330783.678.png',
//             'source': '企微店员',
//             'worker_status': 1,
//             'created_at': 1585562394,
//             'nickname': '李云翔'
//           },
//           {
//             'id': 80466719892160,
//             'worker_no': 'DY20200330996',
//             'name': '侯英才',
//             'mobile': '18629658853',
//             'avatar': 'http://wework.qpic.cn/bizmail/OnjVbBaJ0kZboJOQTJ7sErca5KWzuFgzyjF5l6pyUKHjYT36ia1MicwQ/0',
//             'is_boss': 2,
//             'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330996.633.png',
//             'source': '企微店员',
//             'worker_status': 1,
//             'created_at': 1585562390,
//             'nickname': '侯英才'
//           },
//           {
//             'id': 80463356791488,
//             'worker_no': 'DY20200330935',
//             'name': '胡汉三',
//             'mobile': '18161732940',
//             'avatar': 'http://wework.qpic.cn/bizmail/dmcV8xBkbbyjhsicMvfOZorfu5qbtYS34tWuMchFCTX0hsp52FrP00g/0',
//             'is_boss': 2,
//             'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330935.235.png',
//             'source': '企微店员',
//             'worker_status': 1,
//             'created_at': 1585560748,
//             'nickname': '胡汉三'
//           },
//           {
//             'id': 80463303789248,
//             'worker_no': 'DY20200330588',
//             'name': '李林欢',
//             'mobile': '18092724105',
//             'avatar': 'http://wework.qpic.cn/bizmail/3n8icoXGgbyOnzM23KGeDSydZaeD6QLEiaiabqYFGJSicTmcnoXISOJUXA/0',
//             'is_boss': 1,
//             'qr_code': 'https://img-crs-test.vchangyi.com/download/2020/03/30/DY20200330588.839.png',
//             'source': '企微店员',
//             'worker_status': 1,
//             'created_at': 1585560722,
//             'nickname': '李林欢'
//           }
//         ],
//         'total': 7,
//         'page': 1,
//         'size': 10
//       }
//     }
//   })
//
//   Mock.mock(`${baseUrl}member/statistics`, {
//     statusCode: 200,
//     data: {
//       code: 200,
//       message: 'success',
//       request_id: '050444fbb4ecc0c07a3ee84d3efdc4e412363',
//       data: {
//         'code': 200,
//         'message': 'success',
//         'request_id': 'e098b79deab07a348c7775ca2572946a59317',
//         real_money_total: 8000000000,
//         complete_order_num: 8,
//         complete_money_total: 80000,
//         in_sales_order_num: 8,
//         in_sales_money_total: 80000,
//         sale_all_num: 8,
//         sale_all_money_total: 80000
//       }
//     }
//   })
}
//
module.exports = {
  mock
}
