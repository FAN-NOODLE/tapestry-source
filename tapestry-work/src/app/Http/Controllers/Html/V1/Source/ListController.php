<?php

namespace App\Http\Controllers\Html\V1\Source;

use App\Http\Controllers\Html\V1\BaseController;
use CrsSdk\Console\Media;
use Illuminate\Support\Facades\Redis;

class ListController extends BaseController
{
    /**
     * 获取企业素材
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = request()->only(['search_text', 'type', 'page', 'size',]);

        // 企业ID
        $params['uc_eid'] = $this->domain;

        switchAppHostByDomain($this->domain);

        // 获取企业素材信息
        $data = (new Media())->list($params);

        $url = $this->webUrl();

        if (!empty($data['data'])) {
            foreach ($data['data'] as &$value) {
                $redisKey            = 'source_map_cache:' . $value['scid'];
                $value['media_id']   = Redis::get($redisKey) ?? '';
                $value['source_url'] = $url . '?id=' . $value['scid'];
            }
        }

        return success($data);
    }

    private function webUrl()
    {
        $env = config('app.env');

        if ($env === 'dev') {
            $url = 'https://work-dev-crs.vchangyi.com/' . $this->corpId . '/crs-manu-customer/card-detail';
        } elseif ($env === 'test') {
            $url = 'https://work-test-crs.vchangyi.com/' . $this->corpId . '/crs-manu-customer/card-detail';
        } elseif ($env === 'prod') {
            $url = 'https://work-crs.vchangyi.com/' . $this->corpId . '/crs-manu-customer/card-detail';
        } else {
            $url = 'https://work-dev-crs.vchangyi.com/' . $this->corpId . '/crs-manu-customer/card-detail';
        }
        return $url;
    }
}
