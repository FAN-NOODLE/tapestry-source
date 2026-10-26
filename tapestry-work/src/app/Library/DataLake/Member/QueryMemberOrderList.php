<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/9/2 0002
 * Time: 9:45
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;
use App\Library\DataLake\Config;

class QueryMemberOrderList extends Base
{
    public $errStr = 'query member order list error';

    public $api = '/api/interface/third/member/queryMemberOrderList';

    private $queryType = "";

    private $value = "";

    private $limit = 10;

    private $pageNo = 1;
    private $tradeNo = '';


    public function __construct(string $queryType, string $value, $pageNo = 1, $limit = 10, $tradeNo)
    {
        $this->queryType = $queryType;
        $this->value     = $value;
        $this->pageNo    = $pageNo;
        $this->limit     = $limit;
        $this->limit     = $limit;
        $this->tradeNo   = $tradeNo;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"   => $this->config->getBrandCode(),
            "program_code" => $this->config->getProgramCode(),
            "storeCode"    => "",
            "queryType"    => $this->queryType,
            "value"        => $this->value,
            "orderType"    => 1,
            "tradeNo"      => $this->tradeNo,
            "limit"        => (string)$this->limit,
            "pageNo"       => (string)$this->pageNo
        ];

    }

}