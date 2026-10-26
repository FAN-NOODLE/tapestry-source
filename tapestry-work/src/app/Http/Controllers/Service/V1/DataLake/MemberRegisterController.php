<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2021/4/2 0002
 * Time: 11:02
 */

namespace App\Http\Controllers\Service\V1\DataLake;


use Illuminate\Support\Facades\Config;
use Illuminate\Support\Facades\DB;

class MemberRegisterController extends DataLakeBaseController
{
    /**
     * @throws \Throwable
     */
    public function requestApi()
    {
        return $this->dataLakeService->memberRegister($this->getParam());
    }

    public function getParam()
    {
        $request = $this->request->input();

        return [
            'real_name'     => $request['real_name'] ?? '',
            'mobile'        => $request['mobile'],
            'birthday'      => $request['birthday'] ?? '',
            'sex'           => $request['sex'] ?? 'U',
            'created_at'    => time(),
            'consentStatus' => '0',
            "open_id"       => $request['open_id'],
            "nick_name"     => $request['nick_name'] ?? '',
            "register_at"   => time(),
            "union_id"      => $request['union_id']
        ];
    }

    public function checkParam()
    {

    }

}