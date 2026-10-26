<?php

namespace App\Http\Controllers\App\V1\CustomerInherit;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExternalSearchSw;
use App\Micros\Work\CustomerInherit\CustomerInherit;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class IndexSwController extends BaseController
{

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
        $params['user_id'] = $this->memberAuthInfo['userid'];
        $params = array_merge($params, $this->getRoleRange());
        $params['is_inherit'] = true;
        if ($params['is_group']) {
            $data = (new ExternalSearchSw())->list($params);
            if (empty($data['list'])) {
                return success($data);
            }
            //调用外部联系人关联的导购
            $params['external_user_ids'] = array_column($data['list'], 'id');
            $params['is_relationship'] = true;
            $params['page'] = 1;
            $params['size'] = 999;
            $allData = (new ExternalSearchSw())->list($params);
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
            $data = (new ExternalSearchSw())->list($params);
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
}