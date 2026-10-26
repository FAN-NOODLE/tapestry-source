<?php


namespace App\Library\DataLake\BaBind;

use App\Library\DataLake\Base;

/**
 * Class BaBindQuery
 * @package App\Library\DataLake\BaBind
 *
 * 企业微信BA绑定查询
 */
class BaBindQuery extends Base
{

    public $api = "/api/interface/third/member/wechatBaBindQuery";

    private $queryType = '';
    private $value = '';
    private $bindType = "0";

    /**
     * BaBindQuery constructor.
     * @param string $queryType 查询BA绑定关系的字段类型
     * @param string $value
     * @param string $bindType
     */
    public function __construct(string $queryType, string $value, string $bindType)
    {
        $this->queryType = $queryType;
        $this->value = $value;
        $this->bindType = $bindType;

        parent::__construct();
    }

    protected function setQueryParam()
    {
        $this->queryParam = [
            'queryType' => $this->queryType,
            'value' => $this->value,
            'bindType' => $this->bindType,
            'brand_code' => $this->config->getBrandCode(),
            'program_code' => $this->config->getProgramCode()
        ];
    }

}