<?php


namespace App\Library\DataLake\BaBind;

use App\Library\DataLake\Base;

/**
 * Class BaBind
 * @package App\Library\DataLake\BaBind
 *
 * 企业微信BA绑定解绑
 */
class BaBind extends Base
{
    protected $throttle = true;

    public $api = "/api/interface/third/member/saveWechatBaBind";

    private $cusUnionID;
    private $cusUserID;
    private $baUserID;
    private $bindStatus;
    private $bindType;
    private $opDate;
    private $baCode;
    private $memberCode;

    public function __construct(string $cusUnionID, string $cusUserID, string $baUserID, string $bindStatus, string $bindType, string $opDate, string $baCode, string $memberCode)
    {
        $this->cusUnionID = $cusUnionID;
        $this->cusUserID = $cusUserID;
        $this->baUserID = $baUserID;
        $this->bindStatus = $bindStatus;
        $this->bindType = $bindType;
        $this->opDate = $opDate;
        $this->baCode = $baCode;
        $this->memberCode = $memberCode;

        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'cusUnionId' => $this->cusUnionID,
            'cusUserId' => $this->cusUserID,
            'baUserId' => $this->baUserID,
            'bindStatus' => $this->bindStatus,
            'opDate' => $this->opDate,
            'bindType' => $this->bindType,
            'brand_code' => $this->config->getBrandCode(),
            'program_code' => $this->config->getProgramCode()
        ];

        if (!empty($this->memberCode)) {
            $this->queryParam["memberCode"] = $this->memberCode;
        }

        if (!empty($this->baCode)) {
            $this->queryParam["baCode"] = $this->baCode;
        }
    }
}
