<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/13 0013
 * Time: 10:05
 */

namespace App\Library\DataLake\Member;


use App\Library\DataLake\Base;

class QueryMemberInfo extends Base
{
    public $errStr = 'query member info error';
    public $api = '/api/interface/third/member/queryMemberInfo';

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
            "Brand_code"   => $this->config->getBrandCode(),
            "Program_code" => $this->config->getProgramCode(),
            "queryType"    => $this->queryType,
            "Value"        => $this->value
        ];

    }
}