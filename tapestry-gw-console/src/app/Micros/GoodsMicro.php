<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2020-03-19
 * Time: 16:56
 */

namespace App\Micros;

use CrsSdk\MicroGoodsV1\Goods;
use CrsSdk\MicroGoodsV1\MaterialChannel;
use Illuminate\Support\Facades\Log;

class GoodsMicro
{
    protected $params;
    protected $response;

    public function __construct($params)
    {
        $this->params = $params;
        Log::info('micro-goods-v1-request', $params);
    }

    public function getChannelList()
    {
        $this->response = (new Goods())->list($this->params);

        return $this->response ?? [];
    }

    public function editChannel()
    {
        $this->response = (new MaterialChannel())->edit($this->params);
    }

    public function __destruct()
    {
        Log::info('micro-goods-v1-response', $this->response ? (array)$this->response : []);
    }
}
