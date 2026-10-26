<?php

namespace App\Micros\Work\CustomerInherit;

use CrsSdk\Sdk;

class CustomerInherit
{

    /**
     * @return void
     * 任务列表
     */
    public function InheritList(array $data): array
    {
        return Sdk::work()->setApi("service/v1/customer-inherit/list")->post($data);
    }

    /**
     * @return array|mixed
     * @throws \Exception
     * 任务详情
     */
    public function inheritInfo(array $data): array
    {
        return Sdk::work()->setApi("service/v1/customer-inherit/info")->post($data);
    }

    /**
     * @param array $data
     *
     * @return array|mixed
     * @throws \Exception
     * 顾客转移
     */
    public function inheritTransfer(array $data): array
    {
        return Sdk::work()->setApi("service/v1/customer-inherit/transfer")->post($data);
    }

    /**
     * @param array $data
     *
     * @return void
     * 导购列表
     */
    public function memberList(array $data): array
    {
        return Sdk::work()->setApi("service/v1/customer-inherit/member-list")->post($data);
    }

    /**
     * @param $data
     *
     * @return array|mixed
     * @throws \Exception
     * 外部联系人转接次数
     */
    public function customerInheritRecord($data): array
    {
        return Sdk::work()->setApi("service/v1/customer-inherit/customer-inherit-record")->post($data);
    }

    /**
     * @param $data
     *
     * @return array|mixed
     * @throws \Exception
     * 外部联系人和多个导购的关系
     */
    public function memberRelation($data): array
    {
        return Sdk::work()->setApi("service/v1/customer-inherit/member-record")->post($data);
    }

    /**
     *
     * @param $data
     *
     * @return array
     * @throws \Exception
     * @Time: 2022/11/17
     */
    public function recordForExternalIds($data): array
    {
        return Sdk::work()->setApi("service/v1/customer-inherit/record-external_ids")->post($data);
    }
}