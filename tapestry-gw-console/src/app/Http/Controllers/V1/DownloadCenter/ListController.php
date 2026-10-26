<?php


namespace App\Http\Controllers\V1\DownloadCenter;


use App\Common\Constant;
use App\Micros\DownloadCenter\DownloadCenterMicro;
use App\Http\Controllers\V1\BaseController;

class ListController extends BaseController
{
    public function index()
    {
        $params = $this->request->input();

        // 查询参数处理
        isset($params['search']) && $where['search'] = $params['search'];
        isset($params['search']['status']) && $where['search']['status'] = intval($params['search']['status']);
        $where['page'] = isset($params['page']) ? $params['page'] : Constant::PAGE;
        $where['size'] = isset($params['size']) ? $params['size'] : Constant::LIMIT;
        $where['enterprise_id'] = $this->enterpriseId;

        $list = (new DownloadCenterMicro())->lists($where);

        return success($list);
    }
}
