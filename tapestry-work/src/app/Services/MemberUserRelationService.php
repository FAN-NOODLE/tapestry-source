<?php
/**
 * Created by PhpStorm.
 * User: 大熊
 * Date: 2018/12/3
 * Time: 16:18
 */

namespace App\Services;

use App\Models\GuideUserRelation;
use App\Models\MemberGuide;
use App\Models\MemberUserRelation;
use App\User;
use Illuminate\Support\Facades\DB;

class MemberUserRelationService
{

    /**
     * 专属导购
     */
    const IS_EXCLUSIVE_TRUE = 1;

    /**
     * @var MemberUserRelation
     */
    protected $memberUserRelation;

    public function __construct(MemberUserRelation $memberUserRelation)
    {
        $this->memberUserRelation = $memberUserRelation;
    }

    /**
     * 获取导购某段时间内客户数
     *
     * @param int    $enterpriseId 企业ID
     * @param string $memUid       导购ID
     * @param int    $startTime    开始时间
     * @param int    $endTime      结束时间
     *
     * @return int
     */
    public function timeUserCount($enterpriseId, $memUid, $startTime = 0, $endTime = 0)
    {
        // 查询条件
        $condition = [
            ['enterprise_id', '=', $enterpriseId],
            ['mem_uid', '=', $memUid],
        ];

        if ($startTime) {
            $condition[] = ['created_at', '>=', $startTime];
        }

        if ($endTime) {
            $condition[] = ['created_at', '<=', $endTime];
        }

        $user = $this->memberUserRelation
            ->query()
            ->select(DB::raw('count(*) as count'))
            ->where($condition)
            ->first();

        return $user ? $user['count'] : 0;
    }

    /**
     * 获取导购的客户数
     * @param int $enterpriseId
     * @param array $memUids
     * @return array
     */
    public function timeUserCounts(int $enterpriseId, array $memUids) : array
    {
        $user = $this->memberUserRelation
            ->query()
            ->select(DB::raw('count(*) as count, mem_uid'))
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('mem_uid', $memUids)
            ->groupBy(['mem_uid'])
            ->pluck('count', 'mem_uid')
            ->toArray();

        return $user;
    }


    /**
     * 根据导购ID获取其所有客户信息
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/05 11:16
     * @param int $enterpriseId 企业ID
     * @param string $memUid 导购ID
     * @return mixed
     */
    public function userRelationByMemUid($enterpriseId, $memUid){
        $query = $this->memberUserRelation->query()->where('enterprise_id', $enterpriseId)->where('mem_uid', $memUid)->where('union_id', '<>', '');
        try {
            $userRelation = $query->get();
            return empty($userRelation) ? [] : $userRelation->toArray();
        } catch (\Exception $e) {
            error(110000);
        }
    }

    /**
     * 获取客户信息
     * @author 梁欢 lianghuan@vchangyi.com
     * @date 2018/12/05 11:16
     * @param int $enterpriseId 企业ID
     * @param string $memUid 导购ID
     * @param string $ecUserid 导uc外部联系人userid
     * @return mixed
     */
    public function userRelationDetail($enterpriseId, $memUid, $ecUserid){
        $query = $this->memberUserRelation->query()->where('enterprise_id', $enterpriseId)->where('mem_uid', $memUid)->where('ec_userid', $ecUserid);
        try {
            return $query->first();
        } catch (\Exception $e) {
            error(110000);
        }
    }


    /**
     * 周图表
     * @param $input
     * @return array
     */
    static public function weekChart($input)
    {
        if (isset($input['stores'])) {
            if ($input['type'] === 1) {
                // 自然流量
                $result = DB::select('SELECT FROM_UNIXTIME(created_at,\'%u\') xtime, count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            AND store_id = ""
                            GROUP BY xtime',
                    [
                        'is_exclusive' => self::IS_EXCLUSIVE_TRUE,
                        'enterprise_id' => $input['enterprise_id'],
                        'created_at' => $input['start'],
                        'deleted_at' => 0,
                    ]);
            } else if ($input['type'] === 2) {
                // 门店
                $result = DB::select('SELECT FROM_UNIXTIME(created_at,\'%u\') xtime, count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            AND store_id IN (:store_id)
                            GROUP BY xtime',
                    [
                        'is_exclusive' => self::IS_EXCLUSIVE_TRUE,
                        'enterprise_id' => $input['enterprise_id'],
                        'created_at' => $input['start'],
                        'deleted_at' => 0,
                        'store_id' => join(',', $input['stores'])
                    ]);
            } else {
                // 自然流量+门店
                $result = DB::select('SELECT FROM_UNIXTIME(created_at,\'%u\') xtime, count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            AND (store_id IN (:store_id) OR store_id = "")
                            GROUP BY xtime',
                    [
                        'is_exclusive' => self::IS_EXCLUSIVE_TRUE,
                        'enterprise_id' => $input['enterprise_id'],
                        'created_at' => $input['start'],
                        'deleted_at' => 0,
                        'store_id' => join(',', $input['stores'])
                    ]);
            }

            return $result;
        }

        return DB::select('SELECT FROM_UNIXTIME(created_at,\'%u\') xtime, count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            GROUP BY xtime',
            [
                'is_exclusive'  => self::IS_EXCLUSIVE_TRUE,
                'enterprise_id' => $input['enterprise_id'],
                'created_at'    => $input['start'],
                'deleted_at'    => 0
            ]);
    }

    /**
     * 月图表
     * @param $input
     * @return array
     */
    static public function monthChart($input)
    {
        if (isset($input['stores'])) {
            if ($input['type'] === 1) {
                // 自然流量
                $result = DB::select('SELECT FROM_UNIXTIME(created_at,\'%Y-%m-01\') xtime,count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            AND store_id = ""
                            GROUP BY xtime',
                    [
                        'is_exclusive' => self::IS_EXCLUSIVE_TRUE,
                        'enterprise_id' => $input['enterprise_id'],
                        'created_at' => $input['start'],
                        'deleted_at' => 0,
                    ]);
            } else if ($input['type'] === 2) {
                // 门店
                $result = DB::select('SELECT FROM_UNIXTIME(created_at,\'%Y-%m-01\') xtime,count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            AND store_id IN (:store_id)
                            GROUP BY xtime',
                    [
                        'is_exclusive' => self::IS_EXCLUSIVE_TRUE,
                        'enterprise_id' => $input['enterprise_id'],
                        'created_at' => $input['start'],
                        'deleted_at' => 0,
                        'store_id' => join(',', $input['stores'])
                    ]);
            } else {
                // 自然流量+门店
                $result = DB::select('SELECT FROM_UNIXTIME(created_at,\'%Y-%m-01\') xtime,count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            AND (store_id IN (:store_id) OR store_id = "")
                            GROUP BY xtime',
                    [
                        'is_exclusive' => self::IS_EXCLUSIVE_TRUE,
                        'enterprise_id' => $input['enterprise_id'],
                        'created_at' => $input['start'],
                        'deleted_at' => 0,
                        'store_id' => join(',', $input['stores'])
                    ]);
            }

            return $result;
        }

        return DB::select('SELECT FROM_UNIXTIME(created_at,\'%Y-%m-01\') xtime,count(distinct(mem_uid))
                            number FROM `member_user_relation` 
                            WHERE is_exclusive=:is_exclusive 
                            AND enterprise_id=:enterprise_id 
                            AND created_at>=:created_at 
                            AND deleted_at=:deleted_at 
                            GROUP BY xtime',
            [
                'is_exclusive'  => self::IS_EXCLUSIVE_TRUE,
                'enterprise_id' => $input['enterprise_id'],
                'created_at'    => $input['start'],
                'deleted_at'    => 0
            ]);
    }

    /**
     * 删除导购
     *
     * @param $params
     *          + enterprise_id
     *          + mem_uid()
     *
     * @return bool
     */
    public function deleteGuide($params)
    {

        // 删除导购信息
        MemberGuide::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->when(!empty($params['mem_uid']), function ($query) use ($params) {
                return is_array($params['mem_uid']) ? $query->whereIn('mem_uid', $params['mem_uid']) : $query->where('mem_uid', $params['mem_uid']);
            })
            ->delete();

        return true;
    }
}