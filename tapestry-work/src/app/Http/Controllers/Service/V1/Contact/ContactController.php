<?php
/*
* +------------------------------------+
* | Created: phpstorm
* +------------------------------------+
* | Author : wyg
* +------------------------------------+
* | Date   : 2020/1/21
* +------------------------------------+
* | Time   : 14:13
* +------------------------------------+
* | Email  : <wangyongguo@vchangyi.com>
* +------------------------------------+
*/


namespace App\Http\Controllers\Service\V1\Contact;


use App\Common\Constant;
use App\Exceptions\Error;
use App\Http\Controllers\Service\V1\BaseController;
use App\Library\Uc\Member;
use App\Models\ContactQrCode;
use App\Models\ContactQrCodeParams;
use App\Models\MemberGuide;
use App\Models\WelcomeSpeech;
use App\Services\CommonService;
use App\Services\MemberService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\Rule;

class ContactController extends BaseController
{
    const TYPE_MEMBER = 1; // 员工类型

    const TYPE_DEPART = 2; // 部门类型

    const QR_SINGLE = 1; // 单人二维码

    const QR_MANY = 2; // 多人二维码

    const MAX_LIMIT = 1000; // UC数据查询返回最大条数

    const SHOW_BY_QR_CODE = 2; // 通过二维码联系

    const MAX_MEMBER_COUNT = 100; // 二维码一次最可多应用的人数(总数)

    // 对象类型数组
    const TYPE_ARRAY = [
        self::TYPE_DEPART,
        self::TYPE_MEMBER
    ];

    // 二维码类型数组
    const QR_TYPE_ARRAY = [
        self::QR_SINGLE,
        self::QR_MANY
    ];

    /**
     * 创建二维码
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function createQrCode(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain' => 'required',
            'identifier' => 'required',
            'extend' => 'required|array',
            'extend.channel_id' => 'required',
            'extend.channel' => 'required',
            'extend.name' => 'required',
            'extend.type' => [
                'required',
                Rule::in(self::QR_TYPE_ARRAY),
            ],
            'extend.is_verify' => 'required',
            'extend.members' => 'required|array',
            'extend.members.*.value' => 'required',
            'extend.members.*.label' => 'required',
            'extend.members.*.type' => [
                'required',
                Rule::in(self::TYPE_ARRAY),
            ],
        ];

        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID类型错误',
            'enterprise_id.min' => '企业ID不能为零',
            'domain.required' => '企业域名不能为空',
            'extend.required' => '二维码参数[extend]不能为空',
            'extend.array' => '二维码参数[extend]类型错误',
            'extend.channel_id.required' => '二维码渠道id不能为空',
            'extend.channel.required' => '二维码渠道名称不能为空',
            'extend.name.required' => '二维码名称不能为空',
            'extend.type.required' => '二维码类型不能为空',
            'extend.is_verify.required' => '二维码添加设置名称不能为空',
            'extend.members.required' => '二维码参数[members]不能为空',
            'extend.members.array' => '二维码参数[members]类型错误',
            'extend.members.*.value.required' => '二维码标识数值[value]不能为空',
            'extend.members.*.label.required' => '二维码标识名称[label]不能为空',
            'extend.members.*.type.required' => '二维码标识类型[type]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $enterpriseId = $params['enterprise_id'];

        // 渠道重复检测
//        $isExistChannel = ContactQrCode::query()
//            ->where([
//                'enterprise_id' => $enterpriseId,
//                'channel'       => $params['extend']['channel'],
//            ])->exists();

//        if ($isExistChannel) {
//            throw new Error(60000, '二维码渠道已存在');
//        }

        // 名称重复检测
        $isExistName = ContactQrCode::query()
            ->where([
                'enterprise_id' => $enterpriseId,
                'name' => $params['extend']['name'],
            ])->exists();

        if ($isExistName) {
            throw new Error(60000, '二维码名称不能重复，请重新输入');
        }

        // 查询所有员工
        $memberList = MemberService::memberList($params['domain'], $params['identifier']);

        // 部门对应已加入企业员工数组
        $deptMemUIdArr = [];
        // 员工userId和memUId对应数组
        $allMemUIdsArr = [];
        foreach ($memberList as $member) {
            $dpThirdId = $member['dpName'][0]['dpThirdid'] ?? 0;
            // 排除未加入企业的员工
            if (!is_null($member['memSubscribeTime'])) {
                $memUserId = $member['memUserid'];
                $deptMemUIdArr[$dpThirdId][] = $member['memUid'];
                $allMemUIdsArr[$memUserId] = $member['memUid'];
            }
        }

        // 企业部门数据查询
        $departList = MemberService::departList($params['domain'], $params['identifier'], ['pageSize' => self::MAX_LIMIT]);

        // 部门子部门数组
        $dpThirdIdSubs = [];
        // 部门dpThirdId对应dpId数组
        $allDeptIdArr = [];
        foreach ($departList as $depart) {
            $dpThirdId = $depart['dpThirdid'];
            // 没有员工的部门也需要加入分组
            if (!isset($deptMemUIdArr[$dpThirdId])) {
                $deptMemUIdArr[$dpThirdId] = [];
            }
            $allDeptIdArr[$dpThirdId] = $depart['dpId'];
            $dpThirdIdSubs[$dpThirdId] = CommonService::getSubs($departList, $depart['dpId']);
        }

        // 部门已加入员工数据整合
        foreach ($deptMemUIdArr as $dpThirdId => &$memUIdArr) {
            if (!empty($dpThirdIdSubs[$dpThirdId])) {
                foreach ($dpThirdIdSubs[$dpThirdId] as $dpSubThirdId) {
                    if (isset($deptMemUIdArr[$dpSubThirdId])) {
                        $memUIdArr = array_unique(array_merge($memUIdArr, $deptMemUIdArr[$dpSubThirdId]));
                    }
                }
            }
        }

        // 插入数据准备
        $insert = [];
        $memberPartIds = [];
        $memberUserIds = [];
        $memberUIdsArr = [];
        foreach ($params['extend']['members'] as $member) {
            // 类型装换
            $member['type'] = self::TYPE_MEMBER == $member['type'] ? self::TYPE_DEPART : self::TYPE_MEMBER;
            // 员工数据初始化
            if (self::TYPE_MEMBER == $member['type']) {
                $memberUserId = $member['value'];
                if (isset($allMemUIdsArr[$memberUserId])) {
                    $memberUIdsArr[] = $allMemUIdsArr[$memberUserId];
                }
                // UC员工集合
                $memberUserIds[] = $memberUserId;
                $insert[] = [
                    'enterprise_id' => $enterpriseId,
                    'type' => $member['type'],
                    'object' => $allMemUIdsArr[$memberUserId] ?? '',
                    'object_extra' => $allMemUIdsArr[$memberUserId] ?? '',
                    'created_at' => time(),
                    'updated_at' => time(),
                ];
            }
            // 部门数据初始化
            if (self::TYPE_DEPART == $member['type']) {
                $dpThirdId = $member['value'];
                // UC部门集合
                $memberPartIds[] = $dpThirdId;
                $insert[] = [
                    'enterprise_id' => $enterpriseId,
                    'type' => $member['type'],
                    'object' => $allDeptIdArr[$dpThirdId] ?? '',
                    'object_extra' => implode(',', $deptMemUIdArr[$dpThirdId] ?? []),
                    'created_at' => time(),
                    'updated_at' => time(),
                ];
            }
        }

        try {
            DB::beginTransaction();
            // 本地存储二维码
            $contactQr = ContactQrCode::query()->create([
                'enterprise_id' => $params['enterprise_id'],
                'is_verify' => $params['extend']['is_verify'],
                'channel_id' => $params['extend']['channel_id'],
                'channel' => $params['extend']['channel'],
                'name' => $params['extend']['name'],
                'source_type' => $params['extend']['source_type'] ?? 0,
                'type' => $params['extend']['type'],
            ]);

            // 二维码参数外键更新
            foreach ($insert as &$item) {
                $item['contact_qr_code_id'] = $contactQr['id'];
            }

            // 扩展数据存储
            $insert && ContactQrCodeParams::query()->insert($insert);

            // UC存储并获得二维码id(单人)
            if (self::QR_SINGLE == $params['extend']['type']) {
                $result = (new Member())->getContactQrCode($params['domain'], [
                    "ecwScene" => self::SHOW_BY_QR_CODE,
                    'ecwType' => $params['extend']['type'],
                    'ecwRemark' => $params['extend']['channel'],
                    "ecwSkipVerify" => $params['extend']['is_verify'],
                    "ecwState" => $contactQr['id'],
                    'memUserid' => current($memberUserIds),
                ]);
            } else {
                // UC存储并获得二维码id(多人)
                $result = (new Member())->getContactQrCode($params['domain'], [
                    "ecwScene" => self::SHOW_BY_QR_CODE,
                    'ecwType' => $params['extend']['type'],
                    'ecwRemark' => $params['extend']['channel'],
                    "ecwSkipVerify" => $params['extend']['is_verify'],
                    "ecwState" => $contactQr['id'],
                    "memPartyids" => $memberPartIds,
                    "memUserids" => $memberUserIds
                ]);
            }

            if (!is_array($result)) {
                $result = json_decode($result, true);
                if (isset($result['code'])) {
                    switch ($result['code']) {
                        case 40098;
                            throw new Error(40098, '存在成员尚未实名认证 ');
                            break;
                        case 41041;
                            throw new Error(61000, '二维码一次最大添加人数不超过' . self::MAX_MEMBER_COUNT);
                            break;
                        default:
                            throw new Error($result['code'], 'work we_chart server error!');
                    }
                } else {
                    throw new Error(80000, 'uc server error!');
                }
            }

            // 企业微信二维码链接
            $officialUrl = $result['qrCode'];

            // 企业微信二维码存储到腾讯云
            $transformUrl = CommonService::uploadImage($officialUrl);

            // 更新本地数据
            $contactQr['uc_id'] = $result['id'];
            $contactQr['config_id'] = $result['configId'];
            $contactQr['official_url'] = $officialUrl;
            $contactQr['transform_url'] = $transformUrl;

            $contactQr->save();

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            throw new Error($e->getCode() ?: 10000, $e->getMessage() ?: '二维码创建异常');
        }

        $contactQr['members'] = $params['extend']['members'];

        return success([
            'data' => $contactQr
        ]);
    }

    /**
     * 更新二维码
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function updateQrCode(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain' => 'required',
            'identifier' => 'required',
            'extend' => 'required|array',
            'extend.id' => 'required|integer|min:1',
            'extend.channel' => 'required',
            'extend.name' => 'required',
            'extend.type' => [
                'required',
                Rule::in(self::QR_TYPE_ARRAY),
            ],
            'extend.is_verify' => 'required',
            'extend.members' => 'required|array',
            'extend.members.*.value' => 'required',
            'extend.members.*.label' => 'required',
            'extend.members.*.type' => [
                'required',
                Rule::in(self::TYPE_ARRAY),
            ],
        ];

        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID类型错误',
            'enterprise_id.min' => '企业ID不能为零',
            'domain.required' => '企业域名不能为空',
            'extend.required' => '二维码参数[extend]不能为空',
            'extend.array' => '二维码参数[extend]类型错误',
            'extend.id.required' => '二维码参数ID不能为空',
            'extend.id.integer' => '二维码参数ID类型错误',
            'extend.id.min' => '二维码参数ID不能为零',
            'extend.channel.required' => '二维码渠道名称不能为空',
            'extend.name.required' => '二维码名称不能为空',
            'extend.type.required' => '二维码类型不能为空',
            'extend.is_verify.required' => '二维码添加设置名称不能为空',
            'extend.members.required' => '二维码参数[members]不能为空',
            'extend.members.array' => '二维码参数[members]类型错误',
            'extend.members.*.value.required' => '二维码标识数值[value]不能为空',
            'extend.members.*.label.required' => '二维码标识名称[label]不能为空',
            'extend.members.*.type.required' => '二维码标识类型[type]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $contactQrCodeId = $params['extend']['id'];
        $enterpriseId = $params['enterprise_id'];

        // 渠道重复检测
//        $isExistChannel = ContactQrCode::query()
//            ->where([
//                'enterprise_id' => $enterpriseId,
//                'channel'       => $params['extend']['channel'],
//            ])
//            ->where('id', '<>', $contactQrCodeId)
//            ->exists();

//        if ($isExistChannel) {
//            throw new Error(60000, '二维码渠道已存在');
//        }

        $contactQrCode = ContactQrCode::query()->find($contactQrCodeId);

        if (!$contactQrCode) {
            throw new Error(40000, '二维码不存在');
        }

        //内购渠道不允许编辑和删除
        if ($contactQrCode->source_type) {
            throw new Error(40000, '非手动创建的二维码不允许编辑');
        }

        //如果二维码名称有更新则验重
        $updateName = $params['extend']['name'];
        if ($contactQrCode->name != $updateName) {
            if (ContactQrCode::query()->where('enterprise_id', $params['enterprise_id'])
                ->where('name', $updateName)->count()) {
                throw new Error(40000, '二维码名称不能重复，请重新输入');
            }
        }

        // 查询所有员工
        $memberList = MemberService::memberList($params['domain'], $params['identifier']);

        // 部门对应已加入企业员工数组
        $deptMemUIdArr = [];
        // 员工userId和memUId对应数组
        $allMemUIdsArr = [];
        foreach ($memberList as $member) {
            $dpThirdId = $member['dpName'][0]['dpThirdid'] ?? 0;
            // 排除未加入企业的员工
            if (!is_null($member['memSubscribeTime'])) {
                $memUserId = $member['memUserid'];
                $deptMemUIdArr[$dpThirdId][] = $member['memUid'];
                $allMemUIdsArr[$memUserId] = $member['memUid'];
            }
        }

        // 企业部门数据查询
        $departList = MemberService::departList(
            $params['domain'],
            $params['identifier'],
            ['pageSize' => self::MAX_LIMIT]
        );

        // 部门子部门数组
        $dpThirdIdSubs = [];
        // 部门dpThirdId对应dpId数组
        $allDeptIdArr = [];
        foreach ($departList as $depart) {
            $dpThirdId = $depart['dpThirdid'];
            // 没有员工的部门也需要加入分组
            if (!isset($deptMemUIdArr[$dpThirdId])) {
                $deptMemUIdArr[$dpThirdId] = [];
            }
            if (!isset($deptUnJoinMemUIdArr[$dpThirdId])) {
                $deptUnJoinMemUIdArr[$dpThirdId] = [];
            }
            $allDeptIdArr[$dpThirdId] = $depart['dpId'];
            $dpThirdIdSubs[$dpThirdId] = CommonService::getSubs($departList, $depart['dpId']);
        }

        // 部门已加入员工数据整合
        foreach ($deptMemUIdArr as $dpThirdId => &$memUIdArr) {
            if (!empty($dpThirdIdSubs[$dpThirdId])) {
                foreach ($dpThirdIdSubs[$dpThirdId] as $dpSubThirdId) {
                    if (isset($deptMemUIdArr[$dpSubThirdId])) {
                        $memUIdArr = array_unique(array_merge($memUIdArr, $deptMemUIdArr[$dpSubThirdId]));
                    }
                }
            }
        }

        // 查询本地记录
        $existValue = ContactQrCodeParams::query()
            ->where('contact_qr_code_id', $contactQrCodeId)
            ->pluck('object', 'id')
            ->toArray();

        // 格式化数据
        $insert = [];
        $memberPartIds = [];
        $memberUserIds = [];
        $updateObjects = [];
        $memberUIdsArr = [];
        foreach ($params['extend']['members'] as $member) {
            // 类型转换
            $member['type'] = $member['type'] == self::TYPE_MEMBER ? self::TYPE_DEPART : self::TYPE_MEMBER;
            // 员工数据增量更新
            if (self::TYPE_MEMBER == $member['type']) {
                $memberUserId = $member['value'];
                if (isset($allMemUIdsArr[$memberUserId])) {
                    $memberUIdsArr[] = $allMemUIdsArr[$memberUserId];
                }
                $memberUserIds[] = $memberUserId;
                // memberUserId 转 $memberUId
                $memberUId = $allMemUIdsArr[$memberUserId] ?? '';
                if (in_array($memberUId, $existValue)) {
                    $id = array_search($memberUId, $existValue);
                    // 更新数组初始化
                    $updateObjects[$id] = [
                        'updated_at' => time()
                    ];
                } else {
                    $insert[] = [
                        'enterprise_id' => $params['enterprise_id'],
                        'contact_qr_code_id' => $contactQrCodeId,
                        'type' => $member['type'],
                        'object' => $memberUId,
                        'object_extra' => $memberUId,
                        'created_at' => time(),
                        'updated_at' => time(),
                    ];
                }
            }
            // 部门数据增量更新
            if (self::TYPE_DEPART == $member['type']) {
                $dpThirdId = $member['value'];
                $memberPartIds[] = $dpThirdId;
                // 部门dpThirdId 转部门dpId
                $deptId = $allDeptIdArr[$dpThirdId] ?? '';
                if (in_array($deptId, $existValue)) {
                    // 获取对应主键
                    $id = array_search($deptId, $existValue);
                    // 更新数组初始化
                    $updateObjects[$id] = [
                        'updated_at' => time(),
                        'object_extra' => implode(',', $deptMemUIdArr[$dpThirdId] ?? []),
                    ];
                } else {
                    $insert[] = [
                        'enterprise_id' => $enterpriseId,
                        'contact_qr_code_id' => $contactQrCodeId,
                        'type' => $member['type'],
                        'object' => $deptId,
                        'object_extra' => implode(',', $deptMemUIdArr[$dpThirdId] ?? []),
                        'created_at' => time(),
                        'updated_at' => time(),
                    ];
                }
            }
        }

        // 需要更新的数值
        $updateType = $params['extend']['type'];
        $updateChannel = $params['extend']['channel'];
        $updateVerify = $params['extend']['is_verify'];

        try {
            DB::beginTransaction();
            // 更新UC存储数据(单人)
            if (self::QR_SINGLE == $params['extend']['type']) {
                $result = (new Member())->updateContactQrCode($params['domain'], [
                    "configId" => $contactQrCode['config_id'],
                    "ecwScene" => self::SHOW_BY_QR_CODE,
                    'ecwType' => $updateType,
                    'ecwRemark' => $updateChannel,
                    "ecwSkipVerify" => $updateVerify,
                    'memUserid' => current($memberUserIds),
                ]);
            } else {
                // 更新UC存储数据(多人)
                $result = (new Member())->updateContactQrCode($params['domain'], [
                    "configId" => $contactQrCode['config_id'],
                    "ecwScene" => self::SHOW_BY_QR_CODE,
                    'ecwType' => $updateType,
                    'ecwRemark' => $updateChannel,
                    "ecwSkipVerify" => $updateVerify,
                    "memPartyids" => $memberPartIds,
                    "memUserids" => $memberUserIds
                ]);
            }

            if (!is_array($result)) {
                $result = json_decode($result, true);
                if (isset($result['code'])) {
                    switch ($result['code']) {
                        case 40098;
                            throw new Error(40098, '存在成员尚未实名认证 ');
                            break;
                        case 41041;
                            throw new Error(61000, '二维码一次最大添加人数不超过' . self::MAX_MEMBER_COUNT);
                            break;
                        default:
                            throw new Error($result['code'], 'work we_chart server error!');
                    }
                } else {
                    throw new Error(80000, 'uc server error!');
                }
            }

            // 企业微信二维码链接
            $officialUrl = $result['qrCode'];

            // 企业微信二维码存储到腾讯云
            $transformUrl = CommonService::uploadImage($officialUrl);

            // 二维码数据更新
            $contactQrCode['type'] = $updateType;
            $contactQrCode['name'] = $updateName;
            $contactQrCode['is_verify'] = $updateVerify;
            $contactQrCode['uc_id'] = $result['id'];
            $contactQrCode['config_id'] = $result['configId'];
            $contactQrCode['official_url'] = $officialUrl;
            $contactQrCode['transform_url'] = $transformUrl;

            $contactQrCode->save();

            // 二维码参数更新
            foreach ($updateObjects as $id => $val) {
                ContactQrCodeParams::query()->where('id', $id)->update($val);
            }

            $updateIdArray = array_keys($updateObjects);

            // 删除不需要的
            ContactQrCodeParams::query()
                ->where([
                    'contact_qr_code_id' => $contactQrCodeId
                ])
                ->whereNotIn('id', $updateIdArray)
                ->delete();

            // 新增需要插入的
            $insert && ContactQrCodeParams::query()->insert($insert);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            throw new Error($e->getCode() ?: 50000, $e->getMessage() ?: '二维码更新异常');
        }

        $contactQrCode['members'] = $params['extend']['members'];

        return success([
            'data' => $contactQrCode
        ]);
    }

    /**
     * 二维码详情
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function selectQrCode(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain' => 'required',
            'identifier' => 'required',
            'extend' => 'required|array',
            'extend.id' => 'required|integer|min:1'
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
            'enterprise_id.integer' => '企业[enterprise_id]类型错误',
            'enterprise_id.min' => '企业[enterprise_id]不能为零',
            'domain.required' => '企业[domain]不能为零',
            'extend.required' => '二维码参数[extend]不能为空',
            'extend.array' => '二维码参数[extend]类型错误',
            'extend.id.required' => '二维码[id]不能为空',
            'extend.id.integer' => '二维码[id]类型错误',
            'extend.id.min' => '二维码[id]不能为零',
        ];

        Validator::make($params, $rules, $message)->validate();

        $contactQrCodeId = $params['extend']['id'];

        $contactQrCode = ContactQrCode::query()->find($contactQrCodeId);

        if (!$contactQrCode) {
            throw new Error(10000, '二维码不存在');
        }

        // 部门查询
        $departList = MemberService::departList(
            $params['domain'],
            $params['identifier'],
            ['pageSize' => self::MAX_LIMIT]
        );

        // 部门名称
        $departList = array_column($departList, null, 'dpId');

        // 二维码关联对象
        $localList = ContactQrCodeParams::query()
            ->where([
                'contact_qr_code_id' => $contactQrCodeId
            ])
            ->orderByDesc('id')
            ->get([
                'type',
                'object',
            ])
            ->toArray();

        $memberUIds = array_column(array_filter($localList, function ($member) {
            return Self::TYPE_MEMBER == $member['type'];
        }), 'object');

        // 员工查询
        $memberList = $memberUIds ? MemberService::memberList(
            $params['domain'],
            $params['identifier'],
            ['memUids' => $memberUIds]
        ) : [];

        // 员工集合
        $memberList = array_column($memberList, null, 'memUid');

        // 数据格式化
        $members = [];
        foreach ($localList as $member) {
            if (self::TYPE_MEMBER == $member['type']) {
                $memUid = $member['object'];
                $members[] = [
                    'type' => self::TYPE_DEPART,
                    'label' => $memberList[$memUid]['memUsername'] ?? '',
                    'value' => $memberList[$memUid]['memUserid'] ?? '',
                ];
            } else {
                $departId = $member['object'];
                $members[] = [
                    'type' => self::TYPE_MEMBER,
                    'label' => $departList[$departId]['dpName'] ?? '',
                    'value' => $departList[$departId]['dpThirdid'] ?? '',
                ];
            }
        }

        $contactQrCode['members'] = $members;

        return success([
            'data' => $contactQrCode
        ]);
    }

    /**
     *
     * 删除二维码
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function deleteQrCode(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain' => 'required',
            'extend' => 'required|array',
            'extend.id' => 'required|integer|min:1'
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
            'enterprise_id.integer' => '企业[enterprise_id]类型错误',
            'enterprise_id.min' => '企业[enterprise_id]不能为零',
            'domain.required' => '企业[domain]不能为零',
            'extend.required' => '二维码参数[extend]不能为空',
            'extend.array' => '二维码参数[extend]类型错误',
            'extend.id.required' => '二维码[id]不能为空',
            'extend.id.integer' => '二维码[id]类型错误',
            'extend.id.min' => '二维码[id]不能为零',
        ];

        Validator::make($params, $rules, $message)->validate();

        $contactQrCodeId = $params['extend']['id'];

        $contactQrCode = ContactQrCode::query()->find($contactQrCodeId);

        if (!$contactQrCode) {
            throw new Error(10000, '二维码不存在');
        }
        //非手动创建的不允许编辑和删除
        if ($contactQrCode->source_type) {
            throw new Error(40000, '非手动创建的二维码不允许删除');
        }

        // 检测是否配置欢迎语
//        $isRelatedWelcome = WelcomeSpeech::query()
//            ->where([
//                'enterprise_id'      => $params['enterprise_id'],
//                'contact_qr_code_id' => $contactQrCodeId,
//            ])
//            ->exists();
//
//        if ($isRelatedWelcome) {
//            throw new Error(20000, '该渠道已关联欢迎语,请解除关联后再删除!');
//        }

        try {
            DB::beginTransaction();
            // 从UC删除
            (new Member())->delContactQrCode($params['domain'], ['id' => $contactQrCode['uc_id']]);

            // 从本地删除二维码
            ContactQrCode::query()
                ->where([
                    'enterprise_id' => $params['enterprise_id'],
                    'id' => $contactQrCodeId,
                ])->delete();

            // 删除二维码参数
            ContactQrCodeParams::query()
                ->where([
                    'enterprise_id' => $params['enterprise_id'],
                    'contact_qr_code_id' => $contactQrCodeId,
                ])->delete();

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            throw new Error(50000, '二维码删除异常');
        }

        return success([
            'data' => []
        ]);
    }

    /**
     * 选择员工列表
     *
     * @param Request $request
     *
     * @return array
     * @throws \Exception
     */
    public function memberTree(Request $request)
    {
        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain' => 'required',
            'identifier' => 'required',
        ];

        $message = [
            'enterprise_id.required' => '企业[enterprise_id]不能为空',
            'enterprise_id.integer' => '企业[enterprise_id]类型错误',
            'enterprise_id.min' => '企业[enterprise_id]不能为零',
            'domain.required' => '企业[domain]不能为零',
            'identifier.required' => '企业[identifier]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();

        $hideUser = $params['extend']['type'] ?? 1; // 默认包含员工

        $members = [];
        $notDpMember = [];
        if ($hideUser == 1) {
            // 查询所有员工
            $memberList = MemberService::memberList(
                $params['domain'],
                $params['identifier']
            );
            foreach ($memberList as $member) {
                if (!empty($member['dpName'])) {
                    foreach ($member['dpName'] as $name) {
                        $departId = $name['dpId'];
                        $members[$departId][] = [
                            'user_id' => $member['memUserid'],
                            'mem_uid' => $member['memUid'],
                            'user_name' => $member['memUsername'],
                            'dept_name' => $name['dpName'],
                            'disabled' => $member['memSubscribeStatus'] == 1 ? false : true,
                            'checked' => false
                        ];
                    }
                } else {
                    $notDpMember[] = [
                        'id' => $member['memUid'],
                        'pid' => 0,
                        'lev' => 0,
                        'type' => self::TYPE_DEPART,
                        'label' => $member['memUsername'],
                        'name' => $member['memUsername'],
                        'value' => $member['memUserid'],
                        'disabled' => $member['memSubscribeStatus'] == 1 ? false : true,
                        'checked' => false
                    ];
                }

            }
        }
        // 查询所有部门
        $deptList = MemberService::departList(
            $params['domain'],
            $params['identifier'],
            ['pageSize' => self::MAX_LIMIT]
        );

        /**
         * 树(两种节点类型)字段释义以及和uc字段映射关系
         *
         * memUsername : 员工姓名
         * dpName      : 部门名称
         * memUserid   : 企业微倍用户id
         * dpThirdid   : 企业微倍部门id
         * dpId        : uc部门id
         * dpParentid  : uc父部门id
         * mem_uid     : uc员工mem_uid
         * type        : 节点类型标识: 1=部门, 2=员工
         *
         * 类型   type  id            pid              label             name        value
         * 员工    2    uc(mem_uid)   uc(dpId)         uc(memUsername)   uc(dpName)  uc(memUserid)
         * 部门    1    uc(dpId)      uc(dpParentid)   uc(dpName)        uc(dpName)  uc(dpThirdid)
         */
        $memberList = [];
        foreach ($deptList as $depart) {
            $departId = $depart['dpId'];
            $dpThirdId = $depart['dpThirdid'];
            // 给树添加员工类型
            if (isset($members[$departId])) {
                foreach ($members[$departId] as $member) {
                    $memberList[] = [
                        'id' => $member['mem_uid'],
                        'pid' => $departId,
                        'lev' => $depart['dpLevel'],
                        'type' => self::TYPE_DEPART,
                        'label' => $member['user_name'],
                        'name' => $member['dept_name'],
                        'value' => $member['user_id'],
                        'disabled' => $member['disabled'],
                        'checked' => $member['checked']
                    ];
                }
            }
            // 给树添加部门类型
            $memberList[] = [
                'id' => $departId,
                'pid' => $depart['dpParentid'],
                'lev' => $depart['dpLevel'],
                'type' => self::TYPE_MEMBER,
                'label' => $depart['dpName'],
                'name' => $depart['dpName'],
                'value' => $dpThirdId,
                'disabled' => false,
                'checked' => false
            ];
        }

        // 数组转化tree
        $tree = MemberService::listToTree($memberList);

        $tree = array_merge($tree, $notDpMember);

        return ['data' => $tree];
    }

    /**
     * 获取指定参数二维码
     * @param Request $request
     *
     * @return array
     */
    public function getQrCode(Request $request)
    {

        $params = $request->only(['enterprise_id', 'domain', 'identifier', 'extend']);

        $rules = [
            'enterprise_id' => 'required|integer|min:1',
            'domain' => 'required',
            'identifier' => 'required',
            'extend' => 'required|array',
            'extend.channel_id' => 'required',
            'extend.channel' => 'required',
            'extend.name' => 'required',
            'extend.type' => [
                'required',
                Rule::in(self::QR_TYPE_ARRAY),
            ],
            'extend.is_verify' => 'required',
            'extend.members' => 'required|array',
            'extend.members.*.value' => 'required',
            'extend.members.*.label' => 'required',
            'extend.members.*.mem_uid' => 'required',
            'extend.members.*.type' => [
                'required',
                Rule::in(self::TYPE_ARRAY),
            ],
        ];

        $message = [
            'enterprise_id.required' => '企业ID不能为空',
            'enterprise_id.integer' => '企业ID类型错误',
            'enterprise_id.min' => '企业ID不能为零',
            'domain.required' => '企业域名不能为空',
            'extend.required' => '二维码参数[extend]不能为空',
            'extend.array' => '二维码参数[extend]类型错误',
            'extend.channel_id.required' => '二维码渠道id不能为空',
            'extend.channel.required' => '二维码渠道名称不能为空',
            'extend.name.required' => '二维码名称不能为空',
            'extend.type.required' => '二维码类型不能为空',
            'extend.is_verify.required' => '二维码添加设置名称不能为空',
            'extend.members.required' => '二维码参数[members]不能为空',
            'extend.members.array' => '二维码参数[members]类型错误',
            'extend.members.*.value.required' => '二维码标识数值[value]不能为空',
            'extend.members.*.label.required' => '二维码标识名称[label]不能为空',
            'extend.members.*.type.required' => '二维码标识类型[type]不能为空',
            'extend.members.*.mem_uid.required' => '二维码标识类型[mem_uid]不能为空',
        ];

        Validator::make($params, $rules, $message)->validate();


        $data = ContactQrCodeParams::query()
            ->join('contact_qr_code', 'contact_qr_code_params.contact_qr_code_id', '=', 'contact_qr_code.id')
            ->where([
                'contact_qr_code.enterprise_id' => $params['enterprise_id'],
                'contact_qr_code.source_type' => $params['extend']['source_type'],
                'contact_qr_code.type' => $params['extend']['type'],
                'contact_qr_code.channel_id' => $params['extend']['channel_id'],
                'contact_qr_code_params.enterprise_id' => $params['enterprise_id'],
                'contact_qr_code_params.object' => $params['extend']['members'][0]['mem_uid'],
            ])
            ->first(['contact_qr_code.official_url as official_url', 'contact_qr_code.id']);

        //$transformUrl = $data['official_url'] ?? '';

        return success([
            'data' => $data
        ]);
    }

    /**
     * 更新小程序路径
     *
     * @param Request $request
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     *
     */
    public function updatePoster(Request $request)
    {
        Log::info("更新小程序路径参数：", $request->all());
        $rules = [
            "enterprise_id" => "required|integer",
            "version" => "string"
        ];

        $message = [
            "enterprise_id.*" => "企业ID类型错误",
        ];

        $params = $this->validate($request, $rules, $message);
        $exists = DB::table("poster_release_enterprise")
            ->where("enterprise_id", $params["enterprise_id"])
            ->first();
        if (empty($exists)) {
            DB::table("poster_release_enterprise")->insert([
                "id" => idGenerate(),
                "enterprise_id" => $params["enterprise_id"],
                "marketing_path" => json_encode(Constant::MARKETING_PATH_NEW),
                "created_at" => time(),
                "updated_at" => time()
            ]);
        }

        return success([]);
    }
}