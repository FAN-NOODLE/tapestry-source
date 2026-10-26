<?php
/**
 * Created by PhpStorm.
 * User: Administrator
 * Date: 2020/5/12 0012
 * Time: 17:26
 */

namespace App\Http\Controllers\V1\ActivityCommon;


use App\Common\Constant;
use App\Jobs\Activity\ExportJob;
use CrsSdk\Shop\DownloadCenter;
use App\Exceptions\Error;

class ExportController extends ActivityBaseController
{
    /**
     * 导出路径
     */
    const ROUTE_URL = 'activity-common/export';

    /**
     * 标签
     */
    const TAG = 'activity-export';

    /**
     * 导出类型
     */
    const TYPE_FILE_EXCEL = 1;
    const TYPE_FILE_ZIP = 2;

    /**
     * 按名称搜索
     */
    const TYPE_SEARCH_NAME = 1;

    /**
     * @throws \Exception
     */
    function index()
    {
        $params = $this->params;
        $type = $this->type;
        $name = Constant::ACTIVITY_NAME[$type]['name'] . date('Ymd');
        $header = Constant::ACTIVITY_NAME[$type]['header'];
        $prefix = Constant::ACTIVITY_NAME[$type]['code'];

        $search['page'] = 1;
        $search['size'] = 1000;
        $search['enterprise_id'] = $this->enterpriseId;
        $search['search']['type'] = self::TYPE_SEARCH_NAME;
        $search['search']['name'] = Constant::ACTIVITY_NAME[$type]['name'];

        $fileName = (new DownloadCenter())->lists($search)['list'][0]['name'] ?? "";

        preg_match('/^[^\d]+\d{8}(\d{4})/', $fileName, $match);

        if ($match) {
            $maxNo = $match[1];
            $nowNo = sprintf("%04d", ($maxNo + 1));
        } else {
            $nowNo = '0001';
        }
        $name .= $nowNo;
        $input = [
            'enterprise_id' => $this->enterpriseId,
            'operate_id'    => $this->enterpriseAdminId,
            'operate_name'  => $this->enterpriseAdmin['username'],
            'type'          => self::TYPE_FILE_EXCEL,//1-excel 2-zip
            'route_url'     => self::ROUTE_URL,
            'tag'           => self::TAG,
            'activity_type' => $type,
            'search'        => $params['search']
        ];
        try {
            $this->dispatch(new ExportJob($input, $params, $name, $header, $prefix));
        } catch (\Exception $e) {
            throw new Error($e->getCode(),$e->getMessage());
        }

    }



}
