<?php


namespace App\Http\Controllers\Service\V1\Tag;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Tag;

class SyncTagController extends BaseController
{

    /**
     * 同步企业标签库
     * @api /service/v1/tag/tag-sync
     * @return array
     */
    public function index()
    {
        $qyDomain = $this->request->input('qy_domain', '');

        if (empty($qyDomain)) {
            notice(200018);
        }

        try {
            (new Tag())->syncTag($qyDomain);
        } catch (\Exception $exception) {
            throw new Error($exception->getCode(), $exception->getMessage());
        }

        return success(['data' => []]);
    }
}