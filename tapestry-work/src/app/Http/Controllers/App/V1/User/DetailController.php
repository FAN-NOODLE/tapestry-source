<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/12/3
 * Time: 17:57
 */

namespace App\Http\Controllers\App\V1\User;

use App\Http\Controllers\App\V1\BaseController;
use App\Library\Uc\External;
use App\Models\MemberUserRelation;
use App\Models\UserGroup;
use App\Models\UserGroupRelation;
use CrsSdk\Shop\User;

class DetailController extends BaseController
{
    /**
     * @var MemberUserRelation
     */
    protected $memberUserRelation;

    /**
     * @var UserGroup
     */
    protected $userGroup;

    /**
     * @var UserGroupRelation
     */
    protected $userGroupRelation;

    /**
     * DetailController constructor.
     *
     * @param MemberUserRelation $memberUserRelation
     * @param UserGroup $userGroup
     * @param UserGroupRelation $userGroupRelation
     */
    public function __construct(
        MemberUserRelation $memberUserRelation,
        UserGroup $userGroup,
        UserGroupRelation $userGroupRelation
    )
    {
        parent::__construct();
        $this->memberUserRelation = $memberUserRelation;
        $this->userGroup          = $userGroup;
        $this->userGroupRelation  = $userGroupRelation;
    }

    /**
     * 客户详情
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $params = $this->validates();

        // 获取关系表信息
        $userRelation = $this->memberUserRelation->query()->where([
            'enterprise_id' => $this->enterpriseId,
            'ec_userid' => $params['ec_userid'],
            'mem_uid' => $this->memUid
        ])->first();
        if (!$userRelation) {
            error(200013);
        }

        // 客户基本信息
        if ($userRelation->union_id) {
            $userDetail           = app(User::class)->detail([
                'enterprise_id' => $this->enterpriseId,
                'union_id' => $userRelation->union_id
            ])['data'];
            $userDetail['is_vip'] = true;
        }
        if (!isset($userDetail['id'])) {
            try {
                $userDetail = (new External())->detail($this->qyDomain, ['ecId' => $userRelation->ec_id]);
            } catch (\Exception $exception) {
                $userDetail = [];
            }
            $userDetail = [
                'nick_name' => isset($userDetail['ecName']) ? $userDetail['ecName'] : '',
                'avatar' => isset($userDetail['ecAvatar']) ? $userDetail['ecAvatar'] : '',
                'is_vip' => false
            ];
        }

        $userDetail['stars']            = $userRelation->stars;
        $userDetail['user_remark_name'] = $userRelation->user_remark_name ?? '';
        // 客户所属分组
        $groups = $this->getUserGroup($params['ec_userid']);

        // 客户消费统计
        $count = [
            'year_spend' => 0,
            'year_mean_spend' => 0,
            'year_mean_spend_day' => '',
            'year_category_name' => '',
            'coupon_number' => 0
        ];
        if ($userDetail['is_vip'] && isset($userDetail['id'])) {
            $count = (new User())->count([
                'enterprise_id' => $this->enterpriseId,
                'user_id' => $userDetail['id'],
                'coupon' => 1
            ]);
            $count = [
                'year_spend' => $count['data']['yearSpend'],
                'year_mean_spend' => $count['data']['yearMeanSpend'],
                'year_mean_spend_day' => $count['data']['yearMeanSpendDay'],
                'year_category_name' => $count['data']['yearCategoryName'],
                'coupon_number' => $count['data']['couponNumber'],
            ];
        }

        return success(['data' => [
            'basic' => $userDetail,
            'groups' => $groups,
            'count' => $count
        ]]);
    }

    /**
     * 获取客户所在导购的分组
     *
     * @param $ecUserId
     *
     * @return array
     */
    private function getUserGroup($ecUserId)
    {
        $groups = [];

        // 客户所属所有分组
        $userGroupRelationData = $this->userGroupRelation
            ->query()
            ->where([
                'enterprise_id' => $this->enterpriseId,
                'ec_userid' => $ecUserId,
            ])
            ->get();
        if (!$userGroupRelationData->first()) {
            return $groups;
        }

        // 客户所属所有分组ID集合
        $groupIds = array_column($userGroupRelationData->toArray(), 'group_id');

        // 用户在该导购下分组
        $groups = $this->userGroup
            ->query()
            ->where([
                'enterprise_id' => $this->enterpriseId,
                'mem_uid' => $this->memUid
            ])
            ->whereIn('id', $groupIds)
            ->get()
            ->toArray();

        return $groups;
    }

    /**
     * 数据基础验证
     *
     * @return array|string|null
     * @throws \Illuminate\Validation\ValidationException
     */
    private function validates()
    {
        $rules = [
            'ec_userid' => 'required'
        ];

        $message = [
            'ec_userid.required' => '客户ID不存在'
        ];

        $this->validate(request(), $rules, $message);

        $params = request()->input();

        return $params;
    }
}
