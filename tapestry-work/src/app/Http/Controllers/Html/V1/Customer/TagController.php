<?php


namespace App\Http\Controllers\Html\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Html\V1\BaseController;
use App\Services\TagService;

class TagController extends BaseController
{
    /**
     * 获取企业标签库标签列表
     * @return array
     * @throws \Exception
     */
    public function list()
    {
        // 获取全部标签列表
        $list = (new TagService())->getAllTagList($this->domain);
        $requestId = $list['requestId'];
        unset($list['requestId']);

        foreach ($list as &$value) {
            foreach ($value['tagList'] as &$v) {
                $v['tagId'] = (string)$v['tagId'];
            }
            $value['groupId'] = (string)$value['groupId'];
        }

        unset($v);
        unset($value);
        return success([
            'data' => $list,
            'requestId' => $requestId
        ]);
    }

}