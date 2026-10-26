<?php

/**
 * @var Router $router
 */

/**
 * 小程序授权
 */

use Laravel\Lumen\Routing\Router;

$router->group([
    'prefix'    => 'auth',
    'namespace' => 'Auth'
], function () use ($router) {
    $router->post('shop-auth-detail', [
        'uses'     => 'ShopAuthDetailController@index',
        'describe' => '小程序授权.商城小程序授权详情',
        'page'     => '设置.小程序授权.商城小程序授权'
    ]);

    // todo 废弃
    $router->post('work-auth-detail', [
        'uses'     => 'WorkAuthDetailController@index',
        'describe' => '小程序授权.企业小程序授权详情(废弃)',
        'page'     => '设置.小程序授权.企业微信授权'
    ]);

    $router->post('auth-contact', [
        'uses'     => 'ContactAuthDetailController@index',
        'describe' => '小程序授权.通讯录授权详情',
        'page'     => '设置.小程序授权.客户联系配置'
    ]);

    // todo 废弃
    $router->post('work-auth-complete', [
        'uses'     => 'WorkCompleteController@index',
        'describe' => '小程序授权.企业微信小程序完善信息(废弃)',
        'page'     => '设置.小程序授权.完善资料'
    ]);

    $router->post('shop-auth-complete', [
        'uses'     => 'ShopCompleteController@index',
        'describe' => '小程序授权.商城小程序完善信息',
        'page'     => '设置.小程序授权.完善资料'
    ]);

    // todo 废弃
    $router->post('auth-complete', [
        'uses'     => 'AuthCompleteController@index',
        'describe' => '小程序授权.统一完善信息(废弃)',
        'page'     => ''
    ]);

    $router->post('shop-save', [
        'uses'     => 'ShopSaveController@index',
        'describe' => '小程序授权.商城小程序信息编辑',
        'page'     => ''
    ]);

    // todo 废弃
    $router->post('work-save', [
        'uses'     => 'WorkSaveController@index',
        'describe' => '小程序授权.企业小程序信息编辑(废弃)',
        'page'     => ''
    ]);

    $router->post('mini-detail', [
        'uses'     => 'AuthDetailController@index',
        'describe' => '小程序授权.小程序信息详情',
        'page'     => ''
    ]);
    $router->post('work-auth-complete-new', [
        'uses'     => 'WorkCompleteNewController@index',
        'describe' => '小程序授权.企业微信完善信息-新',
        'page'     => '设置.小程序授权.完善资料'
    ]);

    // todo 废弃
    $router->post('work-save-new', [
        'uses'     => 'WorkSaveNewController@index',
        'describe' => '小程序授权.企业微信信息编辑-新(废弃)',
        'page'     => ''
    ]);

    $router->post('work-mini-list', [
        'uses'     => 'GuideMiniListController@index',
        'describe' => '企业微信小程序授权.小程序列表',
        'page'     => ''
    ]);
    $router->post('install-app', [
        'uses'     => 'InstallAppController@install',
        'describe' => '企业微信.自建应用安装',
        'page'     => ''
    ]);
    $router->post('install-third-app', [
        'uses'     => 'InstallAppController@installThirdApp',
        'describe' => '企业微信.第三方应用安装（智慧导购）',
        'page'     => ''
    ]);
    $router->get('get-work-plugin', [
        'uses'     => 'InstallAppController@getWorkPlugin',
        'describe' => '企业微信.获取自建应用安装信息',
        'page'     => ''
    ]);
    $router->post('register-apply', [
        'uses'     => 'FastRegisterController@apply',
        'describe' => '代注册小程序申请',
        'page'     => ''
    ]);
    $router->post('register-detail', [
        'uses'     => 'FastRegisterController@detail',
        'describe' => '代注册小程序申请详情',
        'page'     => ''
    ]);
    $router->post('register-update-info', [
        'uses'     => 'FastRegisterController@updateInfo',
        'describe' => '代注册小程序更新信息',
        'page'     => ''
    ]);
});

/**
 * 分类 - 刘瑞
 */
$router->group([
    'prefix'    => 'category',
    'namespace' => 'Category'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '分类.列表',
        'page'     => '商品分类.品类.列表'
    ]);
    $router->post('recursion-list', [
        'uses'     => 'RecursionListController@index',
        'describe' => '分类.递归列表',
        'page'     => '商品分类.品类.列表'
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '分类.新增',
        'page'     => '商品分类.品类.新增分类'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '分类.删除',
        'page'     => '商品分类.品类.删除分类'
    ]);
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '分类.编辑',
        'page'     => '商品分类.品类.编辑分类'
    ]);
    $router->get('{id}', [
        'uses'     => 'DetailController@index',
        'describe' => '分类.详情',
        'page'     => '商品分类.品类.分类详情'
    ]);

    $router->post('category-list', [
        'uses'     => 'CategoryListController@index',
        'describe' => '分类.树形商品分类列表',
        'page'     => '商品分类.品类.列表'
    ]);
});

/**
 * 订单 - 刘瑞
 */
$router->group([
    'prefix'    => 'order',
    'namespace' => 'Order'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '订单.列表',
        'page'     => '订单.列表'
    ]);
    $router->post('user-detail', [
        'uses'     => 'UserDetailController@index',
        'describe' => '订单.拼团用户详情',
        'page'     => '订单.拼团用户详情'
    ]);
    $router->get('order-export', [
        'uses'     => 'ExportController@index',
        'describe' => '订单.导出',
        'page'     => '订单.列表'
    ]);
    $router->post('order-export-new', [
        'uses'     => 'ExportNewController@index',
        'describe' => '订单.异步导出',
        'page'     => '订单.列表'
    ]);
    $router->get('order-hand-export', [
        'uses'     => 'HandExportController@index',
        'describe' => '订单.人工导出',
        'page'     => '订单.列表'
    ]);
    $router->post('refund-list', [
        'uses'     => 'RefundListController@index',
        'describe' => '订单.导出',
        'page'     => '订单.列表'
    ]);
    $router->post('fail-order-download', [
        'uses'     => 'FailOrderDownloadController@index',
        'describe' => '订单.下载上传失败的订单',
        'page'     => '订单.下载上传失败的订单'
    ]);

    $router->get('download-send-goods-batch-template', [
        'uses'     => 'SendGoodsBatchController@downloadTemplate',
        'describe' => '订单.下载手工发货模板',
        'page'     => '订单.上传发货物流信息'
    ]);
    $router->post('send-goods-batch', [
        'uses'     => 'SendGoodsBatchController@index',
        'describe' => '订单.上传发货物流信息',
        'page'     => '订单.上传发货物流信息'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '订单.详情',
        'page'     => '订单.详情'
    ]);
    $router->post('remark-add', [
        'uses'     => 'RemarkController@add',
        'describe' => '订单.详情.备注',
        'page'     => '订单.详情.备注'
    ]);
    $router->post('refund-total', [
        'uses'     => 'RefundTotalController@index',
        'describe' => '订单.售后统计',
        'page'     => '订单.售后'
    ]);
    $router->post('address-save', [
        'uses'     => 'AddressSaveController@index',
        'describe' => '订单.修改地址',
        'page'     => '订单.修改地址'
    ]);
    $router->post('after-sale-list', [
        'uses'     => 'AfterSaleListController@index',
        'describe' => '订单.售后列表',
        'page'     => '订单.售后列表'
    ]);
    $router->get('down-apply-list', [
        'uses'     => 'ExportAfterSaleListController@index',
        'describe' => '订单.售后导出',
        'page'     => '订单.售后导出'
    ]);
    $router->get('refund-detail', [
        'uses'     => 'RefundDetailController@index',
        'describe' => '订单.售后详情',
        'page'     => '订单.售后列表'
    ]);
    $router->post('refund-goods-edit', [
        'uses'     => 'RefundGoodsEditController@index',
        'describe' => '订单.售后详情.修改退款商品',
        'page'     => '订单.售后列表.售后详情'
    ]);
    $router->post('logistics-delete', [
        'uses'     => 'LogisticsDeleteController@index',
        'describe' => '订单.详情.删除订单物流',
        'page'     => '订单.详情.删除订单物流'
    ]);
    $router->post('remark-to-read', [
        'uses'     => 'RemarkController@remarkToRead',
        'describe' => '订单.详情.更改用户备注为已读',
        'page'     => '订单.详情.更改用户备注为已读'
    ]);
    $router->post('refund-info-data', [
        'uses'     => 'RefundInfoDataController@index',
        'describe' => '订单.详情.获取售后信息弹窗数据',
        'page'     => '订单.详情.获取售后信息弹窗数据'
    ]);
    $router->post('after-sale-add', [
        'uses'     => 'AfterSaleAddController@index',
        'describe' => '订单.列表.创建手工售后单',
        'page'     => '订单.列表.创建手工售后单'
    ]);
    $router->post('search-by-code', [
        'uses'     => 'TakeGoodsCodeOrderController@index',
        'describe' => '订单.提货.根据自提码搜索订单',
        'page'     => '订单.提货.根据自提码搜索订单'
    ]);
    $router->post('export-order-statement', [
        'uses'     => 'ExportOrderStatementController@index',
        'describe' => '订单.对账导出',
        'page'     => '订单.列表'
    ]);
    // 订单统计
    $router->post('order-statistical', [
        'uses'     => 'OrderStatisticalController@index',
        'describe' => '订单.统计',
        'page'     => ''
    ]);
    // 按日期查看订单统计
    $router->post('order-statistical-by-date', [
        'uses'     => 'OrderStatisticalController@statisticalByDate',
        'describe' => '订单.统计按日期',
        'page'     => ''
    ]);
    // 按导购查看订单统计
    $router->post('order-statistical-by-guide', [
        'uses'     => 'OrderStatisticalController@statisticalByGuide',
        'describe' => '订单.统计按导购',
        'page'     => ''
    ]);
    //导出订单统计
    $router->post('order-statistical-export', [
        'uses'     => 'OrderStatisticalExportController@index',
        'describe' => '订单导出.数据罗盘统计',
        'page'     => ''
    ]);
    $router->post('batch-send-goods', [
        'uses'     => 'BatchSendGoodsController@index',
        'describe' => '批量发货接口',
        'page'     => '订单.待发货列表'
    ]);
    //数据罗盘-订单数据-导出订单明细
    $router->post('order-statistical-detail-export', [
        'uses'     => 'OrderStatisticalDetailExportController@index',
        'describe' => '数据罗盘.导出订单明细',
        'page'     => ''
    ]);
    $router->post('repair-wechat-refund-order', [
        'uses'     => 'RepairWechatRefundOrderController@index',
        'describe' => '修复退款报"订单已全额退款"的订单（运营在微信后台操作了全额退款）',
        'page'     => '订单.修复退款报"订单已全额退款"'
    ]);
});

/**
 * 退货原因-党志强
 */
$router->group([
    'prefix'    => 'refund',
    'namespace' => 'Refund'
], function () use ($router) {
    $router->post('reason-save', [
        'uses'     => 'ReasonSaveController@index',
        'describe' => '退货原因.保存',
        'page'     => '退货原因.保存'
    ]);
    $router->get('reason-list', [
        'uses'     => 'ReasonListController@index',
        'describe' => '退货原因.列表',
        'page'     => '退货原因.列表'
    ]);
    $router->get('reason-delete', [
        'uses'     => 'ReasonDeleteController@index',
        'describe' => '退货原因.删除',
        'page'     => '退货原因.删除'
    ]);
    $router->post('hand-refund', [
        'uses'     => 'HandRefundController@index',
        'describe' => '退货列表.是否同意退款',
        'page'     => '退货列表.是否同意退款'
    ]);
    $router->get('get-white-ids', [
        'uses'     => 'GetWhiteController@index',
        'describe' => '退货列表.可以显示退款按钮的企业ID',
        'page'     => '退货列表.可以显示退款按钮的企业ID'
    ]);

    $router->post('machine-run', [
        'uses'     => 'ChangeRefundApplyStateController@index',
        'describe' => '后台同意退款',
        'page'     => '后台同意退款'
    ]);

    $router->post('refund-to-manual', [
        'uses'     => 'RefundToManualController@index',
        'describe' => '后台一键转人工',
        'page'     => '后台一键转人工'
    ]);

    $router->post('reason-sort', [
        'uses'     => 'ReasonSortController@index',
        'describe' => '售后理由.排序',
        'page'     => '售后理由.列表'
    ]);
});

/**
 * 页面管理
 */
$router->group([
    'prefix'    => 'page',
    'namespace' => 'Page'
], function () use ($router) {
    /**
     * 轮播图
     */
    $router->group([
        'prefix'    => 'banner',
        'namespace' => 'Banner'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '页面管理.轮播图.列表',
            'page'     => '店铺.页面.轮播图'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.轮播图.保存',
            'page'     => '店铺.页面.轮播图保存'
        ]);
    });

    /**
     * 引导区
     */
    $router->group([
        'prefix'    => 'guide',
        'namespace' => 'Guide'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '页面管理.引导区.列表',
            'page'     => '店铺.页面.首页引导区列表'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.引导区.保存',
            'page'     => '店铺.页面.首页引导区保存'
        ]);
    });

    /**
     * 爆品推荐
     */
    $router->group([
        'prefix'    => 'hot',
        'namespace' => 'Hot'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '页面管理.爆品推荐.列表',
            'page'     => '店铺.页面.爆品推荐列表'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.爆品推荐.保存',
            'page'     => '店铺.页面.爆品推荐保存'
        ]);
        $router->post('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '页面管理.爆品推荐.爆品推荐商品列表',
            'page'     => '店铺.页面.爆品推荐选择商品'
        ]);
    });

    /**
     * 商品推荐
     */
    $router->group([
        'prefix'    => 'recommend',
        'namespace' => 'Recommend'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '页面管理.店长推荐.列表',
            'page'     => '店铺.页面.店长推荐列表'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.店长推荐.保存',
            'page'     => '店铺.页面.店长推荐保存'
        ]);
        $router->post('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '页面管理.商品推荐.商品列表',
            'page'     => '店铺.页面.商品推荐选择商品'
        ]);
    });

    /**
     * 猜你喜欢
     */
    $router->group([
        'prefix'    => 'guess',
        'namespace' => 'Guess'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '页面管理.猜你喜欢.列表',
            'page'     => '店铺.页面.猜你喜欢'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.猜你喜欢.保存',
            'page'     => '店铺.页面.猜你喜欢保存'
        ]);
    });

    /**
     * 立减金管理
     */
    $router->group([
        'prefix'    => 'reduce-money',
        'namespace' => 'ReduceMoney'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '页面管理.立减金管理.列表',
            'page'     => '店铺.页面.立减金管理列表'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '页面管理.立减金管理.添加',
            'page'     => '店铺.页面.立减金管理添加'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '页面管理.立减金管理.保存',
            'page'     => '店铺.页面.立减金管理保存'
        ]);
        $router->get('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '页面管理.立减金管理.删除',
            'page'     => '店铺.页面.立减金管理删除'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '页面管理.立减金管理.详情',
            'page'     => '店铺.页面.立减金管理详情'
        ]);
    });

    /**
     * 页底
     */
    $router->group([
        'prefix'    => 'footer',
        'namespace' => 'Footer'
    ], function () use ($router) {
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '页面管理.页底管理.添加',
            'page'     => '店铺.页面.页底管理添加'
        ]);
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '页面管理.页底管理.保存',
            'page'     => '店铺.页面.页底管理保存'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '页面管理.页底管理.详情',
            'page'     => '店铺.页面.页底管理详情'
        ]);
    });

    /**
     * 小程序收藏
     */
    $router->group([
        'prefix'    => 'pop-content',
        'namespace' => 'PopContent'
    ], function () use ($router) {
        $router->post('save', [
            'uses'     => 'PopContentController@contentSave',
            'describe' => '页面管理.页底管理.添加',
            'page'     => '店铺.页面.小程序收藏添加'
        ]);
        $router->post('detail', [
            'uses'     => 'PopContentController@contentDetail',
            'describe' => '页面管理.页底管理.详情',
            'page'     => '店铺.页面.小程序收藏详情'
        ]);

    });
});

/**
 * 礼劵
 */
$router->group([
    'prefix'    => 'coupon',
    'namespace' => 'Coupon'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '礼劵.列表',
        'page'     => '应用.我的应用.礼券.礼券列表'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '礼劵.详情',
        'page'     => '应用.我的应用.礼券.礼券列表'
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '礼劵.新增',
        'page'     => '应用.我的应用.礼券.新增礼券'
    ]);
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '礼劵.编辑',
        'page'     => '应用.我的应用.礼券.编辑礼券'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '礼劵.删除',
        'page'     => '应用.我的应用.礼券.礼券列表'
    ]);
});

/**
 * 任务
 */
$router->group([
    'prefix'    => 'task',
    'namespace' => 'Task'
], function () use ($router) {
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '任务.新增',
        'page'     => '应用.任务.任务中心.新增任务'
    ]);
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '任务.列表',
        'page'     => '应用.任务.任务中心.任务列表'
    ]);
    $router->get('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '任务.删除',
        'page'     => '应用.任务.任务中心.任务列表'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '任务.详情',
        'page'     => '应用.任务.任务中心.任务详情'
    ]);
    $router->get('member-list', [
        'uses'     => 'MemberListController@index',
        'describe' => '任务.员工列表',
        'page'     => '应用.任务.任务中心.新增任务'
    ]);
    $router->post('type-list-page', [
        'uses'     => 'TypePageListController@index',
        'describe' => '任务.类型-分页',
        'page'     => '应用.任务.任务中心.任务列表'
    ]);
});

/**
 * 佣金
 */
$router->group([
    'prefix'    => 'commission',
    'namespace' => 'Commission'
], function () use ($router) {
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '佣金.详情',
        'page'     => '组织.佣金规则'
    ]);
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '佣金.保存',
        'page'     => '组织.佣金规则'
    ]);
});

/**
 * 会员等级
 */
$router->group([
    'prefix'    => 'user-level',
    'namespace' => 'UserLevel'
], function () use ($router) {
    // 列表
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '会员等级.列表',
        'page'     => '客户.会员等级'
    ]);
    // 编辑
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '会员等级.编辑',
        'page'     => '客户.会员等级.编辑会员等级'
    ]);
    // 详情
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '会员等级.详情',
        'page'     => '客户.会员等级.编辑会员等级'
    ]);
    // 状态变更
    $router->post('set-status', [
        'uses'     => 'SetStatusController@index',
        'describe' => '会员等级.状态变更',
        'page'     => '客户.会员等级.会员等级'
    ]);
});

/**
 * 组织管理
 */
$router->group([
    'prefix'    => 'department',
    'namespace' => 'Department'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '组织管理.列表',
        'page'     => '客户.组织'
    ]);
    $router->post('third-list', [
        'uses'     => 'ListController@thirdList',
        'describe' => '组织管理.列表',
        'page'     => '客户.组织'
    ]);
    $router->get('page-list', [
        'uses'     => 'PageListController@index',
        'describe' => '组织管理.分页列表',
        'page'     => '客户.组织'
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '组织管理.新增',
        'page'     => '客户.组织.新增部门'
    ]);
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '组织管理.编辑',
        'page'     => '客户.组织.编辑部门'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '组织管理.删除',
        'page'     => '客户.组织'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '组织管理.详情',
        'page'     => '客户.组织'
    ]);
    $router->post('sort', [
        'uses'     => 'SortController@index',
        'describe' => '客户.组织',
        'page'     => '组织管理.排序'
    ]);
});

/**
 * 客户管理
 */
$router->group([
    'prefix'    => 'customer',
    'namespace' => 'Customer'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '客户管理.列表(暂且废弃)',
        'page'     => ''
    ]);
    $router->get('customer-list', [
        'uses'     => 'ListController@customerList',
        'describe' => '客户管理.列表(新开发2020-04-16)',
        'page'     => ''
    ]);
    $router->get('buy-record', [
        'uses'     => 'ListController@buyRecord',
        'describe' => '客户管理.客户购买记录列表(新开发2020-04-16)',
        'page'     => '客户.客户购买记录列表'
    ]);
    $router->post('member-list', [
        'uses'     => 'MemberListController@index',
        'describe' => '客户管理.会员客户列表',
        'page'     => '客户.会员客户列表'
    ]);
    $router->post('order-member', [
        'uses'     => 'OrderMemberListController@index',
        'describe' => '客户管理.下单客户列表',
        'page'     => '客户.下单客户列表'
    ]);
    $router->post('all-list', [
        'uses'     => 'AllCustomerListController@index',
        'describe' => '客户管理.全部客户列表',
        'page'     => '客户.全部客户列表'
    ]);
    $router->post('friend-list', [
        'uses'     => 'MemberCustomerListController@index',
        'describe' => '客户管理.导购好友客户列表',
        'page'     => '客户.导购好友客户列表'
    ]);
    $router->get('all-list-export', [
        'uses'     => 'AllCustomerListExportController@index',
        'describe' => '客户管理.导出全部客户列表',
        'page'     => '客户.全部客户列表'
    ]);
    $router->get('friend-list-export', [
        'uses'     => 'MemberCustomerListExportController@index',
        'describe' => '客户管理.导出导购好友客户列表',
        'page'     => '客户.导购好友客户列表'
    ]);
    $router->get('order-member-export', [
        'uses'     => 'OrderMemberExportController@index',
        'describe' => '客户管理.导出下单客户列表',
        'page'     => '客户.下单客户列表'
    ]);
    $router->get('member-list-export', [
        'uses'     => 'MemberListExportController@index',
        'describe' => '客户管理.导出会员客户列表',
        'page'     => '客户.会员客户列表'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '客户管理.基础信息',
        'page'     => '客户.客户详情'
    ]);
    $router->get('statistics', [
        'uses'     => 'StatisticsController@index',
        'describe' => '客户管理.统计数据',
        'page'     => '客户.客户详情'
    ]);
    $router->get('order-list', [
        'uses'     => 'OrderListController@index',
        'describe' => '客户管理.订单列表',
        'page'     => '客户.订单列表'
    ]);
    $router->get('order-detail', [
        'uses'     => 'OrderDetailController@index',
        'describe' => '客户管理.订单详情',
        'page'     => '客户.订单详情'
    ]);
    $router->get('addition_list', [
        'uses'     => 'AdditionListController@index',
        'describe' => '客户管理.待添加客户列表',
        'page'     => '客户.待添加客户列表'
    ]);
    $router->get('mark_user', [
        'uses'     => 'MarkUserController@index',
        'describe' => '客户管理.复制手机号标记',
        'page'     => '客户.复制手机号标记'
    ]);
});

/**
 * 员工管理
 */
$router->group([
    'prefix'    => 'user',
    'namespace' => 'User'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '员工管理.列表',
        'page'     => '组织.在职员工列表.在职员工列表'
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '员工管理.新增',
        'page'     => '组织.在职员工列表.添加员工'
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '员工管理.详情',
        'page'     => '组织.在职员工列表.编辑员工'
    ]);
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '员工管理.编辑',
        'page'     => '组织.在职员工列表.编辑员工'
    ]);
    $router->post('leave', [
        'uses'     => 'LeaveController@index',
        'describe' => '员工管理.离职',
        'page'     => '组织.在职员工列表.在职员工列表'
    ]);
    $router->get('handover-list', [
        'uses'     => 'HandoverListController@index',
        'describe' => '员工管理.交接客户列表',
        'page'     => '组织.离职员工列表.已交接客户列表'
    ]);
    $router->post('handover', [
        'uses'     => 'HandoverController@index',
        'describe' => '员工管理.交接客户',
        'page'     => '组织.离职员工列表.已交接客户列表'
    ]);
    $router->post('department-user-all-list', [
        'uses'     => 'DepartmentUserAllListController@index',
        'describe' => '员工管理.部门分层列表',
        'page'     => '组织'
    ]);
    $router->post('change_enterprise_allocation_status', [
        'uses'     => 'ChangeEnterpriseAllocationStatusController@index',
        'describe' => '员工管理.更改企业分配客户状态',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('change_allocation_status', [
        'uses'     => 'ChangeAllocationStatusController@index',
        'describe' => '员工管理.更改员工分配客户状态',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('guide_list', [
        'uses'     => 'GuideController@index',
        'describe' => '员工管理.导购分配员工列表',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('other_guides', [
        'uses'     => 'OtherGuidesController@index',
        'describe' => '员工管理.转交导购列表',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('guide_users', [
        'uses'     => 'GuideUsersController@index',
        'describe' => '员工管理.导购客户列表',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('customer_transfer', [
        'uses'     => 'CustomerTransferController@index',
        'describe' => '员工管理.转交客户',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('guides', [
        'uses'     => 'GuideController@list',
        'describe' => '导购列表(演示用)',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('get_dept_list', [
        'uses'     => 'GuideController@getDeptList',
        'describe' => '门店对应部门列表(演示用)',
        'page'     => '组织.客户分配设置'
    ]);
    $router->post('guide-list-by-store', [
        'uses'     => 'GuideListByStoreOrDepartmentController@index',
        'describe' => '根据门店或部门id查询导购列表',
        'page'     => ''
    ]);
});

/**
 * 员工角色(暂不维护，由前端写死)
 */
$router->group([
    'prefix'    => 'role',
    'namespace' => 'Role'
], function () use ($router) {
    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '员工角色.列表',
        'page'     => ''
    ]);
    $router->post('add', [
        'uses'     => 'AddController@index',
        'describe' => '员工角色.新增',
        'page'     => ''
    ]);
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '员工角色.详情',
        'page'     => ''
    ]);
    $router->post('edit', [
        'uses'     => 'EditController@index',
        'describe' => '员工角色.编辑',
        'page'     => ''
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '员工角色.删除',
        'page'     => ''
    ]);
});

/**
 * 属性 - 朱虹达
 */
$router->group([
    'prefix'    => 'attributes',
    'namespace' => 'Attributes'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '属性.列表',
        'page'     => '商品.规格'
    ]);
    $router->post('save', [
        'uses'     => 'AddController@index',
        'describe' => '属性.新增',
        'page'     => '商品.规格'
    ]);
    $router->post('del', [
        'uses'     => 'DeleteController@index',
        'describe' => '属性.删除',
        'page'     => '商品.规格'
    ]);
    $router->post('category-list', [
        'uses'     => 'CategoryListController@index',
        'describe' => '属性.规格列表',
        'page'     => '商品.规格'
    ]);
});

/**
 * 首页 - 朱虹达
 */
$router->group([
    'prefix'    => 'home',
    'namespace' => 'Home'
], function () use ($router) {
    $router->post('chart', [
        'uses'     => 'ChartController@index',
        'describe' => '首页.下半部分图表',
        'page'     => '首页'
    ]);
    $router->post('task', [
        'uses'     => 'TaskController@index',
        'describe' => '首页.上半部分任务达标率',
        'page'     => '首页'
    ]);
    $router->post('statistical', [
        'uses'     => 'StatisticalController@index',
        'describe' => '首页.客户概况统计',
        'page'     => '首页'
    ]);
    $router->post('order', [
        'uses'     => 'OrderController@index',
        'describe' => '首页.业绩概括统计',
        'page'     => '首页'
    ]);
    $router->get('contact-auth-status', [
        'uses'     => 'AuthDetailController@index',
        'describe' => '首页.获取企业是否授权通讯录',
        'page'     => '首页'
    ]);
    $router->get('statistics', [
        'uses'     => 'StatisticalController@statistics',
        'describe' => '首页.获取数据统计',
        'page'     => '首页'
    ]);
    // 门店列表
    $router->get('store-list', [
        'uses'     => 'StoreListController@index',
        'describe' => '管理后台首页.数据统计.门店列表',
        'page'     => '首页.数据统计.切换门店'
    ]);
    // 企业购买续费提醒
    $router->post('buy-notice', [
        'uses'     => 'BuyNoticeController@index',
        'describe' => '企业购买续费提醒',
        'page'     => '首页'
    ]);
});

/**
 * 商品管理 - 大雄
 */
$router->group([
    'prefix'    => 'goods',
    'namespace' => 'Goods'
], function () use ($router) {
    $router->post('', [
        'uses'     => 'ListController@index',
        'describe' => '商品管理.列表',
        'page'     => '商品.商品列表'
    ]);
    $router->post('skuList', [
        'uses'     => 'SkuListController@index',
        'describe' => '商品管理.SKU列表',
        'page'     => '商品.SKU列表'
    ]);
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '商品管理.详情',
        'page'     => '商品.商品详情'
    ]);
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '商品管理.新增',
        'page'     => '商品.商品新增'
    ]);
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '商品管理.编辑',
        'page'     => '商品.商品编辑'
    ]);
    $router->post('sale', [
        'uses'     => 'SaleController@index',
        'describe' => '商品管理.上下架',
        'page'     => '商品.商品列表'
    ]);
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '商品管理.删除',
        'page'     => '商品.商品列表'
    ]);
    $router->post('sale-check', [
        'uses'     => 'SaleCheckController@index',
        'describe' => '商品管理.下架验证',
        'page'     => '商品.商品列表'
    ]);
    $router->post('cancel-sale', [
        'uses'     => 'CancelSaleController@index',
        'describe' => '商品管理.撤销',
        'page'     => '商品.商品列表'
    ]);
    $router->post('delay-sale', [
        'uses'     => 'DelaySaleController@index',
        'describe' => '商品管理.定时上下架',
        'page'     => '商品.商品列表'
    ]);
    $router->get('export', [
        'uses'     => 'ExportController@index',
        'describe' => '商品管理.导出',
        'page'     => '商品.导出'
    ]);
    $router->post('activity-new-list', [
        'uses'     => 'ActivityNewListController@index',
        'describe' => '商品.活动列表NEW',
        'page'     => '商品.活动列表NEW'
    ]);
    $router->post('goods-sales', [
        'uses'     => 'GoodsSalesController@index',
        'describe' => '商品.月销量列表',
        'page'     => '商品.月销量列表'
    ]);
    $router->post('goods-sales-info', [
        'uses'     => 'GoodsFloatSalesInfoController@index',
        'describe' => '商品.浮动销量详情',
        'page'     => '商品.浮动销量详情'
    ]);
    $router->post('goods-sales-save', [
        'uses'     => 'GoodsFloatSalesSaveController@index',
        'describe' => '商品.浮动销量编辑',
        'page'     => '商品.浮动销量编辑'
    ]);

    // 导入物料编码更改商品库存信息
    $router->post('import-stock',
        [
            'uses'     => 'ImportStockController@index',
            'describe' => '商品.商品库存批量修改',
            'page'     => '商品.商品库存批量修改'
        ]);

    $router->get('download-goods-stock-template', [
        'uses'     => 'ImportStockController@downloadTemplate',
        'describe' => '商品.下载批量更新商品库存模板',
        'page'     => '商品.下载批量更新商品库存模板'
    ]);

    // 导入物料编码更改商品价格
    $router->post('import-sell-price',
        [
            'uses'     => 'ImportSellPriceController@index',
            'describe' => '商品.商品价格批量修改',
            'page'     => '商品.商品价格批量修改'
        ]);

    $router->get('download-goods-sell-price-template', [
        'uses'     => 'ImportSellPriceController@downloadTemplate',
        'describe' => '商品.下载批量更新商品价格模板',
        'page'     => '商品.下载批量更新商品价格模板'
    ]);
    $router->post('fail-goods-download', [
        'uses'     => 'FailGoodsDownloadController@index',
        'describe' => '商品.下载上传失败的商品',
        'page'     => '商品.下载上传失败的商品'
    ]);
    //活动接口
    $router->post('activity-spu-list', [
        'uses'     => 'ActivitySpuListController@index',
        'describe' => '商品.活动spu列表',
        'page'     => '商品.活动spu列表'
    ]);

    $router->post('activity-sku-list', [
        'uses'     => 'ActivitySkuListController@index',
        'describe' => '商品.活动sku列表',
        'page'     => '商品.活动sku列表'
    ]);
});

/**
 * 支付
 * 注意：目前是微信支付
 */
$router->group([
    'prefix'    => 'payment',
    'namespace' => 'PaymentSetting'
], function () use ($router) {
    $router->get('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '支付配置详情',
        'page'     => '设置.支付配置'
    ]);
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '支付配置编辑',
        'page'     => '设置.支付配置编辑'
    ]);
});

/**
 * 应用中心
 */
$router->group([
    'prefix'    => 'app-center',
    'namespace' => 'AppCenter'
], function () use ($router) {
    /**
     * 微信好物圈
     */
    $router->group(['prefix' => 'wx-mall', 'namespace' => 'WxMall'], function () use ($router) {
        // 物品类目导入
        $router->post('category-import', 'CategoryImportController@index');
        // 物品类目列表
        $router->post('category-list', 'CategoryListController@index');
        // 商品列表
        $router->post('goods-list', 'GoodsListController@index');
        // 导出全部商品
        $router->get('goods-export', 'GoodsExportController@index');
        // 商品批量设置好物圈分类
        $router->post('set-category', 'SetCategoryController@index');
    });

    $router->get('list', [
        'uses'     => 'ListController@index',
        'describe' => '应用中心',
        'page'     => '应用.全部应用'
    ]);

    /**
     * 新人优惠 - 梁欢
     */
    $router->group([
        'prefix'    => 'novice',
        'namespace' => 'Novice'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.新人优惠.列表',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->get('all', [
            'uses'     => 'AllController@index',
            'describe' => '应用中心.新人优惠.所有商品SPUID',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.新人优惠.添加',
            'page'     => '应用.我的应用.新人优惠.新增优惠商品'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.新人优惠.详情',
            'page'     => '应用.我的应用.新人优惠.编辑优惠商品'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '应用中心.新人优惠.编辑',
            'page'     => '应用.我的应用.新人优惠.编辑优惠商品'
        ]);
        $router->post('status', [
            'uses'     => 'StatusController@index',
            'describe' => '应用中心.新人优惠.状态修改',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.新人优惠.删除',
            'page'     => '应用.我的应用.新人优惠'
        ]);
        $router->get('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '应用中心.新人优惠.商品列表',
            'page'     => '应用.我的应用.新人优惠'
        ]);
    });

    /**
     * 拼团 - 党志强
     */
    $router->group([
        'prefix'    => 'group',
        'namespace' => 'Group'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.拼团.列表',
            'page'     => '应用.我的应用.拼团'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.拼团.新增',
            'page'     => '应用.我的应用.拼团.新增拼团商品'
        ]);
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.拼团.编辑',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.拼团.删除',
            'page'     => '应用.我的应用.拼团'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.拼团.详情',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('group-goods', [
            'uses'     => 'GroupGoodsController@index',
            'describe' => '应用中心.拼团.团购商品',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('group-setting', [
            'uses'     => 'GroupSettingController@index',
            'describe' => '应用中心.拼团.商品参加活动修改',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
        $router->get('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '应用中心.拼团.团购商品选择列表',
            'page'     => '应用.我的应用.拼团.查看拼团商品'
        ]);
    });

    /**
     * 收藏有礼
     */
    $router->group([
        'prefix'    => 'collection',
        'namespace' => 'Collection'
    ], function () use ($router) {
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.收藏有礼.用户提交列表',
            'page'     => '应用.我的应用.收藏有礼'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '应用中心.收藏有礼.设置',
            'page'     => '应用.我的应用.收藏有礼.设置'
        ]);
        $router->post('audit', [
            'uses'     => 'AuditController@index',
            'describe' => '应用中心.收藏有礼.审核',
            'page'     => '应用.我的应用.收藏有礼.审核'
        ]);
        $router->post('collection_list', [
            'uses'     => 'CollectionListController@index',
            'describe' => '应用中心.收藏有礼.列表',
            'page'     => '应用.我的应用.收藏有礼.列表'
        ]);
    });

    /**
     * 秒杀
     */
    $router->group([
        'prefix'    => 'seckill',
        'namespace' => 'Seckill'
    ], function () use ($router) {
        $router->post('', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.秒杀.列表',
            'page'     => '应用.我的应用.限时购'
        ]);
        // 新增
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '应用中心.秒杀.新增',
            'page'     => '应用.我的应用.限时购.新增活动场次'
        ]);
        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.秒杀.编辑',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.秒杀.详情',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 商品列表详情
        $router->post('goods-list-detail', [
            'uses'     => 'GoodsListDetailController@index',
            'describe' => '应用中心.秒杀.商品列表详情',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 商品详情
        $router->post('goods-detail', [
            'uses'     => 'GoodsDetailController@index',
            'describe' => '应用中心.秒杀.商品详情',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.秒杀.删除',
            'page'     => '应用.我的应用.限时购.限时购'
        ]);
        // 启用/禁用
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.秒杀.启用/禁用',
            'page'     => '应用.我的应用.限时购.限时购'
        ]);
        // 秒杀商品选择列表
        $router->post('goods', [
            'uses'     => 'GoodsController@index',
            'describe' => '应用中心.秒杀.秒杀商品选择列表',
            'page'     => '应用.我的应用.限时购.编辑活动场次'
        ]);
    });

    /**
     * 超低购
     */
    $router->group([
        'prefix'    => 'cheap',
        'namespace' => 'Cheap',
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.超低购.列表',
            'page'     => '应用.我的应用.超低购'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.超低购.新增',
            'page'     => '应用.我的应用.超低购.新增'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.超低购.详情',
            'page'     => '应用.我的应用.超低购.查看'
        ]);
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.超低购.状态修改',
            'page'     => '应用.我的应用.超低购.列表'
        ]);
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.超低购.删除',
            'page'     => '应用.我的应用.超低购.列表'
        ]);

        $router->post('sort', [
            'uses'     => 'SortController@index',
            'describe' => '应用中心.超低购.排序',
            'page'     => '应用.我的应用.超低购.排序'
        ]);
    });

    /**
     * 营销活动
     */
    $router->group([
        'prefix'    => 'activity',
        'namespace' => 'Activity'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.营销活动.列表',
            'page'     => '应用.我的应用.营销活动'
        ]);
        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.营销活动.新增',
            'page'     => '应用.我的应用.营销活动.新增活动'
        ]);
        // 编辑
        $router->post('edit', [
            'uses'     => 'EditController@index',
            'describe' => '应用中心.营销活动.编辑',
            'page'     => '应用.我的应用.营销活动.编辑活动'
        ]);
        // 详情
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.营销活动.详情',
            'page'     => '应用.我的应用.营销活动.详情'
        ]);
        // 删除
        $router->post('del', [
            'uses'     => 'DelController@index',
            'describe' => '应用中心.营销活动.删除',
            'page'     => '应用.我的应用.营销活动.活动删除'
        ]);
        // 启用/禁用
        $router->post('changeStatus', [
            'uses'     => 'ChangeStatusController@index',
            'describe' => '应用中心.营销活动.启用/禁用',
            'page'     => '应用.我的应用.营销活动.启用/禁用'
        ]);

    });

    /**
     * 微信优惠卷
     */
    $router->group([
        'prefix'    => 'wx-coupon',
        'namespace' => 'WxCoupon'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.微信优惠卷.列表',
            'page'     => '应用.微信优惠卷.微信优惠卷'
        ]);
        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.微信优惠卷.新增',
            'page'     => '应用.微信优惠卷.微信优惠卷.新增优惠卷'
        ]);
        // 编辑
        $router->post('edit', [
            'uses'     => 'EditController@index',
            'describe' => '应用中心.微信优惠卷.编辑',
            'page'     => '应用.我的应用.微信优惠卷.编辑优惠卷'
        ]);
        // 详情
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.微信优惠卷.详情',
            'page'     => '应用.我的应用.微信优惠卷.详情'
        ]);

        // 启用/禁用
        $router->post('changeStatus', [
            'uses'     => 'ChangeStatusController@index',
            'describe' => '应用中心.微信优惠卷.启用/禁用',
            'page'     => '应用.我的应用.微信优惠卷.启用/禁用'
        ]);
        // 启用/禁用
        $router->post('del', [
            'uses'     => 'DelController@index',
            'describe' => '应用中心.微信优惠卷.删除',
            'page'     => '应用.我的应用.微信优惠卷.删除'
        ]);

    });

    /**
     * 除去活动已选的商品列表
     */
    $router->get('coupon-spu-list', [
        'uses'     => 'CouponSpuListController@index',
        'describe' => '未参加活动商品列表',
        'page'     => '应用.我的应用.未参加活动商品'
    ]);

    // 抽奖活动
    $router->group(['prefix' => 'lottery', 'namespace' => 'Lottery'], function () use ($router) {
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '抽奖添加', 'page' => '应用,分享抽奖']);
        $router->get('list', ['uses' => 'ListController@index', 'describe' => '抽奖列表', 'page' => '应用,分享抽奖']);
        $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '抽奖编辑', 'page' => '应用,分享抽奖']);
        $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '抽奖详情', 'page' => '应用,分享抽奖']);
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '抽奖删除', 'page' => '应用,分享抽奖']);
        $router->post('status', ['uses' => 'StatusController@index', 'describe' => '抽奖状态', 'page' => '应用,分享抽奖']);
        $router->get('data', ['uses' => 'DataController@index', 'describe' => '抽奖活动数据', 'page' => '应用,分享抽奖']);
    });

    /**
     * 口令优惠券
     */
    $router->group(['prefix' => 'password-coupon', 'namespace' => 'PasswordCoupon'], function () use ($router) {
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '口令券添加', 'page' => '应用,口令优惠券']);
        $router->get('list', ['uses' => 'ListController@index', 'describe' => '口令券列表', 'page' => '应用,口令优惠券']);
        $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '口令券编辑', 'page' => '应用,口令优惠券']);
        $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '口令券详情', 'page' => '应用,口令优惠券']);
        $router->post('download-poster', ['uses' => 'GeneratePostersController@index', 'describe' => '口令券海报下载', 'page' => '应用,口令优惠券']);
        $router->get('download-password', ['uses' => 'DownLoadPasswordListController@index', 'describe' => '口令券下载所有口令', 'page' => '应用,口令优惠券']);
        $router->post('download-sun-code', ['uses' => 'DownLoadSunCodeController@index', 'describe' => '口令券下载太阳码', 'page' => '应用,口令优惠券']);
        $router->post('download-poster-more-status', ['uses' => 'DownLoadPosterMoreStatusController@index', 'describe' => '下载一码多券海报状态变更', 'page' => '应用,口令优惠券']);
    });

    /**
     * 助力
     */
    $router->group([
        'prefix'    => 'help',
        'namespace' => 'Help'
    ], function () use ($router) {

        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.助力.编辑',
            'page'     => '应用.助力.编辑'
        ]);

        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.助力.删除',
            'page'     => '应用.助力.删除'
        ]);

        // 上下架
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.助力.上下架',
            'page'     => '应用.助力.上下架'
        ]);
        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.助力.详情',
            'page'     => '应用.助力.详情'
        ]);
        // 添加
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.助力.添加',
            'page'     => '应用.助力.添加'
        ]);
        // 发起助力用户
        $router->post('user-list', [
            'uses'     => 'UserListController@index',
            'describe' => '应用中心.助力.发起助力用户列表',
            'page'     => '应用.助力.发起助力用户列表'
        ]);
        // 发起助力用户
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.助力.发起助力用户列表',
            'page'     => '应用.助力.发起助力用户列表'
        ]);
    });

    /**
     * 砍价
     */
    $router->group([
        'prefix'    => 'bargain',
        'namespace' => 'Bargain'
    ], function () use ($router) {

        // 列表
        $router->post('', [
            'uses'     => 'ListController@index',
            'describe' => '应用中心.砍价.砍价列表',
            'page'     => '应用.砍价.砍价列表'
        ]);

        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '应用中心.砍价.添加',
            'page'     => '应用.砍价.新增'
        ]);

        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '应用中心.砍价.编辑',
            'page'     => '应用.砍价.编辑'
        ]);

        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '应用中心.砍价.详情',
            'page'     => '应用.砍价.详情'
        ]);

        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '应用中心.砍价.删除',
            'page'     => '应用.砍价.删除'
        ]);

        // 上下架
        $router->post('enable', [
            'uses'     => 'EnableController@index',
            'describe' => '应用中心.砍价.上下架',
            'page'     => '应用.砍价.上下架'
        ]);

        // 发起助力用户
        $router->post('user-list', [
            'uses'     => 'UserListController@index',
            'describe' => '应用中心.砍价.发起砍价用户列表',
            'page'     => '应用.砍价.发起砍价用户列表'
        ]);

    });

    /**
     * 双倍券
     */
    $router->group([
        'prefix'    => 'double-coupon',
        'namespace' => 'DoubleCoupon'
    ], function () use ($router) {

        // 列表
        $router->post('putaway', [
            'uses'     => 'PutawayController@index',
            'describe' => '上下架',
            'page'     => '应用.双倍券.上下架'
        ]);

        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '添加活动',
            'page'     => '应用.双倍券.添加活动'
        ]);
        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '更新活动',
            'page'     => '应用.双倍券.更新活动'
        ]);
        // 列表
        $router->post('list', [
            'uses'     => 'ListController@index',
            'describe' => '活动列表',
            'page'     => '应用.双倍券.活动列表'
        ]);
        // 详情
        $router->post('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '活动详情',
            'page'     => '应用.双倍券.活动详情'
        ]);
        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '活动删除',
            'page'     => '应用.双倍券.活动删除'
        ]);
    });

    /**
     * 社群导流
     */
    $router->group([
        'prefix'    => 'community-diversion',
        'namespace' => 'CommunityDiversion'
    ], function () use ($router) {
        $router->get('get-config', [
            'uses'     => 'GetConfigController@index',
            'describe' => '获取社群导流配置',
            'page'     => '应用.社群导流'
        ]);
        $router->post('save', [
            'uses'     => 'SaveController@index',
            'describe' => '保存社群导流',
            'page'     => '应用.社群导流.编辑'
        ]);
    });

});

/**
 * 管理后台配置
 */
$router->group([
    'prefix'    => 'setting',
    'namespace' => 'Setting'
], function () use ($router) {
    // UC在企业微信配置的回调地址
    $router->get('callback', [
        'uses'     => 'CallbackController@index',
        'describe' => '',
        'page'     => ''
    ]);

    /**
     * 角色管理
     */
    $router->group([
        'prefix'    => 'role',
        'namespace' => 'Role'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '管理后台配置.角色管理.列表',
            'page'     => '设置.账号权限.角色'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '管理后台配置.角色管理.详情',
            'page'     => '设置.账号权限.角色.编辑'
        ]);
        $router->get('power', [
            'uses'     => 'PowerController@index',
            'describe' => '管理后台配置.角色管理.所有权限',
            'page'     => '设置.账号权限.角色.编辑'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '管理后台配置.角色管理.新增',
            'page'     => '设置.账号权限.角色.添加'
        ]);
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '管理后台配置.角色管理.编辑',
            'page'     => '设置.账号权限.角色.编辑'
        ]);
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '管理后台配置.角色管理.删除',
            'page'     => '设置.账号权限.角色'
        ]);
    });

    /**
     * 管理员管理
     */
    $router->group([
        'prefix'    => 'enterprise-admin',
        'namespace' => 'Admin'
    ], function () use ($router) {
        // 列表
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '管理后台配置.管理员管理.列表',
            'page'     => '设置.账号权限.账号'
        ]);
        // 详情
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '管理后台配置.管理员管理.详情',
            'page'     => '设置.账号权限.账号.编辑'
        ]);
        // 新增
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '管理后台配置.管理员管理.新增',
            'page'     => '设置.账号权限.账号.新增账号'
        ]);
        // 编辑
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '管理后台配置.管理员管理.编辑',
            'page'     => '设置.账号权限.账号.编辑'
        ]);
        // 状态修改
        $router->post('status', [
            'uses'     => 'UpdateController@status',
            'describe' => '管理后台配置.管理员管理.状态修改',
            'page'     => '设置.账号权限.账号.编辑'
        ]);
        // 重置密码
        $router->post('reset-password', [
            'uses'     => 'UpdateController@resetPassword',
            'describe' => '管理后台配置.管理员管理.重置密码',
            'page'     => '设置.账号权限.账号'
        ]);
        // 删除
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '管理后台配置.管理员管理.删除',
            'page'     => '设置.账号权限.账号'
        ]);
    });

    /**
     * 操作日志管理
     */
    $router->group([
        'prefix'    => 'operation-log',
        'namespace' => 'OperationLog'
    ], function () use ($router) {
        $router->get('module', [
            'uses'     => 'ModuleController@index',
            'describe' => '管理后台配置.操作日志.全部模块列表',
            'page'     => '设置.账号权限.操作日志'
        ]);
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '管理后台配置.操作日志.日志列表',
            'page'     => '设置.账号权限.操作日志'
        ]);
    });

    /**
     * 商户信息
     */
    $router->group([
        'prefix'    => 'merchant_setting',
        'namespace' => 'Merchant'
    ], function () use ($router) {
        $router->post('detail', [
            'uses'     => 'SettingDetailController@index',
            'describe' => '管理后台配置.商户信息.获取商户信息',
            'page'     => '设置.基本设置.商户信息'
        ]);
        $router->post('save', [
            'uses'     => 'SettingSaveController@index',
            'describe' => '管理后台配置.商户信息.保存商户信息',
            'page'     => '设置.基本设置.商户信息'
        ]);
    });
});

/*
 * 标签
 */
$router->group([
    'prefix'    => 'label',
    'namespace' => 'Label'
], function () use ($router) {
    /**
     * 印象标签
     */
    $router->group([
        'prefix'    => 'impress',
        'namespace' => 'Impress'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '标签管理.印象标签.列表',
            'page'     => '设置.客户.标签'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '标签管理.印象标签.添加',
            'page'     => '设置.客户.标签.新增印象标签'
        ]);
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '标签管理.印象标签.编辑',
            'page'     => '设置.客户.标签.新增印象标签'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '标签管理.印象标签.详情',
            'page'     => '设置.客户.标签.编辑印象标签'
        ]);
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '标签管理.印象标签.删除',
            'page'     => '设置.客户.标签'
        ]);
        $router->post('status', [
            'uses'     => 'StatusController@index',
            'describe' => '标签管理.印象标签.状态修改',
            'page'     => '设置.客户.标签'
        ]);
    });

    /**
     * 系统标签
     */
    $router->group([
        'prefix'    => 'system',
        'namespace' => 'System'
    ], function () use ($router) {
        $router->get('list', [
            'uses'     => 'ListController@index',
            'describe' => '标签管理.系统标签.列表',
            'page'     => '设置.客户.标签'
        ]);
        $router->post('add', [
            'uses'     => 'AddController@index',
            'describe' => '标签管理.系统标签.添加',
            'page'     => '设置.客户.标签.新增标签'
        ]);
        $router->post('update', [
            'uses'     => 'UpdateController@index',
            'describe' => '标签管理.系统标签.编辑',
            'page'     => '设置.客户.标签.编辑标签'
        ]);
        $router->get('detail', [
            'uses'     => 'DetailController@index',
            'describe' => '标签管理.系统标签.详情',
            'page'     => '设置.客户.标签.编辑标签'
        ]);
        $router->post('delete', [
            'uses'     => 'DeleteController@index',
            'describe' => '标签管理.系统标签.删除',
            'page'     => '设置.客户.标签'
        ]);
        $router->get('good', [
            'uses'     => 'GoodController@index',
            'describe' => '标签管理.系统标签.商品',
            'page'     => '设置.客户.标签'
        ]);
        $router->post('status', [
            'uses'     => 'StatusController@index',
            'describe' => '标签管理.系统标签.状态修改',
            'page'     => '设置.客户.标签'
        ]);
    });
});

/**
 * 画像统计
 */
$router->group([
    'prefix'    => 'customer-portrait',
    'namespace' => 'CustomerPortrait'
], function () use ($router) {
    $router->post('age', [
        'uses'     => 'AgeController@index',
        'describe' => '画像统计.年龄统计',
        'page'     => '设置.客户.客户画像'
    ]);
    $router->post('gender', [
        'uses'     => 'GenderController@index',
        'describe' => '画像统计.性别统计',
        'page'     => '设置.客户.客户画像'
    ]);
    $router->post('label-page', [
        'uses'     => 'LabelPageController@index',
        'describe' => '画像统计.标签分页统计',
        'page'     => '设置.客户.客户画像'
    ]);
    $router->post('label', [
        'uses'     => 'LabelController@index',
        'describe' => '画像统计.标签统计',
        'page'     => '设置.客户.客户画像'
    ]);
    $router->post('provinces-page', [
        'uses'     => 'ProvincesPageController@index',
        'describe' => '画像统计.区域分页统计',
        'page'     => '设置.客户.客户画像'
    ]);
    $router->post('provinces', [
        'uses'     => 'ProvincesController@index',
        'describe' => '画像统计.区域不分页统计',
        'page'     => '设置.客户.客户画像'
    ]);
});

/**
 * 积分
 */
$router->group(['prefix' => 'score', 'namespace' => 'Score'], function () use ($router) {

    // 积分规则保存
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '增加积分规则',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则编辑
    $router->post('update', [
        'uses'     => 'UpdateController@index',
        'describe' => '编辑积分规则',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则编辑
    $router->post('update-list', [
        'uses'     => 'RuleDetailController@index',
        'describe' => '编辑积分规则列表',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则删除
    $router->post('delete', [
        'uses'     => 'DeleteController@index',
        'describe' => '删除积分规则',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分统计
    $router->post('count', [
        'uses'     => 'CountController@index',
        'describe' => '积分统计',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分禁用开启
    $router->post('status', [
        'uses'     => 'StatusController@index',
        'describe' => '积分状态',
        'page'     => '客户管理.积分管理'
    ]);

    // 积分规则列表
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '积分规则列表',
        'page'     => '客户管理.积分管理'
    ]);

    // 有效期保存
    $router->post('save-deadline', [
        'uses'     => 'SaveDeadlineController@index',
        'describe' => '积分有效期保存',
        'page'     => '客户管理.积分管理'
    ]);

    // 有效期列表
    $router->post('deadline-list', [
        'uses'     => 'DeadlineListController@index',
        'describe' => '积分有效期列表',
        'page'     => '客户管理.积分管理'
    ]);

    $router->get('user-score-list', [
        'uses'     => 'UserScoreListController@userScoreList',
        'describe' => '客户积分明细列表',
        'page'     => '客户管理.积分管理'
    ]);

    $router->get('user-score-detail', [
        'uses'     => 'UserScoreDetailController@userScoreDetail',
        'describe' => '客户积分明细详情',
        'page'     => '客户管理.积分管理'
    ]);

    $router->post('user-score-change', [
        'uses'     => 'UserScoreChangeController@userScoreChange',
        'describe' => '管理员录入积分',
        'page'     => '客户管理.积分管理'
    ]);
});
// 热门搜索
$router->group(['prefix' => 'hot-search', 'namespace' => 'HotSearch'], function () use ($router) {
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '热词添加', 'page' => '店铺,热门搜索']);
    $router->get('list', ['uses' => 'ListController@index', 'describe' => '热词列表', 'page' => '店铺,热门搜索']);
    $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '热词编辑', 'page' => '店铺,热门搜索']);
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '热词详情', 'page' => '店铺,热门搜索']);
    $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '热词详情', 'page' => '店铺,热门搜索']);
});

/**
 * 页面模板
 */
$router->group(['prefix' => 'template', 'namespace' => 'Template'], function () use ($router) {
    // 首页
    $router->group(['prefix' => 'home', 'namespace' => 'Home'], function () use ($router) {
        // 详情
        $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '首页详情', 'page' => '首页模板']);
        // 优惠券列表（立减金、微信优惠券、站内优惠券）
        $router->get('coupon-list', ['uses' => 'CouponListController@index', 'describe' => '优惠券列表', 'page' => '首页模板']);
        // 新增
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '新增', 'page' => '首页模板']);
        // 编辑
        $router->post('edit', ['uses' => 'UpdateController@index', 'describe' => '编辑', 'page' => '首页模板']);
        // 商品列表
        $router->post('goods-list', ['uses' => 'GoodsController@index', 'describe' => '商品列表', 'page' => '首页模板']);
        // 模板列表
        $router->post('/', ['uses' => 'ListController@index', 'describe' => '模板列表', 'page' => '首页模板']);
        // 删除
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '删除', 'page' => '首页模板删除']);
    });

    // 模板风格
    $router->group(['prefix' => 'shop-style', 'namespace' => 'ShopStyle'], function () use ($router) {
        // 详情
        $router->post('info', ['uses' => 'InfoController@index', 'describe' => '风格详情', 'page' => '风格详情']);
        // 所有企业通用风格配置列表
        $router->post('config', ['uses' => 'ConfigController@index', 'describe' => '企业通用风格配置列表', 'page' => '企业通用风格配置列表']);
        // 设置
        $router->post('set', ['uses' => 'SetController@index', 'describe' => '企业风格设置', 'page' => '企业风格设置']);
        // 添加
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '企业风格添加', 'page' => '企业风格添加']);
        // 删除
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '企业风格删除', 'page' => '企业风格删除']);
    });
    /**
     * 我的模板
     */
    $router->group(['prefix' => 'myTemplate', 'namespace' => 'MyTemplate'], function () use ($router) {
        // 列表
        $router->post('list', 'ListController@index');
        // 新增
        $router->post('add', 'AddController@index');
        // 编辑
        $router->post('update', 'UpdateController@index');
        // 删除
        $router->post('delete', 'DeleteController@index');
        // 详情
        $router->post('detail', 'DetailController@index');
    });

    /**
     * 分类模板
     */
    $router->group(['prefix' => 'category', 'namespace' => 'Category'], function () use ($router) {
        // 列表
        $router->post('', 'ListController@index');
        // 新增
        $router->post('add', 'AddController@index');
        // 编辑
        $router->post('edit', 'EditController@index');
        // 删除
        $router->post('delete', 'DeleteController@index');
        // 详情
        $router->post('detail', 'DetailController@index');
    });

    //版本管理
    $router->group(['prefix' => 'version', 'namespace' => 'Version'], function () use ($router) {
        // 版本更新列表
        $router->get('list', ['uses' => 'ListController@index', 'describe' => '版本更新列表', 'page' => '版本更新列表']);

        // 版本历史记录
        $router->get('history-list', ['uses' => 'HistoryListController@index', 'describe' => '版本历史记录', 'page' => '版本历史记录']);

        // 版本记录详情
        $router->get('history-detail', ['uses' => 'HistoryDetailController@index', 'describe' => '版本记录详情', 'page' => '版本记录详情']);

        // 确认发布版本
        $router->post('release', ['uses' => 'ReleaseController@index', 'describe' => '确认发布版本', 'page' => '确认发布版本']);

        // 查看模板列表数据
        $router->get('template-list', ['uses' => 'TemplateListController@index', 'describe' => '查看模板列表数据', 'page' => '查看模板列表数据']);

        // 选择模板编辑保存
        $router->post('save', ['uses' => 'SaveController@index', 'describe' => '选择模板编辑保存', 'page' => '选择模板编辑保存']);

        // 模拟boss发布版本
        $router->get('save-status', ['uses' => 'SaveStatusController@index', 'describe' => '模拟boss审核版本', 'page' => '模拟boss审核版本']);

        // 自动化打包 提交审核
        $router->post('submit-audit', ['uses' => 'SubmitAuditController@index', 'describe' => '提交审核', 'page' => '提交审核']);

    });

    /**
     * 自动化发布
     */
    $router->group(['prefix' => 'auto-release', 'namespace' => 'AutoRelease'], function () use ($router) {
        // 撤销审核
        $router->post('cancel-audit', ['uses' => 'CancelAuditController@index', 'describe' => '撤销审核', 'page' => '版本管理-']);
        // 定时发布-撤销发布
        $router->post('cancel-release', ['uses' => 'CancelReleaseController@index', 'describe' => '撤销定时发布', 'page' => '版本管理']);
        // 预览二维码
        $router->post('qr-code', ['uses' => 'QrCodeController@index', 'describe' => '体验二维码', 'page' => '版本管理']);
        // 发布
        $router->post('release', ['uses' => 'ReleaseController@index', 'describe' => '发布小程序', 'page' => '版本管理']);
        // 版本回退
        $router->post('revert-release', ['uses' => 'RevertReleaseController@index', 'describe' => '回退小程序版本', 'page' => '版本管理']);
        // 提审
        $router->post('submit-audit', ['uses' => 'SubmitAuditController@index', 'describe' => '提交小程序审核', 'page' => '版本管理']);
        // 版本更新列表
        $router->post('update-list', ['uses' => 'UpdateListController@index', 'describe' => '提审管理详情', 'page' => '版本管理']);
        // 版本管理列表
        $router->post('version-list', ['uses' => 'VersionListController@index', 'describe' => '小程序发布列表', 'page' => '版本管理']);
    });

    /**
     * 小程序开发码
     */
    $router->group(['prefix' => 'develop', 'namespace' => 'Develop'], function () use ($router) {
        // 开发码列表
        $router->get('list', ['uses' => 'ListController@index', 'describe' => '小程序开发码列表', 'page' => '版本管理']);
        // 开发码生成 请求提交
        $router->post('submit', ['uses' => 'SubmitController@index', 'describe' => '提交小程序开发码请求', 'page' => '版本管理']);
    });

    /**
     * 页面版本
     */
    $router->group(['prefix' => 'page-version', 'namespace' => 'PageVersion'], function () use ($router) {
        $router->post('list', ['uses' => 'ListController@index', 'describe' => '版本列表', 'page' => '店铺.店铺模板.版本列表']);
        $router->post('detail', ['uses' => 'DetailController@index', 'describe' => '版本详情', 'page' => '店铺.店铺模板.版本']);
        $router->post('edit', ['uses' => 'EditController@index', 'describe' => '版本编辑', 'page' => '店铺.店铺模板.版本列表.编辑']);
        $router->post('add', ['uses' => 'AddController@index', 'describe' => '版本新增复制', 'page' => '店铺.店铺模板.版本列表']);
        $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '版本删除', 'page' => '店铺.店铺模板.版本列表']);
    });
});

/**
 * 内容管理
 */
$router->group(['prefix' => 'content-manage', 'namespace' => 'ContentManage'], function () use ($router) {
    // 保存内容
    $router->post('save', 'SaveController@index');
    // 详情
    $router->post('info', 'InfoController@index');
    // 活动列表 (解决冲突时候注意不要删掉啊)
    $router->get('activity-list', 'ActivityListController@index');
});

/**
 * 地址管理
 */
$router->group(['prefix' => 'refund-address', 'namespace' => 'RefundAddress'], function () use ($router) {
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '退货地址管理', 'page' => '退货设置.退货配置']);
    $router->post('edit', ['uses' => 'SaveController@index', 'describe' => '支付配置编辑', 'page' => '设置.支付配置编辑']);
});

/**
 * 订单状态扭转新接口
 */
$router->group(['prefix' => 'new-order', 'namespace' => 'NewOrder'], function () use ($router) {
    // 保存内容
    $router->post('update', 'ChangeOrderStateController@index');
    // 查询统计不同状态下订单的个数
    $router->post('statistics-order', 'StatisticsOrderController@index');
    // 订单详情
    $router->post('detail', 'DetailController@index');
    // 手动-补推ERP查不到的订单
    $router->post('push-order-to-erp', 'PushOrderToErpController@index');
});

/**
 * 订单配置
 */
$router->group(['prefix' => 'order-setting', 'namespace' => 'OrderSetting'], function () use ($router) {
    // 保存配置
    $router->post('save', 'SaveController@index');
    // 配置详情
    $router->post('detail', 'DetailsController@index');
});

/**
 * 售后选项接口
 */
$router->group(['prefix' => 'refund-config', 'namespace' => 'RefundConfig'], function () use ($router) {
    $router->post('detail', [
        'uses'     => 'DetailController@index',
        'describe' => '售后选项详情',
        'page'     => '设置.售后'
    ]);
    $router->post('save', [
        'uses'     => 'SaveController@index',
        'describe' => '售后选项保存',
        'page'     => '设置.售后'
    ]);
});

/**
 * 物流接口
 */
$router->group(['prefix' => 'logistics', 'namespace' => 'Logistics'], function () use ($router) {
    $router->post('query', [
        'uses'     => 'QueryController@index',
        'describe' => '批量查询物流',
        'page'     => '订单.详情'
    ]);
    $router->post('key-list', [
        'uses'     => 'KeyListController@index',
        'describe' => '获取支持的快递公司编码',
        'page'     => '订单.详情'
    ]);
});

/**
 * 配置管理
 */
$router->group(['prefix' => 'config', 'namespace' => 'Config'], function () use ($router) {
    //保存配置
    $router->post('save', 'SaveController@index');
    // 详情
    $router->post('info', 'InfoController@index');
});

/**
 * 合并请注意  已经删了两次了
 * 重构模板管理
 */
$router->group(['prefix' => 'new-template', 'namespace' => 'NewTemplate'], function () use ($router) {
    // 列表
    $router->get('list', ['uses' => 'ListController@index', 'describe' => '模板列表', 'page' => '模板列表']);

    // 新增
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '添加模板', 'page' => '添加模板']);

    // 编辑
    $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '修改模板', 'page' => '修改模板']);

    // 详情
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '模板详情', 'page' => '模板详情']);

    // 删除
    $router->post('delete', ['uses' => 'DeleteController@index', 'describe' => '删除模板', 'page' => '删除模板']);

    // 生成物料码
    $router->post('create-template-code', ['uses' => 'CreateAppCodeController@index', 'describe' => '生成物料码', 'page' => '生成物料码']);
});

/**
 * 合并请注意  已经删了两次了
 * 店铺导航
 */
$router->group(['prefix' => 'tab-bar', 'namespace' => 'TabBar'], function () use ($router) {

    // 编辑
    $router->post('update', ['uses' => 'UpdateController@index', 'describe' => '添加编辑底部按钮', 'page' => '添加编辑底部按钮']);

    // 详情
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '底部按钮详情', 'page' => '底部按钮详情']);

});

/**
 * 商品标签
 */
$router->group(['prefix' => 'goods-label', 'namespace' => 'GoodsLabel'], function () use ($router) {
    // 列表
    $router->post('list', 'ListController@index');
    // 新增
    $router->post('add', 'AddController@index');
    // 编辑
    $router->post('update', 'UpdateController@index');
    // 删除
    $router->post('delete', 'DeleteController@index');
});

/**
 * 积分兑券
 */
$router->group(['prefix' => 'score-exchange', 'namespace' => 'ScoreExchange'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ScoreCouponListController', 'describe' => '积分兑券列表', 'page' => '积分兑券']);
    // 设置积分
    $router->post('set-score', ['uses' => 'SetExchangeScoreController', 'describe' => '设置积分券所需积分', 'page' => '积分兑券']);
    // 改变状态
    $router->post('change-status', ['uses' => 'ChangeScoreCouponStatusController', 'describe' => '改变积分券状态', 'page' => '积分兑券']);
    // 选择优惠券
    $router->post('select-coupon', ['uses' => 'SelectCouponController', 'describe' => '选择优惠券', 'page' => '积分兑券']);
    // 创建积分券
    $router->post('create-score-coupon', ['uses' => 'CreateScoreCouponController', 'describe' => '创建积分券', 'page' => '积分兑券']);
});

/**
 * 员工内购券
 */
$router->group(['prefix' => 'staff-coupon', 'namespace' => 'StaffCoupon'], function () use ($router) {
    // 列表
    $router->get('list', [
        'uses'     => 'StaffCouponListController@staffCouponList',
        'page'     => '内购券列表',
        'describe' => '内购券列表',
    ]);
    // 详情
    $router->get('detail', [
        'uses'     => 'StaffCouponDetailController@staffCouponDetail',
        'page'     => '内购券页',
        'describe' => '内购券详情',
    ]);
    // 上下架
    $router->post('change-status', [
        'uses'     => 'StaffCouponChangeStatusController@staffCouponChangeStatus',
        'page'     => '内购券页',
        'describe' => '内购券上下架',
    ]);
    // 删除
    $router->post('delete', [
        'uses'     => 'StaffCouponDeleteController@staffCouponDelete',
        'page'     => '内购券页',
        'describe' => '删除内购券',
    ]);
    // 新增
    $router->post('increase', [
        'uses'     => 'StaffCouponIncreaseController@staffCouponIncrease',
        'page'     => '内购券页',
        'describe' => '创建内购券'
    ]);
    // 更新
    $router->post('update', [
        'uses'     => 'StaffCouponUpdateController@staffCouponUpdate',
        'page'     => '内购券页',
        'describe' => '更新内购券',
    ]);
    // 优惠券列表
    $router->get('wx-coupon-list', [
        'uses'     => 'SelectWxCouponListController@selectWxCouponList',
        'page'     => '内购券页',
        'describe' => '优惠券列表',
    ]);
});

/**
 * 内购商品
 */
$router->group(['prefix' => 'staff-goods', 'namespace' => 'StaffGoods'], function () use ($router) {
    // 列表
    $router->post('list', [
        'uses'     => 'StaffGoodsListController@index',
        'page'     => '内购商品列表',
        'describe' => '内购商品列表',
    ]);
    //列表（过滤空内购价）
    $router->post('filter-list', [
        'uses'     => 'StaffGoodsFilterListController@index',
        'page'     => '内购商品列表',
        'describe' => '内购商品列表',
    ]);
    // 列表
    $router->post('new-list', [
        'uses'     => 'StaffGoodsNewListController@index',
        'page'     => '内购商品新增列表',
        'describe' => '内购商品新增列表',
    ]);
    // 详情
    $router->get('detail', [
        'uses'     => 'StaffGoodsDetailController@index',
        'page'     => '内购商品列表',
        'describe' => '内购商品详情',
    ]);
    // 删除
    $router->post('delete', [
        'uses'     => 'StaffGoodsDeleteController@index',
        'page'     => '内购商品页',
        'describe' => '删除内购商品',
    ]);
    // 保存
    $router->post('increase', [
        'uses'     => 'StaffGoodsIncreaseController@index',
        'page'     => '内购商品页',
        'describe' => '创建内购商品'
    ]);
    // 更新
    $router->post('update', [
        'uses'     => 'StaffGoodsUpdateController@index',
        'page'     => '内购商品页',
        'describe' => '更新内购商品',
    ]);
});

/**
 * 视频上传
 */
$router->group(['prefix' => 'upload-video', 'namespace' => 'UploadVideo'], function () use ($router) {
    // 获取腾讯云签名
    $router->get('signature', 'SignatureController@index');
    // 提交文件基础信息
    $router->post('video-basics', 'VideoBasicsController@index');
    // 文件上传成功或失败通知
    $router->post('upload-video-status', 'UploadVideoStatusController@index');
    // 获取微软云blob上传sas
    $router->get('azure-blob-sas', 'AzureBlobUploadController@getSas');
});

/**
 * 素材管理
 */
$router->group(['prefix' => 'source', 'namespace' => 'Source'], function () use ($router) {
    // 素材列表
    $router->get('list', ['uses' => 'ListController@index', 'describe' => '素材列表', 'page' => '素材列表']);
    // 编辑素材
    $router->post('edit', ['uses' => 'EditController@index', 'describe' => '编辑素材', 'page' => '编辑素材']);
    // 添加素材
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '添加素材', 'page' => '添加素材']);
    // 素材删除
    $router->post('del', ['uses' => 'DelController@index', 'describe' => '素材删除', 'page' => '素材删除']);
    // 素材详情
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '素材详情', 'page' => '素材详情']);
});

/**
 * 搜索穹顶管理
 */
$router->group(['prefix' => 'dome', 'namespace' => 'Dome'], function () use ($router) {
    // 穹顶列表
    $router->get('list', ['uses' => 'ListController@index', 'describe' => '搜索穹顶列表', 'page' => '搜索穹顶']);
    // 编辑穹顶
    $router->post('edit', ['uses' => 'EditController@index', 'describe' => '编辑搜索穹顶', 'page' => '搜索穹顶']);
    // 添加穹顶
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '新增搜索穹顶', 'page' => '搜索穹顶']);
    // 穹顶删除
    $router->post('del', ['uses' => 'DelController@index', 'describe' => '删除搜索穹顶', 'page' => '搜索穹顶']);
    // 穹顶详情
    $router->get('detail', ['uses' => 'DetailController@index', 'describe' => '搜索穹顶详情', 'page' => '搜索穹顶']);
    // 穹顶上下架
    $router->post('status', ['uses' => 'StatusController@index', 'describe' => '搜索穹顶上下架', 'page' => '搜索穹顶']);
});

/**
 * 积分兑券
 */
$router->group(['prefix' => 'score-exchange', 'namespace' => 'ScoreExchange'], function () use ($router) {
    // 列表
    $router->post('list', ['uses' => 'ScoreCouponListController', 'describe' => '积分兑券列表', 'page' => '积分兑券']);
    // 设置积分
    $router->post('set-score', ['uses' => 'SetExchangeScoreController', 'describe' => '设置积分券所需积分', 'page' => '积分兑券']);
    // 改变状态
    $router->post('change-status', ['uses' => 'ChangeScoreCouponStatusController', 'describe' => '改变积分券状态', 'page' => '积分兑券']);
    // 选择优惠券
    $router->post('select-coupon', ['uses' => 'SelectCouponController', 'describe' => '选择优惠券', 'page' => '积分兑券']);
    // 创建积分券
    $router->post('create-score-coupon', ['uses' => 'CreateScoreCouponController', 'describe' => '创建积分券', 'page' => '积分兑券']);
});
// 现场直播
$router->group(['prefix' => 'live-broadcast', 'namespace' => 'LiveBroadcast'], function () use ($router) {
    // 直播列表
    $router->post('list', ['uses' => 'ListController@index', 'describe' => '活动', 'page' => '直播列表']);
    // 直播添加
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '活动', 'page' => '添加直播']);
    // 直播编辑
    $router->post('edit', ['uses' => 'EditController@index', 'describe' => '活动', 'page' => '编辑直播']);
    // 直播删除
    $router->post('del', ['uses' => 'DeleteController@index', 'describe' => '活动', 'page' => '删除直播']);
    // 直播详情
    $router->post('info', ['uses' => 'InfoController@index', 'describe' => '活动', 'page' => '直播详情']);
    // 生成太阳码
    $router->post('create-app-code', ['uses' => 'CreateAppCodeController@index', 'describe' => '活动', 'page' => '生成太阳码']);
    // 更改直播状态
    $router->post('status', ['uses' => 'StatusController@index', 'describe' => '活动', 'page' => '更新直播状态']);
    // 获取小程序直播列表
    $router->get('mini-program-live-list', ['uses' => 'MiNiProgramLiveListController@index', 'describe' => '活动', 'page' => '小程序直播列表']);
    // 直播设置礼券列表
    $router->post('get_coupon_list', ['uses' => 'LiveBroadcastSettingCouponController@getSettingCouponList', 'describe' => '活动', 'page' => '直播设置礼券']);
});

// 公用控制器
$router->group(['prefix' => 'common', 'namespace' => 'Common'], function () use ($router) {
    // 下载文件
    $router->get('download-file', 'DownLoadFileController@index');
    // 上传文件
    $router->post('upload', 'UploadController@upload');
    // 获取企业商城小程序与绑定小程序列表
    $router->get('app_list', 'AppListController@index');
    // 白名单
    $router->post('white-list', 'WhiteListController@index');
    //获得小程序信息
    $router->post('mini-program', 'MiniProgramController@index');
    //上传图片
    $router->post('upload-img', 'UploadImgController@index');
    //上传视频
    $router->post('upload-video', 'UploadVideoController@index');
});

// 企业信息查询
$router->group(['prefix' => 'enterprise', 'namespace' => 'Enterprise'], function () use ($router) {
    // 企业erp对接状况
    $router->post('erp-info', 'ErpInfoController@index');
});

// 邮费模版
$router->group(['prefix' => 'postage-template', 'namespace' => 'PostageTemplate'], function () use ($router) {
    // 列表
    $router->post('list', 'ListController@index');
    // 新增
    $router->post('add', 'AddController@index');
    // 编辑
    $router->post('edit', 'EditController@index');
    // 详情
    $router->post('detail', 'DetailController@index');
    // 删除
    $router->post('delete', 'DeleteController@index');
    // 模版关联商品列表
    $router->post('goods', 'GoodsController@index');
});

/**
 * 欢迎语
 */
$router->group(['prefix' => 'welcome_speech', 'namespace' => 'WelcomeSpeech'], function () use ($router) {
    // 新增欢迎语
    $router->post('add', 'AddController@index');
    // 编辑欢迎语
    $router->post('edit', 'EditController@index');
    // 欢迎语详情
    $router->get('detail', 'DetailController@index');
    // 列表
    $router->get('list', 'ListController@index');
    // 删除
    $router->post('delete', 'DeleteController@index');
});

/**
 * 联系我二维码
 */
$router->group(['prefix' => 'contact', 'namespace' => 'Contact'], function () use ($router) {
    // 联系我二维码列表
    $router->get('list', 'ListController@index');
    // 联系我二维码列表(下拉框选择)
    $router->get('selection_list', 'SelectionListController@index');
    // 员工树形列表
    $router->post('member_tree', 'ContactController@memberTree');
    // 创建二维码
    $router->post('create_qr_code', 'ContactController@createQrCode');
    // 编辑二维码
    $router->post('update_qr_code', 'ContactController@updateQrCode');
    // 查看二维码
    $router->post('select_qr_code', 'ContactController@selectQrCode');
    // 删除二维码
    $router->post('delete_qr_code', 'ContactController@deleteQrCode');
});

/**
 * 客户管理
 */
$router->group([
    'prefix'    => 'member',
    'namespace' => 'Member'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'MemberListController@memberList',
        'describe' => '客户列表',
        'page'     => '企业微信'
    ]);
    $router->post('user-del', [
        'uses'     => 'StaffWithUserDelController@staffWithUserDel',
        'describe' => '删除客户',
        'page'     => '企业微信'
    ]);
    $router->post('user-detail', [
        'uses'     => 'MemberDetailController@memberDetail',
        'describe' => '客户详情',
        'page'     => '企业微信'
    ]);
    $router->post('staff-list', [
        'uses'     => 'StaffListController@staffList',
        'describe' => '员工列表',
        'page'     => '企业微信'
    ]);
    $router->post('member-calc', [
        'uses'     => 'MemberCalcController@memberCalc',
        'describe' => '客户统计',
        'page'     => '企业微信'
    ]);
    $router->post('member-out-data', [
        'uses'     => 'OutDataController@outData',
        'describe' => '客户统计导出',
        'page'     => '企业微信'
    ]);
    $router->post('member-give-tag', [
        'uses'     => 'MemberGiveTagController@memberGiveTag',
        'describe' => '单独批量打标签',
        'page'     => '企业微信'
    ]);
});

/**
 * 标签管理
 */
$router->group(['prefix' => 'tag', 'namespace' => 'Tag'], function () use ($router) {
    // 新增标签标签组
    $router->post('add', 'AddController@index');
    // 编辑标签标签组
    $router->post('edit', 'EditController@index');
    // 删除标签标签组
    $router->post('del', 'DeleteController@index');
    // 标签组列表
    $router->post('list', 'ListController@index');
    // 标签列表-根据标签组id获取
    $router->post('tag-list', 'ListController@tagList');
    // 企业标签库列表
    $router->post('all-list', 'ListController@allList');
    // 标签列表-分页-所有
    $router->post('select-list', 'ListController@selectList');
    // 关联标签和用户分组
    $router->post('tag-customer', 'TagCustomerGroupController@index');
    // 关联标签和用户分组
    $router->post('tag-sync', 'SyncTagController@index');
});


/**
 * 企业微信红包
 */
$router->group(['prefix' => 'red_packet', 'namespace' => 'RedPacket'], function () use ($router) {
    // 红包列表
    $router->get('list', 'RedPacketController@list');
    // 红包对账列表
    $router->get('reconciliation_list', 'RedPacketController@reconciliationList');
    // 红包发送人列表
    $router->post('members', 'RedPacketController@members');
    // 红包上下架
    $router->post('sale', 'RedPacketController@sale');
    // 切换应用开关
    $router->post('switch', 'RedPacketController@switch');
    // 创建红包
    $router->post('create', 'RedPacketController@create');
    // 红包详情
    $router->post('detail', 'RedPacketController@detail');
    // 更新红包
    $router->post('update', 'RedPacketController@update');
    // 删除红包
    $router->post('delete', 'RedPacketController@delete');
});


/**
 * 德芙梦工厂
 */
$router->group(['prefix' => 'dove_dream', 'namespace' => 'DoveDream'], function () use ($router) {
    /**
     * 轮播图
     */
    $router->group(['prefix' => 'banner'], function () use ($router) {
        $router->post('list', [
            'uses'     => 'BannerController@detail',
            'describe' => '轮播图详情',
            'page'     => '页面管理'
        ]);
        $router->post('save', [
            'uses'     => 'BannerController@save',
            'describe' => '轮播图保存',
            'page'     => '页面管理'
        ]);
    });
    /**
     * 基本信息
     */
    $router->group(['prefix' => 'basic'], function () use ($router) {
        $router->post('detail', [
            'uses'     => 'BasicController@detail',
            'describe' => '基本信息详情',
            'page'     => '页面管理'
        ]);
        $router->post('save', [
            'uses'     => 'BasicController@save',
            'describe' => '基本信息保存',
            'page'     => '页面管理'
        ]);
    });
    /**
     * 商品管理
     */
    $router->group(['prefix' => 'spu'], function () use ($router) {
        $router->post('list', [
            'uses'     => 'SpuController@list',
            'describe' => '商品列表',
            'page'     => '商品管理'
        ]);
        $router->post('detail', [
            'uses'     => 'SpuController@detail',
            'describe' => '商品详情',
            'page'     => '商品管理'
        ]);
        $router->post('save', [
            'uses'     => 'SpuController@save',
            'describe' => '商品编辑',
            'page'     => '商品管理'
        ]);
    });
    /**
     * 库存(sku)管理
     */
    $router->group(['prefix' => 'sku'], function () use ($router) {
        $router->post('list', [
            'uses'     => 'SkuController@list',
            'describe' => '库存详情',
            'page'     => '库存管理'
        ]);
        $router->post('save', [
            'uses'     => 'SkuController@save',
            'describe' => '库存编辑',
            'page'     => '库存管理'
        ]);
    });

    /**
     * 德芙梦工厂 - 参加者相关接口
     */
    $router->group(['prefix' => 'order_member'], function () use ($router) {
        // 参加者-详情
        $router->get('detail', [
            'uses'     => 'OrderMemberController@detail',
            'describe' => '梦工厂.参加者详情',
            'page'     => '梦工厂.参加者信息/核销管理'
        ]);
        // 参加者-编辑
        $router->post('save', [
            'uses'     => 'OrderMemberController@save',
            'describe' => '梦工厂.参加者编辑',
            'page'     => '梦工厂.参加者信息/核销管理'
        ]);
        // 参加者-列表+搜索
        $router->get('list', [
            'uses'     => 'OrderMemberController@list',
            'describe' => '梦工厂.参加者列表/搜索',
            'page'     => '梦工厂.参加者信息/核销管理'
        ]);
        // 参加者-导出
        $router->get('export', [
            'uses'     => 'OrderMemberController@export',
            'describe' => '梦工厂.参加者导出',
            'page'     => '梦工厂.参加者信息/核销管理'
        ]);
        // 订单-导出
        $router->post('write', [
            'uses'     => 'OrderMemberController@write',
            'describe' => '梦工厂.参加者核销',
            'page'     => '梦工厂.参加者信息/核销管理'
        ]);
    });

    /**
     * 德芙梦工厂 - 订单相关接口
     */
    $router->group(['prefix' => 'order'], function () use ($router) {
        // 订单-列表+搜索
        $router->get('list', [
            'uses'     => 'OrderController@list',
            'describe' => '梦工厂.订单列表/搜索',
            'page'     => '梦工厂.订单管理'
        ]);
        // 订单-导出
        $router->get('export', [
            'uses'     => 'OrderController@export',
            'describe' => '梦工厂.订单导出',
            'page'     => '梦工厂.订单管理'
        ]);
        // 订单-详情
        $router->get('detail', [
            'uses'     => 'OrderController@detail',
            'describe' => '梦工厂.订单详情',
            'page'     => '梦工厂.订单管理'
        ]);
        // 订单-审核
        $router->get('approval', [
            'uses'     => 'OrderController@approval',
            'describe' => '梦工厂.退款审核',
            'page'     => '梦工厂.订单管理'
        ]);
    });
});


/**
 * 企业微信-客户维护 Customer maintenance
 */
$router->group([
    'prefix'    => 'customer-maintenance',
    'namespace' => 'CustomerMaintenance'
], function () use ($router) {
    /**
     * 企业群发
     */
    $router->group(['prefix' => 'msg', 'namespace' => 'Msg'], function () use ($router) {
        //--添加消息群发--
        $router->post('add', ['uses' => 'AddController@index']);
        //--群发列表----
        $router->post('list', ['uses' => 'ListController@index']);
        //--推送对象全部下拉获取
        $router->get('getObject', ['uses' => 'ListController@getObjectSelect']);
        //--推送对象下拉获取  type:1所属员工|2渠道|3标签
        $router->get('getObject/{type}', ['uses' => 'ListController@getObjectSelect']);
        //--查看h5图文页面 - 素材和消息通用
        $router->get('article/{id}', ['uses' => 'ArticleController@index']);
        //--获取群发结果详情列表
        $router->post('getInfo/{tsid}', ['uses' => 'ListController@getInfo']);
        //--获取群发消息详情用于编辑
        $router->get('getEditInfo/{tsid}', ['uses' => 'ListController@getEditInfo']);
        //--获取预计触达数量
        $router->post('getEst', ['uses' => 'ListController@getEst']);
        //--更新企业群发--
        $router->post('update/{tsid}', ['uses' => 'UpdateController@index']);
        //--更新群发状态和实际触发人数--
        $router->post('updateStatus/{tsid}', ['uses' => 'UpdateStatusController@index']);
        //--删除推送消息-仅草稿状态下
        $router->get('delete/{tsid}', ['uses' => 'DeleteController@index']);
    });
    /**
     * 素材管理
     */
    $router->group(['prefix' => 'media', 'namespace' => 'Media'], function () use ($router) {
        //--用于生成media_id的临时素材
        $router->post('mediaIMG', ['uses' => 'UploadController@upMediaIMG']);
        //--用于富文本编辑的永久素材
        $router->post('upIMG', ['uses' => 'UploadController@upIMG']);
        //--获取商城小程序卡片的默认appid
        $router->get('getShopAppid', ['uses' => 'ListController@getShopAppid']);
        //--添加素材
        $router->post('add', ['uses' => 'AddController@index']);
        //--素材列表 - 所有素材
        $router->post('list', ['uses' => 'ListController@index']);
        //--素材列表 - 指定素材 type:0文本素材 1图片素材|2图文链接|3小程序卡片
        $router->post('list/{type}', ['uses' => 'ListController@index']);
        //--更新素材--
        $router->post('update/{scid}', ['uses' => 'UpdateController@index']);
        //--删除素材
        $router->get('delete/{scid}', ['uses' => 'DeleteController@index']);
    });
    /**
     * 群管理
     */
    $router->group(['prefix' => 'cluster', 'namespace' => 'Cluster'], function () use ($router) {

        $router->post('list', ['uses' => 'ListController@index']);

    });
});


/**
 * 企业微信
 */
$router->group(['prefix' => 'qywx', 'namespace' => 'Qywx'], function () use ($router) {
    /**
     * 渠道管理
     */
    $router->group(['prefix' => 'channel', 'namespace' => 'Channel'], function () use ($router) {
        // 渠道列表
        $router->post('list', 'ListController@index');
        // 新增渠道
        $router->post('add', 'AddController@index');
        // 编辑渠道
        $router->post('edit', 'EditController@index');
        // 渠道统计列表
        $router->post('statistics-list', 'StatisticsListController@index');
        // 渠道统计明细
        $router->post('statistics-detail', 'StatisticsDetailController@index');
    });

    /**
     * 用户分组
     */
    $router->group(['prefix' => 'customer-group', 'namespace' => 'CustomerGroup'], function () use ($router) {
        // 渠道列表
        $router->post('list', 'ListController@index');
        // 新增用户分组
        $router->post('add', 'AddController@index');
        // 用户分组详情
        $router->post('detail', 'DetailController@index');
        // 编辑用户分组
        $router->post('edit', 'EditController@index');
        // 删除用户分组
        $router->post('del', 'DeleteController@index');
    });

    /**
     * 标签管理
     */
    $router->group(['prefix' => 'tag', 'namespace' => 'Tag'], function () use ($router) {
        // 新增标签标签组
        $router->post('add', 'AddController@index');
        // 编辑标签标签组
        $router->post('edit', 'EditController@index');
        // 删除标签标签组
        $router->post('del', 'DeleteController@index');
        // 标签组列表
        $router->post('list', 'ListController@index');
        // 标签列表-根据标签组id获取
        $router->post('tag-list', 'ListController@tagList');
        // 企业标签库列表
        $router->post('all-list', 'ListController@allList');
        // 标签列表-分页-所有
        $router->post('select-list', 'ListController@selectList');
        // 关联标签和用户分组
        $router->post('tag-customer', 'TagCustomerGroupController@index');
        // 关联标签和用户分组
        $router->post('tag-sync', 'SyncTagController@index');
    });
});


/**
 * 企业微信红包
 */
$router->group(['prefix' => 'red_packet', 'namespace' => 'RedPacket'], function () use ($router) {
    // 红包列表
    $router->get('list', 'RedPacketController@list');
    // 红包对账列表
    $router->get('reconciliation_list', 'RedPacketController@reconciliationList');
    // 红包发送人列表
    $router->post('members', 'RedPacketController@members');
    // 红包上下架
    $router->post('sale', 'RedPacketController@sale');
    // 切换应用开关
    $router->post('switch', 'RedPacketController@switch');
    // 创建红包
    $router->post('create', 'RedPacketController@create');
    // 红包详情
    $router->post('detail', 'RedPacketController@detail');
    // 更新红包
    $router->post('update', 'RedPacketController@update');
    // 删除红包
    $router->post('delete', 'RedPacketController@delete');
});

/**
 * 员工内购
 */
$router->group(['prefix' => 'staff-buy', 'namespace' => 'StaffBuyManager'], function () use ($router) {

    // 内购员工列表获取
    $router->post('list', ['uses' => 'ListController@index', 'describe' => '员工亲友内购列表获取', 'page' => '组织']);

    // 内购员工设置获取
    $router->post('setting', ['uses' => 'SettingController@index', 'describe' => '员工亲友内购设置获取', 'page' => '组织']);

    // 内购员工添加
    $router->post('add', ['uses' => 'AddController@index', 'describe' => '员工亲友内购添加', 'page' => '组织']);

    // 内购员工修改
    $router->post('save', ['uses' => 'UpdateController@index', 'describe' => '员工亲友内购修改', 'page' => '组织']);

    // 运营设置数据添加
    $router->post('fill', ['uses' => 'FillController@index', 'describe' => '员工亲友内购运营设置数据添加', 'page' => '组织']);

    // 运营设置数据详情
    $router->post('ext-detail', ['uses' => 'StaffExtDetailController@index', 'describe' => '员工亲友内购运营设置数据详情', 'page' => '组织']);
});


/**
 * 内购订单
 */
$router->group([
    'prefix'    => 'staff-order',
    'namespace' => 'StaffOrder'
], function () use ($router) {
    $router->post('list', [
        'uses'     => 'ListController@index',
        'describe' => '内购订单.列表',
        'page'     => '内购订单.列表'
    ]);
    $router->post('statistics', [
        'uses'     => 'StatisticsController@index',
        'describe' => '内购订单.统计数据',
        'page'     => '内购订单.统计数据'
    ]);
    $router->post('export-order-list', [
        'uses'     => 'ExportOrderListController@index',
        'describe' => '创建内购订单.列表导出任务',
        'page'     => '创建内购订单.列表导出任务'
    ]);

    $router->post('achievement-list', [
        'uses'     => 'AchievementListController@index',
        'describe' => '员工业绩.列表',
        'page'     => '员工业绩.列表'
    ]);

    $router->post('export-achievement-list', [
        'uses'     => 'ExportAchievementListController@index',
        'describe' => '创建内购业绩列表.导出任务',
        'page'     => '创建内购业绩列表.导出任务'
    ]);

    $router->post('department-achievement-list', [
        'uses'     => 'AchievementDpListController@index',
        'describe' => '部门业绩.列表',
        'page'     => '部门业绩.列表'
    ]);

    $router->post('export-department-achievement-list', [
        'uses'     => 'ExportAchievementDpListController@index',
        'describe' => '创建内购部门业绩列表.导出任务',
        'page'     => '创建内购部门业绩列表.导出任务'
    ]);
});

$router->group([
    'prefix'    => 'staff-right',
    'namespace' => 'StaffRight'
], function () use ($router) {
    $router->post('staff-bind-list', [
        'uses'     => 'StaffBindListController@index',
        'describe' => '导购列表.获取导购绑定顾客信息',
        'page'     => '导购列表.获取导购绑定顾客信息'
    ]);
    $router->post('partial-inherit', [
        'uses'     => 'PartialInheritController@index',
        'describe' => '导购列表导购业绩关系部分继承',
        'page'     => '导购列表.导购业绩关系继承'
    ]);
    $router->post('all-inherit', [
        'uses'     => 'AllInheritController@index',
        'describe' => '导购列表导购业绩关系全部继承',
        'page'     => '导购列表.导购业绩关系继承'
    ]);
});

/**
 * 物料
 */
$router->group([
    'prefix'    => 'material',
    'namespace' => 'Material'
], function () use ($router) {
    // 新增物料接口（暂时不用）
    $router->post('add', 'AddController@index');

    // 获取企业物料列表
    $router->post('list', 'ListController@index');

    // 编辑物料
    $router->post('edit', 'EditController@index');

    // 物料详情
    $router->post('detail', 'DetailController@index');

    // 物料库存修改
    $router->post('update-stock', 'UpdateStockController@index');
});


/**
 * 素材中心
 */
$router->group(['prefix' => 'material-center', 'namespace' => 'MaterialCenter'], function () use ($router) {
    // 标签列表
    $router->post('tag-list', [
        'uses'     => 'TagListController@index',
        'describe' => '标签列表',
        'page'     => '素材中心'
    ]);
    // 编辑标签
    $router->post('tag-edit', [
        'uses'     => 'TagEditController@index',
        'describe' => '编辑标签',
        'page'     => '素材中心'
    ]);
    // 素材标签列表
    $router->post('material-tag-list', [
        'uses'     => 'MaterialTagListController@index',
        'describe' => '素材标签列表',
        'page'     => '素材中心'
    ]);
    // 回收站素材列表
    $router->post('recycle-bin-list', [
        'uses'     => 'RecycleBinListController@index',
        'describe' => '回收站素材列表',
        'page'     => '素材中心'
    ]);
    // 从素材中心删除素材
    $router->post('del-from-center', [
        'uses'     => 'MaterialDelFromCenterController@index',
        'describe' => '从素材中心删除素材',
        'page'     => '素材中心'
    ]);
    // 从回收站彻底删除素材
    $router->post('del-from-bin', [
        'uses'     => 'MaterialDelFromBinController@index',
        'describe' => '从回收站彻底删除素材',
        'page'     => '素材中心'
    ]);
    // 从回收站恢复素材
    $router->post('recover-from-bin', [
        'uses'     => 'MaterialRecoverFromBinController@index',
        'describe' => '从回收站恢复素材',
        'page'     => '素材中心'
    ]);
    // 回收站恢复判断原路径是否完整
    $router->post('recover-from-bin-check-parent-path', [
        'uses'     => 'MaterialCheckParentPathController@index',
        'describe' => '回收站恢复判断原路径是否完整',
        'page'     => '素材中心'
    ]);
    // 素材类型下所有文件夹列表
    $router->post('folder-list-tree', [
        'uses'     => 'MaterialFolderListTreeController@index',
        'describe' => '素材类型下所有文件夹列表',
        'page'     => '素材中心'
    ]);
    // 素材变更父级文件夹
    // 1、回收站恢复时存在父级文件夹被变更或被删除
    $router->post('exchange-parent-folder', [
        'uses'     => 'MaterialExchangeParentFolderController@index',
        'describe' => '素材变更父级文件夹',
        'page'     => '素材中心'
    ]);
    // 判断素材是否存在(综合素材、商品素材列表)
    $router->post('check-material-valid', [
        'uses'     => 'MaterialCheckValidController@index',
        'describe' => '判断素材是否存在(综合素材、商品素材列表)',
        'page'     => '素材中心'
    ]);
    // 素材列表接口
    $router->post('material-list', [
        'uses'     => 'MaterialListController@materialList',
        'describe' => '素材列表',
        'page'     => '素材中心'
    ]);
    // 新建文件夹
    $router->post('create-dir', [
        'uses'     => 'CreateDirController@createDir',
        'describe' => '创建文件夹',
        'page'     => '素材中心'
    ]);
    // 编辑文件夹名称
    $router->post('edit-dir-name', [
        'uses'     => 'EditDirNameController@editDirName',
        'describe' => '编辑文件夹名字',
        'page'     => '素材中心'
    ]);
    // 文件夹列表
    $router->post('dir-list', [
        'uses'     => 'DirListController@dirList',
        'describe' => '文件夹列表接口',
        'page'     => '素材中心'
    ]);
    // 下载文件
    $router->post('download-file', [
        'uses'     => 'DownLoadFileController@downloadFile',
        'describe' => '下载文件',
        'page'     => '素材中心'
    ]);
    // 上传文件
    $router->post('upload-file', [
        'uses'     => 'UploadFileController@uploadFile',
        'describe' => '上传文件',
        'page'     => '素材中心'
    ]);
    // 腾讯临时秘钥
    $router->post('tmp-secret', [
        'uses'     => 'TmpSecretController@tmpSecret',
        'describe' => '获取临时秘钥',
        'page'     => '素材中心'
    ]);
    // 商品上传文件
    $router->post('goods-file', [
        'uses'     => 'GoodsFileController@goodsFile',
        'describe' => '商品上传文件',
        'page'     => '素材中心'
    ]);
    $router->post('goods-material-list', [
        'uses'     => 'GoodsMaterialListController@goodsMaterialList',
        'describe' => '商品素材中心',
        'page'     => '素材中心'
    ]);
});

// 团购活动迁移
$router->group(['prefix' => 'activity-common', 'namespace' => 'ActivityCommon'], function () use ($router) {
    $router->post('activity', [
        'uses'     => 'ActivityController@index',
        'describe' => '活动通用(列表/新增/删除/详情/编辑/上架/下架)',
        'page'     => '应用'
    ]);
    $router->post('export', [
        'uses'     => 'ExportController@index',
        'describe' => '活动列表导出',
        'page'     => '应用'
    ]);
    $router->post('increase-stock', [
        'uses'     => 'IncreaseStockController@index',
        'describe' => '活动库存增加',
        'page'     => '应用'
    ]);
    $router->post('generate-sun-code', [
        'uses'     => 'GenerateSunCodeController@index',
        'describe' => '生成活动太阳码',
        'page'     => '应用'
    ]);
    $router->get('download-sun-code', [
        'uses'     => 'DownloadSunCodeController@index',
        'describe' => '下载活动太阳码',
        'page'     => '应用'
    ]);
    $router->post('activity-list', [
        'uses'     => 'ActivityListController@index',
        'describe' => '获取活动列表接口',
        'page'     => '应用'
    ]);
    $router->post('activity-spu-list', [
        'uses'     => 'ActivitySpuListController@index',
        'describe' => '获取活动商品列表接口',
        'page'     => '应用'
    ]);
    $router->post('activity-type-list', [
        'uses'     => 'ActivityTypeListController@index',
        'describe' => '获取活动列表接口',
        'page'     => '应用'
    ]);
    $router->post('activity-sku', [
        'uses'     => 'ActivitySkuController@index',
        'describe' => '获取活动占用sku列表接口',
        'page'     => '应用'
    ]);
});
$router->group(['prefix' => 'withdrawal-setting', 'namespace' => 'WithdrawaSetting'], function () use ($router) {
    //提现设置
    $router->post('add', ['uses' => 'AddController@index']);
    //提现设置查看
    $router->post('detail', ['uses' => 'DetailController@index']);
    //自动提现设置
    $router->post('auto_add', ['uses' => 'AutoAddController@index']);
    //自动提现设置查看
    $router->post('auto_detail', ['uses' => 'AutoDetailController@index']);
    //根据项目编码获取项目信息
    $router->post('check', ['uses' => 'CheckController@index']);
});

$router->group(['prefix' => 'withdrawal', 'namespace' => 'Withdrawal'], function () use ($router) {

    //审批列表
    $router->post('list', ['uses' => 'ApprovalListController@index']);
    //审批列表
    $router->post('export_list', ['uses' => 'ListExportController@index']);
    //审批
    $router->post('approval-check', ['uses' => 'ApprovalCheckController@index']);
});

$router->group(['prefix' => 'partner', 'namespace' => 'Partner'], function () use ($router) {
    // 合伙人列表
    $router->post('list', ['uses' => 'listController@index']);
    // 合伙人详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 合伙人启用禁用
    $router->post('switch', ['uses' => 'SwitchController@index']);
    // 合伙人导出
    $router->post('export_list', ['uses' => 'exportController@index']);
    // 合伙人余额明细导出
    $router->post('export_detail', ['uses' => 'exportDetail@index']);
});

$router->group(['prefix' => 'partner-distribution', 'namespace' => 'PartnerSettlement'], function () use ($router) {
    // 结算记录
    $router->post('list', ['uses' => 'listController@index']);
    // 结算记录详情
    $router->post('detail', ['uses' => 'DetailController@index']);
    // 结算记录导出
    $router->post('export', ['uses' => 'exporListController@index']);
});

$router->group(['prefix' => 'marketing-plan', 'namespace' => 'MarketingPlan'], function () use ($router) {
    //营销计划列表
    $router->post('list', ['uses' => 'ListController@index']);
    //详情/复制
    $router->post('detail', ['uses' => 'DetailController@index']);
    //创建
    $router->post('create', ['uses' => 'CreateController@index']);
    //删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
    //未完成名单导出
    $router->post('export', ['uses' => 'ExportController@index']);
    //更新
    $router->post('update', ['uses' => 'UpdateController@index']);
    //获得太阳码
    $router->post('sun-code', ['uses' => 'GetSunCodeController@index']);
});

$router->group(['prefix' => 'mission_center', 'namespace' => 'MissionCenter'], function () use ($router) {
    //营销计划列表
    $router->post('list', ['uses' => 'ListController@index']);
    //详情/复制
    $router->post('detail', ['uses' => 'DetailController@index']);
    //创建
    $router->post('create', ['uses' => 'CreateController@index']);
    //删除
    $router->post('delete', ['uses' => 'DeleteController@index']);
    //未完成名单导出
    $router->post('export', ['uses' => 'ExportController@index']);
    //更新
    $router->post('update', ['uses' => 'UpdateController@index']);
    //获得太阳码
    $router->post('sun-code', ['uses' => 'GetSunCodeController@index']);

});

$router->group(['prefix' => 'data_compass', 'namespace' => 'DataCompass'], function () use ($router) {
    // 顾客数据按导购导出
    $router->post('export_by_worker', ['uses' => 'DataExportByWorkerController@index']);
    // 顾客数据按日期导出
    $router->post('export_by_date', ['uses' => 'DataExportByDateController@index']);
    // 顾客数据按日期导出
    $router->post('export_binding_customer_detail', ['uses' => 'ExportBindingCustomerDetailController@index']);
});


$router->group(['prefix' => 'data-lake', 'namespace' => 'DataLake'], function () use ($router) {
    // 添加配置
    $router->post('add-config', ['uses' => 'AddConfigController@index']);
    // 修改配置
    $router->post('edit-config', ['uses' => 'EditConfigController@index']);

    // 修改配置
    $router->post('get-config', ['uses' => 'GetConfigController@index']);

});

// 标签
$router->group(['prefix' => 'tag', 'namespace' => 'Tag'], function () use ($router) {
    // 人群包导出
    $router->get('export-crowd-template', ['uses' => 'TagController@exportCrowdTemplate']);
});