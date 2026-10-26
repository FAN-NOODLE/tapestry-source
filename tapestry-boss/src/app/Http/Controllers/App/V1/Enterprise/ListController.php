<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/21
 * Time: 8:16 PM
 */

namespace App\Http\Controllers\App\V1\Enterprise;

use App\Models\Enterprise;
use Exception;
use CrsSdk\Work\EnterprisePayment;
use Illuminate\Http\Request;
use App\Http\Controllers\App\V1\BaseController;
use App\Services\EnterpriseService;

class ListController extends BaseController
{
    /**
     * 企业service
     *
     * @var EnterpriseService
     */
    protected $enterpriseService;

    /**
     * ListController constructor.
     * @param EnterpriseService $enterpriseService
     * @throws Exception
     */
    public function __construct(EnterpriseService $enterpriseService)
    {
        parent::__construct();

        // 企业service赋值
        $this->enterpriseService = $enterpriseService;
    }

    /**
     * 企业列表
     *
     * @author 梁欢 lianghuan@vchangyi.com
     * @author 赵奇
     * @date   2018/11/25 11:16
     *
     * @param Request $request
     *
     * @return array
     */
    public function index(Request $request)
    {
        // 查询参数处理
        $params = $request->input();
        $limit = isset($params['limit']) ? $params['limit'] : 10;

        // 列表查询条件
        $where = [
            'version' => Enterprise::NORMAL_VERSION
        ];
        isset($params['name']) && $where['name'] = $params['name'];

        // 列表数据查询
        $query = Enterprise::query();

        // 条件组合
        foreach ($where as $key => $value) {
            if (!empty($value) || $value == 0) {
                switch ($key) {
                    case 'name':
                        $query = $query->where($key, 'like', '%' . $value . '%');
                        break;
                    default:
                        $query = $query->where($key, '=', $value);
                }
            }
        }

        // 排序
        $query = $query->orderBy('id', 'desc');
        $enterprise = $query->paginate((int)$limit)->toArray();

        $enterpriseIds = array_column($enterprise['data'], 'id');

        if ($enterpriseIds) {
            $list = (new EnterprisePayment())->budgetByEnterprise(['enterprise_ids' => $enterpriseIds]);
            $list = array_column($list['data'] ?? [], null, 'enterprise_id');
            foreach ($enterprise['data'] as &$item) {
                $item['balance'] = 0;
                if (isset($list[$item['id']])) {
                    $item['balance'] = $list[$item['id']]['balance'] ?? 0;
                }
            }
            unset($item);
        }
        return success([
            'data' => $enterprise['data'],
            'meta' => [
                'total' => $enterprise['total'],
                'size' => (int)$enterprise['per_page'],
                'page' => $enterprise['current_page']
            ]
        ]);
    }
}
