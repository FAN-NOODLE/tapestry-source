<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models;

use App\Common\Constant;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class MemberUserRelation extends Base
{
    /**
     * 是否专属导购：是
     */
    const IS_EXCLUSIVE_TRUE = 1;
    /**
     * 是否专属导购：否
     */
    const IS_EXCLUSIVE_FALSE = 2;

    /**
     * 星标客户 1：是
     */
    const STARS_TRUE = 1;

    /**
     * 星标客户 2：否
     */
    const STARS_FALSE = 2;

    /**
     * 是否是绑定客户：是
     */
    const BINDING_CUSTOMER_YES = 1;
    /**
     * 是否是绑定客户：否
     */
    const BINDING_CUSTOMER_NO = 2;

    /**
     * 外部联系人性别 男
     */
    const EC_GENDER_MAN = 1;

    /**
     * 外部联系人性别 女
     */
    const EC_GENDER_WOMAN = 2;

    /**
     * 外部联系人性别 未知
     */
    const EC_GENDER_DEFAULT = 3;

    /**
     * 外部联系人性别 未知
     */
    const EC_AS_GENDER_DEFAULT = 0;

    /**
     * 未进行互相删除 (员工客户未互删)
     */
    const STAFF_USER_NOT_DEL = 0;

    /**
     * 未删除
     */
    const DEL_DEFAULT = 0;

    /**
     * 删除外部联系人 (员工删除客户)
     */
    const DEL_EXTERNAL_CONTACT = 1;

    /**
     * 删除跟进成员 (客户删除员工)
     */
    const DEL_FOLLOW_USER = 2;

    /**
     * 导购主动添加
     */
    const SOURCE_TAKE_SHOPPERS_ADD = 1;

    /**
     * 支付即导购添加
     */
    const SOURCE_PAY_SHOPPERS_ADD = 2;

    /**
     * 社群添加
     */
    const SOURCE_COMMUNITY_ADD = 3;

    // 外部联系人正常状态
    const NORMAL_STATUS = 0;

    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'member_user_relation';

    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';

    /**
     * @var array
     */
    protected $fillable = [
        'id',
        'enterprise_id',
        'store_id',
        'mem_uid',
        'ec_id',
        'ec_userid',
        'ec_name',
        'ec_avatar',
        'ec_gender',
        'user_remark_name',
        'user_name_index',
        'user_description',
        'union_id',
        'relation_at',
        'stars',
        'is_exclusive',
        'source_type',
        'user_id',
        'delete_type',
        'is_binding_customer',
        'cancel_relation_at',
        'union_code',
        'mobile',
        'grade',
        'point',
        'birthday'
    ];

    public function channel()
    {
        return $this->hasOne(Channel::class, 'id', 'source_type');
    }

    /**
     * getMemberListWithOrder
     * @desc 通过条件筛选外部联系人，需要关联订单数据
     * @param     $param
     * @param int $skip $skip = ($page - 1) * $size
     * @param int $limit $limit = $size
     *
     * @return array
     * @date 2020-07-15
     * @author XuHuitao
     */
    public static function getMemberListWithOrder($param, $skip = 0, $limit = 10, $forPage = false) {
        if (isset($param['page'], $param['size'])) {
            $skip = ($param['page'] - 1) * $param['size'];
            $limit = $param['size'];
        }

        $select = ['ec_name', 'ec_avatar', 'ec_userid', 'ec_gender', 'relation_at',
            'member_user_relation.user_id', 'mem_uid', 'mos.total_order_num',
            'mos.total_order_amount', 'member_user_relation.union_id', 'member_user_relation.source_type', 'ec_id',
            'member_user_relation.id as relation_id'];

        $query = MemberUserRelation::query()
            ->leftJoin('member_order_statistics as mos', function ($join) {
                $join->on([
                    ['member_user_relation.union_id', '=', 'mos.union_id'],
                    ['member_user_relation.union_id', '!=', DB::raw("''")]
                ]);
            })
            ->where('member_user_relation.enterprise_id', $param['enterprise_id'])
            ->where('mem_uid', $param['mem_uid'])
            ->where('delete_type', 0);
        if (isset($param['relation_id'])) {
            if (empty($param['relation_id'])) {
                return [];
            }
            $query = $query->whereIn('member_user_relation.id', $param['relation_id']);
        }
        if (isset($param['search_union_id']) || isset($param['search_relation_id'])) {
            if (empty($param['search_union_id']) && empty($param['search_relation_id'])) {
                return [];
            }
            $query = $query->where(function ($query) use ($param) {
                if (!empty($param['search_union_id'])) {
                    $query->whereIn('member_user_relation.union_id', $param['search_union_id']);
                }
                if (!empty($param['search_relation_id'])) {
                    $query->orWhereIn('member_user_relation.id', $param['search_relation_id']);
                }
            });
        }

        if (isset($param['tag_union_id'])) {
            $query = $query->whereIn('member_user_relation.union_id', $param['tag_union_id']);
        }
        if (isset($param['gender'])) {
            $query = $query->where('ec_gender', $param['gender']);
        }

        if ($param['add_type'] == Constant::ADD_TYPE_BUSINESS_CARD) {
            $query = $query->where('source_type', $param['qr_code_id']);
        } elseif ($param['add_type'] == Constant::ADD_TYPE_BUSINESS_CARD_OTHER) {
            $query = $query->where('source_type', '<>', $param['qr_code_id']);
        }

        if (isset($param['order_num_min'], $param['order_num_max'])) {
            if ($param['order_num_min'] == 0) {
                $query = $query->where(function ($query) use ($param){
                    $query->whereNull('mos.total_order_num');
                    $query->orWhereBetween('mos.total_order_num', [$param['order_num_min'], $param['order_num_max']]);
                });
            } else {
                $query = $query->whereBetween('mos.total_order_num', [$param['order_num_min'], $param['order_num_max']]);
            }
        }
        if (isset($param['order_amount_min'], $param['order_amount_max'])) {
            if ($param['order_amount_min'] == 0) {
                $query = $query->where(function ($query) use ($param){
                    $query->whereNull('mos.total_order_amount');
                    $query->orWhereBetween('mos.total_order_amount', [$param['order_amount_min'], $param['order_amount_max']]);
                });
            } else {
                $query = $query->whereBetween('mos.total_order_amount', [$param['order_amount_min'], $param['order_amount_max']]);
            }
        }
        if (isset($param['join_time_start'], $param['join_time_end'])) {
            $query = $query->whereBetween('relation_at', [$param['join_time_start'], $param['join_time_end']]);
        }
        if (isset($param['buy_time_start'], $param['buy_time_end'])) {
            $query = $query->whereBetween('mos.latest_buy_time', [$param['buy_time_start'], $param['buy_time_end']]);
        }
        $query = $query->select($select)->orderByDesc('member_user_relation.relation_at');
        if ($forPage) {
            $query = $query->paginate($limit);
        }else{
            $query = $query->get();
        }

        $data = $query->toArray();

        return $data;
    }

    /**
     * getCustomerNumOfGuide
     * @desc 查询当前导购关联的顾客数量
     * @param $memUid
     * @param $enterpriseId
     * @param $storeId
     *
     * @return int
     * @date 2020-07-15
     * @author XuHuitao
     */
    public static function getCustomerNumOfGuide($memUid, $enterpriseId, $storeId) {
        return self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->where('relation_at', '>', 0)
            ->count();
    }

    /**
     * getCustomerLatestGuide
     * @desc 查询顾客最新绑定的导购信息
     * @param $memUids
     * @param $unionIds
     *
     * @return \Illuminate\Database\Eloquent\Model|null|object|static
     * @date 2020-07-16
     * @author XuHuitao
     */
    public static function getCustomerLatestGuide($memUids, $unionIds) {
        $data = self::query()
            ->whereIn('mem_uid', $memUids)
            ->whereIn('union_id', $unionIds)
            ->select(['*'])
            ->get()->toArray();
        $latestGuide = [];
        foreach ($data as $item) {
            if (isset($latestGuide[$item['union_id']])) {
                if ($item['relation_at'] > $latestGuide[$item['union_id']]['relation_at']) {
                    $latestGuide[$item['union_id']] = $item;
                }
            } else {
                $latestGuide[$item['union_id']] = $item;
            }
        }
        return $latestGuide;
    }

    /**
     * searchCustomerByName
     * @desc 通过姓名搜索顾客
     * @param $enterpriseId
     * @param $searchKey
     * @param $fields
     *
     * @return \Illuminate\Database\Eloquent\Collection|static[]
     * @date 2020-07-26
     */
    public static function searchCustomerByName($enterpriseId, $searchKey, $fields = ['*']) {
        return self::query()->where('enterprise_id', $enterpriseId)
            ->where('ec_name', 'like', '%' . $searchKey . '%')
            ->select($fields)
            ->get();
    }

    /**
     * getCustomerUnionidsOfMember
     * @desc 获取导购关联的所有顾客的unionid列表
     * @param $enterpriseId
     * @param $memUid
     *
     * @return array
     * @date 2020-07-27
     */
    public static function getCustomerUnionidsOfMember($enterpriseId, $memUid) {
        $result = self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('mem_uid', $memUid)
            ->select(['union_id'])
            ->get()
            ->toArray();
        return array_column($result, 'union_id');
    }

    /**
     * getDataById
     * @desc 通过id查询导购顾客关联关系
     * @param $relationId
     *
     * @return \Illuminate\Database\Eloquent\Model|null|object|static
     * @date 2020-08-06
     */
    public static function getDataById($relationId)
    {
        return self::query()
            ->where('id', $relationId)
            ->first();
    }


    /**
     * 根据主键id查询客户信息
     *
     * @param $id
     *
     * @return array|\Illuminate\Database\Eloquent\Model|null|object|static
     * @author  baizhongyang <baizhongyang@vchangyi.com>
     * @date    2020/8/6 16:06
     */
    public static function getCustomerById($id)
    {
        $result = self::query()
            ->where('id', $id)
            ->first();
        if (!is_null($result)) {
            $result = $result->toArray();
        } else {
            $result = [];
        }

        return $result;
    }

    /**
     * statisticCustomerAndMember
     * @desc
     * @param $enterpriseId
     * @param $codeIdList
     * @param $timeStart
     * @param $timeEnd
     *
     * @return array
     * @date 2020-08-14
     */
    public static function statisticCustomerAndMemberCount($enterpriseId, $codeIdList, $timeStart, $timeEnd, $forAll)
    {
        $query = self::query()
            ->where('member_user_relation.enterprise_id', $enterpriseId);
        if ($forAll) {
            $query = $query->where('source_type', '>', 1);
        } else {
            $query = $query->whereIn('source_type', $codeIdList);
        }
        $customerNum = $query
            ->where('delete_type', 0)
            ->when(!empty($timeStart)  && !empty($timeEnd), function ($query) use ($timeStart,$timeEnd) {
                return  $query->whereBetween('member_user_relation.created_at', [$timeStart, $timeEnd]);
            })
            ->count();
        return $customerNum;
    }


    /**
     * 查询顾客列表
     * @param int   $enterpriseId
     * @param array $codeIdList
     * @param int   $timeStart
     * @param int   $timeEnd
     * @param bool  $forAll
     *
     * @return array
     */
    public static function statisticCustomerList($enterpriseId=0, $memUidList,$timeStart=0, $timeEnd=0)
    {
        $query = self::query()
            ->where('member_user_relation.enterprise_id', $enterpriseId);
        $customerNum = $query
            ->where('delete_type', 0)
            ->whereIn('mem_uid', $memUidList)
            ->when(!empty($timeEnd), function ($query) use ($timeStart,$timeEnd) {
                return  $query->whereBetween('member_user_relation.created_at', [$timeStart, $timeEnd]);
            })
            ->select(['union_id','source_type','created_at as created_date'])
            ->get()->toArray();

        return $customerNum;
    }


    /**
     * statisticCustomerDaily
     * @desc 按天统计顾客数据
     * @param $enterpriseId
     * @param $codeIdList
     * @param $timeStart
     * @param $timeEnd
     *
     * @return array
     * @date 2020-08-14
     */
    public static function statisticCustomerDaily($enterpriseId, $codeIdList, $timeStart, $timeEnd, $forAll)
    {
        $query = self::query()
            ->where('member_user_relation.enterprise_id', $enterpriseId);
        if ($forAll) {
            $query = $query->where('source_type', '>', 1);
        } else {
            $query = $query->whereIn('source_type', $codeIdList);
        }
        return $query->whereBetween('member_user_relation.created_at', [$timeStart, $timeEnd])
            ->where('delete_type', 0)
            ->select(['created_at as created_date','union_id'])
            ->get()
            ->toArray();
    }


    /**
     * statisticCustomerByGuide
     * @desc 获取导购顾客数据
     * @param $enterpriseId
     * @param $codeIdList
     * @param $timeStart
     * @param $timeEnd
     *
     * @return array
     * @date 2020-08-22
     */
    public static function statisticCustomerByGuide($enterpriseId, $memUidList, $timeStart, $timeEnd)
    {
        $timeData = self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('source_type', '>',0)
            ->where('delete_type', 0)
            ->whereIn('mem_uid', $memUidList)
            ->whereBetween('created_at', [$timeStart, $timeEnd])
            ->select(['source_type','mem_uid'])
            ->get()
            ->toArray();
        $returnTimeData = [];
        foreach ($timeData as $item) {
            $mem_uid = $item['mem_uid'];
            $returnTimeData[$mem_uid][] = [
                'source_type'       => $item['source_type'],
                'mem_uid'           => $item['mem_uid']
            ];
        }

        $returnTotalData = [];
        $totalData = self::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('source_type', '>',0)
            ->whereIn('mem_uid', $memUidList)
            ->where('delete_type', 0)
            ->select(['source_type','mem_uid'])
            ->get()
            ->toArray();

        foreach ($totalData as $item) {
            $mem_uid = $item['mem_uid'];
            $returnTotalData[$mem_uid][] = [
                'source_type'       => $item['source_type'],
                'mem_uid'           => $item['mem_uid']
            ];
        }
        return [
            'time_data'  => $returnTimeData,
            'total_data' => $returnTotalData,
        ];
    }
}
