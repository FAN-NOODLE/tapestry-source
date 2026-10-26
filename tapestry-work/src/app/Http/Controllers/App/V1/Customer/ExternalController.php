<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\External;
use App\Library\Uc\Member;

class ExternalController extends BaseController
{
    /**
     * 测试uc接口，获取员工下的外部联系人列表
     * @param External $external
     *
     * @return array
     * @throws \Exception
     */
    public function index(External $external)
    {
        $list = $external->listForMember($this->qyDomain, ['memUid' => $this->memUid]);

        return success([$list]);
    }

    /**
     * 测试uc接口，获取员工列表
     * @param Member $member
     *
     * @return array
     * @throws \Exception
     */
    public function member(Member $member)
    {
        $list = $member->list($this->qyDomain, $this->appIdentifier);

        return success([$list]);
    }
}
