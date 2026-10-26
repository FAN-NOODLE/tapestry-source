<?php

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;
use Overtrue\Pinyin\Pinyin;

class EditController extends BaseController
{
    /**
     * 编辑客户资料
     *
     * @param Request $request
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index(Request $request)
    {
        $this->validate($request, [
            'ec_userid'        => 'required',
            'user_remark_name' => 'string',
            'user_description' => 'string'
        ], [
            'ec_userid.required' => '外部联系人id不能为空',
        ]);

        $params = $request->only(['user_remark_name', 'user_description']);

        // 更新客户备注名时，同时更新客户姓名索引，索引为拼音
        if ($request->input('user_remark_name')) {
            $pinyin = app(Pinyin::class)->abbr($request->input('user_remark_name'));

            // 如果昵称是表情符号，user_name_index使用汉字"井",因为如果使用符号'#'，客户列表排序时'#'会在最上面，按照微信习惯，'#'应该最下面
            $params['user_name_index'] =
                (preg_match('/^[a-zA-Z\x{4e00}-\x{9fa5}]+/u', $request->input('user_remark_name')))
                    ? $pinyin : '井';
        }

        $data = MemberUserRelation::where([
            'ec_userid' => $request->input('ec_userid'),
            'mem_uid'   => $this->memUid
        ])->update($params);

        return success(['data' => $data]);
    }
}
