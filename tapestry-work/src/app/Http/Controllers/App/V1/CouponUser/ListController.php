<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\CouponUser;

use App\Services\MemberUserRelationService;
use CrsSdk\Shop\CouponUser;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;

class ListController extends BaseController
{
    /**
     * 标记外部联系Service
     *
     * @var MemberUserRelationService
     */
    protected $memberUserRelationService;

    /**
     * ListController constructor.
     * @param MemberUserRelationService $memberUserRelationService
     */
    public function __construct(MemberUserRelationService $memberUserRelationService)
    {
        parent::__construct();

        // 标记外部联系Service赋值
        $this->memberUserRelationService = $memberUserRelationService;
    }

    /**
     * 用户礼劵
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/25 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->input();

        if (!isset($params['ec_userid']) || empty($params['ec_userid'])) {
            error(200008);
        }

        $where['enterprise_id'] = $this->enterpriseId;
        $where['page'] = isset($params['page']) ? $params['page'] : 1;
        $where['limit'] = isset($params['limit']) ? $params['limit'] : 10;

        //是否是客户条件判断
        $userRelation = $this->memberUserRelationService->userRelationDetail($this->enterpriseId, $this->memUid, $params['ec_userid']);
        if (isset($userRelation->union_id)) {
            $where['union_id'] = $userRelation->union_id;
            try {
                $data = (new CouponUser())->lists($where);
            } catch (\Exception $e) {
                error(710004);
            }

            //列表处理
            array_walk($data['data'], function (&$value) use ($userRelation) {
                $value['ec_userid'] = $userRelation->ec_userid;
                $value['surplus_day'] = ceil(($value['end_time'] - time()) / 86400) ?? 0;
            });
        } else {
            $data = [
                'data' => [],
                'meta' => [
                    'total' => 0,
                    'size' => $where['page'],
                    'page' => $where['limit']
                ]
            ];
        }

        return success($data);
    }
}
