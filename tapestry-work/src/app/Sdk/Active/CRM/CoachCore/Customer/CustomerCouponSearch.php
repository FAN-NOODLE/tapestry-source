<?php

namespace App\Sdk\Active\CRM\CoachCore\Customer;

use App\Sdk\Active\CRM\CoachCore\Base;

/**
 * Customer Coupon Search
 * 根据customerInfo查询会员优惠卷记录
 */
class CustomerCouponSearch extends Base
{
    public $errStr = 'query member info error';
    public $api = '/v2/cdp/coupon/queryCustomerCoupons';


    private $pageNum;
    private $pageSize;
    private $customerInfo;
    private $status;        // 优惠券状态1未使用 2已使用 3过期 0 all
    private $channel_type;  // 渠道 RETAIL: 零售店 OUTLET：奥莱店
    private $source;        // 请求来源： 1会员中心 2 零售小程序3 奥莱小程序 4 企业微信

    public function __construct(array $customerInfo, int $status, string $channel_type, int $source, int $page = 1, int $size = 10)
    {
        parent::__construct();
        $this->customerInfo = $customerInfo;
        $this->status = $status;
        $this->channel_type = $channel_type;
        $this->source = $source;
        $this->pageNum = $page;
        $this->pageSize = $size;
        $this->setQueryParam();

    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'pageNum' => $this->pageNum,
            'pageSize' => $this->pageSize,
            'customerInfo' => $this->customerInfo,
            'status' => $this->status,
            'channel_type' => $this->channel_type,
            'source' => $this->source,
        ];

    }
}