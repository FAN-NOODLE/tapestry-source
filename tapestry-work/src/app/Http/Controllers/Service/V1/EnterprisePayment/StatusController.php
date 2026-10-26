<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/12/3
 * Time: 3:12 PM
 */

namespace App\Http\Controllers\Service\V1\EnterprisePayment;

use App\Http\Controllers\Service\V1\BaseController;
use App\Http\Requests\EnterprisePaymentRequest;
use App\Exceptions\Error;
use App\Library\RandomRedPacket;
use App\Models\EnterprisePayment;
use App\Models\EnterprisePaymentRandom;
use App\Services\EnterprisePaymentService;
use Illuminate\Support\Facades\DB;
use function PHPSTORM_META\type;

class StatusController extends BaseController
{
    // 删除活动
    const DELETED_YES = 5;

    /**
     * 修改企业微信红包状态
     *
     * @return array
     */
    public function index()
    {
        $params = $this->request->input();
        // 验证数据
        (new EnterprisePaymentRequest())->status($params);

        $detail = EnterprisePayment::query()
            ->where('enterprise_id', $params['enterprise_id'])
            ->where('id', $params['id'])
            ->first();
        if (empty($detail)) {
            throw new Error(110009, '活动信息不存在');
        }

        // 验证状态
        $this->operationStatus($params, $detail);

        $data = [
            'is_sale' => $params['status']
        ];

        // 已发送红包数量缓存key
        $recordCount = EnterprisePayment::SEND_PACK_KEY . $params['id'];
        // 随机红包缓存key
        $lkey = EnterprisePayment::RANDOM_LIST_KEY . $params['id'];

        // 如果是上架且随机红包
        if ($params['status'] == EnterprisePayment::IS_SALE_UP) {
            (new EnterprisePaymentService())->setCount($params['id'], $recordCount);
            if ($detail->type == EnterprisePayment::RANDOM_AMOUNT_TYPE) {
                $this->createRandomAmount($detail);

                (new EnterprisePaymentService())->setRandomList($params['id'], $lkey);
            }
        }

        // 如果是已删除
        if ($params['status'] == self::DELETED_YES) {
            $data = [
                'is_deleted' => EnterprisePayment::IS_DELETED_YES
            ];
        }

        // 如果是已结束
        if ($params['status'] == EnterprisePayment::IS_SALE_OVER) {
            (new EnterprisePaymentService())->overActivity($detail);
            delRedis([$lkey, $recordCount]);
        }

        // 更新数据
        EnterprisePayment::query()->where('enterprise_id', $params['enterprise_id'])
            ->where('id', $params['id'])
            ->update($data);

        return success(['data' => []]);
    }

    /**
     * @param $params
     * @param $detail
     */
    protected function operationStatus($params, $detail)
    {

        // 原始活动已删除
        if ($detail->is_deleted == EnterprisePayment::IS_DELETED_YES) {
            throw new Error(110014, '活动已删除不能进行操作');
        }

        // 原始数据未上架
        if ($detail->is_sale == EnterprisePayment::IS_SALE_DEFAULT) {
            if ($params['status'] == EnterprisePayment::IS_SALE_DOWN) {
                throw new Error(110010, '活动未上架不能下架');
            }
        }

        // 原始数据已上架
        if ($detail->is_sale == EnterprisePayment::IS_SALE_UP) {
            if ($params['status'] == EnterprisePayment::IS_SALE_UP) {
                throw new Error(110011, '活动不能重复上架');
            }

            // 已上架不能删除
            if ($params['status'] == self::DELETED_YES) {
                throw new Error(110011, '活动已上架不能删除');
            }
        }

        // 原始数据已下架
        if ($detail->is_sale == EnterprisePayment::IS_SALE_DOWN) {
            if ($params['status'] == EnterprisePayment::IS_SALE_DOWN) {
                throw new Error(110012, '活动不能重复下架');
            }
        }

        // 原始数据已结束
        if ($detail->is_sale == EnterprisePayment::IS_SALE_OVER) {
            // 如果不是删除操作则提示错误
            if ($params['status'] !== self::DELETED_YES) {
                throw new Error(110013, '活动已结束,不能进行操作');
            }
        }

        // 如果是已删除
        if ($params['status'] == self::DELETED_YES) {
            // 活动未结束
            if ($detail->is_sale !== EnterprisePayment::IS_SALE_OVER) {
                throw new Error(110014, '活动未结束,不能进删除');
            }
        }
    }

    /*
     * 生成随机红包
     */
    protected function createRandomAmount($detail)
    {
        $total = EnterprisePaymentRandom::query()
            ->where('enterprise_payment_id', $detail->id)
            ->count();
        if (!empty($total)) {
            return true;
        }

        $bonus_total = $detail->enterprise_payment_amount;
        $bonus_count = $detail->quantity;
        $bonus_max   = $detail->max_amount;//此算法要求设置的最大值要大于平均值
        $bonus_min   = $detail->min_amount;

        $result_bonus = (new RandomRedPacket())->getBonus($bonus_total, $bonus_count, $bonus_max, $bonus_min);

        $list = array_chunk($result_bonus, 5000);
        // 添加
        DB::beginTransaction();
        try {
            foreach ($list as $value) {
                $data = [];
                foreach ($value as $item) {
                    $data[] = [
                        'amount'                => $item,
                        'enterprise_id'         => $detail->enterprise_id,
                        'enterprise_payment_id' => $detail->id,
                        'is_receive'            => EnterprisePaymentRandom::IS_RECEIVE_NO,
                        'created_at'            => time()
                    ];
                }

                if (!empty($data)) {
                    $sql = $this->createdAllInSql($data);
                    DB::insert($sql);
                }
            }

            DB::commit();
        } catch (\Exception $e) {
            DB::rollBack();
            notice('100000', '更改活动状态失败');
        } catch (\Throwable $e) {
            DB::rollBack();
            notice('100000', '更改活动状态失败');
        }

        return true;
    }


    /**
     * @title createdAllInSql
     * @description 将所有的数据拼接成sql语句
     * @author postbird
     *
     * @param $data
     *
     * @return bool|string
     */
    protected function createdAllInSql($data)
    {
        $sql = "INSERT INTO `enterprise_payment_random` (`amount`, `enterprise_id`, `enterprise_payment_id`, `is_receive`, `created_at`) VALUES";
        foreach ($data as $item) {

            $sql .= "('{$item['amount']}', '{$item['enterprise_id']}', '{$item['enterprise_payment_id']}', '{$item['is_receive']}','{$item['created_at']}'),";
        }
        $sql = substr($sql, 0, strlen($sql) - 1);
        return $sql;
    }
}
