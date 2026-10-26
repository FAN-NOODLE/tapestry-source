<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/11/26
 * Time: 16:24
 */

namespace App\Models\ScrmEnterprise;


class WxUserCorpExternal extends ScrmEnterpriseInit
{


    /**
     * 主键
     *
     * @var string
     */
    public $primaryKey = 'id';
    /**
     * 表名
     *
     * @var string
     */
    protected $table = 'wx_user_corp_external';
    /**
     * @var array
     */
    protected $guarded = [];


    /**
     * 潜在客户
     */
    const CARD_LEVEL_POTENTIAL=20;

    /**
     * 优质客户
     */
    const CARD_LEVEL_HIGH =10;

    /**
     * 切库根据外部联系人user_id查询
     *
     * @param $param
     * @param $workExternalUserIds
     *
     * @return array
     */
    public static function getExternalUserByIds($param, $workExternalUserIds): array
    {

        $query = static::query()
            ->setConn($param['enterprise_id'])
            ->select('wx_user_corp_external.id', 'wx_user_corp_external.mobile', 'wx_user_corp_external.openid', 'wx_user_corp_external.unionid', 'wx_user_corp_external.euid', 'wx_user_corp_external.nick', 'wx_user_corp_external.name', 'wx_user_corp_external.position', 'wx_user_corp_external.avatar', 'wx_user_corp_external.union_code', 'wx_user_corp_external.point', 'wx_user_corp_external.grade', 'wx_user_corp_external.store_code', 'wx_user_corp_external.ba_code', 'wx_user_corp_external.birthday', 'wx_user_corp_external.corp_name', 'wx_user_corp_external.corp_full_name', 'wx_user_corp_external.type', 'wx_user_corp_external.sex', 'wx_user_corp_external.external_profile', 'enterprise_user.mobile', 'enterprise_user.id as eu_id', 'enterprise_user.tags as eu_tags')
            ->leftJoin('enterprise_user',
                'wx_user_corp_external.euid',
                '=',
                'enterprise_user.id')
            ->whereIn('wx_user_corp_external.id', $workExternalUserIds)
            ->when(isset($param['gender']) && !empty($param['gender']), function ($query) use ($param) {
                return $query->where('wx_user_corp_external.sex', $param['gender']);
            })
            ->when(isset($param['search_key']) && !empty($param['search_key']), function ($query) use ($param) {

                $query->where(function ($query) use ($param) {

                    return $query->where('wx_user_corp_external.name', 'like', '%' . $param['search_key'] . '%')
                        ->orWhere('wx_user_corp_external.nick', 'like', '%' . $param['search_key'] . '%')
                        ->orWhere('enterprise_user.mobile', 'like', '%' . $param['search_key'] . '%');

                });

            });

        $count = $query->count();
        $list  = $query->forPage($param['page'], $param['size'])
            ->orderBy('wx_user_corp_external.created_at', 'desc')
            ->get()
            ->toArray();

        return [$list, $count];

    }


    /**
     * 切库获取外部联系人详情
     *
     * @param $enterpriseId
     * @param $workExternalUserIds
     *
     * @return array
     */
    public static function getExternalUserInfoById($enterpriseId, $workExternalUserId): array
    {

        return static::query()
            ->setConn($enterpriseId)
            ->where('id', $workExternalUserId)
            ->first()
            ->toArray();

    }

    /**
     * 根据openid获取外部联系人信息
     *
     * @param int $enterpriseId
     * @param array $openIds
     * @param array $fields
     *
     * @return array
     * @Date: 2021/6/29
     * @Time: 14:50
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public static function getExternalUserInfoByOpenId(int $enterpriseId, array $openIds, array $fields = ['*']): array
    {
        return static::query()
            ->setConn($enterpriseId)
            ->select($fields)
            ->whereIn('openid', $openIds)
            ->get()
            ->toArray();
    }

    /**
     * 根据主键id获取数据
     *
     * @param int $enterpriseId
     * @param array $ids
     * @param array $fields
     *
     * @return array
     * @Date: 2021/7/5
     * @Time: 20:38
     */
    public static function getExternalUserInfoByIds(int $enterpriseId, array $ids, array $fields = ['*']): array
    {
        return static::query()
            ->setConn($enterpriseId)
            ->select($fields)
            ->whereIn('id', $ids)
            ->get()
            ->toArray();
    }

    /**
     * 根据$unionId获取外部联系人
     *
     * @param int $enterpriseId
     * @param array $unionId
     * @param array $fields
     *
     * @return array
     * @Date: 2021/7/5
     * @Time: 21:49
     * @Author: 烽火连城 <zhaijunliang@vchangyi.com>
     */
    public static function getExternalUserInfoByUnionId(int $enterpriseId, array $unionId, array $fields = ['*']): array
    {
        return static::query()
            ->setConn($enterpriseId)
            ->select($fields)
            ->whereIn('unionid', $unionId)
            ->where('union_code', '<>', '')
            ->get()
            ->toArray();
    }

    public static function updateExternalUserInfoByUnionId(int $enterpriseId, string $unionId, array $updateData)
    {
        return static::query()
            ->setConn($enterpriseId)
            ->where('unionid', $unionId)
            ->update($updateData);
    }

    public static function getExternalUserInfo(int $enterpriseId, string $unionId)
    {
        return static::query()
            ->setConn($enterpriseId)
            ->where('unionid', $unionId)
            ->first();
    }

    /**
     * @param int $enterpriseId
     * @param string $unionId
     * @param array $fields
     * @return array
     */
    public static function getExternalUserInfoByUnionIdNoCode(int $enterpriseId, string $unionId, array $fields = ['*']): array
    {
        return static::query()
            ->setConn($enterpriseId)
            ->select($fields)
            ->where('unionid', $unionId)
            ->get()
            ->toArray();
    }

}

