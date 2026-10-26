<?php
/**
 * Created by :PhpStorm
 * User: lvshuaishuai
 * Date: 2020/11/16
 * Time: 3:00 下午
 */

namespace App\Http\Controllers\Service\V1\DistributionPartner;


use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\Contact\ContactController;
use App\Models\DistributionOrderSettlement;
use App\Models\PartnerEnterpriseMap;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Validator;

class ListController extends ContactController
{
    private $input = [];

    public function index(Request $request)
    {
        // 接收参数
        $params = $request->only([
            'enterprise_id',
            'worker_id',
            'worker_ids',  //合伙人数组
            'partner_status',
            'worker_name',
            'partner_name',
            'partner_mobile',
            'partner_ids',  //合伙人id数组
            'user_id',  //合伙人id数组
            'page',
            'size'
        ]);
        // 验证规则
        $rules = [
            'enterprise_id'  => 'integer',
            'partner_name'   => 'string',
            'partner_mobile' => 'integer',
            'worker_id'      => 'integer'
        ];

        // 提示语
        $messages = [
            'partner_mobile.integer' => '手机号需为int类型',
            'worker_id.integer'      => '导购id需int类型'
        ];
        // 参数验证
        $validator = Validator::make($params, $rules, $messages);

        if ($validator->fails()) {
            throw new Error(500, $validator->errors()->first());
        }
        $params['size'] = $params['size'] ?? 10;
        $params['page'] = $params['page'] ?? 1;
        $this->input    = $params;

        $list  = $this->getPartnerList();
        $total = count($list);
        $data  = [
            'page'  => intval($this->input['page']),
            'size'  => intval($this->input['size']),
            'total' => $total,
            'list'  => $list ?? []
        ];

        return success($data);

    }

    /**
     *  合伙人列表
     */
    public function getPartnerList()
    {
        $params = $this->input;

        $columns = [
            'partner_enterprise_map.*',
            'partner_member.id',
            'partner_member.partner_sn',
            'partner_member.partner_name',
            'partner_member.partner_sex',
            'partner_member.partner_mobile',
            'partner_member.open_id',
            'partner_member.union_id',
            'partner_member.mobile_short',
            'partner_member.avatar'];
        $partnerList = PartnerEnterpriseMap::query()
            ->leftJoin('partner_member', 'partner_enterprise_map.partner_id', '=', 'partner_member.id')
            ->when(!empty($params['enterprise_id']), function ($query) use ($params) {
                return $query->where('partner_enterprise_map.enterprise_id', $params['enterprise_id']);
            })
            ->when(!empty($params['worker_id']), function ($query) use ($params) {
                return $query->where('partner_enterprise_map.worker_id', $params['worker_id']);
            })
            ->when(!empty($params['worker_ids']), function ($query) use ($params) {
                return $query->whereIn('partner_enterprise_map.worker_id', $params['worker_ids']);
            })
            ->when(isset($params['partner_status']), function ($query) use ($params) {
                return $query->where('partner_enterprise_map.partner_status', $params['partner_status']);
            })
            ->when(!empty($params['worker_name']), function ($query) use ($params) {
                return $query->where("partner_enterprise_map.worker_name", 'like', '%' . $params['worker_name'] . '%');
            })
            ->when(!empty($params['partner_ids']), function ($query) use ($params) {
                return $query->whereIn("partner_member.id", $params['partner_ids']);
            })
            ->when(!empty($params['partner_mobile']), function ($query) use ($params) {
                return $query->where("partner_member.partner_mobile", $params['partner_mobile']);
            })
            ->when(!empty($params['user_id']), function ($query) use ($params) {
                if (is_array($params['user_id'])) {
                    return $query->whereIn("partner_enterprise_map.user_id", $params['user_id']);
                } else {
                    return $query->where("partner_enterprise_map.user_id", $params['user_id']);
                }
            })
            ->when(!empty($params['partner_name']), function ($query) use ($params) {
                $query->where(function ($query) use ($params) {
                    $query->where("partner_member.partner_name",  'like', '%' . $params['partner_name'] . '%')
                        ->orWhere("partner_member.partner_mobile",  'like', '%' . $params['partner_name'] . '%');
                });
            })
            ->paginate($this->input['size'], $columns, '', $this->input['page'])
            ->toArray();

        // 转数组
        foreach ($partnerList['data'] as $ids => &$orders) {
            unset($partnerList['data'][$ids]['partner_sn']);
            unset($partnerList['data'][$ids]['partner_sex']);
            unset($partnerList['data'][$ids]['open_id']);
            unset($partnerList['data'][$ids]['union_id']);
            if($partnerList['data'][$ids]['mobile_sensitive'] == 1 && !empty($partnerList['data'][$ids]['partner_mobile'])){
                $partnerList['data'][$ids]['partner_mobile'] = substr_replace(
                    $partnerList['data'][$ids]['partner_mobile'], '****', 3, 4);
            }

            $orders['created_at']     = (int)$orders['created_at'];
            $orders['updated_at']     = (int)$orders['updated_at'];
            $orders['deleted_at']     = (int)$orders['deleted_at'];
//            $orders['mobile_show']    = (int)$orders['mobile_show'];
//            $orders['partner_mobile'] = (int)$orders['partner_mobile'];
            // 分销订单数
            $partnerList['data'][$ids]['orders'] = DistributionOrderSettlement::query()
                ->where('enterprise_id', '=', $partnerList['data'][$ids]['enterprise_id'])
                ->where('partner_id', '=', $partnerList['data'][$ids]['id'])
                ->where('order_status', '>=', 30) // 已支付
                ->count();
        }

        return $partnerList['data'];


    }
}