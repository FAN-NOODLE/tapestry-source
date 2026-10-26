<?php

namespace App\Http\Controllers\V1\Store;

use App\Http\Controllers\V1\BaseController;
use App\Traits\ExcelTrait;
use CrsSdk\Sdk;

class DownloadStoreListController extends BaseController
{
    use ExcelTrait;

    /**
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function index()
    {
        $params = [
            'enterprise_id' => $this->enterpriseId,
            'size'  =>  9999
        ];
        $result = Sdk::work()
            ->setApi('service/v1/store/department-list')
            ->post($params);
        if (empty($result)) {
            return success([]);
        }
        $dataArr = [];
        foreach ($result['list'] as $val) {
            $dataArr[] = [
                $val['id'],
                $val['wechat_department_id'],
                $val['third_store_code'],
                $val['store_name'],
                $val['store_type'] ?: '',
                $val['address'],
                $val['inventory_name'],
                $val['city'],
                $val['phone'],
                $val['sa_share_order_url'],
                $val['sa_check_order_url'],
                $val['popup_store_keys'],
                $val['email'],
            ];
        }
        $header = ['系统id','企微部门id','门店编号','门店对外名称','门店属性:1奥莱 2零售',
            '门店地址','库存查询门店名','门店所属城市','门店电话','SA分享订单链接','店长查看门店链接','签名密钥','邮箱'];

        $this->excelExport($header, $dataArr, 'store-department-template');

    }
}
