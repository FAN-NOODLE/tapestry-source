<?php

namespace App\Services;

use App\Common\Constant;
use App\Exceptions\Error;
use App\Jobs\StoreWorkerGenerateQrCodeJob;
use App\Jobs\SyncHistoryExternalJob;
use App\Library\ImageTools;
use App\Library\Uc\Member;
use App\Models\Channel;
use App\Models\MemberUserRelation;
use App\Models\PartnerEnterpriseMap;
use App\Models\PartnerMember;
use App\Models\Scrm\Members;
use App\Models\StaffShare;
use App\Models\StoreWorker;
use App\Models\UpdateWorkerRecord;
use App\Models\UserStoreWorkerRelation;
use App\Services\ShareChain\StaffShareStrategy;
use App\Services\ShareChain\Substance;
use CrsSdk\Boss\Auth;
use CrsSdk\Boss\Enterprise;
use CrsSdk\MicroAuthority\Account;
use CrsSdk\MicroAuthority\BuiltInRole;
use CrsSdk\Shop\Store;
use Exception;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Log;
use Illuminate\Support\Facades\Redis;

class StoreWorkerService extends BaseService
{
    /**
     * 新增店员
     *
     * @var int
     */
    const ADD_WORKER = 1;

    /**
     * 编辑
     *
     * @var int
     */
    const EDIT_WORKER = 2;

    /**
     * 导购渠道名称
     */
    const GUIDE_ARRAY = [
        Constant::STORE_WORKER_TYPE_COMMUNITY  => '社群导购',
        Constant::STORE_WORKER_TYPE_ENTERPRISE => '企业导购',
        Constant::STORE_WORKER_TYPE_MANAGER    => '企业管理'
    ];


    /**
     * 获取门店导购列表
     * @param $params
     * @param bool $isPage
     * @param array $selectFields
     * @return array
     * @throws \Exception
     */
    public static function list($params, $isPage = true, $selectFields = ['*'])
    {
        $page = $params['page'] ?? Constant::PAGE;
        $size = $params['size'] ?? Constant::SIZE;

        $where = [];
        if (!empty($params['worker_no'])) {
            $where[] = ['mem_userid', '=', $params['worker_no']];
        }
        if (!empty($params['name'])) {
            $where[] = ['name', 'like', '%' . $params['name'] . '%'];
        }
        if (!empty($params['mobile'])) {
            $where[] = ['mobile', 'like', '%' . $params['mobile'] . '%'];
        }
        if (!empty($params['worker_status'])) {
            $where[] = ['worker_status', '=', $params['worker_status']];
        }
        if (!empty($params['start_date']) && !empty($params['end_date'])) {
            $where[] = ['created_at', '>=', $params['start_date']];
            $where[] = ['created_at', '<=', $params['end_date']];
        }
        // 是否店长筛选
        if (!empty($params['is_boss'])) {
            $where[] = ['is_boss', '=', $params['is_boss']];
        }
        if (!empty($params['store_id'])) {
            $where[] = ['store_id', '=', $params['store_id']];
        }
        if (!empty($params['type'])) {
            $where[] = ['type', '=', $params['type']];
        }


        $sql   = StoreWorker::query()
            ->when(!empty($params['enterprise_id']), function ($query) use ($params) {
                return $query->where('enterprise_id', $params['enterprise_id']);
            })
            ->when(!empty($params['enterprise_ids']), function ($query) use ($params) {
                return $query->whereIn('enterprise_id', $params['enterprise_ids']);
            })
            ->when(!empty($params['mem_uids']), function ($query) use ($params) {
                return $query->whereIn('mem_uid', $params['mem_uids']);
            })
            ->when(!empty($params['ids']), function ($query) use ($params) {
                return $query->whereIn('id', $params['ids']);
            })
            ->when(!empty($params['store_ids']), function ($query) use ($params) {
                return $query->whereIn('store_id', $params['store_ids']);
            })
            ->when(isset($params['contact_qr_code_id']) && !empty($params['contact_qr_code_id']),
                function ($query) use ($params) {
                    return $query->where('contact_qr_code_id', $params['contact_qr_code_id']);
                })
            ->when(isset($params['contact_qr_code_ids']) && !empty($params['contact_qr_code_ids']),
                function ($query) use ($params) {
                    return $query->whereIn('contact_qr_code_id', $params['contact_qr_code_ids']);
                })
            ->when(!empty($params['user_ids']), function ($query) use ($params) {
                return $query->whereIn('user_id', $params['user_ids']);
            })
            ->where($where);
        $total = $sql
            ->count();

        if ($isPage) {
            $sql = $sql->forPage($page, $size);
        }

        $list = $sql->orderBy('created_at', 'desc')
            ->get($selectFields)
            ->toArray();

        $storeIds = array_values(array_unique(array_filter(array_column($list, 'store_id'))));
        if (!empty($storeIds)) {
            switchAppHost(['enterprise_id' => $params['enterprise_id'] ?? current($params['enterprise_ids'])]);
            $storeList = (new Store())->listByIds(['store_ids' => $storeIds, 'is_store' => 1]);
            $storeList = array_column($storeList, 'short_name', 'id');
        } else {
            $storeList = [];
        }
        // 本期没有微信昵称，赋值姓名
        foreach ($list as $key => $value) {
            if (isset($value['source']) && $value['source'] == StoreWorker::SOURCE_WORK) {
                $list[$key]['source'] = '企微店员';
            }
            if (!empty($value['store_id']) && isset($storeList[$value['store_id']])) {
                $list[$key]['store_name'] = $storeList[$value['store_id']];
            } else {
                $list[$key]['store_name'] = '';
            }
        }
        return [
            'list'  => $list,
            'total' => $total,
            'page'  => (int)$page,
            'size'  => (int)$size
        ];
    }

    /**
     * 获取导购详情
     * @param $params
     * @return array
     */
    public static function detail($params): array
    {
        $columns = ['id',
                    'enterprise_id',
                    'type',
                    'recently_logged_time',
                    'recently_logged_store_id',
                    'user_id',
                    'store_id',
                    'worker_no',
                    'name',
                    'mobile',
                    'mem_uid',
                    'mem_userid',
                    'contact_qr_code_id',
                    'avatar',
                    'nickname',
                    'sms_mobile',
                    'is_send_sms',
                    'qr_code',
                    'store_sun_code',
                    'is_boss',
                    'source',
                    'worker_status',
                    'created_at',
                    'updated_at'];
        if (empty($params['worker_id']) && empty($params['enterprise_id']) && empty($params['store_id']) && empty($params['user_id']) && empty($params['mobile']) && empty($params['worker_status'] && empty($params['type']))) {
            return [];
        }
        $storeWorker = StoreWorker::query()
            ->when(!empty($params['worker_id']), function ($query) use ($params) {
                return $query->where('id', '=', $params['worker_id']);
            })
            ->when(!empty($params['enterprise_id']), function ($query) use ($params) {
                return $query->where('enterprise_id', '=', $params['enterprise_id']);
            })
            ->when(!empty($params['store_id']), function ($query) use ($params) {
                return $query->where('store_id', $params['store_id']);
            })
            ->when(!empty($params['user_id']), function ($query) use ($params) {
                return $query->where('user_id', '=', $params['user_id']);
            })
            ->when(!empty($params['mobile']), function ($query) use ($params) {
                return $query->where('mobile', '=', $params['mobile']);
            })
            ->when(!empty($params['mem_uid']), function ($query) use ($params) {
                return $query->where('mem_uid', '=', $params['mem_uid']);
            })
            ->when(!empty($params['worker_status']), function ($query) use ($params) {
                return $query->where('worker_status', '=', $params['worker_status']);
            })
            ->when(!empty($params['type']), function ($query) use ($params) {
                if (is_array($params['type'])) {
                    return $query->whereIn('type', $params['type']);
                } else {
                    return $query->where('type', $params['type']);
                }
            })
            ->when(!empty($params['is_boss']), function ($query) use ($params) {
                return $query->where('is_boss', $params['is_boss']);
            })
            ->when(!empty($params['mem_userid']), function ($query) use ($params) {
                return $query->where('mem_userid', $params['mem_userid']);
            })
            ->orderBy('id', 'desc')
            ->first($columns);

        return !empty($storeWorker) ? $storeWorker->toArray() : [];
    }

    /**
     * 更新导购信息
     * @param array $where
     * @param array $data
     * @return int
     */
    public static function update(array $where, array $data)
    {
        return StoreWorker::query()
            ->where($where)
            ->update($data);
    }

    /**
     * 变更管理员
     * @param $worker
     * @param $workerStatus
     */
    public static function changeAdmin($worker, $workerStatus)
    {
        if ($workerStatus == StoreWorker::IN_SERVICE) {
            // 添加登录权限
            self::addAdminByWorker($worker);
        } else {
            // 登录账号信息发生变化，重置门店助手登录
            self::removeAssistantLoginId($worker['mobile']);

            // 如果该手机号不存在员工了，删除管理员
            $count = StoreWorker::query()
                ->where('enterprise_id', $worker['enterprise_id'])
                ->where('mobile', $worker['mobile'])
                ->where('worker_status', StoreWorker::IN_SERVICE)
                ->count();
            if ($count == 0) {
                self::addAdminByWorker($worker, StoreWorkerService::EDIT_WORKER);
            }
        }
    }

    public static function removeWorkerRedis($workerId)
    {
        $key   = RedisKey('STORE_WORKER_INFO', ['worker_id' => $workerId]);
        $redis = Redis::connection();
        if ($redis->exists($key)) {
            $redis->del([$key]);
        }
    }

    /**
     * 添加管理员
     * @param array $data
     * @param int $isAdd 是否新增
     */
    public static function addAdminByWorker(array $data, int $isAdd = self::ADD_WORKER)
    {
        $params       = [
            'enterprise_id' => $data['enterprise_id'],
            'username'      => $data['name'],
            'mobile'        => !empty($data['mobile']) ? $data['mobile'] : '',
            'operator_id'   => isset($data['operator_id']) ? $data['operator_id'] : 0,
        ];
        $method       = '';
        $data['type'] = $data['type'] ?? Constant::STORE_WORKER_TYPE_STORE;
        Log::info("变更权限参数:", $data);
        switch ($data['type']) {
            case Constant::STORE_WORKER_TYPE_STORE:
                if (StoreWorker::BOSS_YES == $data['is_boss']) {
                    $params['before_mobile'] = !empty($data['before_mobile']) ? $data['before_mobile'] : '';
                    $params['store_id']      = $data['store_id'];
                    $method                  = 'storeManger';
                } else if (StoreWorker::BOSS_NO == $data['is_boss']) {
                    $params['type']     = $isAdd;
                    $params['store_id'] = $data['store_id'];
                    $method             = 'storeMember';
                }
                break;
            case Constant::STORE_WORKER_TYPE_COMMUNITY:
                $params['type'] = $isAdd;
                $method         = 'communityGuide';
                break;
            case Constant::STORE_WORKER_TYPE_ENTERPRISE:
                $params['type'] = $isAdd;
                $method         = 'enterpriseGuide';
                break;
            case Constant::STORE_WORKER_TYPE_MANAGER:
                $params['type']    = $isAdd;
                $params['range']   = $data['range'] ?? [];
                $params['noCheck'] = $data['noCheck'] ?? false;
                $method            = 'enterpriseAdmin';
                break;
        }
        if (!empty($method)) {
            app(BuiltInRole::class)->$method($params);
        }
    }

    /**
     * 重置员工门店助手单点登录
     *
     * @param $mobile
     * author: zhaojie<zhaojie@vchangyi.com>
     * create: 2020/5/15 13:53
     */
    public static function removeAssistantLoginId($mobile)
    {
        $key   = 'assistant:wx_login_' . $mobile;
        $redis = Redis::connection();
        $redis->select(3);
        if ($redis->exists($key)) {
            $redis->del([$key]);
        }
    }

    public static function delUserStoreWorkerRelation(array $where)
    {
        return UserStoreWorkerRelation::query()
            ->where($where)
            ->delete();
    }

    /**
     * 获取导购分享链路信息
     *
     * @param string $memUid
     * @param int $sceneUserId
     * @param int $storeId
     * @return array
     */
    public static function getWorkerByShareLink($enterpriseId, $memUid = '', $sceneUserId = 0, $storeId = 0)
    {
        $workerInfo = StoreWorker::query()
            ->where('enterprise_id', '=', $enterpriseId)
            ->when(!empty($memUid), function ($query) use ($memUid) {
                $query->where('mem_uid', '=', $memUid);
            })
            ->when(!empty($sceneUserId), function ($query) use ($sceneUserId) {
                $query->where('user_id', '=', $sceneUserId);
            })
            ->where('worker_status', '=', StoreWorker::IN_SERVICE)
            ->orderBy('id', 'desc')
            ->get(['id', 'mem_uid', 'store_id', 'type', 'recently_logged_store_id', 'user_id'])
            ->toArray();
        if ($workerInfo) {
            $currentInfo = current($workerInfo);
            $memUid      = $currentInfo['mem_uid'];
            $userId      = $currentInfo['user_id'];
            $workerId    = $currentInfo['id'];
            if ($currentInfo['type'] == Constant::STORE_WORKER_TYPE_STORE) {
                // 获取门店id
                $storeIds = array_column($workerInfo, 'store_id');
                $storeId  = $storeId ? $storeId : (in_array($currentInfo['recently_logged_store_id'], $storeIds) ? $currentInfo['recently_logged_store_id'] : $currentInfo['store_id']);

                $workerInfoByStoreId = array_column($workerInfo, null, 'store_id');
                $workerId            = $workerInfoByStoreId[$storeId]['id'] ?? 0;
            }
        }

        return [
            'mem_uid'   => $memUid,
            'user_id'   => $userId ?? 0,
            'worker_id' => $workerId ?? 0,
            'store_id'  => $storeId
        ];
    }

    /**
     * 查询导购mem_uid
     * @param int $enterpriseId 企业ID
     * @param array $ids 导购ID
     * @return array
     */
    public static function getMemUid(int $enterpriseId, array $ids)
    {
        return StoreWorker::query()
            ->select(DB::Raw("distinct(mem_uid) as mem_uid,id,name,mobile"))
            ->where('enterprise_id', $enterpriseId)
            ->whereIn('id', $ids)
            ->get()
            ->toArray();
    }

    /**
     * 生成企微员工联系我二维码海报
     * @param $data
     * @throws Exception
     */
    public static function generateWorkQrCode($data)
    {
        try {
            // 生成联系我二维码
            $qrCode = self::getContactQrCode($data);
            // 生成海报
            $imgPath = self::generateWorkPoster($qrCode['qr_code'], $data);
            if (!empty($imgPath)) {
                $result = CommonService::uploadImage($imgPath);
                // 保存生成的海报
                if ($result) {
                    unlink($imgPath);
                    $upWhere = [
                        'id'            => $data['id'],
                        'enterprise_id' => $data['enterprise_id']
                    ];
                    $upData  = [
                        'qr_poster_code'            => $result,
                        'qr_poster_code_id' => $qrCode['contact_qr_code_id']];
//                    self::update($upWhere, $upData);
                    Members::updateMember($upWhere, $upData);
                    self::removeWorkerRedis($data['id']);
                }
            }
        } catch (Exception $e) {
            Log::info('生成企业微信员工联系二维码海报失败' . $e);
            throw new Error($e->getCode(), $e->getMessage());
        }

    }

    /**
     * 生成联系我二维码
     * @param $data
     * @return array
     * @throws Exception
     */
    public static function getContactQrCode($data)
    {
        if ($data['type'] == Constant::STORE_WORKER_TYPE_STORE) {
            $store = (new Store())->detail([
                'id'            => $data['store_id'],
                'enterprise_id' => $data['enterprise_id']
            ]);
            $channelId = $store['channel_id'];
            if (!$channelId) {
                // 创建渠道
                $channelInsert = [
                    'channel'       => '门店' . $store['store_no'],
                    'enterprise_id' => $data['enterprise_id'],
                    'source_type'   => 2,
                    'admin_id'      => $data['operator_id'] ?? 1,
                    'admin_name'    => $data['operator_name'] ?? '系统'
                ];
                $channelId     = (int)self::addChannel($channelInsert);
                // 更新渠道id
                if ($channelId) {
                    (new Store())->changeChannel([
                        'store_id'      => $data['store_id'],
                        'enterprise_id' => $data['enterprise_id'],
                        'channel_id'    => $channelId
                    ]);
                }
            }
            $channel = '门店' . $store['store_no'];
            $rqName  = '门店' . $store['store_no'] . $data['name'];
        } else {
            $channel       = self::GUIDE_ARRAY[$data['type']];
            $rqName        = $data['name'] . $data['memUserId'];
            $channelInsert = [
                'enterprise_id' => $data['enterprise_id'],
                'admin_id'      => $data['operator_id'] ?? 0,
                'admin_name'    => $data['operator_name'] ?? '系统',
                'channel'       => $channel,
                'source_type'   => Constant::STORE_WORKER_TYPE_COMMUNITY,
            ];
            $channelId     = (int)self::addChannel($channelInsert);
        }

        // 组装二维码参数
        $params = [
            'enterprise_id' => $data['enterprise_id'],
            'domain'        => $data['domain'],
            'identifier'    => 'contact',
            'extend'        => [
                'id'          => $data['id'],
                'channel_id'  => $channelId,
                'channel'     => $channel,
                'name'        => $rqName,
                'type'        => 1, // 二维码类型:1.单人, 2.多人
                'is_verify'   => 1, // 是否需要验证:0、是, 1、否
                'source_type' => 2, // 渠道类型 2=门店商城
                'members'     => [
                    [
                        'value'   => $data['memUserId'],
                        'label'   => $data['name'],
                        'type'    => 2, // 对象类型:1 、员工, 2、部门
//                        'mem_uid' => $data['mem_uid']
                        'mem_userid' => $data['memUserId']
                    ],
                ],
            ],
        ];
        // 生成联系我二维码
        $workerQrCode    = (new ContactQrCodeParamsService())->getQrCode($params);
        $qrCode          = $workerQrCode['official_url'] ?? '';
        $contactQrCodeId = $workerQrCode['id'] ?? 0;
        if (empty($qrCode)) {
            try {
                $result = (new ContactQrCodeParamsService())->createQrCode($params);;
                $qrCode = $data['memQrCode'];
                if (!empty($result)) {
                    $qrCode          = $result['transform_url'];
                    $contactQrCodeId = $result['id'];
                }
            } catch (Exception $e) {
                Log::info('生成二维码失败getContactQrCode:' . $e->getCode() . $e->getMessage());
                throw new Error($e->getCode(), $e->getMessage());
            }
        }
        return [
            'contact_qr_code_id' => $contactQrCodeId,
            'qr_code'            => $qrCode
        ];
    }

    /**
     * 生成联系我海报
     * @param string $qrCode
     * @param array $worker
     * @return string url
     */
    public static function generateWorkPoster($qrCode, $worker)
    {
        $imgTmpPath      = '';
        $isDefaultAvatar = false;
        $imageObj        = [];
        try {
            $imageTools = new ImageTools();

            // 如果没有qrcode，取消生成
            $qrCode = str_replace(' ', '', $qrCode);
            if (empty($qrCode)) {
                return $imgTmpPath;
            }
            // 没有头像设置默认头像
            if (empty($worker['avatar']) || strpos($worker['avatar'], 'DefaultAvatar') !== false) {
                $isDefaultAvatar = true;
                // $worker['avatar'] = base_path('public/sourceImages/store_worker_header.png');
                $worker['avatar'] = imagecreatefrompng('public/sourceImages/store_worker_header.png');
            } else {
                // $worker['avatar'] = $imageTools->getImage($worker['avatar']);
                list($worker['avatar'], $httpcode) = $imageTools->getImageSource($worker['avatar']);
            }
            // 存储公共元素
//            $appInfo['logo'] = $imageTools->getImage($appInfo['logo']);
            // $qrCode          = $imageTools->getImage($qrCode);
            list($qrCode, $httpcode) = $imageTools->getImageSource($qrCode);
            // 头像和logo处理成圆形
//            $appInfo['logo'] = $imageTools->toCircular($appInfo['logo']);
            if (!$isDefaultAvatar) {
                $worker['avatar'] = $imageTools->toCircularBySource($worker['avatar']);
            }

            // 海报背景图片
            $bg_file = base_path('public/sourceImages/store_worker_' . $worker['enterprise_id'] . '.png');
            if (!is_file($bg_file)) {
                $bg_file = base_path('public/sourceImages/store_worker_default.png');
            }
            $imageObj['img'] = imagecreatefrompng($bg_file);
            $width           = imagesx($imageObj['img']);
            //字体文件
            $font       = base_path('public/fonts/msyh.ttc');
            $fontWeight = base_path('public/fonts/msyhbd.ttc');

            // 员工名称
            $fontColor = imagecolorallocate($imageObj['img'], 51, 51, 51);
            $fontSize  = 30;
            $fontBox   = imagettfbbox($fontSize, 0, $fontWeight, $worker['name']);
            imagettftext($imageObj['img'], $fontSize, 0, ceil(($width - $fontBox[2]) / 2), 315, $fontColor, $fontWeight, $worker['name']);
            // 导购展示
            /*   $fontSize = 22;
               $jobName  = $appInfo['name'] . '专属导购';
               $fontBox  = imagettfbbox($fontSize, 0, $font, $jobName);
               imagettftext($imageObj['img'], $fontSize, 0, ceil(($width - $fontBox[2]) / 2), 360, $fontColor, $font, $jobName);
            */
            // 小程序名称
            $fontSize = 30;
//            $fontBox  = imagettfbbox($fontSize, 0, $font, $appInfo['name']);
            $x = ceil(($width - $fontBox[2] - 130) / 2);
//            imagettftext($imageObj['img'], $fontSize, 0, $x + 180, 1200, $fontColor, $font, $appInfo['name']);


//            $newFontColor = imagecolorallocate($imageObj['img'], 225, 225, 225);
            // 竖线
//            imagettftext($imageObj['img'], 30, 0, $x + 130, 1200, $newFontColor, $font, '|  ');

            // Logo
//            list($l_w, $l_h) = getimagesize($appInfo['logo']);
//            $imageObj['logo'] = $imageTools->createImageFromFile($appInfo['logo']);
//            imagecopyresampled($imageObj['img'], $imageObj['logo'], $x, 1140, 0, 0, 100, 100, $l_w, $l_h);
            // 头像
            // list($g_w, $g_h) = getimagesize($worker['avatar']);
            $g_w = imagesx($worker['avatar']);
            $g_h = imagesy($worker['avatar']);

            //$imageObj['avatar'] = $imageTools->createImageFromFile($worker['avatar']);
            imagecopyresampled($imageObj['img'], $worker['avatar'], 306, 96, 0, 0, 140, 140, $g_w, $g_h);

            //二维码
            // list($code_w, $code_h) = getimagesize($qrCode);
            $code_w = imagesx($qrCode);
            $code_h = imagesy($qrCode);

            //$imageObj['codeImg'] = $imageTools->createImageFromFile($qrCode);
            imagecopyresampled($imageObj['img'], $qrCode, 228, 490, 0, 0, 300, 300, $code_w, $code_h);


            //输出图片
            $path = resource_path('tempimg/');
            if (!is_dir($path)) {
                mkdir($path, 0777, true);
            }
            $imgTmpPath = $path . $worker['userid'] . '.png';
            imagepng($imageObj['img'], $imgTmpPath);
            return $imgTmpPath;
        } catch (Exception $e) {
            Log::info('生成海报失败' . $e);
        } finally {
            //释放空间
            foreach ($imageObj as $key => $value) {
                imagedestroy($value);
            }
//            if (!empty($appInfo)) {
//                if (file_exists($appInfo['logo'])) {
//                    unlink($appInfo['logo']);
//                }
//            }
//            if (!$isDefaultAvatar && file_exists($worker['avatar'])) {
//                unlink($worker['avatar']);
//            }
//            if (file_exists($qrCode)) {
//                unlink($qrCode);
//            }
        }
        return $imgTmpPath;
    }

    /**
     * 获取小程序名称
     * @param $enterpriseId
     * @return array|void
     * @throws Exception
     */
    public static function getMiniAppInfo($enterpriseId)
    {
        // 尝试获取小程序简称
        $enterprise = (new Enterprise())->enterpriseMerchantSettingInfo(['enterprise_id' => $enterpriseId]);
        if (empty($enterprise)) {
            return [];
        }

        if (empty($enterprise['data']['short_name'])) {
            // 无小程序简称，获取小程序全名
            $miniApp = (new Auth())->authDetail(['auth_type' => Constant::APP_TYPE, 'enterprise_id' => $enterpriseId]);
            if (empty($miniApp) || empty($miniApp['data'])) {
                return [];
            }
            $appInfo = [
                'name' => $miniApp['data']['name'],
                'logo' => $miniApp['data']['logo'],
            ];
        } else {
            // 有小程序简称
            $appInfo = [
                'name' => $enterprise['data']['short_name'],
                'logo' => $enterprise['data']['logo'],
            ];
        }
        return $appInfo;
    }

    /**
     * 根据手机号获取门店ID
     * @param string $mobile
     * @param int $enterpriseId
     * @return array
     */
    public static function getStoreIdsByMobile($mobile, $enterpriseId)
    {
        return StoreWorker::query()
            ->where('mobile', $mobile)
            ->where('enterprise_id', $enterpriseId)
            ->get(['store_id'])
            ->toArray();
    }

    /**
     * 生成导购编号
     * @param $enterpriseId
     * @param string $identifier
     * @return false|string
     */
    private static function generateWorkerNo($enterpriseId, $identifier = 'DG')
    {
        while (true) {
            $no     = $identifier . date('YmdHis') . mt_rand(100, 999);
            $coupon = StoreWorker::query()->where('worker_no', $no)->where('enterprise_id', $enterpriseId)->first(['id']);
            if (empty($coupon)) {
                return $no;
            }
        }

        return false;
    }

    /**
     * 添加一个导购信息
     * @param array $data
     * @param string $mobile
     * @return array
     */
    public static function addStoreWorker(array $data, $mobile = '')
    {
        // 添加失败的用户信息
        $lostMemUid = [];
        try {
            // 获取企业微信员工信息
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $data['enterprise_id']])['data'];
            if (!is_array($data['mem_uid'])) {
                $data['mem_uid'] = [$data['mem_uid']];
            }
            // 添加成功的workerId集合
            $workerIds = [];
            foreach ($data['mem_uid'] as $item) {
                $data['mem_uid'] = $item;
                $result          = self::addOneWorker($data, $enterprise['domain'], $mobile);
                if (isset($result['name'])) {
                    $lostMemUid[] = $result;
                }
                if (isset($result['workerId'])) {
                    $workerIds[] = $result['workerId'];
                }
            }

            if (!empty($workerIds)) {
                // 导购商城用户分享链路解绑
                $staffShareStrategy = new Substance(new StaffShareStrategy());
                $staffShareStrategy->unbindingChange($data['enterprise_id'], [], $workerIds, Constant::GUIDE_TYPE);

                $jobData = [
                    'worker_ids'    => $workerIds,
                    'enterprise_id' => $data['enterprise_id'],
                    'store_id'      => $data['store_id'],
                    'domain'        => $enterprise['domain']
                ];
                // 添加店长，只能走异步生成二维码和太阳码
                if ($mobile) {
                    dispatch((new StoreWorkerGenerateQrCodeJob($jobData))->delay(5));
                } else {
                    dispatch(new StoreWorkerGenerateQrCodeJob($jobData));
                }
            }
            return ['list' => $lostMemUid];
        } catch (Exception $e) {
            Log::info('添加员工失败' . $e);
            throw new Error($e->getCode(), $e->getMessage(), 200);
        }
    }

    /**
     * 依次添加店员
     * @param $input
     * @param $domain
     * @param $inputMobile
     * @return array|void
     */
    private static function addOneWorker($input, $domain, $inputMobile)
    {
        try {
            // 获取企业微信员工信息
            $member = (new Member())->detail($domain, 'contact', ['memUid' => $input['mem_uid']]);
            if (empty($member)) {
                return [];
            }
            // 根据手机号或mem_uid判断是否重复添加员工
            $mobile = $member['memMobile'];
            if ($input['is_boss'] == StoreWorker::BOSS_YES && !empty($inputMobile)) {
                $mobile = $inputMobile;
            }
            $memUid = $input['mem_uid'];
            $worker = StoreWorker::query()
                ->where('enterprise_id', '=', $input['enterprise_id'])
                ->where('store_id', '=', $input['store_id'])
                ->where(function ($query) use ($mobile, $memUid) {
                    $query->where('mobile', '=', $mobile)->orWhere('mem_uid', '=', $memUid);
                })->first();
            // 如果添加普通员工，有相同手机号或mem_uid，取消添加
            if (!empty($worker) && $input['is_boss'] == StoreWorker::BOSS_NO) {
                return ['id' => $input['mem_uid'], 'name' => $member['memUsername']];
            }

            $input['mobile'] = $mobile;
            $input['name']   = $member['memUsername'];
            // 如果该员工是店长的话，修改该门店之前的店长为普通职员
            if ($input['is_boss'] == StoreWorker::BOSS_YES) {
                // 添加的店长,走店长变更逻辑
                $existWorker = is_null($worker) ? [] : $worker->toArray();
                self::changeStoreBoss($input, $existWorker);
                if (!empty($worker)) {
                    StoreWorker::query()
                        ->where('id', '=', $worker->id)
                        ->update([
                            'is_boss'       => StoreWorker::BOSS_YES,
                            'mem_userid'    => $member['memUserid'],
                            'worker_status' => StoreWorker::IN_SERVICE,
                            'mobile'        => $mobile,
                        ]);
                    StoreWorker::query()
                        ->where('enterprise_id', $input['enterprise_id'])
                        ->where('mem_uid', $worker->mem_uid)
                        ->update(['user_id' => $input['user_id'] ?? 0]);

                    self::removeWorkerRedis($worker->id);
                    return [];
                }
            }
            // 如果店员不存在，新增
            $workerNo = self::generateWorkerNo($input['enterprise_id']);
            $id       = idGenerate();
            $data     = [
                'id'            => $id,
                'enterprise_id' => $input['enterprise_id'],
                'user_id'       => $input['user_id'] ?? 0,
                'store_id'      => $input['store_id'],
                'worker_no'     => $workerNo,
                'name'          => $member['memUsername'],
                'nickname'      => $member['memUsername'],
                'mobile'        => $mobile,
                'sms_mobile'    => $mobile,
                'mem_uid'       => $input['mem_uid'],
                'mem_userid'    => $member['memUserid'],
                'avatar'        => $member['memFace'],
                'is_boss'       => $input['is_boss'],
                'source'        => $input['source'],
                'worker_status' => StoreWorker::IN_SERVICE
            ];
            StoreWorker::query()->create($data);
            self::addAdminByWorker($input);
            $data['operator_id'] = isset($input['operator_id']) ?? 0;

            return ['workerId' => $data['id']];
        } catch (Exception $e) {
            Log::info('添加员工失败addOneWorker' . $e, $input);
            throw new Error($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 变更门店店长
     * @param array $data
     * @param array $worker
     * @throws Exception
     */
    private static function changeStoreBoss(array &$data, array $worker)
    {
        // 查询之前店长信息
        $beforeStoreManager = StoreWorker::query()
            ->where('store_id', $data['store_id'])
            ->where('enterprise_id', $data['enterprise_id'])
            ->where('is_boss', StoreWorker::BOSS_YES)
            ->first(['id', 'mobile']);
        if ($beforeStoreManager && $beforeStoreManager->mobile != $data['mobile']) {
            try {
                self::checkGuideAccount($data['enterprise_id'], $beforeStoreManager->mobile, Constant::GUIDE_STORE_MEMBER);
                $data['before_mobile'] = $beforeStoreManager->mobile;
            } catch (Exception $exception) {
                Log::debug('旧手机号添加店员角色互斥');
            }

            // 变更之前的店长为普通店员
            $beforeStoreManager->update(['is_boss' => StoreWorker::BOSS_NO]);
        }
        // 查询要变更的人是否已存在此店中，如果存在，更新类型为店长
        Log::debug('去除变更人的导购角色-debug', $worker);
        if (!empty($worker)) {
            $worker['is_boss'] = StoreWorker::BOSS_YES;
            if (!empty($data['before_mobile'])) {
                $worker['before_mobile'] = $data['before_mobile'];
            }
            self::addAdminByWorker($worker);
        }
    }

    /**
     * 验证导购角色
     * @param int $enterpriseId
     * @param string $mobile
     * @param string $wantAddRole
     * @param int $storeId
     * @throws Exception
     */
    public static function checkGuideAccount(int $enterpriseId, string $mobile, string $wantAddRole, int $storeId = 0): void
    {
        $searchAdminRoleParams = [
            'enterprise_id' => $enterpriseId,
            'mobile'        => $mobile
        ];
        // 查询管理员所有的角色
        try {
            $adminRoles = (new Account())->adminRoles($searchAdminRoleParams);
            if (!empty($adminRoles)) {
                $guideAdminRoleMap = array_column($adminRoles, 'mark');
                // 管理员拥有的内置角色
                $hadGuideRole = array_intersect($guideAdminRoleMap, Constant::GUIDE_DEFAULT_ROLE);
                // 想添加店长或者店员，判断有无企业导购、社群导购、企业管理角色
                if (in_array($wantAddRole, Constant::GUIDE_IN_STORE_ROLE)) {
                    $hadOtherRole = array_intersect($hadGuideRole, Constant::GUIDE_NOT_STORE_ROLE);
                    // 说明有企业导购、社群导购、企业管理角色，不能添加店长或者店员成功
                    if (!empty($hadOtherRole)) {
                        notice(700001);
                    }
                } else {
                    // 想添加企业导购、社群导购、企业管理角色，不能有门店类角色
                    $hadStoreRole = array_intersect($hadGuideRole, Constant::GUIDE_IN_STORE_ROLE);
                    if (!empty($hadStoreRole)) {
                        notice(700001);
                    }
                }
            }
        } catch (\App\Exceptions\Error $error) {
            Log::error("调用接口MicroAuthority:Account->adminRoles接口出错，错误信息：{$error->getMessage()}, 错误code: {$error->getCode()}", $searchAdminRoleParams);
            if ($error->getCode() != 700012) {
                throw $error;
            }
        }
    }

    /**
     * 编辑导购信息
     * @param array $data
     * @return mixed
     */
    public static function editStoreWorker(array $data)
    {
        $count = StoreWorker::query()
            ->where('enterprise_id', $data['enterprise_id'])
            ->where('store_id', $data['store_id'])
            ->where('mobile', $data['mobile'])
            ->where([['id', '!=', $data['id']]])
            ->count();
        if ($count > 0) {
            notice(2300005);
        }

        // 门店店长不能修改手机号
        $worker = StoreWorker::query()
            ->where('enterprise_id', $data['enterprise_id'])
            ->where('store_id', $data['store_id'])
            ->where('id', $data['id'])
            ->first();
        if (empty($worker)) {
            return false;
        }
        if ($worker->is_boss == StoreWorker::BOSS_YES) {
            return false;
        }

        try {
            StoreWorker::query()
                ->where('enterprise_id', $data['enterprise_id'])
                ->where('store_id', $data['store_id'])
                ->where('id', $data['id'])
                ->where('is_boss', StoreWorker::BOSS_NO)
                ->update(['mobile' => $data['mobile'], 'name' => $data['name']]);
            self::removeWorkerRedis($data['id']);
            // 如果店员名称发生修改，重新生成海报
            if ($worker->name != $data['name']) {
                $enterprise = (new Enterprise())->detail(['enterprise_id' => $data['enterprise_id']])['data'];
                $jobData    = [
                    'worker_ids'    => ['worker_id' => $data['id']],
                    'enterprise_id' => $data['enterprise_id'],
                    'store_id'      => $data['store_id'],
                    'domain'        => $enterprise['domain']
                ];
                dispatch(new StoreWorkerGenerateQrCodeJob($jobData));
            }
            // 暂时去掉管理员新增逻辑 打开
            if ($worker->mobile != $data['mobile']) {
                $beforeMobile = '';
                $detail       = StoreWorker::query()
                    ->where('enterprise_id', $data['enterprise_id'])
                    ->where('mobile', $worker->mobile)
                    ->where('worker_status', StoreWorker::IN_SERVICE)
                    ->where('id', '!=', $data['id'])
                    ->first();
                if (!empty($detail)) {
                    $params = $detail->toArray();
                    self::addAdminByWorker($params);
                } else {
                    $beforeMobile = $worker->mobile;
                }
                $params                  = $worker->toArray();
                $params['mobile']        = $data['mobile'];
                $params['name']          = $data['name'];
                $params['before_mobile'] = $beforeMobile;
                self::addAdminByWorker($params);
                if (empty($beforeMobile)) {
                    self::addAdminByWorker($worker->toArray());
                }
                // 登录账号信息发生变化，重置门店助手登录
                if (!empty($beforeMobile)) {
                    self::removeAssistantLoginId($beforeMobile);
                }
            }
            return true;
        } catch (Exception $e) {
            Log::info('修改员工失败' . $e, $data);
            notice(2300002);
        }
    }

    public static function selectWorkerDetailList(array $params)
    {
        return StoreWorker::query()
            ->where($params)
            ->get()
            ->toArray();
    }

    /**
     * 添加店员
     * @param $input
     * @return bool
     */
    public static function addGuide($input)
    {
        try {
            // 获取企业信息
            $enterprise = (new Enterprise())->detail(['enterprise_id' => $input['enterprise_id']])['data'];

            $memUids = array_column($input['member'], 'mem_uid');
            // 查询离岗导购
            $allGuideList      = StoreWorker::query()
                ->where('enterprise_id', $input['enterprise_id'])
                ->where('store_id', 0)
                ->whereIn('mem_uid', $memUids)
//                ->where('type', $input['type'])
                ->get(['id', 'mem_uid', 'type', 'worker_no', 'worker_status'])
                ->toArray();
            $guideList         = [];
            $existGuideMemuids = [];
            foreach ($allGuideList as $value) {
                if ($value['worker_status'] == StoreWorker::OUT_SERVICE) {
                    $guideList[$value['mem_uid']] = $value;
                } else {
                    $existGuideMemuids[] = $value['mem_uid'];
                }
            }


            $userList = self::getUserIds([
                'enterprise_id' => $input['enterprise_id'],
                'mem_uid'       => $memUids
            ]);
            $userList = array_column($userList, 'user_id', 'mem_uid');

            // 添加成功的workerId集合
            $workerIds = [];
            $memberMobile = [];
            foreach ($input['member'] as $item) {
                if (in_array($item['mem_uid'], $existGuideMemuids)) {
                    StoreWorker::query()
                        ->where('enterprise_id', $input['enterprise_id'])
                        ->where('mem_uid', $item['mem_uid'])
                        ->where('user_id', $item['user_id'])
                        ->update(['mobile' => $item['mobile'], 'name' => $item['user_name']]);
                } else {
                    $data           = $input;
                    $data['member'] = $item;
                    $result         = self::addOneGuide($data, $guideList, $userList);
                    if (isset($result['workerId'])) {
                        $workerIds[] = $result['workerId'];
                    }
                }
                $memberMobile[] = $item['mobile'];

            }

            if (!empty($workerIds)) {
//                $data = [
//                    'worker_ids'    => $workerIds,
//                    'enterprise_id' => $input['enterprise_id'],
//                    'store_id'      => 0,
//                    'domain'        => $enterprise['domain'],
//                    'operator_id'   => $input['operator_id'] ?? 0,
//                    'operator_name' => $input['operator_name'] ?? '系统'
//                ];
//                //异步生成二维码和太阳码
//                dispatch(new StoreWorkerGenerateQrCodeJob($data));
                //同步该导购的外部联系人
                dispatch(new SyncHistoryExternalJob([
                    'enterprise_id' => $input['enterprise_id'],
                    'worker_ids'    => $workerIds
                ]));
                if (!empty($memberMobile)){
                    $memberListIds = Members::query()
                        ->whereIn('mobile', $memberMobile)
                        ->where('enterprise_id', $input['enterprise_id'])
                        ->pluck('id')
                        ->toArray();
                    $data = [
                        'worker_ids'    => $memberListIds,
                        'enterprise_id' => $input['enterprise_id'],
                        'store_id'      => 0,
                        'domain'        => $enterprise['domain'],
                        'operator_id'   => $input['operator_id'] ?? 0,
                        'operator_name' => $input['operator_name'] ?? '系统'
                    ];
                    //异步生成二维码和太阳码
                    dispatch(new StoreWorkerGenerateQrCodeJob($data));
                }
            }
        } catch (Exception $e) {
            Log::info('添加店员失败' . $e);
            throw new Error($e->getCode(), $e->getMessage(), 200);
        }

        return true;
    }

    /**
     * 依次添加导购
     * @param $input
     * @param $guideList
     * @param $userList
     * @return array
     * @author zhangjunchang<zhangjunchang@vchangyi.com>
     * @date 2020/11/5 18:00
     */
    private static function addOneGuide($input, $guideList, $userList)
    {
        try {
            DB::beginTransaction();

            $data = [
                'enterprise_id' => $input['enterprise_id'],
                'store_id'      => 0,
                'name'          => $input['member']['user_name'],
                'nickname'      => $input['member']['user_name'],
                'mobile'        => $input['member']['mobile'],
                'sms_mobile'    => $input['member']['mobile'],
                'mem_uid'       => $input['member']['mem_uid'],
                'mem_userid'    => $input['member']['user_id'] ?? '',
                'avatar'        => $input['member']['avatar'],
                'type'          => $input['type'],
                'range'         => $input['range'] ?? [],
                'user_id'       => $userList[$input['member']['mem_uid']] ?? 0,
                'is_boss'       => StoreWorker::BOSS_NO,
                'source'        => StoreWorker::SOURCE_WORK,
                'worker_status' => StoreWorker::IN_SERVICE
            ];

            if (!empty($guideList[$input['member']['mem_uid']])) {
                StoreWorker::query()
                    ->where('id', $guideList[$input['member']['mem_uid']]['id'])
                    ->update([
                        'worker_status' => StoreWorker::IN_SERVICE,
                        'type'          => $input['type'],
                        'mobile'        => $data['mobile'],
                        'created_at'    => time()
                    ]);

                $data['id']        = $guideList[$input['member']['mem_uid']]['id'];
                $data['worker_no'] = $guideList[$input['member']['mem_uid']]['worker_no'];
            } else {
                $workerNo          = self::generateWorkerNo($input['enterprise_id'], 'DG');
                $id                = idGenerate();
                $data['id']        = $id;
                $data['worker_no'] = $workerNo;
                StoreWorker::query()->create($data);
            }

            $return              = ['workerId' => $data['id']];
            $data['operator_id'] = isset($input['operator_id']) ?? 0;

            self::addAdminByWorker($data);

            DB::commit();
            return $return;
        } catch (Exception $e) {
            DB::rollBack();
            Log::info('添加员工失败addOneGuide' . $e, $input);
            throw new Error($e->getCode(), $e->getMessage());
        }
    }

    /**
     * 获取导购的user_id
     * @param $params
     * @return array
     */
    public static function getUserIds($params)
    {
        return StoreWorker::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->whereIn('mem_uid', $params['mem_uid'])
            ->where('user_id', '<>', 0)
            ->get(['mem_uid', 'user_id'])
            ->toArray();
    }

    /**
     * 移除导购
     * @param $params
     * @return bool
     * @throws Exception
     */
    public static function deleteGuide($params)
    {
        $worker = StoreWorker::query()
            ->where('id', $params['id'])
            ->first();
        if (!$worker) {
            notice(2300006);
        }
        $worker->update([
            'worker_status' => StoreWorker::OUT_SERVICE
        ]);
        $data                = $worker->toArray();
        $data['operator_id'] = $params['operator_id'];
        self::addAdminByWorker($data, self::EDIT_WORKER);
        // 登录账号信息发生变化，重置门店助手登录
        self::removeAssistantLoginId($data['mobile']);
        return true;
    }

    public static function updateWorkerInfo($input)
    {
        $record  = UpdateWorkerRecord::query()
            ->where('domain', $input['domain'])
            ->pluck('worker_id')
            ->toArray();
        $account = new Account();
        foreach ($input['member_list'] as $member) {
            try {
                if ($member['mobile'] != $member['before_mobile'] && !empty($member['mobile'])) {
                    // 更新管理员手机号
                    $account->infoChange([
                        'enterprise_id' => $member['enterprise_id'],
                        'mobile'        => $member['before_mobile'],
                        'new_mobile'    => $member['mobile']
                    ]);
                    self::removeAssistantLoginId($member['before_mobile']);
                }

                $where        = [
                    'id' => $member['id']
                ];
                $updateParams = [
                    'name'   => $member['name'],
                    'mobile' => $member['mobile'],
                    'avatar' => $member['avatar'],
                ];
                self::update($where, $updateParams);

                if (!in_array($member['id'], $record)) {
                    $data = [
                        'worker_ids'    => [$member['id']],
                        'enterprise_id' => $member['enterprise_id'],
                        'store_id'      => $member['store_id'],
                        'domain'        => $input['domain'],
                        'operator_id'   => 1,
                        'operator_name' => '系统'
                    ];
                    //异步生成太阳码
                    dispatch(new StoreWorkerGenerateQrCodeJob($data));
                }
                return true;
            } catch (Exception $e) {
                Log::info('定时修改员工信息失败' . $e, $input);
                UpdateWorkerRecord::query()
                    ->create([
                        'domain'        => $input['domain'],
                        'enterprise_id' => $member['enterprise_id'],
                        'worker_id'     => $member['id'],
                        'result'        => json_encode([
                            'file'    => $e->getFile(),
                            'line'    => $e->getLine(),
                            'code'    => $e->getCode(),
                            'message' => $e->getMessage(),
                            'trace'   => $e->getTraceAsString(),
                        ])
                    ]);
            }
        }
    }

    public static function workerBindUser($params)
    {
        // 绑定员工和商城客户关系
        StoreWorker::query()
            ->where('mem_uid', $params['mem_uid'])
            ->where('enterprise_id', $params['enterprise_id'])
            ->update([
                'user_id' => $params['user_id']
            ]);

        // 导购商城用户分享链路解绑
        $staffShareStrategy = new Substance(new StaffShareStrategy());
        $staffShareStrategy->unbindingChange($params['enterprise_id'], [$params['mem_uid']], [], Constant::GUIDE_TYPE);

        return true;
    }

    /**
     * 生成企微员工联系我二维码海报
     * @param $data
     * @throws Exception
     */
    public static function generateOldWorkQrCode($data)
    {
        try {
            // 生成联系我二维码
            $qrCode = self::getContactQrCode($data);
            // 生成海报
            $path = self::generateWorkPoster($qrCode['qr_code'], $data);
            if (!empty($path)) {
                $result = CommonService::uploadImage($path);
//                $result = (new AzureBlobService())->uploadByFilePath($path);
                // 保存生成的海报
                if ($result) {
                    unlink($path);

                    StoreWorker::query()
                        ->where('id', $data['id'])
                        ->where('store_id', $data['store_id'])
                        ->where('enterprise_id', $data['enterprise_id'])
                        ->update(['qr_code' => $result, 'contact_qr_code_id' => $qrCode['contact_qr_code_id']]);
                    self::removeWorkerRedis($data['id']);
                }
            }
        } catch (Exception $e) {
            echo 'StoreWorkerID:' . $data['id'] . '企业ID:' . $data['enterprise_id'];
            echo '生成企业微信员工联系二维码海报失败' . $e->getMessage();
        }
    }

    public static function editStoreWorkerByMemId(array $params)
    {
        // 修改的手机号不能重复
        $count = StoreWorker::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('store_id', $params['store_id'])
            ->where('mobile', $params['mobile'])
            ->where([['mem_uid', '!=', $params['mem_uid']]])
            ->count();
        if ($count > 0) {
            notice(2300005, [], [], 400);
        }
        // 门店店长不能修改手机号
        $worker = StoreWorker::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('store_id', $params['store_id'])
            ->where('mem_uid', $params['mem_uid'])
            ->first(['id', 'mobile', 'is_boss', 'name']);
        if (empty($worker)) {
            notice(2300006, [], [], 400);
        }
        if ($worker->is_boss == StoreWorker::BOSS_YES) {
            return true;
        }
        try {
            StoreWorker::query()
                ->where('enterprise_id', $params['enterprise_id'])
                ->where('store_id', $params['store_id'])
                ->where('id', $worker->id)
                ->where('is_boss', StoreWorker::BOSS_NO)
                ->update(['mobile' => $params['mobile'], 'name' => $params['name']]);
            self::removeWorkerRedis($worker->id);
            // 如果店员名称发生修改，重新生成海报
            if ($worker->name != $params['name']) {
                $enterprise = (new Enterprise())->detail(['enterprise_id' => $params['enterprise_id']])['data'];
                $data       = [
                    'worker_ids'    => ['worker_id' => $worker->id],
                    'enterprise_id' => $params['enterprise_id'],
                    'store_id'      => $params['store_id'],
                    'domain'        => $enterprise['domain']
                ];
                dispatch(new StoreWorkerGenerateQrCodeJob($data));
            }
        } catch (Exception $e) {
            Log::info('修改员工失败' . $e->getMessage(), $params);
            notice(2300002, [], [], 400);
        }
        return true;
    }

    /**
     * 根据条件获取导购列表
     * @param $params
     * @return array
     */
    public static function allList($params)
    {
        $where = [
            ['user_id', '<>', 0]
        ];
        if (!empty($params['mobile'])) {
            $where[] = ['mobile', '=', $params['mobile']];
        }
        if (!empty($params['worker_status'])) {
            $where[] = ['worker_status', '=', $params['worker_status']];
        }
        // 是否店长筛选
        if (!empty($params['is_boss'])) {
            $where[] = ['is_boss', '=', $params['is_boss']];
        }
        if (!empty($params['store_id'])) {
            $where[] = ['store_id', '=', $params['store_id']];
        }
        if (!empty($params['type'])) {
            $where[] = ['type', '=', $params['type']];
        }
        // 字段筛选
        if (empty($params['columns'])) {
            $selectFields = ['id', 'worker_no', 'name', 'nickname', 'mobile', 'avatar', 'is_boss', 'qr_code', 'source', 'worker_status', 'created_at', 'type',
                             'mem_uid', 'store_id', 'user_id', 'is_send_sms', 'contact_qr_code_id'];
        } else {
            $selectFields = $params['columns'];
        }

        return StoreWorker::query()
            ->when(!empty($params['enterprise_id']), function ($query) use ($params) {
                return $query->where('enterprise_id', $params['enterprise_id']);
            })
            ->when(!empty($params['mem_uids']), function ($query) use ($params) {
                return $query->whereIn('mem_uid', $params['mem_uids']);
            })
            ->when(!empty($params['ids']), function ($query) use ($params) {
                return $query->whereIn('id', $params['ids']);
            })
            ->when(!empty($params['store_ids']), function ($query) use ($params) {
                return $query->whereIn('store_id', $params['store_ids']);
            })
            ->where($where)
            ->orderBy('created_at', 'desc')
            ->get($selectFields)
            ->toArray();
    }

    /**
     * 获取导购信息
     * @param $workerId
     * @return array
     * @author zhangjunchang<zhangjunchang@vchangyi.com>
     * @date 2020/11/9 17:52
     */
    public static function detailByIdOrMemUid($workerId)
    {
        $storeWorker = StoreWorker::query()
            ->where('id', '=', $workerId)
            ->orWhere('mem_uid', '=', $workerId)
            ->first();
        return !empty($storeWorker) ? $storeWorker->toArray() : [];
    }

    public static function getWorkerIdByName($enterpriseId, $name)
    {
        return StoreWorker::query()
            ->where('enterprise_id', $enterpriseId)
            ->where('name', "like", '%' . $name . '%')
            ->pluck('name', 'id')
            ->toArray();
    }

    public static function getWorkersInfo(array $params)
    {
        $workerInfo = self::list($params, false);
        $data       = [];
        if (empty($workerInfo)) {
            return $data;
        }
        $storeIds = array_values(array_unique(array_column($workerInfo, 'store_id')));
        switchAppHost(['enterprise_id' => $params['enterprise_id']]);
        $storeList = (new Store())->listByIds(['store_ids' => $storeIds, 'is_store' => 1]);
        $storeList = array_column($storeList, 'name', 'id');
        foreach ($workerInfo as $val) {
            $data[] = [
                'store_name'  => $storeList[$val['store_id']],
                'worker_name' => $val['name'],
                'worker_id'   => $val['id'],
                'avatar'      => $val['avatar']
            ];
        }
        return $data;
    }

    /**
     * 在职导购或手机号查询门店ID
     * @param array $param
     * @return array
     */
    public static function getStoreIdListByMobile(array $param)
    {
        return StoreWorker::query()
            ->select(DB::raw('distinct store_id as store_id'))
            ->where([
                'enterprise_id' => $param['enterprise_id'],
                'worker_status' => $param['worker_status']
            ])
            ->when(!empty($param['mobile']), function ($query) use ($param) {
                return $query->where(['mobile' => $param['mobile']]);
            })
            ->orderByDesc('created_at')
            ->get()
            ->toArray();
    }

    public static function addChannel(array $data)
    {
        $existChannel = Channel::query()
            ->where([
                'enterprise_id' => $data['enterprise_id'],
                'channel'       => $data['channel']
            ])->first(['id']);

        if ($existChannel) {
            return $existChannel->id;
        }
        $ret = Channel::query()->create($data);
        return $ret->id;
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
     * 获取用户绑定的导购或归属导购
     * @param int $enterpriseId
     * @param int $userId
     *
     * @return array
     * @throws Exception
     */
    public static function getShareWorkerIdByUserId($enterpriseId = 0, $userId = 0, $ascription = false)
    {
        $guideSettingKey = redisKey('ENTERPRISE_GUIDE_SETTING', ['enterprise_id' => $enterpriseId, 'type' => Constant::LOGIN_IDENTITY_GUIDE]);
        $guideSetting    = getOrSetData($guideSettingKey);

        if (empty($guideSetting)) {
            $guideSetting = (new StaffManagerService())->getStaffSetting($enterpriseId, Constant::LOGIN_IDENTITY_GUIDE) ?? [];
            if (!empty($guideSetting)) {
                getOrSetData($guideSettingKey, $guideSetting);
            }
        }

        $commissionType = $guideSetting['commission_type'] ?? 0;
        $fissionType    = $guideSetting['fission_type'] ?? 0;

        $workerId          = $updatedAt = 0;
        $partnerMemberInfo = PartnerEnterpriseMap::query()
            ->where([
                    'enterprise_id'  => $enterpriseId,
                    'user_id'        => $userId,
                    'partner_status' => Constant::PARTNER_STATUS_ENABLE]
            )->first(['partner_id as id']);
        if (!empty($partnerMemberInfo)) {
            $partnerMemberInfo   = $partnerMemberInfo->toArray();
            $worker['worker_id'] = $partnerMemberInfo['id'];
        } else {
            $worker = self::getWorkerByShareLink($enterpriseId, '', $userId, 0);
        }

        if ($worker['worker_id']) {
            // 导购购买商品有佣金
            if ($commissionType == Constant::GUIDE_COMMISSION_TYPE_YES) {
                $workerId = $worker['worker_id'];
            }

            // 归属导购
            if ($ascription) {
                $workerId = $worker['worker_id'];
            }
        } else {
            $shareLink = StaffShare::query()
                ->where('enterprise_id', $enterpriseId)
                ->where('user_id', $userId)
                ->where('binding_type', StaffShare::BINDING_TYPE_YES)
                ->whereIn('type', [Constant::LOGIN_IDENTITY_GUIDE, Constant::LOGIN_IDENTITY_PARTNER])
                ->first(['worker_id', 'level', 'updated_at']);
            if ($shareLink) {
                $shareLink = $shareLink->toArray();
                // level等于1表示直接分享 || level大于1表示间接分享,同时后台裂变分享开关为计算业绩
                if ($shareLink['level'] == 1 || ($shareLink['level'] > 1 && $fissionType == Constant::GUIDE_FISSION_TYPE)) {
                    $workerId = $shareLink['worker_id'];
                }

                // 归属导购
                if ($ascription && $shareLink['level'] >= 1) {
                    $workerId  = $shareLink['worker_id'];
                    $updatedAt = $shareLink['updated_at'];
                }
            }
        }

        return [$workerId, $updatedAt];
    }

    public static function getWorkerListByIds($ids)
    {
        return StoreWorker::query()
            ->whereIn('id', $ids)
            ->get(['id', 'name', 'mobile', 'mem_uid', 'store_id'])
            ->toArray();
    }

    /**
     * 绑定门店员工关系(不带员工参数、顾客参数访问门店)
     *
     * @param $params
     * @return mixed
     */
    public static function storeMemBind($params)
    {
        // 当前顾客是否绑定过门店员工
        // A、绑定过则不重新绑定，延用之前绑定的门店员工关系(结束)
        // B、没绑定过则绑定门店店长(结束)

        // 初始化返回数据
        $returnData = [
            'store_id'  => $params['store_id'], // 门店ID
            'worker_id' => 0, // 员工ID
        ];

        // 获取用户当前绑定导购信息
        $shareRelation = StaffShare::query()
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'user_id'       => $params['user_id'],
                'type'          => Constant::LOGIN_IDENTITY_GUIDE,
            ])
            ->whereIn('binding_type', [StaffShare::BINDING_TYPE_YES, StaffShare::BINDING_TYPE_PARENT_UNBINDING])
            ->first();

        if (!empty($shareRelation->worker_id)) {
            $returnData['worker_id'] = $shareRelation->worker_id;
            return $returnData;
        }

        // 查询门店的店长信息
        $storeBoss = StoreWorker::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('store_id', $params['store_id'])
            ->where('is_boss', StoreWorker::BOSS_YES)
            ->first();

        if (empty($storeBoss)) {
            return [];
        }
        $returnData['worker_id'] = $storeBoss['id'];

        return $returnData;
    }

    /**
     * 绑定门店员工关系(带参数访问门店)
     * @param $params
     * @return array
     */
    public static function storeMemBindWithScene($params)
    {
        // A、带员工参数，绑定门店员工(结束)
        // B、带顾客参数，查询顾客绑定的门店员工进行绑定(结束)
        // C、只带门店参数，调用绑定门店员工关系(不带员工参数、顾客参数访问门店)(结束)

        // 初始化返回数据
        $returnData = [
            'store_id'  => $params['scene_store_id'], // 门店ID
            'worker_id' => 0, // 员工ID
        ];

        // 带员工参数
        if (!empty($params['scene_worker_id'])) {
            return $returnData['worker_id'] = $params['scene_worker_id'];
        }

        // 获取用户当前绑定导购信息
        $shareRelation = StaffShare::query()
            ->where([
                'enterprise_id' => $params['enterprise_id'],
                'user_id'       => $params['scene_user_id'],
                'type'          => Constant::LOGIN_IDENTITY_GUIDE,
            ])
            ->whereIn('binding_type', [StaffShare::BINDING_TYPE_YES, StaffShare::BINDING_TYPE_PARENT_UNBINDING])
            ->first();

        if (!empty($shareRelation->worker_id)) {
            $returnData['worker_id'] = $shareRelation->worker_id;
            return $returnData;
        }

        // 查询门店的店长信息
        $storeBoss = StoreWorker::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('store_id', $params['store_id'])
            ->where('is_boss', StoreWorker::BOSS_YES)
            ->first();

        if (empty($storeBoss)) {
            return [];
        }

        $returnData['worker_id'] = $storeBoss['id'];
        return $returnData;
    }

    /**
     * 获取绑定导购信息
     * @param int $enterpriseId
     * @param int $userId
     * @return array
     * @throws Exception
     */
    public static function getBindingWorker($enterpriseId = 0, $userId = 0)
    {
        $workerBindInfo = [];
        // 查询是否导购
        list($workerId, $updatedAt) = self::getShareWorkerIdByUserId($enterpriseId, $userId, true);
        if (!empty($workerId)) {
            $workerInfo = StoreWorker::query()
                ->where('id', $workerId)
                ->first(['name', 'mem_uid', 'user_id']);

            if (empty($workerInfo)) {
                $workerInfo = PartnerMember::query()
                    ->leftJoin('partner_enterprise_map', 'partner_member.id', '=', 'partner_enterprise_map.partner_id')
                    ->where(['partner_member.id' => $workerId, 'partner_enterprise_map.enterprise_id' => $enterpriseId])
                    ->first(['partner_member.partner_name as name', 'partner_enterprise_map.user_id']);
            }
            $workerBindInfo = [
                'name'      => $workerInfo->name ?? '',
                'mem_uid'   => $workerInfo->mem_uid ?? '',
                'worker_id' => $workerId,
                'user_id'   => $workerInfo->user_id ?? '',
                'bing_time' => $updatedAt
            ];
        }

        return $workerBindInfo;
    }


    /**
     * 获取所有企业id
     * @return array
     */
    public static function getAllEnterpriseId()
    {
        return StoreWorker::query()
            ->select('enterprise_id')
            ->distinct()
            ->pluck('enterprise_id')
            ->toArray();
    }
}
