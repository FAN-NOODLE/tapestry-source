<?php

namespace App\Micros;

use CrsSdk\Shop\StoreArea;

class StoreMicro
{

    public function getStoreIds($params)
    {
        $data = (new StoreArea())->storeListByArea($params);

        if (empty($data)) {
            return [];
        }
        $storeList = [];
        foreach ($data as $v) {
            $storeList = array_merge($storeList, $v['store_list']);
        }

        return array_values(array_unique(array_column($storeList,'id')));
    }
}