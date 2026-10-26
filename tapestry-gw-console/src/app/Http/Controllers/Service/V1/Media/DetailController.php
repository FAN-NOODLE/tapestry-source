<?php

namespace App\Http\Controllers\Service\V1\Media;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\QYwxMedia;

class DetailController extends BaseController
{

    /**
     * 获取企业素材信息
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->request->only(
            'scid'
//            , 'uc_eid'
        );
        if (!$input['scid']) return error(1100002);

        $media = QYwxMedia::query()->where($input)->first();
        if (!$media) {
            throw new Error(10000, '素材不存在');
        }
        return success(['data' => $media]);
    }
}
