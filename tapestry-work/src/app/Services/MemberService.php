<?php

namespace App\Services;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Library\Uc\Department;
use App\Library\Uc\Member;
use App\Micros\Common;
use App\Models\MemberUserRelation;
use App\Models\Scrm\Departments;
use App\Models\Scrm\Members;
use CrsSdk\Boss\Enterprise;
use Exception;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;

class MemberService
{

    // UC查询删除人员状态标识
    const MEMBER_DELETE_REMARK = 3;

    // 部门类型
    const IS_DEPART = 1;

    // UC列表数据较大时,多次查询每次返回最大记录数
    const PAGE_SIZE = 500;

    /**
     * 获取企业微信mediaId
     *
     * @param string $corpId
     * @param string $appIdentifier
     * @param string $imageUrl
     *
     * @return mixed
     * @throws Exception
     */
    public static function getWorkWeChartMediaId($corpId = '', $appIdentifier = '', $imageUrl = '')
    {
        if (!$imageUrl) {
            throw new Error(10000, '图片地址缺失!');
        }

        $path = base_path('/resources/upload/tmp');

        if (!is_dir($path)) {
            mkdir($path, 0777, true);
        }

        $fileName = $path . '/' . md5(time() . rand(100000, 999999)) . '.jpg';

        $content = file_get_contents($imageUrl);

        if (false === file_put_contents($fileName, $content)) {
            throw new Error(20000, '本地图片存储失败!');
        }

        // 文件resource(此处如果直接能打开网络图片地址,则不需要本地储存)
        $file = fopen($fileName, 'r');

        $multipart = [
            [
                'name'     => 'file',
                'contents' => $file,
            ],
            [
                'name'     => 'fileType',
                'contents' => 1, // 图片
            ],
            [
                'name'     => 'filename',
                'contents' => rand(10000, 99999),
            ],
        ];

        $result = (new Member())->imageToMediaId($corpId, $appIdentifier, $multipart);

        Log::debug('uc-media-id-result' . var_export($result, true));

        if (empty($result['mediaId'])) {
            throw new Error(20000, '企业微信mediaId获取失败!');
        }
        unlink($fileName);
        return $result['mediaId'];
    }

    /**
     * 列表转数
     *
     * @param array $list
     * @param string $pk
     * @param string $pid
     * @param string $child
     * @param int $root
     *
     * @return array
     */
    public static function listToTree(array $list = [], $pk = 'id', $pid = 'pid', $child = 'children', $root = '')
    {
        $refer = [];
        foreach ($list as $key => $val) {
            $refer[$val[$pk]] = &$list[$key];
        }

        $tree = [];
        foreach ($list as $key => $val) {
            // level等于1是顶级
            if ($val['lev'] == 1 && self::IS_DEPART == $val['type']) {
                $tree[] = &$list[$key];
            }
            if (isset($refer[$val[$pid]])) {
                $refer[$val[$pid]][$child][] = &$list[$key];
            }
        }

        return $tree;
    }

    /**
     * 员工列表
     *
     * @param $domain
     * @param $identifier
     * @param $where
     *
     * @return array
     * @throws Exception
     */
    public static function memberList($domain, $identifier, $where = [])
    {
        $page       = 1;
        $maxPages   = 0;
        $memberList = [];
        do {
            // 参数初始化
            $where['pageNum']  = $page;
            $where['pageSize'] = self::PAGE_SIZE;

            // 结果处理
            $result     = (new Member())->ignorePermissionsList($domain, $identifier, $where);
            $memberList = array_merge($memberList, $result['list']);
            $maxPages   = $maxPages ?: ceil($result['total'] / self::PAGE_SIZE);

            // 步程自增
            $page++;
        } while ($page <= $maxPages);

        return $memberList;
    }

    /**
     * 部门列表
     *
     * @param $domain
     * @param $identifier
     * @param $where
     *
     * @return array
     * @throws Exception
     */
    public static function departList($domain, $identifier, $where = [])
    {
        $result = (new Department())->list($domain, $identifier, $where);

        return $result['list'] ?? [];
    }

    /**
     * 部门所有列表
     *
     * @param $domain
     * @param $identifier
     * @param $where
     *
     * @return array
     * @throws Exception
     */
    public static function allDepList($domain, $identifier, $where = [])
    {
        $page     = 1;
        $maxPages = 0;
        $list     = [];
        do {
            // 参数初始化
            $where['pageNum']  = $page;
            $where['pageSize'] = self::PAGE_SIZE;

            // 结果处理
            $result   = (new Department())->list($domain, $identifier, $where);
            $list     = array_merge($list, $result['list']);
            $maxPages = $maxPages ?: ceil($result['total'] / self::PAGE_SIZE);

            // 步程自增
            $page++;
        } while ($page <= $maxPages);

        return $list;
    }

    /**
     * 无部门员工列表
     *
     * @param $domain
     * @param $identifier
     * @param $where
     *
     * @return array
     * @throws Exception
     */
    public static function noDpMemberList($domain, $identifier, $where = [])
    {
        $page       = 1;
        $maxPages   = 0;
        $memberList = [];
        do {
            // 参数初始化
            $where['pageNum']  = $page;
            $where['pageSize'] = self::PAGE_SIZE;

            // 结果处理
            $result     = (new Member())->noDpList($domain, $identifier, $where);
            $memberList = array_merge($memberList, $result['list']);
            $maxPages   = $maxPages ?: ceil($result['total'] / self::PAGE_SIZE);

            // 步程自增
            $page++;
        } while ($page <= $maxPages);

        return $memberList;
    }

    public static function memberListNoAuth($qyDomain = '', $appIdentifier = '', $params = [], $enterpriseId = 0)
    {
        $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, $params);
        if (isset($result['list']) && $enterpriseId) {
            foreach ($result['list'] as &$item) {
                if (!is_null($item['memSubscribeTime'])) {
                    $item['user_count'] = (new MemberUserRelationService(new MemberUserRelation()))->timeUserCount($enterpriseId, $item['memUid']);
                } else {
                    $item['user_count'] = 0;
                }
            }
        }
        return $result;
    }

    /**
     * 获取members 导购列表
     *
     * @param $params
     *
     * @return array
     * @throws Exception
     * @Time: 2022/10/19
     */
    public static function getMemberList($params): array
    {
        $page = $params['page'] ?? Constant::PAGE;
        $size = $params['size'] ?? Constant::SIZE;

        $where = [];
        if (!empty($params['member_name'])) {
            $where[] = ['members.name', 'like', '%' . $params['member_name'] . '%'];
        }
        if (!empty($params['mobile'])) {
            $where[] = ['members.mobile', 'like', '%' . $params['mobile'] . '%'];
        }
        if (!empty($params['worker_status'])) {
            $where[] = $params['worker_status'] == 1 ?
                ['members.status', '=', Members::IN_SERVICE] :
                ['members.status', '<>', Members::IN_SERVICE];
        }
        if (!empty($params['type'])) {
            $where[] = ['type', '=', $params['type']];
        }
        if (!empty($params['dp_third_id'])) {
            $where[] = ['d.third_id', '=', $params['dp_third_id']];
        }
        if (!empty($params['worker_no'])) {
            $where[] = ['members.userid', '=', $params['worker_no']];
        }
        $sql = Members::query()
            ->leftJoin('member_departments as md', 'members.id', 'md.member_id')
            ->leftJoin('departments as d', 'md.department_id', 'd.id')
            ->where($where)
            ->where(['members.enterprise_id' => $params['enterprise_id'],
                     'md.deleted_at'         => 0,
                     'd.deleted_at'          => 0,]);
//            ->when(!empty($params['worker_status']) && $params['worker_status'] == 1, function ($query) use ($params) {
//                return $query->where('members.status', 1);
//            })
//            ->when(!empty($params['worker_status']) && $params['worker_status'] == 2, function ($query) use ($params) {
//                return $query->where('members.status', '<>', 1);
//            });

        $total = $sql->orderBy('members.created_at', 'desc')
            ->distinct()->count('members.id');

        if ($page && $size) {
            $sql = $sql->forPage($page, $size);
        }

        $memberList = $sql
            ->groupBy(['members.id'])
            ->get(['members.id', 'members.userid as mem_userid', 'members.name as mem_name', 'members.mobile as mem_mobile', 'members.position', 'members.avatar as mem_avatar', 'members.thumb_avatar', 'members.qr_poster_code as qr_code', 'members.status as worker_status', 'members.type', 'members.created_at',
                   'd.id as department_id', 'd.third_id',
                   DB::raw("GROUP_CONCAT(DISTINCT `d`.`name` ) dp_name"),
                   DB::raw("GROUP_CONCAT(DISTINCT `d`.`third_id` ) dp_third_ids"),])
            ->toArray();
        //获取列表管理范围
        $memberList = self::getRangeName($memberList, $params['enterprise_id']);
        return [
            'list'  => $memberList,
            'total' => $total,
            'page'  => (int)$page,
            'size'  => (int)$size
        ];

    }

    /**
     * 获取导购列表管理范围
     *
     * @param $memberList
     * @param $enterpriseId
     *
     * @return mixed
     * @throws Exception
     */
    private static function getRangeName($memberList, $enterpriseId)
    {
        //获取列表所有人的管理范围
        foreach ($memberList as &$val) {

            $val['position']   = '店员';
            $val['range_info'] = [];
            if ($val['type'] == Members::STORE_SM) {
                $val['position'] = '店长';
                $mobiles[]       = $val['mem_mobile'];
            }
        }
        if (empty($mobiles)) {
            return $memberList;
        }
        $rangeList      = Common::authority()->setApi('micro-authority/v1/account/range-list-by-mobiles')
            ->post(['enterprise_id' => $enterpriseId, 'mobiles' => $mobiles]);
        $allRangesDpIds = [];
        foreach ($rangeList as $dpIds) {
            $allRangesDpIds = array_merge($allRangesDpIds, $dpIds);
        }
        $dpNameByIdArr = Departments::query()->where('enterprise_id', $enterpriseId)
            ->whereIn('third_id', array_unique($allRangesDpIds))
            ->get(['third_id', 'name'])
            ->keyBy('third_id')
            ->toArray();
        foreach ($memberList as &$val) {
            $rangeDpNameArr = [];
            if (isset($rangeList[$val['mem_mobile']])) {
                foreach ($rangeList[$val['mem_mobile']] as $dpId) {
                    if (isset($dpNameByIdArr[$dpId])) {
                        $rangeDpNameArr[] = $dpNameByIdArr[$dpId]['name'];
                    }
                }
                $val['range_info'] = $rangeDpNameArr;
            }

        }
        return $memberList;
    }

    /**
     * 根据用户uid获取用户所在部门的所有人员列表
     *
     * @param string $qyDomain
     * @param string $appIdentifier
     * @param string $memUid 用户uid
     *
     * @return array|void
     * @throws Exception
     */
    public function getMemberListByUid($qyDomain = '', $appIdentifier = '', $memUid = '')
    {
        // 获取当前用户的详情信息
        $ucResult = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $memUid]);

        if (!isset($ucResult)) {
            return error(60103);
        }

        $dpId = '';
        if (isset($ucResult['dpName'])) {
            // 取出所属部门id
            $dpId = $ucResult['dpName'][0]['dpId'];
        }

        if (empty($dpId)) {
            // 部门信息不存在
            return error(600102);
        }

        $where = [
            'dpId'   => $dpId,
            'memAll' => 1
        ];

        // 获取人员列表
        $list = $this->listAllMember($qyDomain, $appIdentifier, $where);

        return $list;
    }

    /**
     * 根据条件获取人员列表
     *
     * @param string $qyDomain
     * @param string $appIdentifier
     * @param array $condition
     *
     * @return array
     * @throws Exception
     */
    public function listAllMember($qyDomain = '', $appIdentifier = '', $condition = [])
    {

        $page     = 1;
        $page_max = 0;
        $list     = [];

        // 获取所有成员
        do {
            $condition['pageNum']  = $page;
            $condition['pageSize'] = 50;

            $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, $condition);

            if (isset($result['list']) && !empty($result['list'])) {
                $list = array_merge($list, $result['list']);
            }

            // 计算总页数
            if ($page_max === 0) {
                $page_max = ceil($result['total'] / 50);
            }

            $page++;
        } while ($page <= $page_max);

        return $list;
    }

    /**
     * 获取外部联系人呼起数据
     *
     * @param string $memUid
     * @param string $unionId
     * @param int $enterpriseId
     *
     * @return mixed|string
     * @throws Exception
     */
    public function getUserEcUserId($memUid = '', $unionId = '', $enterpriseId = 0)
    {

        //根据导购二维码id统计外部联系人数
        $memberUserInfo = MemberUserRelation::query()
            ->where('mem_uid', $memUid)
            ->where('union_id', $unionId)
            ->where('delete_type', '=', 0)
            ->orderBy('created_at', 'desc')
            ->first(['union_id', 'ec_userid']);

        $enterpriseIdentifierList = EnterpriseService::getEnterpriseIdentifierList($enterpriseId);

        $identifierList = array_column($enterpriseIdentifierList, 'identifier');

        if (in_array(Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT, $identifierList)) {
            $identifier = Constant::ENTERPRISE_IDENTIFIER_STORE_ASSISTANT;
            if (env('APP_ENV') == 'dev') {
                $identifier = 'crs_store_assistant_dev';
            } else {
                $list = array_column($enterpriseIdentifierList, null, 'identifier');
                if ($list[$identifier]['created_at'] > env('ICLK_APP_CREATED_AT')) {
                    $identifier = 'iclk_crs_store_assistant';
                }
            }
            try {
                $enterpriseInfo = (new Enterprise())->detail(['enterprise_id' => $enterpriseId])['data'];

                $pltUserid = CustomerService::ecUserIdToPltUserId($memberUserInfo['ec_userid'], $enterpriseInfo['domain'], $identifier);

                $memberUserInfo['ec_userid'] = $pltUserid ?? $memberUserInfo['ec_userid'];
            } catch (Exception $exception) {
            }
        }

        return $memberUserInfo['ec_userid'] ?? '';
    }

}
