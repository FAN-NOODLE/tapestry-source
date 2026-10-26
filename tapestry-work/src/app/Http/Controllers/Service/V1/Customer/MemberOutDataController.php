<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Models\Channel;
use App\Library\Uc\Member;
use App\Models\ContactQrCode;
use App\Models\MemberUserRelation;
use App\Http\Controllers\Service\V1\BaseController;
use App\Services\StaffOrderService;

class MemberOutDataController extends BaseController
{
    /**
     * 按日期查
     */
    const TIME_DATA = 1;

    /**
     * 按员工查
     */
    const STAFF_DATA = 2;

    /**
     * 客户统计导出
     */
    public function memberOutData()
    {
        $params = $this->request->all();

        $qyDomain      = $params['qy_domain'] ?? null;
        $enterpriseId  = $params['enterprise_id'] ?? null;
        $appIdentifier = $params['app_identifier'] ?? null;
        //       $type          = $params['type'] ?? 4;
        $memIds     = [];
        $conditions = [];

        if (empty($qyDomain) || empty($appIdentifier)) {
            notice(110007);
        }

        if (empty($enterpriseId)) {
            notice(600000);
        } else {
            $conditions[] = ['enterprise_id', '=', $enterpriseId];
        }

        // 部门ID 和 员工ID 都在 - 只给这个员工的客户,不验证部门
        if (!empty($params['mem_uid']) && !empty($params['dep_uid'])) {
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

        // 添加渠道
        $contactQrcodeIds = [];
        if (!empty($params['channel_id'])) {
            $contactQrcodeIds = ContactQrCode::query()
                ->where('channel_id', $params['channel_id'])
                ->pluck('id')
                ->toArray();
        }

        $endTime   = $params['end_time'] ?? '';
        $startTime = $params['start_time'] ?? '';

        // 处理数据类型
        $dataType = self::TIME_DATA;
        if (!empty($params['data_type']) && in_array($params['data_type'], [self::TIME_DATA, self::STAFF_DATA])) {
            $dataType = $params['data_type'];
        }

        // 以下取到的时间多了最后一天, 按今天和明天的时间段算好数据, 存储好后, 去掉最后一个时间
        $next             = true;
        $data['x_data'][] = $startTime;
        $tmpTime          = $startTime;
        do {
            $tmpTime = $tmpTime + 86400;
            if ($tmpTime < $endTime) {
                $data['x_data'][] = $tmpTime;
            } else {
                $next = false;
            }
        } while ($next);
        $data['x_data'][] = $endTime;
        sort($data['x_data'], SORT_NUMERIC);

        /*
        开始处理数据
        - 总添加人数(截止到某个时间的总添加人数)
        - 净增人数(不包含互删的客户)
        - 添加人数(包含互删的客户)
        - 删除人数(员工删除客户)
        - 被删除人数(客户删除员工)
        */

        $outData = [];

        // 按日期
        if ($dataType == self::TIME_DATA) {

            $num       = count($data['x_data']);
            $outData[] = ['日期', '客户总数', '净增客户总数', '流失客户总数', '新增客户数', '删除人数', '被删除人数'];

            for ($i = 0; $i < $num - 1; $i++) {

                // 客户总数
                $total = MemberUserRelation::query()
                    ->when($memIds, function ($query, $memIds) {
                        return $query->whereIn('mem_uid', $memIds);
                    })
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })->where('relation_at', '<=', $data['x_data'][$i + 1])
                    ->count();

                // 净增客户总数
                $totally = MemberUserRelation::query()
                    ->when($memIds, function ($query, $memIds) {
                        return $query->whereIn('mem_uid', $memIds);
                    })
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })->where('relation_at', '<=', $data['x_data'][$i + 1])
                    ->where('delete_type', MemberUserRelation::STAFF_USER_NOT_DEL)
                    ->count();


                // 新增客户数
                $add = MemberUserRelation::query()
                    ->when($memIds, function ($query, $memIds) {
                        return $query->whereIn('mem_uid', $memIds);
                    })
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('relation_at', '>=', $data['x_data'][$i])
                    ->where('relation_at', '<=', $data['x_data'][$i + 1])
                    ->count();

                // 删除人数
                $del = MemberUserRelation::query()
                    ->when($memIds, function ($query, $memIds) {
                        return $query->whereIn('mem_uid', $memIds);
                    })
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('cancel_relation_at', '>=', $data['x_data'][$i])
                    ->where('cancel_relation_at', '<=', $data['x_data'][$i + 1])
                    ->where('delete_type', '=', MemberUserRelation::DEL_EXTERNAL_CONTACT)
                    ->count();

                // 被删除人数
                $passiveDel = MemberUserRelation::query()
                    ->when($memIds, function ($query, $memIds) {
                        return $query->whereIn('mem_uid', $memIds);
                    })
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('cancel_relation_at', '>=', $data['x_data'][$i])
                    ->where('cancel_relation_at', '<=', $data['x_data'][$i + 1])
                    ->where('delete_type', '=', MemberUserRelation::DEL_FOLLOW_USER)
                    ->count();

                // 流失客户总数
                $lossTotal = $total - $totally;

                $outData[] = [date('Y-m-d', $data['x_data'][$i]), $total, $totally, $lossTotal, $add, $del, $passiveDel];
            }
        }

        // 按员工
        if ($dataType == self::STAFF_DATA) {

            // 获取这个时间段有客户变动的员工mem_uid - 有指定员工和没员工的情况不一样, 时间问题
            if (empty($memIds)) {
                $members = MemberUserRelation::query()
                    ->where($conditions)
                    ->where([
                        ['relation_at', '>=', $startTime],
                        ['relation_at', '<=', $endTime],
                    ])
                    ->orWhere([
                        ['cancel_relation_at', '>=', $startTime],
                        ['cancel_relation_at', '<=', $endTime],
                    ])
                    ->pluck('mem_uid')
                    ->toArray();
            } else {
                // 当天有没有新增客户
                $members = MemberUserRelation::query()
                    ->where($conditions)
                    ->whereIn('mem_uid', $memIds)
                    ->where([
                        ['relation_at', '>=', $startTime],
                        ['relation_at', '<=', $endTime],
                    ])
                    ->pluck('mem_uid')
                    ->toArray();
                // 没有查找当天有没有删除客户
                if (empty($members)) {
                    $members = MemberUserRelation::query()
                        ->where($conditions)
                        ->whereIn('mem_uid', $memIds)
                        ->where([
                            ['cancel_relation_at', '>=', $startTime],
                            ['cancel_relation_at', '<=', $endTime],
                        ])
                        ->pluck('mem_uid')
                        ->toArray();
                }
            }

            $members   = array_unique($members);
            $outData[] = ['员工', '客户总数', '净增客户总数', '流失客户总数', '新增客户数', '删除人数', '被删除人数'];

            // 根据mem_uid计算每个mem_id的5项数据
            foreach ($members as $member) {
                // 获取员工名字
                $memUser = (new Member())->detail($qyDomain, $appIdentifier, ['memUid' => $member,'memAll'=>1]);

                // todo::洽洽出现了, 客户关联关系的 员工数据 ,不是属于 非我们boss, 的企业标识的员工。 注意解决
                if (empty($memUser['memUsername'])) {
                    continue;
                }

                // 客户总数
                $total = MemberUserRelation::query()
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('relation_at', '<=', $endTime)
                    ->where('mem_uid', '=', $member)
                    ->count();

                // 净增客户总数
                $totally = MemberUserRelation::query()
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('relation_at', '<=', $endTime)
                    ->where('delete_type', MemberUserRelation::STAFF_USER_NOT_DEL)
                    ->where('mem_uid', '=', $member)
                    ->count();


                // 新增客户数
                $add = MemberUserRelation::query()
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('relation_at', '>=', $startTime)
                    ->where('relation_at', '<=', $endTime)
                    ->where('mem_uid', '=', $member)
                    ->count();

                // 删除人数
                $del = MemberUserRelation::query()
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('cancel_relation_at', '>=', $startTime)
                    ->where('cancel_relation_at', '<=', $endTime)
                    ->where('delete_type', '=', MemberUserRelation::DEL_EXTERNAL_CONTACT)
                    ->where('mem_uid', '=', $member)
                    ->count();

                // 被删除人数
                $passiveDel = MemberUserRelation::query()
                    ->where($conditions)
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
                    ->where('cancel_relation_at', '>=', $startTime)
                    ->where('cancel_relation_at', '<=', $endTime)
                    ->where('delete_type', '=', MemberUserRelation::DEL_FOLLOW_USER)
                    ->where('mem_uid', '=', $member)
                    ->count();

                // 流失客户总数
                $lossTotal = $total - $totally;

                if ($total == 0 && $totally == 0 && $lossTotal == 0 && $add == 0 && $del == 0 && $passiveDel == 0) {
                    continue;
                }

                $outData[] = [$memUser['memUsername'] ?? '', $total, $totally, $lossTotal, $add, $del, $passiveDel];
            }
        }

        $result = [
            'type' => $dataType,
            'data' => $outData,
        ];

        return success($result);
    }
}