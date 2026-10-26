<?php

namespace App\Http\Controllers\V1\Customer;

use App\Exceptions\Error;
use CrsSdk\Work\Customer;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class AdditionListController extends BaseController
{
    /**
     * 待添加客户列表
     * @return JsonResponse
     */
    public function index()
    {
        $page   = $this->request->input('page', 1);
        $size   = $this->request->input('size', 10);
        $type   = $this->request->input('type', 0);
        $memUid = $this->request->input('mem_uid', '');

        // 是否有手机号 1 是 2否
        $haveMobile = $this->request->input('have_mobile', 1);

        // 是否复制手机号 0 否 1 是
        $isCopyMobile = $this->request->input('is_copy_mobile', 0);
        if (!$memUid) {
            throw new Error(19112615, '员工memUid不能为空!');
        }

        $result = (new Customer())->additionList([
            'page'           => $page,
            'size'           => $size,
            'mem_uid'        => $memUid,
            'have_mobile'    => $haveMobile,
            'is_copy_mobile' => $isCopyMobile,
            'type'           => $type,
            'enterprise_id'  => $this->enterpriseId
        ]);

        return success([
            'list'  => $result['data'],
            'page'  => $result['meta']['page'],
            'size'  => $result['meta']['size'],
            'total' => $result['meta']['total'],
        ]);
    }
}
