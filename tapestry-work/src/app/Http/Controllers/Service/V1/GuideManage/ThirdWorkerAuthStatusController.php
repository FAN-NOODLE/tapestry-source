<?php


namespace App\Http\Controllers\Service\V1\GuideManage;


use App\Http\Controllers\Service\V1\BaseController;
use App\Models\CrsWorkerThirdPartRelation;
use App\Models\PartnerEnterpriseMap;
use App\Models\PartnerMember;
use App\Services\PartnerMemberService;
use Illuminate\Http\Request;

/**
 * 获取第三方导购认证状态
 * Class ThirdWorkerAuthStatusController
 * @package App\Http\Controllers\Service\V1\GuideManage
 */
class ThirdWorkerAuthStatusController extends BaseController
{

    public function index(Request $request)
    {
        $rules = [
            'enterprise_id' => 'required',
            'source_type' => 'required',
            'worker_id' => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
            'source_type.required' => '企业[source_type]不能为空',
            'worker_id.required' => '企业[worker_id]不能为空',
        ];

        $input = $this->validate($request, $rules, $message);

        $where = [
            'enterprise_id' => $input['enterprise_id'],
            'store_staff_id' => $input['worker_id'],
            'sync_source' => $input['source_type'],
        ];
        $relationExisted = CrsWorkerThirdPartRelation::query()->where($where)->exists();

        //查询合伙人关系表
        if (!$relationExisted) {
            $info = (new PartnerMemberService())->getPartnerInfo([
                'id' => $input['worker_id'],
                'enterprise_id' => $input['enterprise_id']
            ]);
            if (!empty($info) && $info['user_id'] > 0) {
                $relationExisted = true;
            }
        }

        return success([
            'result' => $relationExisted
        ]);
    }
}