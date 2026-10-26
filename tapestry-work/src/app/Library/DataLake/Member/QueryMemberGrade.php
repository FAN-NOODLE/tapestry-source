<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/14 0014
 * Time: 16:58
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;
use App\Library\DataLake\Config;

class QueryMemberGrade extends Base
{
    public $errStr = 'query member grade error';
    public $api = '/api/interface/third/member/queryMemberGrade';

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
            "brand_code"   => $this->config->getBrandCode(),
            "program_code" => $this->config->getProgramCode(),
            "queryType"    => $this->queryType,
            "value"        => $this->value,
            "category"     => "default"
        ];

    }
}