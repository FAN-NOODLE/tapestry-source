<?php

namespace App\Services;

use App\Exceptions\Error;
use App\Library\Uc\Member;
use App\Models\ContactQrCode;
use App\Models\ContactQrCodeParams;
use App\Models\Scrm\Members;
use Exception;
use Illuminate\Support\Facades\DB;

class ContactQrCodeParamsService extends BaseService
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

    public function getQrCode(array $params)
    {
        $ret = ContactQrCodeParams::query()
            ->join('contact_qr_code', 'contact_qr_code_params.contact_qr_code_id', '=', 'contact_qr_code.id')
            ->where([
                'contact_qr_code.enterprise_id'        => $params['enterprise_id'],
                'contact_qr_code.source_type'          => $params['extend']['source_type'],
                'contact_qr_code.type'                 => $params['extend']['type'],
                'contact_qr_code.channel_id'           => $params['extend']['channel_id'],
                'contact_qr_code.name'                 => $params['extend']['name'],
                'contact_qr_code_params.enterprise_id' => $params['enterprise_id'],
                'contact_qr_code.deleted_at'           => 0,
            ])
            ->first(['contact_qr_code.official_url as official_url', 'contact_qr_code.id']);
        return is_null($ret) ? [] : $ret->toArray();
    }

    public function createQrCode(array $params)
    {
        // 名称重复检测
        $isExistName = ContactQrCode::query()
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'name'          => $params['extend']['name'] . rand(100, 999),
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
                $memUserId                   = $member['memUserid'];
                $deptMemUIdArr[$dpThirdId][] = $member['memUid'];
                $allMemUIdsArr[$memUserId]   = $member['memUid'];
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
            $allDeptIdArr[$dpThirdId]  = $depart['dpId'];
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
        $insert        = [];
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
                $insert[]        = [
                    'enterprise_id' => $params['enterprise_id'],
                    'type'          => $member['type'],
                    'object'        => $allMemUIdsArr[$memberUserId] ?? '',
                    'object_extra'  => $allMemUIdsArr[$memberUserId] ?? '',
                    'created_at'    => time(),
                    'updated_at'    => time(),
                ];
            }
            // 部门数据初始化
            if (self::TYPE_DEPART == $member['type']) {
                $dpThirdId = $member['value'];
                // UC部门集合
                $memberPartIds[] = $dpThirdId;
                $insert[]        = [
                    'enterprise_id' => $params['enterprise_id'],
                    'type'          => $member['type'],
                    'object'        => $allDeptIdArr[$dpThirdId] ?? '',
                    'object_extra'  => implode(',', $deptMemUIdArr[$dpThirdId] ?? []),
                    'created_at'    => time(),
                    'updated_at'    => time(),
                ];
            }
        }

        try {
            DB::beginTransaction();
            // 本地存储二维码
            $contactQr = ContactQrCode::query()->create([
                'enterprise_id' => $params['enterprise_id'],
                'is_verify'     => $params['extend']['is_verify'],
                'channel_id'    => $params['extend']['channel_id'],
                'channel'       => $params['extend']['channel'],
                'name'          => $params['extend']['name'],
                'source_type'   => $params['extend']['source_type'] ?? 0,
                'type'          => $params['extend']['type'],
            ]);
            // 二维码参数外键更新
            foreach ($insert as &$item) {
                $item['contact_qr_code_id'] = $contactQr['id'];
            }
            // 扩展数据存储
            $insert && ContactQrCodeParams::query()->insert($insert);
            // UC存储并获得二维码id(单人)
            $ucData = [
                "ecwScene"      => self::SHOW_BY_QR_CODE,
                'ecwType'       => $params['extend']['type'],
                'ecwRemark'     => $params['extend']['channel'],
                "ecwSkipVerify" => $params['extend']['is_verify'],
                "ecwState"      => $contactQr['id']
            ];
            if (self::QR_SINGLE == $params['extend']['type']) {
                $ucData['memUserid'] = current($memberUserIds);
            } else {
                // UC存储并获得二维码id(多人)
                $ucData['memPartyids'] = $memberPartIds;
                $ucData['memUserids']  = $memberUserIds;
            }
            // 成员名称
            $name   = $params['extend']['members'][0]['label'];
            $result = (new Member())->getContactQrCode($params['domain'], $ucData);
            if (!is_array($result)) {
                $result = json_decode($result, true);
                if (isset($result['code'])) {
                    $upWhere = [
                        'id'            => $params['extend']['id'],
                        'enterprise_id' => $params['enterprise_id']
                    ];
                    $upData  = ['open_userid'            => $result['code']];
                    Members::updateMember($upWhere, $upData);

                    switch ($result['code']) {

                        case 40098:
                            throw new Error(40098, '成员' . $name . '尚未实名认证');
                        case 41041:
                            throw new Error(41041, '二维码一次最大添加人数不超过' . self::MAX_MEMBER_COUNT);
                        case 41045:
                            throw new Error(41045, 'API添加「联系我」达到数量上限');
                        case 41054:
                            throw new Error(41054, '成员' . $name . '尚未激活');
                        case 43009:
                            throw new Error(43009, '企业未验证主体');
                        case 50002:
                            throw new Error(50002, '成员' . $name . '不在权限范围');
                        case 61000:
                            throw new Error(61000, '永久二维码超过每个员工的限制');
                        case 60129:
                            throw new Error(60129, '成员' . $name . '手机和邮箱都为空');
                        case 84074:
                            throw new Error(84074, '没有外部联系人权限');
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
            $contactQr['uc_id']         = $result['id'];
            $contactQr['config_id']     = $result['configId'];
            $contactQr['official_url']  = $officialUrl;
            $contactQr['transform_url'] = $transformUrl;
            $contactQr->save();
            DB::commit();
        } catch (Exception $e) {
            DB::rollBack();
            throw new Error($e->getCode() ?: 10000, $e->getMessage() ?: '二维码创建异常');
        }
        $contactQr['members'] = $params['extend']['members'];
        return $contactQr;
    }
}
