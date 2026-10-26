const Mock = require('../mock.js')

const mock = baseUrl => {
  Mock.mock(`${baseUrl}withdrawal/withdrawal-home`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'repr',
      'request_id': 'consectetur id',
      'data': {
        'info': {
          'allow_withdrawal': true,
          'id': 33153379.03423053,
          'enterprise_id': 46484644.45617947,
          'worker_id': 42877931.504592806,
          'real_name': '杨森',
          'mobile': 'laborum magna',
          'id_type': 795790.6670084596,
          'id_code': 'eu Excepteur sit adipisicing quis',
          'bank_card_username': 'p',
          'bank_card_code': 'consequat',
          'bank_name': 'officia sit',
          'agree_deal': -97281753.92541867,
          'auth_status': 13095781.950182945,
          'status': -13609138.226661384,
          'created_at': 'do enim minim',
          'updated_at': 'in',
          'deleted_at': 'qui commodo eiusmod eu aliquip',
          'available_cash_amount': 676784.093771413,
          'today_withdrawaled': -37285259
        },
        'allow_withdrawal': true,
        'transfer_no': 'Duis',
        'need_certification': false
      }
    }
  })
  Mock.mock(`${baseUrl}withdrawal/withdrawal-logs`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'amet',
      'request_id': 'laboris in quis consequat',
      'data': {
        data: [{
          'id': 97249784.2350815,
          'operator_id': 51971294.59295225,
          'worker_id': -99439501.57694922,
          'application_amount': 12345,
          'request_initiated_at': -40236390.87081539,
          'request_handle_at': 1894227.698005274,
          'status': 1,
          'reason': 'mollit amet voluptate adipisicing'
        }]
      }
    }
  })
  Mock.mock(`${baseUrl}distribution-goods/goods-list`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'success',
      'request_id': '916fb5c688374cc2f33cb2c2ee8956ed93531',
      'data': {
        'list': [
          {
            'spu_id': 116515474108224,
            'spu_img': 'https://img-crs.vchangyi.com/2020/04/22/6ad9e47b9c41e3e86633105ab7988c17.jpg',
            'spu_name': 'Dior\u53e3\u7ea2',
            'sales': '20',
            'commission': 980,
            'new_goods': 1,
            'sku_id': 116515474665280,
            'sign_price': 20000,
            'sell_price': 9800
          },
          {
            'spu_id': 116515474108224,
            'spu_img': 'https://img-crs.vchangyi.com/2020/04/22/6ad9e47b9c41e3e86633105ab7988c17.jpg',
            'spu_name': 'Dior\u53e3\u7ea2',
            'sales': '20',
            'commission': 980,
            'new_goods': 1,
            'sku_id': 116515474665280,
            'sign_price': 20000,
            'sell_price': 9800
          },
          {
            'spu_id': 116515474108224,
            'spu_img': 'https://img-crs.vchangyi.com/2020/04/22/6ad9e47b9c41e3e86633105ab7988c17.jpg',
            'spu_name': 'Dior\u53e3\u7ea2',
            'sales': '20',
            'commission': 980,
            'new_goods': 1,
            'sku_id': 116515474665280,
            'sign_price': 20000,
            'sell_price': 9800
          }
        ],
        'page': 1,
        'size': 10,
        'total': 1
      }
    }
  })
  Mock.mock(`${baseUrl}withdrawal/certification-audit`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'success',
      'request_id': '6a4004544e6fcc61c86912140064c90313317',
      'data': {'status': true, 'message': '认证成功'}
    }
  })
  Mock.mock(`${baseUrl}withdrawal/withdrawal-add`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'success',
      'request_id': '6a4004544e6fcc61c86912140064c90313317',
      'data': {'status': true, 'message': '认证成功'}
    }
  })
  Mock.mock(`${baseUrl}common/login`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'success',
      'request_id': 'fc15730b2e3424a98a5c333609e56ecd89356',
      'data': {
        'token': 'REDACTED_FOR_SOURCE_SHARING',
        'enterprise': [
          {
            'enterprise_id': 2,
            'hash': '054edce1d4cec91c76001d968841c08f',
            'version': 2,
            'name': '畅移测试企业-聚骧',
            'logo': '',
            'contact_name': '韩志刚',
            'contact_mobile': '18229062631',
            'address': '陕西省西安市高新区财富中心Ⅱ期D座',
            'crs_status': 1,
            'store_status': 1,
            'domain': '72D5430F7F000001381C7A561A35C79D',
            'corp_id': 'wwb1a6952b5c909031',
            'account_start_at': 0,
            'account_expire_at': 1687023409,
            'is_split_flow': 1,
            'admin_id': 33328,
            'account': '18700465991',
            'account_name': '黄小米',
            'account_global_roaming': '+86',
            'account_mobile': '18700465991',
            'account_email': '',
            'admin_available_system': '',
            'role': 0,
            'is_super': 2,
            'range_type': 1,
            'serve_type': 0,
            'shop_name': '聚骧',
            'enterprise_role_list': []
          },
          {
            'enterprise_id': 2,
            'hash': '054edce1d4cec91c76001d968841c08f',
            'version': 2,
            'name': '畅移测试企业-聚骧',
            'logo': '',
            'contact_name': '韩志刚',
            'contact_mobile': '18229062631',
            'address': '陕西省西安市高新区财富中心Ⅱ期D座',
            'crs_status': 1,
            'store_status': 1,
            'domain': '72D5430F7F000001381C7A561A35C79D',
            'corp_id': 'wwb1a6952b5c909031',
            'account_start_at': 0,
            'account_expire_at': 1687023409,
            'is_split_flow': 1,
            'admin_id': 33328,
            'account': '18700465991',
            'account_name': '黄小米',
            'account_global_roaming': '+86',
            'account_mobile': '18700465991',
            'account_email': '',
            'admin_available_system': '',
            'role': 0,
            'is_super': 2,
            'range_type': 1,
            'serve_type': 0,
            'shop_name': '聚骧',
            'enterprise_role_list': []
          }
        ],
        'is_first_login': 0,
        'login_cache_key': 'wx_login_18700465991'
      }
    }
  })
  Mock.mock(`${baseUrl}distribution-settlement/list`, {
    statusCode: 200,
    data: {
      'code': 200,
      'message': 'veniam commodo aute ad ut',
      'request_id': 'fugiat aliqua dolore minim enim',
      'data': {
        'list': [
          {
            'id': -73063265.05617839,
            'settlement_status': 86246258.84179342,
            'order_no': 'laborum non Lorem',
            'order_status': 16254574.382783994,
            'order_pay_time': -14436592.137329727,
            'goods_info': [
              {
                'order_goods_price': 874826.9491889477,
                'promotion_commission': 45066094.60712066,
                'service_commission': 41172426.660030484,
                'purchase_quantity': -41173758.899779946,
                'material_no': 'qui est inci',
                'spu_img': 'incididunt ullamco Excepteur dolor',
                'spu_name': 'incididunt eiusmod',
                'attribute': [
                  'Duis'
                ],
                'commission': 60398218.805211514
              },
              {
                'order_goods_price': 39743090.91590005,
                'promotion_commission': -87913767.82231736,
                'service_commission': -93711225.39425083,
                'purchase_quantity': 16131972.998414978,
                'material_no': 'consectetur magna id',
                'spu_img': 'eiusmod non aute',
                'spu_name': 'commodo',
                'attribute': [
                  'c',
                  'non esse tempor'
                ],
                'commission': -84365447.92327023
              },
              {
                'order_goods_price': 58480080.599983364,
                'promotion_commission': -79598441.68959264,
                'service_commission': -61432262.48703933,
                'purchase_quantity': 27542635.290418312,
                'material_no': 'magna tempor velit ',
                'spu_img': 'laboris pariatur',
                'spu_name': 'laboris sit fugiat consequat est',
                'attribute': [
                  'deserunt occaecat'
                ],
                'commission': 59282570.31912252
              },
              {
                'order_goods_price': -99380132.4081939,
                'promotion_commission': 11151094.572259009,
                'service_commission': 29089933.210079134,
                'purchase_quantity': -60717644.80129298,
                'material_no': 'ut in esse',
                'spu_img': 'ut',
                'spu_name': 'nulla irure',
                'attribute': [
                  'ad exercitation',
                  'occae',
                  'culpa amet exercitation magna'
                ],
                'commission': 44091884.27665791
              }
            ],
            'settlement_time': 'l',
            'settlement_money': -52889435.71778515
          },
          {
            'id': 82398148.17929855,
            'settlement_status': 29456335.08530055,
            'order_no': 'dolor consectetur sunt tempor',
            'order_status': -41254244.1521246,
            'order_pay_time': -16640296.89989151,
            'goods_info': [
              {
                'order_goods_price': 40499060.883199126,
                'promotion_commission': -3372724.6766011864,
                'service_commission': -87730947.96384013,
                'purchase_quantity': -66702823.127481595,
                'material_no': 'incididunt aliqua cillum labore fugiat',
                'spu_img': 'eiusmod nulla fugiat',
                'spu_name': 'deserunt nisi',
                'attribute': [
                  'consectetur dolore',
                  'in aliquip quis ipsum in'
                ],
                'commission': 81676147.91045636
              },
              {
                'order_goods_price': -40062570.59693214,
                'promotion_commission': -91132999.13371086,
                'service_commission': -1238625.3057869077,
                'purchase_quantity': -93540416.03202412,
                'material_no': 'aliquip fugiat',
                'spu_img': 'fugiat in ad laboris',
                'spu_name': 'id ',
                'attribute': [
                  'mollit',
                  'occaeca'
                ],
                'commission': -62101750.61273486
              },
              {
                'order_goods_price': -15911102.079266697,
                'promotion_commission': -3006817.758140117,
                'service_commission': -49884731.77786355,
                'purchase_quantity': 27205685.651297003,
                'material_no': 'minim amet velit eu sed',
                'spu_img': 'sed dolore nisi',
                'spu_name': 'Ut tempor ex',
                'attribute': [
                  'amet magna',
                  'ut',
                  'fugiat',
                  'sint deserunt laborum elit'
                ],
                'commission': 52879785.6023151
              }
            ],
            'settlement_time': 'minim adipisicing velit in cillum',
            'settlement_money': -47287288.30500692
          }
        ],
        'total': 20,
        'page': 1,
        'size': 5
      }
    }
  })
}

module.exports = {
  mock
}
