<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/20
 * Time: 3:53 下午
 */

namespace App\Http\Controllers\Service\V1\DistributionPartner;


use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\PartnerEnterpriseMap;
use App\Models\PartnerMember;
use App\Services\StoreWorkerService;
use CrsSdk\Boss\Enterprise;
use CrsSdk\Boss\MiniProgram;
use CrsSdk\Shop\StoreWorker;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;

class MyController extends BaseController
{
    public function index(Request $request)
    {
        $params = $request->only([
            'enterprise_id',
            'partner_id',
        ]);

        // 验证规则
        $rules = [
            'enterprise_id' => 'required|integer',
            'partner_id'    => 'required|integer',
        ];

        // 提示语
        $messages = [
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID需为int类型',
            'partner_id.required' => '合伙人ID不能为空',
            'partner_id.integer' => '合伙人ID需为int类型',
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }

        // 条件
        $where = [
            'enterprise_id' => $params['enterprise_id'],
            'partner_id'    => $params['partner_id']
        ];
        switchAppHost(['enterprise_id' => $params['enterprise_id']]);
        // 字段筛选
        $selectFields = [
            'partner_member.partner_name', 'partner_member.avatar', 'partner_enterprise_map.worker_avatar',
            'partner_enterprise_map.worker_id', 'partner_enterprise_map.worker_name', 'partner_enterprise_map.enterprise_id','partner_member.partner_mobile',
            'partner_enterprise_map.mobile_sensitive','partner_enterprise_map.sms_switch'
        ];
        // 合伙人信息
        $partnerInfo = PartnerMember::query()
            ->leftJoin('partner_enterprise_map', 'partner_member.id', '=', 'partner_enterprise_map.partner_id')
            ->where($where)
            ->first($selectFields);
        if (!empty($partnerInfo)) {
            $partnerInfo = $partnerInfo->toArray();
            $workerInfo = StoreWorkerService::detail([
                'enterprise_id' => $params['enterprise_id'],
                'worker_id' => $partnerInfo['worker_id'],
            ]);
            if (!empty($workerInfo)){
                $partnerInfo['worker_status'] = $workerInfo['worker_status'];
            }
        } else {
            return success();
        }
        // 多品牌
        $partnerEnterpriseInfo = PartnerEnterpriseMap::query()
            ->where('partner_id', '=', $params['partner_id'])
            ->get()
            ->toArray();
        // 合伙人所在企业
        $enterpriseIds = array_column($partnerEnterpriseInfo, 'enterprise_id');
        // 合伙人不再组织架构，所以查询所有的企业信息
        //$enterpriseInfo = (new Enterprise())->listsOfficial($enterpriseIds);
        $enterpriseInfo = (new Enterprise())->listNoPage(["enterprise_ids"=>$enterpriseIds]);
        $einfo = array_column($enterpriseInfo['data'], null, 'id');
        unset($partnerEnterpriseInfo);
        unset($enterpriseInfo);

        //$enterprise = (new Enterprise());
        foreach ($enterpriseIds as $info) {
            $partnerInfo['enterprise_info'][] = $einfo[$info] ?? [];
        }
        foreach ($partnerInfo['enterprise_info'] as $key => $value) {
            /*$params = [
                'enterprise_id' => $value['id']
            ];*/
            //$partnerInfo['enterprise_info'][$key]['logo'] = $enterprise->enterpriseMerchantSettingInfo($params)['data']['logo'];
            //企业名称和logo 使用小程序商城的名称和logo
            $miniProgramInfo = (new MiniProgram())->detail(['enterprise_id'=>$value['id']]);
            $value['shop_name'] = $miniProgramInfo['data']['name'] ?? '';
            $value['logo'] = $miniProgramInfo['data']['logo'] ?? '';
            $partnerInfo['enterprise_info'][$key] = $value;
        }
        return success($partnerInfo);
    }
}