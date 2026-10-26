<?php

namespace App\Http\Controllers\Service\V1\Source;

use App\Http\Controllers\Service\V1\BaseController;
use App\Services\SourceServices;

class UpDataMediaIdController extends BaseController
{
    /**
     * 更新素材临时Id
     *
     * @return array
     * @throws \Exception
     */
    public function index()
    {
        $params = request()->only(['uc_eid', 'image_media_url', 'scid']);

        (new SourceServices())->sysSourceMediaId([$params]);

        return success([]);
    }
}
