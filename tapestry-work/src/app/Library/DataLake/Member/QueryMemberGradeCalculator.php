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

class QueryMemberGradeCalculator extends Base
{
    public $errStr = 'query member grade calculator error';
    public $api = '/api/interface/third/member/queryMemberGradeCalculator';

    private $queryType = "";
    private $value = "";
    //当前等级
    private $grade = "";
    //要升级到的等级
    private $upGrade = "";


    function __construct(string $queryType, string $value, string $grade, string $upGrade)
    {
        $this->queryType = $queryType;
        $this->value     = $value;
        $this->grade     = $grade;
        $this->upGrade   = $upGrade;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"   => $this->config->getBrandCode(),
            "program_code" => $this->config->getProgramCode(),
            "queryType"    => $this->queryType,
            "value"        => $this->value,
            "grade"        => $this->grade,
            "upGrade"      => $this->upGrade
        ];

    }
}