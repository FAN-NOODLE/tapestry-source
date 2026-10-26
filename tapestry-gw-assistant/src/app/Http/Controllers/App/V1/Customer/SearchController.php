<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/14
 * Time: 10:09
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExternalSearch;
use App\Micros\Work\CustomerInherit\CustomerInherit;
use App\Micros\Work\Work;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SearchController extends BaseController
{

    /**
     * 外部联系人筛选
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     * @Time: 2022/10/21
     */
    public function getChoiceOption(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params = array_merge($params, $this->getRoleRange());
        return success((new ExternalSearch())->choiceOption($params));
    }

    /**
     * 外部联系人列表
     *
     * @param Request $request
     *
     * @return JsonResponse
     * @throws Exception
     * @Time: 2022/10/21
     */
    public function list(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params = array_merge($params, $this->getRoleRange());
        $data = (new ExternalSearch())->list($params);

        if (empty($data['list'])) {
            return success($data);
        }
        $recordForExternal = (new CustomerInherit())->recordForExternalIds([
            'enterprise_id' => $this->enterpriseId,
            'external_ids' => array_column($data['list'], 'id'),
        ]);
        foreach ($data['list'] as $key => $info) {
            foreach ($recordForExternal as $item) {
                $data['list'][$key]['is_new_inherit'] = $item['takeover_member_id'] == $this->memberAuthInfo['member_id'] && $item['external_id'] == $info['work_external_user_id'];
            }
        }
        return success($data);

    }


    /**
     * 客户详情
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function detail(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        return success((new ExternalSearch())->detail($params));
    }

    /**
     * CRM三方客户详情
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function detailByCrm(Request $request): JsonResponse
    {
        $res = (new Work())->customerDetailByCoachCrm([
            'union_id' => $request->input('unionid', ''),
            'enterprise_id' => $this->enterpriseId
        ]);
        if (!empty($res)) {
            $res['expiry_date'] = !empty($res['expiry_date']) ? date('Y-m-d', strtotime($res['expiry_date'])) : '--';
            $res['upgrade_date_time'] = !empty($res['upgrade_date_time']) ? date('Y-m-d', strtotime($res['upgrade_date_time'])) : '';
            if (!empty($res['vipStatus'])) {
                $res['vipStatus']['issueDate'] = date('Y-m-d', strtotime($res['vipStatus']['issueDate']));
                $res['vipStatus']['expiryDate'] = date('Y-m-d', strtotime($res['vipStatus']['expiryDate']));
            }
            $res['registration_date'] = !empty($res['registration_date']) ? date('Y-m-d', strtotime($res['registration_date'])) : '';
        }
        return success($res);
    }

    /**
     * CDP推荐商品
     *
     * @param Request $request
     * @return JsonResponse
     */
    public function cdpInterestedProductsByCdp(Request $request): JsonResponse
    {
        $unionId = $request->input('unionid', '');
        $res = [];
        if (!empty($unionId)) {
            $cdpParams = [
                "userId" => $unionId,
                "userKey" => "union_id",
                "offset" => 0,
                "limit" => 15,
                "channelId" => "weclient"
            ];
             $res = ((new Work())->getCustomerInterestedProducts($cdpParams));
        }
        return success($res);

    }


    /**
     * 全选的外部联系人openid集合
     * @param Request $request
     * @return JsonResponse
     */
    public function getOpenIds(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params = array_merge($params, $this->getRoleRange());
        return success((new ExternalSearch())->getOpenIds($params));
    }

}