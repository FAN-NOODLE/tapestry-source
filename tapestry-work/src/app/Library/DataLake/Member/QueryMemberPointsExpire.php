<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/8/14 0014
 * Time: 17:11
 */

namespace App\Library\DataLake\Member;

use App\Library\DataLake\Base;
use App\Library\DataLake\Config;

class QueryMemberPointsExpire extends Base
{
    public $errStr = 'query member PointsExpire error';
    public $api = '/api/interface/third/member/queryMemberPointsExpire';

    private $queryType = "";
    private $value = "";
    private $pointTypeGroup;
    private $Startdate;
    private $Enddate;


    public function __construct(string $queryType, string $value, string $pointTG, string $start, string $end)
    {
        $this->queryType      = $queryType;
        $this->value          = $value;
        $this->pointTypeGroup = $pointTG;
        $this->Startdate      = $start;
        $this->Enddate        = $end;
        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            "brand_code"     => $this->config->getBrandCode(),
            "program_code"   => $this->config->getProgramCode(),
            "queryType"      => $this->queryType,
            "value"          => $this->value,
            "pointTypeGroup" => $this->pointTypeGroup,
            "Startdate"      => $this->Startdate,
            "Enddate"        => $this->Enddate
        ];

    }

}