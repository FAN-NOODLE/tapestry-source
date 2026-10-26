<?php
/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\App\V1\Coupon;

use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use CrsSdk\Shop\Coupon;
use App\Services\MemberUserRelationService;

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
     * 礼劵列表
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/11/25 11:16
     * @param Request $request
     * @return array
     */
    public function index(Request $request)
    {
        $params = $request->input();

        // 查询参数处理
        isset($params['status']) && $where['status'] = $params['status'];
        $where['page'] = isset($params['page']) ? $params['page'] : 1;
        $where['limit'] = isset($params['limit']) ? $params['limit'] : 10;
        $where['enterprise_id'] = $this->enterpriseId;
        $where['mem_uid'] = $this->memUid;

        $userRelation = $this->memberUserRelationService->userRelationByMemUid($this->enterpriseId, $this->memUid);
        $where['union_id'] = array_column($userRelation, 'union_id');
        try {
            $list = (new Coupon())->guide($where);
        } catch (\Exception $e) {
             error(710001);
        }
        if (isset($list['data']) && !empty($list['data'])) {
            $resultList = [];
            foreach ($list['data'] as $key => $value) {
                $resultList[$key]['id'] = $value['id'];
                $resultList[$key]['name'] = $value['name'];
                $resultList[$key]['condition_value'] = $value['condition_value'];
                $resultList[$key]['discount_amount'] = $value['discount_amount'];
                $resultList[$key]['start_time'] = $value['start_time'];
                $resultList[$key]['end_time'] = $value['end_time'];
                $resultList[$key]['status'] = $value['status'];
                $resultList[$key]['range'] = $value['range'];
                $resultList[$key]['receive_count'] = $value['receive_count'];
                $resultList[$key]['use_count'] = $value['use_count'];
                if (isset($value['categorys'])) {
                    $resultList[$key]['categorys'] = $value['categorys'];
                }
            }
            $list['data'] = $resultList;
        }
        return success($list);
    }
}
