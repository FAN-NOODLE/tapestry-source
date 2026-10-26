<?php


namespace App\Http\Controllers\V1\Qywx\Tag;

use App\Http\Controllers\V1\BaseController;
use CrsSdk\Work\Tag;

class SyncTagController extends BaseController
{
    /**
     * 同步企业标签库
     * @api /app/v1/qywx/tag/tag-sync
     * @return mixed
     * @throws \Exception
     */
    public function index()
    {
        $input['qy_domain'] = $this->qyDomain;

        return (new Tag())->tagSync($input);
    }

}