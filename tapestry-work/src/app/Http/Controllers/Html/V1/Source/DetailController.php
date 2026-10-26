<?php

namespace App\Http\Controllers\Html\V1\Source;

use App\Http\Controllers\Html\V1\BaseController;
use CrsSdk\Console\Media;

class DetailController extends BaseController
{
    protected $requireLogin = false;

    /**
     * 获取企业素材
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = request()->only(['scid']);

        // 企业ID
//        $params['uc_eid'] = $this->domain;

        switchAppHostByDomain($this->domain);

        // 获取企业素材信息
        $data = (new Media())->detail($params);

        return success($data);
    }
}
