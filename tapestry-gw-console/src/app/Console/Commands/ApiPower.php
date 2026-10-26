<?php
/*
 * Created by PhpStorm.
 * User: wangce
 * Date Time: 2019/2/25 11:49
 */

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class ApiPower extends Command
{

    protected $signature = 'initialize:api_power';

    protected $description = '初始化API数据';

    private $insertData = [];

    public function handle()
    {
        $data = [
            [
                'id' => 1,
                'pid' => 0,
                'page' => '系统',
                'route' => 'admin',
                'describe' => '管理员',
                'child' => [
                    [
                        'id' => 2,
                        'pid' => 1,
                        'page' => '登陆页',
                        'route' => 'admin/login',
                        'describe' => '登陆',
                        'child' => []
                    ],
                    [
                        'id' => 3,
                        'pid' => 1,
                        'page' => '',
                        'route' => 'admin/logout',
                        'describe' => '退出',
                        'child' => []
                    ],
                    [
                        'id' => 4,
                        'pid' => 1,
                        'page' => '',
                        'route' => 'admin/mock-login',
                        'describe' => '模拟登陆',
                        'child' => []
                    ],
                    [
                        'id' => 5,
                        'pid' => 1,
                        'page' => '',
                        'route' => 'admin/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 6,
                        'pid' => 1,
                        'page' => '',
                        'route' => 'admin/password',
                        'describe' => '修改密码',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 7,
                'pid' => 0,
                'page' => '设置.小程序授权',
                'route' => 'auth',
                'describe' => '小程序授权',
                'child' => [
                    [
                        'id' => 8,
                        'pid' => 7,
                        'page' => '商城小程序授权',
                        'route' => 'auth/shop-auth-detail',
                        'describe' => '商城小程序授权详情',
                        'child' => []
                    ],
                    [
                        'id' => 9,
                        'pid' => 7,
                        'page' => '企业微信授权',
                        'route' => 'auth/work-auth-detail',
                        'describe' => '企业微信小程序授权详情',
                        'child' => []
                    ],
                    [
                        'id' => 10,
                        'pid' => 7,
                        'page' => '客户联系配置',
                        'route' => 'auth/auth-contact',
                        'describe' => '通讯录授权详情',
                        'child' => []
                    ],
                    [
                        'id' => 11,
                        'pid' => 7,
                        'page' => '完善资料',
                        'route' => 'auth/work-auth-complete',
                        'describe' => '企业微信小程序完善信息',
                        'child' => []
                    ],
                    [
                        'id' => 12,
                        'pid' => 7,
                        'page' => '完善资料',
                        'route' => 'auth/shop-auth-complete',
                        'describe' => '商城小程序完善信息',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 13,
                'pid' => 0,
                'page' => '品类',
                'route' => 'category',
                'describe' => '商品分类',
                'child' => [
                    [
                        'id' => 14,
                        'pid' => 13,
                        'page' => '品类.列表',
                        'route' => 'category/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 15,
                        'pid' => 13,
                        'page' => '品类.列表',
                        'route' => 'category/recursion-list',
                        'describe' => '递归列表',
                        'child' => []
                    ],
                    [
                        'id' => 16,
                        'pid' => 13,
                        'page' => '品类.新增分类',
                        'route' => 'category/add',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 17,
                        'pid' => 13,
                        'page' => '品类.删除分类',
                        'route' => 'category/delete',
                        'describe' => '删除',
                        'child' => []
                    ],
                    [
                        'id' => 18,
                        'pid' => 13,
                        'page' => '品类.编辑分类',
                        'route' => 'category/update',
                        'describe' => '编辑',
                        'child' => []
                    ],
                    [
                        'id' => 19,
                        'pid' => 13,
                        'page' => '品类.分类详情',
                        'route' => 'category/{id}',
                        'describe' => '详情',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 20,
                'pid' => 0,
                'page' => '订单',
                'route' => 'order',
                'describe' => '订单',
                'child' => [
                    [
                        'id' => 21,
                        'pid' => 20,
                        'page' => '订单.列表',
                        'route' => 'order/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 22,
                        'pid' => 20,
                        'page' => '订单.订单详情',
                        'route' => 'order/{id}',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 23,
                        'pid' => 20,
                        'page' => '订单.拼团用户详情',
                        'route' => 'order/user-detail',
                        'describe' => '拼团用户详情',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 24,
                'pid' => 0,
                'page' => '店铺.页面',
                'route' => 'page',
                'describe' => '页面',
                'child' => [
                    [
                        'id' => 25,
                        'pid' => 24,
                        'page' => '店铺.页面.轮播图',
                        'route' => 'page/banner',
                        'describe' => '轮播图',
                        'child' => [
                            [
                                'id' => 26,
                                'pid' => 25,
                                'page' => '店铺.页面.轮播图列表',
                                'route' => 'page/banner/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 27,
                                'pid' => 25,
                                'page' => '店铺.页面.轮播图保存',
                                'route' => 'page/banner/save',
                                'describe' => '保存'
                            ]
                        ]
                    ],
                    [
                        'id' => 28,
                        'pid' => 24,
                        'page' => '店铺.页面.首页引导区',
                        'route' => 'page/guide',
                        'describe' => '首页引导区',
                        'child' => [
                            [
                                'id' => 29,
                                'pid' => 28,
                                'page' => '店铺.页面.首页引导区列表',
                                'route' => 'page/guide/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 30,
                                'pid' => 28,
                                'page' => '店铺.页面.首页引导区保存',
                                'route' => 'page/guide/save',
                                'describe' => '保存'
                            ]
                        ]
                    ],
                    [
                        'id' => 31,
                        'pid' => 24,
                        'page' => '店铺.页面.爆品推荐',
                        'route' => 'page/hot',
                        'describe' => '爆品推荐',
                        'child' => [
                            [
                                'id' => 32,
                                'pid' => 31,
                                'page' => '店铺.页面.爆品推荐列表',
                                'route' => 'page/hot/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 33,
                                'pid' => 31,
                                'page' => '店铺.页面.爆品推荐保存',
                                'route' => 'page/hot/save',
                                'describe' => '保存'
                            ],
                            [
                                'id' => 34,
                                'pid' => 31,
                                'page' => '店铺.页面.爆品推荐选择商品',
                                'route' => 'page/hot/goods',
                                'describe' => '商品列表'
                            ]
                        ]
                    ]
                ]
            ],
            [
                'id' => 35,
                'pid' => 0,
                'page' => '应用.我的应用.礼券',
                'route' => 'coupon',
                'describe' => '礼券',
                'child' => [
                    [
                        'id' => 36,
                        'pid' => 35,
                        'page' => '应用.我的应用.礼券.礼券列表',
                        'route' => 'coupon/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 37,
                        'pid' => 35,
                        'page' => '应用.我的应用.礼券.礼券列表',
                        'route' => 'coupon/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 38,
                        'pid' => 35,
                        'page' => '应用.我的应用.礼券.新增礼券',
                        'route' => 'coupon/add',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 39,
                        'pid' => 35,
                        'page' => '应用.我的应用.礼券.编辑礼券',
                        'route' => 'coupon/save',
                        'describe' => '编辑',
                        'child' => []
                    ],
                    [
                        'id' => 40,
                        'pid' => 35,
                        'page' => '应用.我的应用.礼券.礼券列表',
                        'route' => 'coupon/delete',
                        'describe' => '删除',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 41,
                'pid' => 0,
                'page' => '应用.任务.任务中心',
                'route' => 'task',
                'describe' => '任务',
                'child' => [
                    [
                        'id' => 42,
                        'pid' => 41,
                        'page' => '应用.任务.任务中心.任务列表',
                        'route' => 'task/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 43,
                        'pid' => 41,
                        'page' => '应用.任务.任务中心.任务详情',
                        'route' => 'task/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 44,
                        'pid' => 41,
                        'page' => '应用.任务.任务中心.新增任务',
                        'route' => 'task/add',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 45,
                        'pid' => 41,
                        'page' => '应用.任务.任务中心.新增任务',
                        'route' => 'task/member-list',
                        'describe' => '员工列表',
                        'child' => []
                    ],
                    [
                        'id' => 46,
                        'pid' => 41,
                        'page' => '应用.任务.任务中心.任务列表',
                        'route' => 'task/delete',
                        'describe' => '删除',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 47,
                'pid' => 0,
                'page' => '组织.佣金',
                'route' => 'commission',
                'describe' => '佣金',
                'child' => [
                    [
                        'id' => 48,
                        'pid' => 46,
                        'page' => '组织.佣金规则',
                        'route' => 'commission/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 49,
                        'pid' => 46,
                        'page' => '组织.佣金规则',
                        'route' => 'commission/save',
                        'describe' => '保存',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 50,
                'pid' => 0,
                'page' => '客户.会员等级.会员等级',
                'route' => 'user-level',
                'describe' => '会员等级',
                'child' => [
                    [
                        'id' => 51,
                        'pid' => 50,
                        'page' => '客户.会员等级.会员等级',
                        'route' => 'user-level/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 52,
                        'pid' => 50,
                        'page' => '客户.会员等级.编辑会员等级',
                        'route' => 'user-level/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 53,
                        'pid' => 50,
                        'page' => '客户.会员等级.编辑会员等级',
                        'route' => 'user-level/update',
                        'describe' => '编辑',
                        'child' => []
                    ],
                    [
                        'id' => 54,
                        'pid' => 50,
                        'page' => '客户.会员等级.会员等级',
                        'route' => 'user-level/set-status',
                        'describe' => '状态变更',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 55,
                'pid' => 0,
                'page' => '客户.组织',
                'route' => 'department',
                'describe' => '部门',
                'child' => [
                    [
                        'id' => 56,
                        'pid' => 55,
                        'page' => '客户.组织',
                        'route' => 'department/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 57,
                        'pid' => 55,
                        'page' => '客户.组织.新增部门',
                        'route' => 'department/add',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 58,
                        'pid' => 55,
                        'page' => '客户.组织.编辑部门',
                        'route' => 'department/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 59,
                        'pid' => 55,
                        'page' => '客户.组织.编辑部门',
                        'route' => 'department/update',
                        'describe' => '编辑',
                        'child' => []
                    ],
                    [
                        'id' => 60,
                        'pid' => 55,
                        'page' => '客户.组织',
                        'route' => 'department/delete',
                        'describe' => '删除',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 61,
                'pid' => 0,
                'page' => '客户.客户列表',
                'route' => 'customer',
                'describe' => '客户',
                'child' => [
                    [
                        'id' => 62,
                        'pid' => 61,
                        'page' => '',
                        'route' => 'customer/list',
                        'describe' => '列表(暂且废弃)',
                        'child' => []
                    ],
                    [
                        'id' => 63,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/member-list',
                        'describe' => '会员客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 64,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/order-member',
                        'describe' => '下单客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 65,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/all-list',
                        'describe' => '全部客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 66,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/friend-list',
                        'describe' => '导购好友客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 67,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/all-list-export',
                        'describe' => '导出全部客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 68,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/friend-list',
                        'describe' => '导购好友客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 69,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/friend-list-export',
                        'describe' => '导出导购好友客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 70,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/order-member-export',
                        'describe' => '导出下单客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 71,
                        'pid' => 61,
                        'page' => '客户.客户列表',
                        'route' => 'customer/member-list-export',
                        'describe' => '导出会员客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 72,
                        'pid' => 61,
                        'page' => '客户.客户详情',
                        'route' => 'customer/detail',
                        'describe' => '基础信息',
                        'child' => []
                    ],
                    [
                        'id' => 73,
                        'pid' => 61,
                        'page' => '客户.客户详情',
                        'route' => 'customer/statistics',
                        'describe' => '统计数据',
                        'child' => []
                    ],
                    [
                        'id' => 74,
                        'pid' => 61,
                        'page' => '客户.客户详情',
                        'route' => 'customer/order-list',
                        'describe' => '订单列表',
                        'child' => []
                    ],
                    [
                        'id' => 75,
                        'pid' => 61,
                        'page' => '客户.订单详情',
                        'route' => 'customer/order-detail',
                        'describe' => '订单详情',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 76,
                'pid' => 0,
                'page' => '组织',
                'route' => 'user',
                'describe' => '员工',
                'child' => [
                    [
                        'id' => 77,
                        'pid' => 76,
                        'page' => '组织.在职员工列表.在职员工列表',
                        'route' => 'user/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 78,
                        'pid' => 76,
                        'page' => '组织.在职员工列表.添加员工',
                        'route' => 'user/add',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 79,
                        'pid' => 76,
                        'page' => '组织.在职员工列表.编辑员工',
                        'route' => 'user/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 80,
                        'pid' => 76,
                        'page' => '组织.在职员工列表.编辑员工',
                        'route' => 'user/update',
                        'describe' => '更新',
                        'child' => []
                    ],
                    [
                        'id' => 81,
                        'pid' => 76,
                        'page' => '组织.在职员工列表.在职员工列表',
                        'route' => 'user/leave',
                        'describe' => '离职',
                        'child' => []
                    ],
                    [
                        'id' => 82,
                        'pid' => 76,
                        'page' => '组织.离职员工列表.已交接客户列表',
                        'route' => 'user/handover-list',
                        'describe' => '交接客户列表',
                        'child' => []
                    ],
                    [
                        'id' => 83,
                        'pid' => 76,
                        'page' => '组织.离职员工列表.已交接客户列表',
                        'route' => 'user/handover',
                        'describe' => '交接客户',
                        'child' => []
                    ],
                    [
                        'id' => 84,
                        'pid' => 77,
                        'page' => '组织',
                        'route' => 'user/department-user-all-list',
                        'describe' => '员工部门分层列表',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 85,
                'pid' => 0,
                'page' => '设置.账号权限.角色',
                'route' => 'role',
                'describe' => '员工角色',
                'child' => [
                    [
                        'id' => 86,
                        'pid' => 85,
                        'page' => '设置.账号权限.角色.角色',
                        'route' => 'role/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 87,
                        'pid' => 85,
                        'page' => '设置.账号权限.角色.新增角色',
                        'route' => 'role/add',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 88,
                        'pid' => 85,
                        'page' => '设置.账号权限.角色.编辑角色',
                        'route' => 'role/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 89,
                        'pid' => 85,
                        'page' => '设置.账号权限.角色.编辑角色',
                        'route' => 'role/edit',
                        'describe' => '编辑',
                        'child' => []
                    ],
                    [
                        'id' => 90,
                        'pid' => 85,
                        'page' => '设置.账号权限.角色',
                        'route' => 'role/delete',
                        'describe' => '删除',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 91,
                'pid' => 0,
                'page' => '商品.规格',
                'route' => 'attributes',
                'describe' => '属性',
                'child' => [
                    [
                        'id' => 92,
                        'pid' => 91,
                        'page' => '商品.规格',
                        'route' => 'attributes/list',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 93,
                        'pid' => 91,
                        'page' => '商品.规格',
                        'route' => 'attributes/save',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 94,
                        'pid' => 91,
                        'page' => '商品.规格',
                        'route' => 'attributes/del',
                        'describe' => '删除',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 95,
                'pid' => 0,
                'page' => '首页',
                'route' => 'home',
                'describe' => '首页',
                'child' => [
                    [
                        'id' => 96,
                        'pid' => 95,
                        'page' => '首页',
                        'route' => 'home/chart',
                        'describe' => '下半部分图表',
                        'child' => []
                    ],
                    [
                        'id' => 97,
                        'pid' => 95,
                        'page' => '首页',
                        'route' => 'home/task',
                        'describe' => '上半部分任务达标率',
                        'child' => []
                    ],
                    [
                        'id' => 98,
                        'pid' => 95,
                        'page' => '首页',
                        'route' => 'home/statistical',
                        'describe' => '客户概况统计',
                        'child' => []
                    ],
                    [
                        'id' => 99,
                        'pid' => 95,
                        'page' => '首页',
                        'route' => 'home/order',
                        'describe' => '业绩概括统计',
                        'child' => []
                    ],
                    [
                        'id' => 100,
                        'pid' => 95,
                        'page' => '首页',
                        'route' => 'home/contact-auth-status',
                        'describe' => '获取企业是否授权通讯录',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 101,
                'pid' => 0,
                'page' => '商品',
                'route' => 'goods',
                'describe' => '商品',
                'child' => [
                    [
                        'id' => 102,
                        'pid' => 101,
                        'page' => '商品.商品列表',
                        'route' => 'goods',
                        'describe' => '列表',
                        'child' => []
                    ],
                    [
                        'id' => 103,
                        'pid' => 101,
                        'page' => '商品.商品详情',
                        'route' => 'goods/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 104,
                        'pid' => 101,
                        'page' => '商品.商品新增',
                        'route' => 'goods/save',
                        'describe' => '新增',
                        'child' => []
                    ],
                    [
                        'id' => 105,
                        'pid' => 101,
                        'page' => '商品.商品编辑',
                        'route' => 'goods/update',
                        'describe' => '编辑',
                        'child' => []
                    ],
                    [
                        'id' => 106,
                        'pid' => 101,
                        'page' => '商品.商品列表',
                        'route' => 'goods/sale',
                        'describe' => '上下架',
                        'child' => []
                    ],
                    [
                        'id' => 107,
                        'pid' => 101,
                        'page' => '商品.商品列表',
                        'route' => 'goods/delete',
                        'describe' => '删除',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 108,
                'pid' => 0,
                'page' => '设置.支付配置',
                'route' => 'payment',
                'describe' => '支付配置',
                'child' => [
                    [
                        'id' => 109,
                        'pid' => 108,
                        'page' => '设置.支付配置',
                        'route' => 'payment/detail',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 110,
                        'pid' => 108,
                        'page' => '设置.支付配置',
                        'route' => 'payment/save',
                        'describe' => '编辑',
                        'child' => []
                    ]
                ]
            ],
            [
                'id' => 111,
                'pid' => 0,
                'page' => '应用',
                'route' => 'app-center',
                'describe' => '应用',
                'child' => [
                    [
                        'id' => 112,
                        'pid' => 111,
                        'page' => '应用.全部应用',
                        'route' => 'app-center/list',
                        'describe' => '详情',
                        'child' => []
                    ],
                    [
                        'id' => 113,
                        'pid' => 111,
                        'page' => '应用.我的应用.新人优惠',
                        'route' => 'app-center/novice',
                        'describe' => '新人优惠',
                        'child' => [
                            [
                                'id' => 114,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠',
                                'route' => 'app-center/novice/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 115,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠',
                                'route' => 'app-center/novice/all',
                                'describe' => '优惠所有商品SPUID'
                            ],
                            [
                                'id' => 116,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠.新增优惠商品',
                                'route' => 'app-center/novice/add',
                                'describe' => '新增'
                            ],
                            [
                                'id' => 117,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠.编辑优惠商品',
                                'route' => 'app-center/novice/detail',
                                'describe' => '详情'
                            ],
                            [
                                'id' => 118,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠.编辑优惠商品',
                                'route' => 'app-center/novice/save',
                                'describe' => '编辑'
                            ],
                            [
                                'id' => 119,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠',
                                'route' => 'app-center/novice/status',
                                'describe' => '状态修改'
                            ],
                            [
                                'id' => 120,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠',
                                'route' => 'app-center/novice/delete',
                                'describe' => '删除'
                            ],
                            [
                                'id' => 121,
                                'pid' => 112,
                                'page' => '应用.我的应用.新人优惠',
                                'route' => 'app-center/novice/goods',
                                'describe' => '商品列表'
                            ]
                        ]
                    ],
                    [
                        'id' => 122,
                        'pid' => 111,
                        'page' => '应用.我的应用.拼团',
                        'route' => 'app-center/group',
                        'describe' => '拼团',
                        'child' => [
                            [
                                'id' => 123,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团',
                                'route' => 'app-center/group/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 124,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团.新增拼团商品',
                                'route' => 'app-center/group/add',
                                'describe' => '新增'
                            ],
                            [
                                'id' => 125,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团.查看拼团商品',
                                'route' => 'app-center/group/update',
                                'describe' => '编辑'
                            ],
                            [
                                'id' => 126,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团',
                                'route' => 'app-center/group/delete',
                                'describe' => '删除'
                            ],
                            [
                                'id' => 127,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团.查看拼团商品',
                                'route' => 'app-center/group/detail',
                                'describe' => '详情'
                            ],
                            [
                                'id' => 128,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团.查看拼团商品',
                                'route' => 'app-center/group/group-goods',
                                'describe' => '团购商品'
                            ],
                            [
                                'id' => 129,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团.查看拼团商品',
                                'route' => 'app-center/group/group-setting',
                                'describe' => '商品参加活动修改'
                            ],
                            [
                                'id' => 130,
                                'pid' => 122,
                                'page' => '应用.我的应用.拼团.查看拼团商品',
                                'route' => 'app-center/group/goods',
                                'describe' => '商品选择列表'
                            ]
                        ]
                    ],
                    [
                        'id' => 131,
                        'pid' => 111,
                        'page' => '应用.我的应用.限时购',
                        'route' => 'app-center/seckill',
                        'describe' => '秒杀',
                        'child' => [
                            [
                                'id' => 132,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购',
                                'route' => 'app-center/seckill/',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 133,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.新增活动场次',
                                'route' => 'app-center/seckill/save',
                                'describe' => '新增'
                            ],
                            [
                                'id' => 134,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.编辑活动场次',
                                'route' => 'app-center/seckill/update',
                                'describe' => '编辑'
                            ],
                            [
                                'id' => 135,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.编辑活动场次',
                                'route' => 'app-center/seckill/detail',
                                'describe' => '详情'
                            ],
                            [
                                'id' => 136,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.编辑活动场次',
                                'route' => 'app-center/seckill/goods-list-detail',
                                'describe' => '列表详情'
                            ],
                            [
                                'id' => 137,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.编辑活动场次',
                                'route' => 'app-center/seckill/goods-detail',
                                'describe' => '商品详情'
                            ],
                            [
                                'id' => 138,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.限时购',
                                'route' => 'app-center/seckill/delete',
                                'describe' => '删除'
                            ],
                            [
                                'id' => 139,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.限时购',
                                'route' => 'app-center/seckill/enable',
                                'describe' => '启用/禁用'
                            ],
                            [
                                'id' => 140,
                                'pid' => 131,
                                'page' => '应用.我的应用.限时购.编辑活动场次',
                                'route' => 'app-center/seckill/goods',
                                'describe' => '商品选择列表'
                            ]
                        ]
                    ]
                ]
            ],
            [
                'id' => 141,
                'pid' => 0,
                'page' => '设置',
                'route' => 'setting',
                'describe' => '设置',
                'child' => [
                    [
                        'id' => 142,
                        'pid' => 141,
                        'page' => '设置.支付配置',
                        'route' => 'setting/callback',
                        'describe' => 'UC在企业微信配置的回调地址',
                        'child' => []
                    ],
                    [
                        'id' => 143,
                        'pid' => 141,
                        'page' => '设置.账号权限.角色',
                        'route' => 'setting/role',
                        'describe' => '角色',
                        'child' => [
                            [
                                'id' => 144,
                                'pid' => 143,
                                'page' => '设置.账号权限.角色',
                                'route' => 'setting/role/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 145,
                                'pid' => 143,
                                'page' => '设置.账号权限.角色.编辑',
                                'route' => 'setting/role/detail',
                                'describe' => '详情'
                            ],
                            [
                                'id' => 146,
                                'pid' => 143,
                                'page' => '设置.账号权限.角色.编辑',
                                'route' => 'setting/role/power',
                                'describe' => '所有权限'
                            ],
                            [
                                'id' => 147,
                                'pid' => 143,
                                'page' => '设置.账号权限.角色.添加',
                                'route' => 'setting/role/add',
                                'describe' => '添加'
                            ],
                            [
                                'id' => 148,
                                'pid' => 143,
                                'page' => '设置.账号权限.角色.编辑',
                                'route' => 'setting/role/update',
                                'describe' => '更新'
                            ],
                            [
                                'id' => 149,
                                'pid' => 143,
                                'page' => '设置.账号权限.角色',
                                'route' => 'setting/role/delete',
                                'describe' => '删除'
                            ]
                        ]
                    ],
                    [
                        'id' => 150,
                        'pid' => 141,
                        'page' => '设置.账号权限.账号',
                        'route' => 'setting/enterprise-admin',
                        'describe' => '管理员',
                        'child' => [
                            [
                                'id' => 151,
                                'pid' => 150,
                                'page' => '设置.账号权限.账号',
                                'route' => 'setting/enterprise-admin/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 152,
                                'pid' => 150,
                                'page' => '设置.账号权限.账号.编辑',
                                'route' => 'setting/enterprise-admin/detail',
                                'describe' => '详情'
                            ],
                            [
                                'id' => 153,
                                'pid' => 150,
                                'page' => '设置.账号权限.账号',
                                'route' => 'setting/enterprise-admin/reset-password',
                                'describe' => '重置密码'
                            ],
                            [
                                'id' => 154,
                                'pid' => 150,
                                'page' => '设置.账号权限.账号.新增账号',
                                'route' => 'setting/enterprise-admin/add',
                                'describe' => '添加'
                            ],
                            [
                                'id' => 155,
                                'pid' => 150,
                                'page' => '设置.账号权限.账号.编辑',
                                'route' => 'setting/enterprise-admin/update',
                                'describe' => '更新'
                            ],
                            [
                                'id' => 156,
                                'pid' => 150,
                                'page' => '设置.账号权限.账号',
                                'route' => 'setting/enterprise-admin/delete',
                                'describe' => '删除'
                            ]
                        ]
                    ],
                    [
                        'id' => 157,
                        'pid' => 141,
                        'page' => '设置.账号权限.操作日志',
                        'route' => 'setting/operation-log',
                        'describe' => '操作日志',
                        'child' => [
                            [
                                'id' => 158,
                                'pid' => 157,
                                'page' => '设置.账号权限.操作日志',
                                'route' => 'setting/operation-log/module',
                                'describe' => '模块列表'
                            ],
                            [
                                'id' => 159,
                                'pid' => 157,
                                'page' => '设置.账号权限.操作日志',
                                'route' => 'setting/operation-log/list',
                                'describe' => '列表'
                            ]
                        ]
                    ]
                ]
            ],
            [
                'id' => 160,
                'pid' => 0,
                'page' => '设置.客户',
                'route' => 'label',
                'describe' => '标签',
                'child' => [
                    [
                        'id' => 161,
                        'pid' => 160,
                        'page' => '设置.客户.标签',
                        'route' => 'setting/label/impress',
                        'describe' => '印象标签',
                        'child' => [
                            [
                                'id' => 162,
                                'pid' => 161,
                                'page' => '设置.客户.标签',
                                'route' => 'setting/label/impress/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 163,
                                'pid' => 161,
                                'page' => '设置.客户.标签.新增印象标签',
                                'route' => 'setting/label/impress/add',
                                'describe' => '新增'
                            ],
                            [
                                'id' => 164,
                                'pid' => 161,
                                'page' => '设置.客户.标签.编辑印象标签',
                                'route' => 'setting/label/impress/update',
                                'describe' => '更新'
                            ],
                            [
                                'id' => 165,
                                'pid' => 161,
                                'page' => '设置.客户.标签.编辑印象标签',
                                'route' => 'setting/label/impress/detail',
                                'describe' => '详情'
                            ],
                            [
                                'id' => 166,
                                'pid' => 161,
                                'page' => '设置.客户.标签',
                                'route' => 'setting/label/impress/delete',
                                'describe' => '删除'
                            ],
                            [
                                'id' => 167,
                                'pid' => 161,
                                'page' => '设置.客户.标签',
                                'route' => 'setting/label/impress/status',
                                'describe' => '状态'
                            ]
                        ]
                    ],
                    [
                        'id' => 168,
                        'pid' => 160,
                        'page' => '设置.客户.标签',
                        'route' => 'setting/label/system',
                        'describe' => '印象标签',
                        'child' => [
                            [
                                'id' => 169,
                                'pid' => 168,
                                'page' => '设置.客户.标签',
                                'route' => 'setting/label/system/list',
                                'describe' => '列表'
                            ],
                            [
                                'id' => 170,
                                'pid' => 168,
                                'page' => '设置.客户.标签.新增标签',
                                'route' => 'setting/label/system/add',
                                'describe' => '新增'
                            ],
                            [
                                'id' => 171,
                                'pid' => 168,
                                'page' => '设置.客户.标签.编辑标签',
                                'route' => 'setting/label/system/update',
                                'describe' => '更新'
                            ],
                            [
                                'id' => 172,
                                'pid' => 168,
                                'page' => '设置.客户.标签.编辑标签',
                                'route' => 'setting/label/system/detail',
                                'describe' => '详情'
                            ],
                            [
                                'id' => 173,
                                'pid' => 168,
                                'page' => '设置.客户.标签',
                                'route' => 'setting/label/system/delete',
                                'describe' => '删除'
                            ],
                            [
                                'id' => 174,
                                'pid' => 168,
                                'page' => '设置.客户.标签',
                                'route' => 'setting/label/system/good',
                                'describe' => '商品'
                            ],
                            [
                                'id' => 175,
                                'pid' => 168,
                                'page' => '设置.客户.标签',
                                'route' => 'setting/label/system/status',
                                'describe' => '状态'
                            ]
                        ]
                    ]
                ]
            ],
            [
                'id' => 176,
                'pid' => 0,
                'page' => '设置.客户',
                'route' => 'customer-portrait',
                'describe' => '画像统计',
                'child' => [
                    [
                        'id' => 177,
                        'pid' => 176,
                        'page' => '设置.客户.客户画像',
                        'route' => 'customer-portrait/age',
                        'describe' => '年龄统计'
                    ],
                    [
                        'id' => 178,
                        'pid' => 176,
                        'page' => '设置.客户.客户画像',
                        'route' => 'customer-portrait/gender',
                        'describe' => '性别统计'
                    ],
                    [
                        'id' => 179,
                        'pid' => 176,
                        'page' => '设置.客户.客户画像',
                        'route' => 'customer-portrait/label-page',
                        'describe' => '标签分页统计'
                    ],
                    [
                        'id' => 180,
                        'pid' => 176,
                        'page' => '设置.客户.客户画像',
                        'route' => 'customer-portrait/label',
                        'describe' => '标签统计'
                    ],
                    [
                        'id' => 181,
                        'pid' => 176,
                        'page' => '设置.客户.客户画像',
                        'route' => 'customer-portrait/provinces-page',
                        'describe' => '区域分页统计'
                    ],
                    [
                        'id' => 182,
                        'pid' => 176,
                        'page' => '设置.客户.客户画像',
                        'route' => 'customer-portrait/provinces',
                        'describe' => '区域不分页统计'
                    ]
                ]
            ]
        ];

        $data = $this->mergeData($data);

        foreach ($data as $key => $value) {
            $data[$key]['created_at'] = time();
            $data[$key]['updated_at'] = time();
        }

        DB::table('api_power')->insert($data);
    }

    private function mergeData($data)
    {
        foreach ($data as &$v) {

            if (empty($v['child'])) {
                unset($v['child']);
                $this->insertData[] = $v;

            } else {

                $arr = $v['child'];
                $this->mergeData($arr);
                unset($v['child']);
                $this->insertData[] = $v;
            }
        }

        return $this->insertData;
    }

}