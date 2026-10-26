<?php
/**
 * Created by PhpStorm.
 * User: wangpengfei
 * Date: 2019/11/15
 * Time: 2:11 PM
 */

namespace App\StateMachine\ActivityCenter\Group;

use App\Caches\ActivityCenter\Common\ActivityInitiateCache;
use App\Models\ActivitySku;
use App\Models\ActivityStatus;
use App\Models\Error;
use App\Models\GroupInitiate;
use App\StateMachine\ActivityCenter\Machine;
use Illuminate\Support\Facades\DB;

class InitiateTranstion
{
    const ACTION_PAY = 'pay';
    const ACTION_CANCEL = 'cancel';
    const ACTION_SUCCESS = 'success';
    const ACTION_FAIL = 'fail';
    const ACTION_DOWN = 'down';

    private static function transtion()
    {
        $transtion = [
            GroupInitiate::STATUS_DEFAULT => [
                self::ACTION_PAY    => GroupInitiate::STATUS_ING,
                self::ACTION_CANCEL => GroupInitiate::STATUS_CANCEL,
                self::ACTION_DOWN   => GroupInitiate::STATUS_FAIL,
            ],

            GroupInitiate::STATUS_ING => [
                self::ACTION_SUCCESS => GroupInitiate::STATUS_SUCCESS,
                self::ACTION_FAIL    => GroupInitiate::STATUS_FAIL,
                self::ACTION_DOWN    => GroupInitiate::STATUS_FAIL,
            ],
        ];

        return $transtion;
    }

    /**
     * @param $enterpriseId int 企业id
     * @param $groupInitiateId int 主键
     * @param $action string 动作
     * @param $to
     *
     * @return mixed
     */
    public static function run($enterpriseId, $groupInitiateId, $action, $to, $extUpdate = [])
    {
        $transtions = self::transtion();

        $history = ActivityStatus::getHistory($enterpriseId, 'GroupInitiate', $groupInitiateId);

        $fsm = Machine::create($transtions, $history);

        $assertEqualsTo = $fsm->doAction($action);

        // 成功
        if ($to != $assertEqualsTo) {
            notice('group_initiate_status_abnormal');
        }

        $activityInitiateCache = new ActivityInitiateCache();

        try {
            DB::beginTransaction();

            // 保存历史
            ActivityStatus::setHistory($enterpriseId, 'GroupInitiate', $groupInitiateId, $fsm->getHistory());

            $update = array_merge(['status' => $to], $extUpdate);
            $activityInitiateCache->updateInitiateInfoCache($enterpriseId, $groupInitiateId, ActivitySku::TYPE_GROUP, $update);

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();

            // 更新失败，删除缓存
            $activityInitiateCache->deleteInitiateInfoCache($enterpriseId, $groupInitiateId, ActivitySku::TYPE_GROUP);
            Error::log('更新团状态', $e);
            notice('group_initiate_status_abnormal');
        }

        return $to;
    }
}
