<?php

namespace App\Http\Controllers\Service\V1\EnterprisePlugin;

use App\Library\Uc\EnterprisePlugin;
use App\Http\Controllers\Service\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 查询企业安装的应用
     *
     * @return array|\Psr\Http\Message\StreamInterface
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        $this->validate($this->request, [
            'qy_domain' => 'required'
        ]);
        $input = $this->request->all();

        // CRS渠道类型:5
        $data = (new EnterprisePlugin())->list($input['qy_domain'], ['channelType' => 5]);

        unset($data['requestId']);

        // 过滤掉没安装的
        foreach ($data as $key => $value) {
            if ($value['available'] !== 1) {
                unset($data[$key]);
            }
        }

        return $data;
    }
}
