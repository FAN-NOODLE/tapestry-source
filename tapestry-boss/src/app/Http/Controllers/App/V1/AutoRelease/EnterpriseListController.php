<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2019/11/8
 * Time: 14:34
 */
namespace App\Http\Controllers\App\V1\AutoRelease;

use App\Models\Enterprise;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;

class EnterpriseListController extends BaseController
{
    /**
     * 小程序列表
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        $rules = [
            'page'  => 'integer',
            'limit' => 'integer',
            'type'  => 'integer'
        ];
        $messages = [
            'page.integer'    => '页码必须为正整数',
            'limit.integer'   => '每页数量必须为正整数',
            'type.integer' => '小程序类型必须为正整数',
        ];
        $this->validate($request, $rules, $messages);

        $page = $request->post('page', 1);
        $size = $request->post('limit', 10);
        $type = $request->post('type', EnterpriseThirdPartAccount::RELEASE_STABLE);

        $enterpriseObj = EnterpriseThirdPartAccount::query();

        if ($request->filled('name')) {
            $name = $request->post('name');
            $enterpriseObj = $enterpriseObj->where('name', 'like', '%' . $name . '%');
        }
        $enterpriseObj->where('type', '=', EnterpriseThirdPartAccount::TYPE_SHOP);
        $enterpriseObj->where('status', '=', EnterpriseThirdPartAccount::STATUS_ACTIVE);
        $enterpriseObj->where('verify_status', '=', EnterpriseThirdPartAccount::VERIFIED);
        $enterpriseObj->where('enterprise_type', '=', $type);

        $total = $enterpriseObj->count('id');
        $enterpriseObj = $enterpriseObj->forpage($page, $size)
            ->get(['enterprise_id as id', 'name']);

        $enterpriseList = $enterpriseObj->toArray() ?? [];

        return [
            'data' => $enterpriseList,
            'meta' => [
                'total' => $total,
                'size'  => $size,
                'page'  => $page
            ]
        ];
    }
}