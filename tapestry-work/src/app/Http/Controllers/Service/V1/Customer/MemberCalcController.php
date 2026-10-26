<?php

namespace App\Http\Controllers\Service\V1\Customer;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\Channel;
use App\Models\ContactQrCode;
use App\Models\MemberUserRelation;
use Carbon\Carbon;

class MemberCalcController extends BaseController
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
     * @return array
     * @throws \Exception
     */
    public function memberCalc()
    {
        $params = $this->request->all();

        $qyDomain      = $params['qy_domain'] ?? null;
        $enterpriseId  = $params['enterprise_id'] ?? null;
        $appIdentifier = $params['app_identifier'] ?? null;
        $page          = $params['page'] ?? Constant::DEFAULT_PAGE;
        $size          = $params['size'] ?? Constant::DEFAULT_SIZE;
        //    $type          = $params['type'] ?? 4;
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
            if (is_string($params['mem_uid'])) {
                $memIds[] = $params['mem_uid'];
            } else {
                $memIds = array_merge($memIds, $params['mem_uid']);
            }
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

        // 以下取到的时间多了最后一天, 按今天和明天的时间段算好数据, 存储好后, 去掉最后一个时间(日期时间轴)
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
        $num            = count($data['x_data']); // 一共几天
        $data['y_data'] = [];
        for ($i = 0; $i < $num - 1; $i++) {
            // 客户总数
            $total = MemberUserRelation::query()
                ->when($memIds, function ($query) use ($memIds) {
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
            $data['y_data'][] = [$total, $totally, $lossTotal, $add, $del, $passiveDel];
        }
        unset($data['x_data'][$num - 1]);

        $numData  = [];
        $numTotal = 0;

        // 按日期
        if ($dataType == self::TIME_DATA) {

            $datum = [];
            foreach ($data['x_data'] as $key => $item) {
                $datum[] = array_merge([date('Y-m-d', $item)], $data['y_data'][$key]);
            }

            $numTotal = count($datum);
            $datum    = array_reverse($datum);
            $numData  = array_slice($datum, $size * ($page - 1), $size);
        }

        // 按员工
        if ($dataType == self::STAFF_DATA) {

            // 获取这个时间段有客户变动的员工mem_uid - 有指定员工和没员工的情况不一样, 时间问题
            if (empty($memIds)) {
                $members = MemberUserRelation::query()
                    ->where($conditions)
                    ->where(function ($query) use ($startTime, $endTime) {
                        $query->where([
                            ['relation_at', '>=', $startTime],
                            ['relation_at', '<=', $endTime],
                        ])->orWhere([
                            ['cancel_relation_at', '>=', $startTime],
                            ['cancel_relation_at', '<=', $endTime],
                        ]);
                    })
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
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
                    ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                        $query->whereIn('source_type', $contactQrcodeIds);
                    })
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
                        ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                            $query->whereIn('source_type', $contactQrcodeIds);
                        })
                        ->pluck('mem_uid')
                        ->toArray();
                }
            }

            $members = array_unique($members);

            // 根据mem_uid计算每个mem_id的5项数据
            $staffDatums = [];
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
                $staffDatums[] = [$memUser['memUsername'] ?? '', $total, $totally, $lossTotal, $add, $del, $passiveDel];
            }

            $numTotal = count($staffDatums);
            $numData  = array_slice($staffDatums, $size * ($page - 1), $size);
        }

        $todayTimestamp = strtotime(date('Y-m-d 23:59:59', time()));

        // 顾客总数
        $customerTotal = MemberUserRelation::query()
            ->when($memIds, function ($query) use ($memIds) {
                return $query->whereIn('mem_uid', $memIds);
            })
            ->where($conditions)
            ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                $query->whereIn('source_type', $contactQrcodeIds);
            })->where('relation_at', '<=', $todayTimestamp)
            ->count();

        // 净增客户总数
        $pureCustomerTotal = MemberUserRelation::query()
            ->when($memIds, function ($query, $memIds) {
                return $query->whereIn('mem_uid', $memIds);
            })
            ->where($conditions)
            ->when(!empty($contactQrcodeIds), function ($query) use ($contactQrcodeIds) {
                $query->whereIn('source_type', $contactQrcodeIds);
            })->where('relation_at', '<=', $todayTimestamp)
            ->where('delete_type', MemberUserRelation::STAFF_USER_NOT_DEL)
            ->count();

        foreach ($numData as &$datum) {
            $datum = [
                'date_name'  => $datum[0],
                'total'      => $datum[1],
                'totally'    => $datum[2],
                'loss_total' => $datum[3],
                'add'        => $datum[4],
                'del'        => $datum[5],
                'passive'    => $datum[6]
            ];
        }

        $result = [
            'chart'    => $data,
            'paginate' => [
                'type' => $dataType,
                'data' => $numData,
                'meta' => [
                    'page'  => (int)$page,
                    'size'  => (int)$size,
                    'total' => $numTotal,
                ]
            ],
            'data_overview' => [
                'customer_total' => $customerTotal,
                'pure_customer_total' => $pureCustomerTotal
            ],
        ];

        return success(['data' => $result]);
    }
}