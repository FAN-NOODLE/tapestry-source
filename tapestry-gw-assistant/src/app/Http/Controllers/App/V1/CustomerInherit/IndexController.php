<?php

namespace App\Http\Controllers\App\V1\CustomerInherit;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExternalSearch;
use App\Micros\MicroScrm\GoMicroScrm;
use App\Micros\Work\CustomerInherit\CustomerInherit;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class IndexController extends BaseController
{
    /**
     * @var $inherit
     */
    private $inherit;

    public function __construct()
    {
        parent::__construct();
        $this->inherit = new CustomerInherit();
    }

    /**
     * @return void
     * 任务列表
     */
    public function inheritList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $dataRange = $this->dataRangeList;
        $third_id = array_column($dataRange, "object_id");
        $params['third_ids'] = $third_id;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $result = $this->inherit->InheritList($params);
        return success($result);
    }

    /**
     * @return JsonResponse
     * 任务详情
     */
    public function inheritInfo(): JsonResponse
    {
        $params = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;

        $result = $this->inherit->inheritInfo($params);
        return success($result);
    }

    /**
     * 顾客转移
     *
     * @return void
     * @throws Exception
     */
    public function inheritTransfer(): JsonResponse
    {
        $params = request()->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params['member_name'] = $this->memberAuthInfo['name'];
        $dataRange = $this->dataRangeList;
        $third_id = array_column($dataRange, "object_id");
        $params['third_ids'] = $third_id;
        $result = $this->inherit->inheritTransfer($params);
        return success($result);
    }

    /**
     * @return void
     * 导购列表
     */
    public function memberList(Request $request): JsonResponse
    {
//        $params                  = request()->all();
//        $params['enterprise_id'] = $this->enterpriseId;
//        $params['mobile']        = $this->mobile;
//
//        $admin_id = $this->loginAdminEnterprise['admin_id'];
//        // 获取企业管理者权限范围数据
//        $dataRange        = (new Account())->dataRangeList(['admin_id' => $admin_id]);
//        $dpIds            = collect($dataRange)->pluck('object_id')->all();


        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $dataRange = $this->dataRangeList;
        $third_id = array_column($dataRange, "object_id");
        $params['third_ids'] = $third_id;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params['dp_ids'] = $third_id;
        $params['mobile'] = $this->mobile;
        $result = $this->inherit->memberList($params);
        $res = (new GoMicroScrm())->getTakeoverMemberList($params);
        if (!empty($res)) {
            foreach ($result['list'] as &$item) {
                if (in_array($item['member_id'], $res)) {
                    $item['check'] = false;
                }
            }
            unset($item);
        }
        return success($result);
    }

    /**
     * @return JsonResponse
     * @throws Exception
     * 外部联系人接口
     */
    public function customerList(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['mobile'] = $this->mobile;

        $admin_id = $this->loginAdminEnterprise['admin_id'];
        // 获取企业管理者权限范围数据
        $dataRange = $this->dataRangeList;
        $dpIds = array_column($dataRange, "object_id");
        $params['third_ids'] = $dpIds;
        $params['member_type'] = 3;// 3导购 4管理
        if ($this->isManager == 1) {
            $params['member_type'] = 4;
        }
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params['is_group'] = (bool)$params['is_group'];
        if (!empty($params['birthday'])) {
            $params["birthday_start"] = $params['birthday']["birthday_start"];
            $params["birthday_end"] = $params['birthday']["birthday_end"];
        }
        if (!empty($params['time'])) {
            $params["last_purchase_start_time"] = $params['time']["last_purchase_start_time"];
            $params["last_purchase_end_time"] = $params['time']["last_purchase_end_time"];
        }

        $result = (new GoMicroScrm())->getExternalUserList($params);
        //调三方添加标签数据 scrm 已处理
        //$result = $this->splicingExternalUser($result, $params);

        if ($params['is_group']) {
            //调用外部联系人关联的导购
            $work_external_user_id = collect($result['list'])->pluck('work_external_user_id')->all();

            $members = (new CustomerInherit())->memberRelation([
                'enterprise_id' => $this->enterpriseId,
                'work_external_user_id' => $work_external_user_id,
                'third_ids' => $dpIds,
                'params' => $params
            ]);
            //组装数据
            foreach ($result['list'] as &$item) {

                $itemArr = $item;
                $item['more_member'] = [];
                foreach ($members as $value) {
                    if ($value['work_external_user_id'] == $item['work_external_user_id']) {
                        if ($value['member_id'] != $item['member_id']) {
                            $itemArr['id'] = $value['id'];
                            $itemArr['member_name'] = $value['member_name'];
                            $itemArr['member_id'] = $value['member_id'];
                            $itemArr['remark'] = $value['remark'];
                            $itemArr['work_external_user_id'] = $value['work_external_user_id'];
                            $itemArr['created_at'] = $value['created_at'];
                            $itemArr['is_key_sa'] = $value['is_key_sa'];
                            $item['more_member'][] = $itemArr;
                        }
                    }
                }
            }

        } else {

            //调用外部联系人被分配的次数
            $external_id = collect($result['list'])->pluck('work_external_user_id')->all();

            $res = (new CustomerInherit())->customerInheritRecord(['enterprise_id' => $this->enterpriseId, 'external_id' => $external_id]);

            $record = collect($res['external_to_total'])->keyBy('external_id')->all();
            $waitExternalIdArr = $res['wait_external_arr'] ?? [];
            //组装数据
            foreach ($result['list'] as &$item) {
                $item['transfer_num'] = $record[$item['work_external_user_id']]['num'] ?? 0;
                $item['transfer_num'] = max(2 - $item['transfer_num'], 0);
                if (in_array($item['work_external_user_id'], $waitExternalIdArr)) {
                    $item['transfer_num'] = -1; //分配中
                }
            }
        }

        return success($result);
    }


    /**
     * 外部联系人接口（新版）
     *
     * @param Request $request
     * @return JsonResponse
     * @throws Exception
     */
    public function searchList(Request $request): JsonResponse
    {
        $params = $request->input();
        $params['enterprise_id'] = $this->enterpriseId;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $params = array_merge($params, $this->getRoleRange());
        if ($params['is_group']) {
            $data = (new ExternalSearch())->list($params);
            if (empty($data['list'])) {
                return success($data);
            }
            //调用外部联系人关联的导购
            $params['external_user_ids'] = array_column($data['list'], 'id');
            $params['is_relationship'] = true;
            $params['page'] = 1;
            $params['size'] = 999;
            $allData = (new ExternalSearch())->list($params);
            //组装数据
            $externalList = [];
            $unionIds = [];
            foreach ($allData['list'] as $value) {
                if (in_array($value['unionid'],$unionIds)) {
                    $externalList[$value['unionid']]['more_member'][] = $value;
                }else{
                    $externalList[$value['unionid']] = $value;
                    $unionIds[] = $value['unionid'];
                }
            }
            $data['list'] = array_values($externalList);
        } else {
            $params['is_relationship'] = true;
            $data = (new ExternalSearch())->list($params);
            if (empty($data['list'])) {
                return success($data);
            }

            $res = (new CustomerInherit())->customerInheritRecord(['enterprise_id' => $this->enterpriseId, 'external_id' => array_column($data['list'], 'work_external_user_id')]);
            $record = collect($res['external_to_total'])->keyBy('external_id')->all();
            $waitExternalIdArr = $res['wait_external_arr'] ?? [];
            //组装数据
            foreach ($data['list'] as &$item) {
                $item['transfer_num'] = $record[$item['work_external_user_id']]['num'] ?? 0;
                $item['transfer_num'] = max(2 - $item['transfer_num'], 0);
                if (in_array($item['work_external_user_id'], $waitExternalIdArr)) {
                    $item['transfer_num'] = -1; //分配中
                }
            }
        }

        return success($data);
    }

    /**
     * @return void
     * 拼接三方接口数据
     */
    public function splicingExternalUser($data, $params): array
    {
        $cardDescription = ['O' => 0,
            'PROSPECT' => 1,
            'C' => 2,
            'D' => 3,
            'GOLD' => 4,
            'SILVER' => 5,
            'DIAMOND' => 6];
        foreach ($data['list'] as $users => $user) {

            $data['list'][$users]['customer_card'] = $cardDescription[$user['card_level']] ?? 0;
            $data['list'][$users]['last_12m_purchase_amount'] = $user['amount_total'] / 100;
            $data['list'][$users]['client_type'] = $user['clt_type'];

            /*$data['list'][$users]['customer_type']            = '';
            $data['list'][$users]['customer_card']            = 0;
            $data['list'][$users]['name']                     = '';
            $data['list'][$users]['mobile']                   = '';
            $data['list'][$users]['last_12m_purchase_amount'] = 0;
            $data['list'][$users]['client_type']              = '';

            if ($user['unionid']) {
                // 查询crm
                $userCrm = (new Work())->getCrmUserInfo($params['enterprise_id'], $user['unionid'], $user['euid']);
                if ($userCrm['status']) {
                    $data['list'][$users]['customer_type']            = $userCrm['customer_type'];
                    $data['list'][$users]['customer_card']            = $userCrm['customer_card'];
                    $data['list'][$users]['name']                     = $userCrm['name'];
                    $data['list'][$users]['mobile']                   = $userCrm['mobile'];
                    $data['list'][$users]['last_12m_purchase_amount'] = $userCrm['last_12m_purchase_amount'];
                    $data['list'][$users]['client_type']              = $userCrm['client_type'];
                }
            }*/
        }

        return $data;
    }

    /**
     * 在职分配的顾客列表全选的人数
     * @param Request $request
     * @return JsonResponse
     */
    public function selectAllNum(Request $request): JsonResponse
    {
        $params = $request->all();
        $params['enterprise_id'] = $this->enterpriseId;
        $dataRange = $this->dataRangeList;
        $third_id = array_column($dataRange, "object_id");
        $params['third_ids'] = $third_id;
        $params['member_id'] = $this->memberAuthInfo['member_id'];
        $result = (new GoMicroScrm())->getOnInheritSelectAllNum($params);
        return success(["num" => $result]);
    }
}