<?php
/**
 * 商品列表导出
 * User: 何岳龙
 * Date: 2018/11/26
 * Time: 10:04
 */

namespace App\Http\Controllers\V1\AppCenter\WxMall;
use App\Http\Controllers\V1\BaseController;
use App\Jobs\Export\WxMallGoodsListExportJob;

class GoodsExportController extends BaseController
{
    /**
     * 文件名称（列表显示和导出名称）
     * @var
     */
    const FILE_NAME = '好物圈-商品导出';

    /**
     * 下载中心route_url
     */
    const ROUTE_URL = 'gw-console/v1/app-center/wx-mall/goods-export';

    /**
     * 总库列表出tag
     */
    const TAG = 'material-list';
    /**
     * 好物圈商品导出
     * @throws \Box\Spout\Common\Exception\IOException
     * @throws \Box\Spout\Common\Exception\InvalidArgumentException
     * @throws \Box\Spout\Common\Exception\UnsupportedTypeException
     * @throws \Box\Spout\Writer\Exception\InvalidSheetNameException
     * @throws \Box\Spout\Writer\Exception\WriterNotOpenedException
     */
    public function index()
    {
        $input['enterprise_id'] = $this->enterpriseId;
        $input['size']          = 1000;
        $input['page']          = 1;
        $input['tag'] = $input['tag'] ?? self::TAG;
        $input['operate_id']  = $this->enterpriseAdminId;
        $input['operate_name'] = $this->enterpriseAdmin['username'];
        $input['route_url'] = $input['route_url'] ?? self::ROUTE_URL;
        (new WxMallGoodsListExportJob($input, $input, self::FILE_NAME))->handle();
        return success();
    }

}
