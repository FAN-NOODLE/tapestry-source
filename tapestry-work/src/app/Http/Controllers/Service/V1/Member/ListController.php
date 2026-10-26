<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 15:47
 */

namespace App\Http\Controllers\Service\V1\Member;

use App\Library\Uc\Member;
use Illuminate\Http\Request;
use App\Services\MemberUserRelationService;
use App\Http\Controllers\Service\V1\BaseController;

class ListController extends BaseController
{
    /**
     * @var MemberUserRelationService
     */
    protected $memberUserRelationService;

    /**
     * CountController constructor.
     *
     * @param MemberUserRelationService $memberUserRelationService
     */
    public function __construct(
        MemberUserRelationService $memberUserRelationService
    ) {
        parent::__construct();
        $this->memberUserRelationService = $memberUserRelationService;
    }

    /**
     * 员工列表（忽略权限）
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function index(Request $request)
    {
        $qyDomain      = $request->post('qy_domain', '');
        $appIdentifier = $request->post('app_identifier', '');
        $params        = $request->post('params', []);
        $enterpriseId  = $request->post('enterprise_id', 0);

        $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, $params);
        $result = $this->putUserCount($result, $enterpriseId);

        return success($result);
    }

    /**
     * 获取导购的会员总数
     * @param $result
     * @param $enterpriseId
     * @return mixed
     */
    public function putUserCount($result, $enterpriseId)
    {
        if (isset($result['list']) && $enterpriseId) {
            $memUids = [];
            foreach ($result['list'] as $item) {
                if (!is_null($item['memSubscribeTime'])) {
                   $memUids[] = $item['memUid'];
                }
            }

            // v2.3.1 -> task-1081230 优化代码结构
            $userCountBy = [];
            if (!empty($memUids)) {
                $userCountBy = $this->memberUserRelationService->timeUserCounts($enterpriseId, $memUids);
            }

            foreach ($result['list'] as &$v) {
                if (!is_null($item['memSubscribeTime'])) {
                    $v['user_count'] = $userCountBy[$v['memUid']] ?? 0;
                } else {
                    $v['user_count'] = 0;
                }
            }
        }

        return $result;
    }
}
