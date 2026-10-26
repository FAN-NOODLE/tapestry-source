<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Models\Channel;
use App\Common\Constant;
use App\Library\Uc\Member;
use App\Models\ContactQrCode;
use Illuminate\Http\Request;
use App\Models\MemberUserRelation;
use App\Http\Controllers\Service\V1\BaseController;

class StaffUserDelController extends BaseController
{
    /**
     * 员工删客户-客户删员工
     */
    const DEL_RELATION = [
        MemberUserRelation::DEL_EXTERNAL_CONTACT,
        MemberUserRelation::DEL_FOLLOW_USER,
    ];

    /**
     * 添加类型 - 1 导购主动添加  2 支付即导购添加  3 社群添加
     */
    const SOURCE_TYPES = [
        MemberUserRelation::SOURCE_TAKE_SHOPPERS_ADD,
        MemberUserRelation::SOURCE_PAY_SHOPPERS_ADD,
        MemberUserRelation::SOURCE_COMMUNITY_ADD,
    ];

    /**
     * 删除客户接口
     * @param Request $request
     * @return array|void
     * @throws \Exception
     */
    public function staffUserDel(Request $request)
    {
        $conditions = [];
        $params     = $request->all();

        if (empty($params['enterprise_id'])) {
            notice(600000);
        }

        $page = $params['page'] ?? Constant::PAGE;
        $size = $params['size'] ?? Constant::SIZE;

        $qyDomain      = $params['qy_domain'];
        $appIdentifier = $params['app_identifier'];
        $memIds        = [];

        $conditions[] = ['enterprise_id', '=', $params['enterprise_id']];

        // 部门ID 和 员工ID 都在 - 只给这个员工的客户,不验证部门
        if (!empty($params['dep_uid']) && !empty($params['mem_uid'])) {
            $memIds[] = $params['mem_uid'];
        }

        // 只有部门ID - 找这个部门的所有员工的客户
        if (!empty($params['dep_uid']) && empty($params['mem_uid'])) {
            $isTrue  = true;
            $pageNum = 1;
            $sizeNum = 500;
            do {
                $result = (new Member())->ignorePermissionsList($qyDomain, $appIdentifier, [
                    'dpId'     => $params['dep_uid'],
                    'pageNum'  => $pageNum,
                    'pageSize' => $sizeNum,
                ]);

                if (isset($result['list']) && count($result['list']) > 0) {
                    foreach ($result['list'] as $item) {
                        $memIds[] = $item['memUid'];
                    }
                }

                if ($result['total'] < $pageNum * $sizeNum) {
                    $isTrue = false;
                } else {
                    $pageNum = $pageNum + 1;
                }
            } while ($isTrue);
        }

        // 只有员工ID - 只找这个员工的客户
        if (!empty($params['mem_uid']) && empty($params['dep_uid'])) {
            $memIds[] = $params['mem_uid'];
        }

        // 有开始结束时间段
        if (!empty($params['start_time']) && !empty($params['end_time'])) {
            $conditions[] = ['cancel_relation_at', '>=', $params['start_time']];
            $conditions[] = ['cancel_relation_at', '<=', $params['end_time']];
        }

        // 客户姓名
        if (!empty($params['name'])) {
            $conditions[] = ['ec_name', 'like', "%{$params['name']}%"];
        }

        // 渠道来源
        $sourceTypes = [];
        if (!empty($params['channel_id'])) {
            $channel = Channel::query()->where('id', '=', $params['channel_id'])->first();
            if (empty($channel)) {
                notice(610001);
            }
            $contactQrcode = ContactQrCode::query()->where('channel_id', $channel->id)->get()->toArray();

            $sourceTypes = array_column($contactQrcode, 'id');
        }

        // 删除类型 - 员工删 还是 客户删
        if (isset($params['type']) && !in_array($params['type'], self::DEL_RELATION)) {
            notice(900001);
        } else {
            $conditions[] = ['delete_type', '=', $params['type']];
        }

        // 以下三段代码为, 是为了把部门员工的条件加入进去
        $query = MemberUserRelation::query()
            ->with('channel')
            ->where($conditions)
            ->when(!empty($sourceTypes), function($query) use ($sourceTypes) {
                $query->whereIn('source_type', $sourceTypes);
            });

        if (!empty($memIds)) {
            $query = $query->whereIn('mem_uid', $memIds);
        }

        $members = $query->orderBy('cancel_relation_at', 'desc')
            ->paginate($size, ['*'], '', $page);

        // 最终的响应数据组装
        $data = [];
        $sourceType = array_column($members->items(), 'source_type');
        $channelByContact = ContactQrCode::query()->with('channel')->whereIn('id', $sourceType)->get()->toArray();
        $channelByContact = array_column($channelByContact, null, 'id');
        foreach ($members->items() as $member) {
            $memUser = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $member['mem_uid']]);
            $data[]  = [
                'id'         => $member['id'],
                'name'       => $member['ec_name'],
                'avatar'     => $member['ec_avatar'],
                'union_id'   => $member['union_id'],
                'user_id'    => $member['user_id'],
                'channel'    => $channelByContact[$member['source_type']]['channel']['channel'] ?? '',
                'staff_name' => $memUser['memUsername'] ?? '',
                'created_at' => date('Y-m-d H:i:s', $member['relation_at']),
                'deleted_at' => date('Y-m-d H:i:s', $member['cancel_relation_at']),
            ];
        }

        return success([
            'data' => $data,
            'meta' => [
                'page'  => (int)$page,
                'size'  => (int)$size,
                'total' => (int)$members->total(),
            ]
        ]);
    }
}