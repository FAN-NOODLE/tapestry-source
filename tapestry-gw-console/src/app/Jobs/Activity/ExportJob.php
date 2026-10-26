<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/12 0012
 * Time: 19:44
 */

namespace App\Jobs\Activity;

use App\Common\Constant;
use App\Jobs\Export\ExportJob as BaseJob;
use App\Traits\ActivityTrait;
use CrsSdk\MicroActivity\Activity;
use Illuminate\Support\Facades\Log;
use App\Exceptions\Error;

class ExportJob extends BaseJob
{
    use ActivityTrait;
    /**
     * excel 标题头
     *
     * @var array
     */
    private $header;

    /**
     * 文件名前缀
     *
     * @var string
     */
    private $prefix;

    public function __construct(array $params, array $paramJson = [], string $name = '活动列表导出', $header = [], $prefix = '')
    {

        parent::__construct($params, $paramJson, $name);
        $this->header = $header;
        $this->prefix = $prefix;
    }

    public function handle()
    {
        $params           = $this->params;
        $params['method'] = 'list';

        try {
            $this->excelExportInit($this->header)->listExport($params);
        } catch (\Exception $exception) {
            Log::info('导出失败' . $exception->getMessage(), $params);
            // 修改状态为失败
            $this->jobStatus(false, $exception->getMessage().':'.getDownloadPath());
            return true;
        }
        // 导出完成 修改状态为成功
        $this->jobStatus();
        return true;

    }

    /**
     * @param $params
     *
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    function listExport($params)
    {
        $search = [
            'page'          => Constant::PAGE,
            'size'          => Constant::MAX_DEFAULT_SIZE,
            "enterprise_id" => $params['enterprise_id'],
            "method"        => "list",
            'type'          => $params['activity_type'],
            'search'        => $params['search']
        ];

        // 查询获取物料库存管理列表
        try {
            $data = (new Activity())->getActivity($search);

        } catch (\Exception $exception) {
            Log::info('获取活动列表失败' . $exception->getMessage(), $params);
            throw new Error(10000, '获取活动列表失败'.$exception->getMessage());
        }

        // 组装数据
        switch ($params['activity_type']) {
            case 'group':
                $data = $this->groupFormData($data['list']);
                break;
            case 'flash-sale':
                $data = $this->flashSaleFormData($data['list']);
                break;
            case 'haggle':
                $data = $this->haggleFormData($data['list']);
                break;
            case 'help':
                $data = $this->helpFormData($data['list']);
                break;
        }

        // 插入excel
        $this->exportExcel($data);

        $this->exportClose();
    }

    /**
     * 组装限时购导出数据
     *
     * @param $data
     *
     * @return array
     */
    private function flashSaleFormData($data)
    {
        $rows = [];
        foreach ($data as $key => $rule) {

            $rows[$key]['id']     = $this->prefix . $rule['id'];
            $rows[$key]['name']   = $rule['name'];
            $rows[$key]['status'] = $this->activityStatus($rule['activity_status']);
            $rows[$key]['time']              = date('Y-m-d H:i:s', $rule['start_at']) . ' 至 ' . date('Y-m-d H:i:s', $rule['end_at']);
            $rows[$key]['order_num']         = $rule['activity_detail']['order_num'];
            $rows[$key]['customer_num']      = $rule['activity_detail']['customer_num'];
            $rows[$key]['order_total_money'] = $rule['activity_detail']['order_total_money'] > 0 ? round($rule['activity_detail']['order_total_money'] / 100, 2) : 0;
        }

        return $rows;
    }


    private function groupFormData($data)
    {

        $rows = [];
        foreach ($data as $key => $rule) {

            $rows[$key]['id']     = $this->prefix . $rule['id'];
            $rows[$key]['name']   = $rule['name'];
            $rows[$key]['status'] = $this->activityStatus($rule['activity_status']);
            $rows[$key]['time'] = date('Y-m-d H:i:s', $rule['start_at']) . ' 至 ' . date('Y-m-d H:i:s', $rule['end_at']);
            $this->ActivityArrayToJson($rule, 'activity_sku', 'spu_name', true);
            $rows[$key]['spu_name'] = $rule['spu_name'];
            $rows[$key]['need_num'] = $rule['need_num'];
            unset($rule['activity_detail']['activity_rule_id']);
            $rows[$key] = array_merge($rows[$key], $rule['activity_detail']); //todo 如何保证数据与标题一致
            $this->ActivityArrayToJson($rule, 'spu_list', 'path');
            $rows[$key]['path'] = $rule['path'];
            $this->ActivityArrayToJson($rule, 'spu_list', 'qr_code', true);
            $rows[$key]['qr_code'] = $rule['qr_code'];
        }

        return $rows;
    }


    private function haggleFormData($data)
    {
        $rows = [];
        foreach ($data as $key => $rule) {

            $rows[$key]['id']     = $this->prefix . $rule['id'];
            $rows[$key]['name']   = $rule['name'];
            $rows[$key]['status'] = $this->activityStatus($rule['activity_status']);
            $rows[$key]['time'] = date('Y-m-d H:i:s', $rule['start_at']) . ' 至 ' . date('Y-m-d H:i:s', $rule['end_at']);
            $this->ActivityArrayToJson($rule, 'activity_sku', 'spu_name', true);
            $rows[$key]['init_num']  = $rule['activity_detail']['init_num'];
            $rows[$key]['help_num']  = $rule['activity_detail']['help_num'];
            $rows[$key]['order_num'] = $rule['activity_detail']['order_num'];
            unset($rule['activity_detail']['activity_rule_id']);
            $rows[$key] = array_merge($rows[$key], $rule['activity_detail']); //todo 如何保证数据与标题一致

        }

        return $rows;
    }

    /**
     * 好友助力组装数据
     *
     * @param $data
     *
     * @return array
     * @date    2020/6/16 18:24
     */
    private function helpFormData($data)
    {
        $rows = [];
        foreach ($data as $key => $rule) {
            $rows[$key]['id'] = $this->prefix . $rule['id'];
            $rows[$key]['name'] = $rule['name'];
            $rows[$key]['status'] = $this->activityStatus($rule['activity_status']);
            $rows[$key]['time'] = date('Y-m-d H:i:s', $rule['start_at']) . ' 至 ' . date('Y-m-d H:i:s', $rule['end_at']);
            $this->ActivityArrayToJson($rule, 'activity_sku', 'spu_name');
            $this->ActivityArrayToJson($rule, 'activity_sku', 'sku_name', true);
            $rows[$key]['spu_name'] = !empty($rule['spu_name']) ? $rule['spu_name'] : $rule['sku_name'] ;
            $rows[$key]['init_num'] = $rule['activity_detail']['init_num'];
            $rows[$key]['help_num'] = $rule['activity_detail']['help_num'];
            $rows[$key]['success_number'] = $rule['activity_detail']['success_number'];
            $rows[$key]['order_num'] = $rule['activity_detail']['order_num'];
            $rows[$key]['activity_stock'] = $rule['activity_detail']['activity_stock'];
            unset($rule['activity_detail']['activity_rule_id']);
            $rows[$key] = array_merge($rows[$key], $rule['activity_detail']);
        }

        return $rows;
    }

    protected function activityStatus($status)
    {
        $activityStatus = [
            '1' => "未上架",
            '2' => "已上架",
            '3' => "进行中",
            '4' => "已下架",
            '5' => "已结束",
        ];
        return $activityStatus[$status] ?? '未知';
    }

}
