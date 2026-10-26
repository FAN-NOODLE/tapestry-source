<?php
/**
 * Created by PhpStorm.
 * User: XuHuitao
 * Date: 2020/7/14
 * Time: 10:09
 */

namespace App\Http\Controllers\App\V1\Customer;

use App\Http\Controllers\App\V1\BaseController;
use App\Micros\MicroScrm\ExternalSearchSw;
use App\Micros\Work\CustomerInherit\CustomerInherit;
use Exception;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

class SearchSwController extends BaseController
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
        return success((new ExternalSearchSw())->choiceOption($params));
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
        $params['user_id'] = $this->memberAuthInfo['userid'];
        $params = array_merge($params, $this->getRoleRange());
        $data = (new ExternalSearchSw())->list($params);

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
        return success((new ExternalSearchSw())->getOpenIds($params));
    }

}