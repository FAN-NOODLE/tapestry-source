<?php
/**
 * Created by PhpStorm.
 * User: huangxm
 * Date: 2020-08-31
 * Time: 14:18
 */

namespace App\Jobs\Export;

use App\Common\Constant;
use CrsSdk\Shop\Order;
use CrsSdk\Work\Customer;
use Illuminate\Support\Facades\Log;

class CustomerExportJob extends ExportJob
{
    /**
     * 导出表header 头
     *
     * @var string[]
     */
    protected $header
        = [
            'name'              => '客户',
            'gender'            => '性别',
            'phone'             => '手机号',
            'tag'               => '企业标签',
            'group'             => '用户分组',
            'order_total'       => '订单数',
            'order_price_total' => '订单金额',
            'mem_name'          => '所属员工',
            'dep_name'          => '所属部门',
            'channel_name'      => '渠道',
            'created_at'        => '添加时间',
            'mem_mobile'        => '所属员工手机号',
            'is_wechat'         => '是否授权小程序'
        ];

    public function __construct(array $params, array $paramsJson, string $name = '客户列表')
    {
        parent::__construct($params, $paramsJson, $name);
    }

    /**
     * @return bool
     */
    public function handle()
    {
        $params = $this->params;
        $header = array_values($this->header);
        // 查询顾客列表
        try {
            $this->excelExportInit($header)->customers($params);
        } catch (\Exception $exception) {
            Log::info('导出失败' . $exception->getMessage(), $params);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage());
            return true;
        }

        // 导出完成 修改状态为成功
        $this->jobStatus();
        return true;
    }

    /**
     * @param array $params
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    private function customers(array $params)
    {
        $customers = $this->yieldCustomer($params);

        foreach ($customers as $customer) {
            $orders    = $this->searchOrder($customer['ids']);
            $excelData = [];
            foreach ($customer['list']['data'] as $customerItem) {
                $orderPrice = 0;
                if (isset($orders[$customerItem['union_id']]['order_price_total'])) {
                    $orderPrice = $orders[$customerItem['union_id']]['order_price_total'] / 100;
                }
                $isWechat = "否";
                if (isset($orders[$customerItem['union_id']]) && $orders[$customerItem['union_id']]['user_id'] > 0) {
                    $isWechat = "是";
                }
                $excelData[] = [
                    $customerItem['name'],
                    $customerItem['gender'],
                    $orders[$customerItem['union_id']]['phone'] ?? '',
                    implode(";", $customerItem['tag']),
                    implode(';', $customerItem['group']),
                    $orders[$customerItem['union_id']]['order_total'] ?? 0,
                    $orderPrice,
                    $customerItem['mem_name'],
                    $customerItem['dep_name'],
                    $customerItem['channel_name'],
                    $customerItem['created_at'],
                    $customerItem['mem_mobile'],
                    $isWechat
                ];
            }
            $this->exportExcel($excelData);
        }

        $this->exportClose();
    }

    /**
     * @param array $params
     * @return \Generator
     * @throws \Exception
     */
    private function yieldCustomer(array $params)
    {
        $params['page'] = Constant::PAGE;
        $params['size'] = Constant::BATCH_UPDATE_LIMIT;
        $customerObject = new Customer();
        while (true) {
            $data = $customerObject->customerList($params);
            if (empty($data['list']['data'])) {
                break;
            }
            yield $data;
            $params['page'] += 1;
        }
    }

    /**
     * 订单信息
     *
     * @param array $unionIds
     * @return array
     * @throws \Exception
     */
    private function searchOrder(array $unionIds)
    {
        $unionIds = array_filter($unionIds);
        if (empty($unionIds)) {
            return [];
        }
        // 获取客户结果对应的订单数据
        $order = (new Order)->calcUserOrder([
            'ids'           => $unionIds,
            'enterprise_id' => $this->params['enterprise_id'],
        ]);
        return array_column($order['data'], null, 'union_id');
    }
}
