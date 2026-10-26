<?php

namespace App\Http\Controllers\Service\V1\Admin;

use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Admin;
use App\Library\Uc\Member;
use App\Library\Util\Logger;
use App\Models\StoreWorker;

class DetailController extends BaseController
{
    /**
     * 管理员详情
     *
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'qy_domain' => 'required|string',
            'admin_id'  => 'required|string'
        ]);
        $input = $this->request->all();

        return (new Admin())->detail($input['qy_domain'], [
            'eaId' => $input['admin_id']
        ]);
    }

    public function list()
    {
        $this->validate($this->request, [
            'qy_domain' => 'required|string'
        ]);
        $input = $this->request->all();
        $args  = [
//            'eqMemStatus' => 3,
            'pageNum'     => 1,
            'pageSize'    => 999,
            'mobiles'     => $input['mobiles'] ?? []
        ];

        $res   = (new Member())->
        ignorePermissionsList($input['qy_domain'], null, $args);

        if($res['list']){
            return array_unique(array_column($res['list'],'memMobile'));
        }
        return [];
    }


}
