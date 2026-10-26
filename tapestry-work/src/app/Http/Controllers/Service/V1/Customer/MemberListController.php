<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Library\Uc\Member;
use App\Library\Uc\Tag;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\CustomerGroup;
use App\Models\CustomerGroupRelation;
use App\Models\MemberUserRelation;
use Illuminate\Http\Request;
use App\Http\Controllers\Service\V1\BaseController;
use Illuminate\Support\Facades\Log;

class MemberListController extends BaseController
{
    /**
     * 客户和员工关系正常
     */
    const MS_NORMAL = 0;

    /**
     * 包含 客户删除员工 的客户
     */
    const CONDITION_MDS = 2;

    /**
     * 性别集合
     */
    const GENDERS = [
        MemberUserRelation::EC_GENDER_MAN,
        MemberUserRelation::EC_GENDER_WOMAN,
        MemberUserRelation::EC_GENDER_DEFAULT,
        MemberUserRelation::EC_AS_GENDER_DEFAULT,
    ];

    /**
     * 性别映射
     */
    const GENDER = [
        MemberUserRelation::EC_GENDER_MAN        => '男',
        MemberUserRelation::EC_GENDER_WOMAN      => '女',
        MemberUserRelation::EC_GENDER_DEFAULT    => '未知',
        MemberUserRelation::EC_AS_GENDER_DEFAULT => '未知',
    ];

    public function memberList(Request $request)
    {
        $params = $request->all();

        if (empty($params['enterprise_id'])) {
            notice(600000);
        }

        $page = $params['page'] ?? Constant::PAGE;
        $size = $params['size'] ?? Constant::SIZE;
        $memberUserRelation = MemberUserRelation::query();
        $conditions = [];

        $qyDomain      = $params['qy_domain'];
        $appIdentifier = $params['app_identifier'];

        $conditions[] = ['enterprise_id', '=', $params['enterprise_id']];

        if (!empty($params['name'])) {
            $conditions[] = ['ec_name', 'like', '%' . $params['name'] . '%'];
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

        // 性别筛选
        if (!empty($params['gender']) && in_array($params['gender'], self::GENDERS)) {
            $conditions[] = ['ec_gender', '=', $params['gender']];
        }

        // 有开始结束时间段
        if (!empty($params['start_time']) && !empty($params['end_time'])) {
            $conditions[] = ['relation_at', '>=', $params['start_time']];
            $conditions[] = ['relation_at', '<=', $params['end_time']];
        }

        // 已被客户删除
        if (!empty($params['type']) && $params['type'] == self::CONDITION_MDS) {
            $conditions[] = ['delete_type', '=', self::CONDITION_MDS];
        }

        //以下三段代码为, 是为了把部门员工的条件加入进去
        $query = $memberUserRelation
            ->where($conditions)
            ->when(!empty($sourceTypes), function($query) use ($sourceTypes) {
                $query->whereIn('source_type', $sourceTypes);
            });

        $memIds = [];
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

        // 只传递了部门且部门下没有人
        if (empty($params['mem_uid']) && !empty($params['dep_uid']) && empty($memIds)) {
            return success([
                'ids'  => [],
                'list' => [
                    'data' => [],
                    'meta' => [
                        'page'  => (int)$page,
                        'size'  => (int)$size,
                        'total' => 0,
                    ]
                ]
            ]);
        }

        // 分组筛选 等cdp
        $ecUserIds = [];
        if (!empty($params['group_id']) && CustomerGroup::query()->where('id', '=', $params['group_id'])->exists()) {
            // 变量名不改了，实际上为ec_userid
            $ecUserIds = CustomerGroupRelation::query()
                ->where('enterprise_id', '=', $params['enterprise_id'])
                ->where('customer_group_id', '=', $params['group_id'])
                ->pluck('ec_userid')
                ->toArray();
        }

        if (!empty($params['group_id']) && empty($ecUserIds)) {
            return success([
                'ids'  => [],
                'list' => [
                    'data' => [],
                    'meta' => [
                        'page'  => (int)$page,
                        'size'  => (int)$size,
                        'total' => 0,
                    ]
                ]
            ]);
        }

        $datum  = [];
        $repeat = [];
        if (!empty($params['tag_id'])) {
            $tag = (new Tag)->tagDetail($qyDomain, ['tagId' => $params['tag_id']]);
            if (!empty($tag)) {
                foreach ($tag as $key => $item) {
                    if (is_numeric($key)) {
                        if (empty($repeat[$item['ecUserid'] . '/' . $item['memUid']])) {
                            // 去重
                            $repeat[$item['ecUserid'] . '/' . $item['memUid']] = 1;
                            $datum[]                                           = [
                                'ec_userid' => $item['ecUserid'],
                                'mem_uid'   => $item['memUid']
                            ];
                        }
                    }
                }
            }
        }

        if (!empty($params['tag_id']) && empty($datum)) {
            return success([
                'ids'  => [],
                'list' => [
                    'data' => [],
                    'meta' => [
                        'page'  => (int)$page,
                        'size'  => (int)$size,
                        'total' => 0,
                    ]
                ]
            ]);
        }

        $endMemUids  = [];
        $endEcUserIs = [];

        if ($memIds && $datum) {
            foreach ($datum as $k => $v) {
                if (!in_array($v['mem_uid'], $memIds)) {
                    unset($datum[$k]);
                }
            }
        }

        if ($ecUserIds && $datum) {
            foreach ($datum as $k => $v) {
                if (!in_array($v['ec_userid'], $ecUserIds)) {
                    unset($datum[$k]);
                }
            }
        }

        if (empty($datum) && $ecUserIds) {
            $endEcUserIs = $ecUserIds;
        }

        if (empty($datum) && $memIds) {
            $endMemUids = $memIds;
        }

        if ($endEcUserIs) {
            $query = $query->whereIn('ec_userid', $endEcUserIs);
        }

        if ($endMemUids) {
            $query = $query->whereIn('mem_uid', $endMemUids);
        }

        if ($datum) {
            $count = count($datum);
            $datum = array_values($datum);
            $query->where(function ($query) use ($count, $datum) {
                if ($count > 1) {
                    $query->where([
                        ['ec_userid', '=', $datum[0]['ec_userid']],
                        ['mem_uid', '=', $datum[0]['mem_uid']],
                    ]);
                    unset($datum[0]);
                    foreach ($datum as $item) {
                        $query->orWhere([
                            ['ec_userid', '=', $item['ec_userid']],
                            ['mem_uid', '=', $item['mem_uid']],
                        ]);
                    }
                } else {
                    $query->where([
                        ['ec_userid', '=', $datum[0]['ec_userid']],
                        ['mem_uid', '=', $datum[0]['mem_uid']],
                    ]);
                }
            });
        }

        $members = $query->orderBy('relation_at', 'desc')
            ->paginate($size, ['*'], '', $page);

        // 最终的响应数据组装
        $ids  = [];
        $data = [];
        $sourceType = array_column($members->items(), 'source_type');
        $channelByContact = ContactQrCode::query()->with('channel')->whereIn('id', $sourceType)->get()->toArray();
        $channelByContact = array_column($channelByContact, null, 'id');
        $memUids = array_column($members->items(), 'mem_uid');
        $memberByUc = (new Member())->list($qyDomain, $appIdentifier, [
            'memUids' => $memUids
        ])['list'];
        $memUsers = array_column($memberByUc, null, 'memUid');
        unset($memberByUc);
        $tagObject = new Tag();
        foreach ($members->items() as $member) {
            $memUser = $memUsers[$member->mem_uid] ?? [];
            $dep_name = "未知部门";
            if (!empty($memUser['dpName'])) {
                $memberUserDp = array_column($memUser['dpName'], null, 'dpId');
                $dep_name = current($memberUserDp)['dpName'];
                if (!empty($params['dep_uid'])) {
                    $dep_name = $memberUserDp[$params['dep_uid']]['dpName'];
                }
            }
            $mem_name = !empty($memUser['memUsername']) ? $memUser['memUsername'] : '';

            $tagIds   = [];
            $tagNames = [];
            $tagList  = $tagObject->tagDetail($qyDomain, ['ecUserid' => $member['ec_userid']]);
            unset($tagList['requestId']);
            // 获取这个员工给这个客户打的标签
            if (!empty($tagList)) {
                $tagCollect = [];
                foreach ($tagList as $key => $item) {
                    if ($item['memUid'] == $member['mem_uid']) {
                        $tagCollect[$item['tagId']] = $item['tagName'];
                    }
                }
                $tagIds   = array_keys($tagCollect);
                $tagNames = array_values($tagCollect);
                foreach ($tagIds as &$tagId) {
                    $tagId = (string)$tagId;
                }
            }

            // 客户员工关系正常, 显示该员工给客户打的标签
            if ($member['delete_type'] > 0) {
                // 查询是否与其他客户关系正常
                $existsOtherRelation = $memberUserRelation
                    ->where('ec_id', '=', $member['ec_id'])
                    ->where('delete_type', '=', 0)
                    ->where('mem_uid', '!=', $member['mem_uid'])
                    ->exists();

                // 存在就啥都不做, 返回空
                if ($existsOtherRelation) {
                    $tagIds   = [];
                    $tagNames = [];
                }
            }

            // 获取分组
            $group = CustomerGroupRelation::query()
                ->leftJoin('customer_group AS cg', 'cg.id', '=', 'customer_group_relation.customer_group_id')
                ->where('customer_group_relation.enterprise_id', $params['enterprise_id'])
                ->where('ec_userid', '=', $member['ec_userid'])
                ->pluck('name')
                ->toArray();

            $ids[]  = $member['union_id'];
            $data[] = [
                'id'                => $member['id'],
                'name'              => $member['ec_name'],
                'avatar'            => $member['ec_avatar'],
                'gender'            => self::GENDER[$member['ec_gender']],
                'ec_userid'         => $member['ec_userid'],
                'user_id'           => $member['user_id'],
                'union_id'          => $member['union_id'],
                "group"             => $group,
                "tag"               => $tagNames,
                "tag_ids"           => $tagIds,
                'dep_name'          => $dep_name,
                'mem_name'          => $mem_name,
                'mem_mobile'        => $memUser['memMobile'] ?? '',
                'mem_uid'           => $member['mem_uid'],
                'channel_name'      => $channelByContact[$member['source_type']]['channel']['channel'] ?? '',
                'phone'             => '',
                'order_total'       => 0,
                'order_price_total' => 0,
                'created_at'        => date('Y-m-d H:i:s', $member['relation_at']),
                'deleted_at'        => date('Y-m-d H:i:s', $member['cancel_relation_at']),
            ];
        }

        return success([
            'ids'  => $ids,
            'list' => [
                'data' => $data,
                'meta' => [
                    'page'  => (int)$page,
                    'size'  => (int)$size,
                    'total' => (int)$members->total(),
                ]
            ]
        ]);
    }
}