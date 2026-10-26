<?php

/**
 * Created by PhpStorm.
 * User: wanghuan
 * Date: 2018/11/27
 * Time: 21:34
 */

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\MemberUserRelation;
use CrsSdk\Shop\User;
use CrsSdk\Shop\UserLevel;
use App\Library\Uc\Member;
use Illuminate\Http\Request;

class ListController extends BaseController
{
    /**
     * 客户列表
     *
     * @param Request $request
     *
     * @return array|mixed
     * @throws \Exception
     */
    public function index(Request $request)
    {
        // 分页参数
        $page = $request->get('page', Constant::DEFAULT_PAGE);
        $size = $request->get('size', Constant::DEFAULT_SIZE);
        // 查询条件
        $search = $request->get('search', []);

        // 企业id
        $enterpriseId = $request->get('enterprise_id', 0);
        if (!$enterpriseId) {
            error(600000);
        }

        // 企业domain
        $qyDomain = $request->get('qy_domain', '');
        if (!$qyDomain) {
            error(600003);
        }

        // 应用标识
        $appIdentifier = $request->get('app_identifier', '');
        if (!$appIdentifier) {
            error(600004);
        }

        // 关键字
        $keywords = '';
        if (isset($search['keywords'])) {
            $keywords = $search['keywords'];
        }

        // 调用sdk获取用户数据
        $params   = [
            'enterprise_id' => $enterpriseId,
            'search'        => $search
        ];
        $userList = (new User())->users($params);
        $userList = $userList['data'];

        // 用户union_id集合
        $unionIds = [];
        if ('' != $keywords) {
            $unionIds = array_filter(array_column($userList, 'union_id'));
        }

        // 查询条件
        $where = [
            'enterprise_id' => $enterpriseId,
            'is_exclusive'  => MemberUserRelation::IS_EXCLUSIVE_TRUE // 专属导购
        ];
        // 总数据量
        $total = MemberUserRelation::where($where)
            ->when('' != $keywords, function ($query) use ($unionIds) {
                return $query->whereIn('union_id', $unionIds);
            })->count();

        // 待返数据
        $data = [];
        if ($total) {
            $relationList = MemberUserRelation::where($where)
                ->when('' != $keywords, function ($query) use ($unionIds) {
                    return $query->whereIn('union_id', $unionIds);
                })->forpage($page, $size)
                ->orderBy('id', 'DESC')
                ->get()
                ->toArray();

            if ($relationList) {
                // 数据格式化
                $data = $this->format([
                    'enterprise_id'  => $enterpriseId,
                    'qy_domain'      => $qyDomain,
                    'app_identifier' => $appIdentifier,
                    'relation_list'  => $relationList,
                    'user_list'      => $userList
                ]);
            }
        }

        return success([
            'data' => $data,
            'meta' => [
                'total' => intval($total),
                'page'  => intval($page),
                'size'  => intval($size)
            ]
        ]);
    }

    /**
     * 列表数据格式化
     *
     * @param array $params 参数列表
     *                      + int enterprise_id 企业id
     *                      + string qy_domain 企业domain
     *                      + string app_identifier 应用标识
     *                      + array relation_list 员工标记外部联系人列表
     *                      + array user_list 用户列表
     *
     * @return array
     * @throws \Exception
     */
    protected function format($params)
    {
        // 初始化返回数据
        $data = [];

        // 企业id
        $enterpriseId = $params['enterprise_id'];
        // 企业domain
        $qyDomain = $params['qy_domain'];
        // 应用标识
        $appIdentifier = $params['app_identifier'];

        // 用户列表
        $userList = $params['user_list'];
        $userList = array_column($userList, null, 'union_id');

        // 员工标记外部联系人列表
        $relationList = $params['relation_list'];
        // 导购mem_uid集合
        $memUids = array_unique(array_column($relationList, 'mem_uid'));
        sort($memUids);

        // 参数列表
        $params = [
            'memUids'  => $memUids,
            'pageNum'  => Constant::DEFAULT_PAGE,
            'pageSize' => count($memUids)
        ];
        // 调用uc sdk获取导购信息
        $members = (new Member())->list($qyDomain, $appIdentifier, $params);
        $members = array_column($members['list'], null, 'memUid');

        // 会员等级
        $userLevel = (new UserLevel())->userLevel(['enterprise_id' => $enterpriseId]);
        $userLevel = array_column($userLevel['data'], null, 'id');

        // 组装返回数据
        foreach ($relationList as $relation) {
            // 用户信息
            $user = $userList[$relation['union_id']] ?? [];
            // 专属导购id
            $mem_uid = $relation['mem_uid'];

            // 会员名称
            $levelName = '';
            if ($user) {
                $levelName = $userLevel[$user['level_id']]['name'] ?? '';
            }

            // 返回数据赋值
            $data[] = [
                'ec_id'          => $relation['ec_id'],
                'user_id'        => $user['id'] ?? 0,
                'nick_name'      => $user['nick_name'] ?? $relation['ec_name'],
                'real_name'      => $user['real_name'] ?? '',
                'avatar'         => $user['avatar'] ?? $relation['ec_avatar'],
                'level_name'     => $levelName,
                'birthday'       => $user['birthday'] ?? '',
                'mobile'         => $user['mobile'] ?? '',
                'sales_username' => $members[$mem_uid]['memUsername'] ?? '',
                'sales_dp_name'  => $members[$mem_uid]['dpName'][0]['dpName'] ?? ''
            ];
        }

        return $data;
    }

    /**
     * @notes 说明
     * @params ec_id [array] [外部联系Id列表]
     * @return array
     */
    public function userIdList(Request $request)
    {
        $ec_ids = $request->post('ec_id', []);

        if (empty($ec_ids)) {
            error(600008);
        }

        $relationList = MemberUserRelation::whereIn('ec_id', $ec_ids)->groupBy('ec_id', 'user_id')->orderBy('id')->get(['ec_id', 'union_id', 'id'])->toArray();

        return success(['data' => $relationList]);
    }
}
