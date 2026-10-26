<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/3/25 0025
 * Time: 17:28
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;
use App\Library\DataLake\Config;

class QueryOrderDetail extends Base
{
    public $errStr = 'query order detail error';

    public $api = '/api/interface/third/member/queryMemberOrderDetailList';


    private $queryType = "";
    private $value = "";
    private $tradeNo;

    public function __construct(string $queryType,string $value,string $tradeNo)
    {
        $this->queryType = $queryType;
        $this->value = $value;
        $this->tradeNo = $tradeNo;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"   => (string)$this->config->getBrandCode(),
            "program_code" => $this->config->getProgramCode(),
            "queryType"    => $this->queryType,
            "value"         => $this->value,
            "tradeNo"       => $this->tradeNo,
        ];

    }
}