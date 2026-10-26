<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/17 0017
 * Time: 17:11
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;
use App\Library\DataLake\Config;

class QueryMemberPoints extends Base
{
    public $errStr = 'query member point error';

    public $api = '/api/interface/third/member/queryMemberPoints';

    private $queryType = "";
    private $value = "";


    function __construct(string $queryType,string $value)
    {
        $this->queryType = $queryType;
        $this->value = $value;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"     => $this->config->getBrandCode(),
            "program_code"   => $this->config->getProgramCode(),
            "queryType"      => $this->queryType,
            "value"          => $this->value,
            "pointTypeGroup" => "BP" // BP 累计积分 RP 累计兑换（扣减）积分
        ];
    }

}