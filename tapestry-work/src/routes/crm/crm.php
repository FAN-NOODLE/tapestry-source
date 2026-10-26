<?php

/**
 * CRM-CDP
 */
/** @var $router */
$router->group([
], function () use ($router) {

    $router->post('getToken', [
        'uses' => 'CrmController@getToken',
        'desc' => '获取token'
    ]);
    $router->post('customer-info', [
        'uses' => 'CrmController@getCustomerInfo',
        'desc' => '通用获取第三方会员信息'
    ]);
    $router->post('get-member-info', [
        'uses' => 'CrmController@getMemberInfo',
        'desc' => '获取第三方会员信息'
    ]);
    $router->post('get-point-history',[
        'uses' => 'CrmController@getCustomerPointHistory',
        'desc' => '查询会员积分记录'
    ]);
    $router->post('get-customer-coupon',[
        'uses' => 'CrmController@getCustomerCoupon',
        'desc' => '查询会员优惠卷记录'
    ]);
    $router->post('get-sw-customer-coupon',[
        'uses' => 'CrmController@getSwCustomerCoupon',
        'desc' => '查询Sw会员优惠卷记录'
    ]);
    $router->post('get-customer-purchase',[
        'uses' => 'CrmController@getCustomerPurchase',
        'desc' => '查询会员购买记录'
    ]);
    $router->post('get-customer-tag',[
        'uses' => 'CrmController@getCustomerTag',
        'desc' => '查询会员标签记录'
    ]);
    $router->post('check-member',[
        'uses' => 'CrmController@checkMember',
        'desc' => '用户是否是会员'
    ]);
    $router->post('welfare-give',[
        'uses' => 'CrmController@actionWelfareGive',
        'desc' => 'coach发放优惠卷'
    ]);
    $router->post('welfare-haved',[
        'uses' => 'CrmController@actionWelfareHaved',
        'desc' => 'coach用户是否发放优惠卷'
    ]);

    $router->post('user-join',[
        'uses' => 'CrmController@actionUserJoin',
        'desc' => 'ks用户是否发放优惠卷'
    ]);

    $router->post('sa-assignment',[
        'uses' => 'CrmController@saAssignment',
        'desc' => 'coach导购添加顾客'
    ]);
    $router->post('sa-un-assignment',[
        'uses' => 'CrmController@saUnAssignment',
        'desc' => 'coach导购删除顾客'
    ]);

    $router->post('customer-option',[
        'uses' => 'CrmController@customerOption',
        'desc' => 'coach聚合查询会员列表筛选条件'
    ]);

    $router->post('customer-list',[
        'uses' => 'CrmController@customerList',
        'desc' => 'coach根据门店或门店导购批量查询会员列表'
    ]);

    $router->post('customer-detail',[
        'uses' => 'CrmController@customerDetail',
        'desc' => 'coach查询会员详情'
    ]);

    $router->post('customer-detail-by-coach-crm',[
        'uses' => 'CrmController@getCoachCustomerInfo',
        'desc' => 'coach查询会员详情'
    ]);
});

