<?php


namespace App\Http\Controllers\V1\Goods;


use App\Http\Controllers\V1\BaseController;
use App\Services\GoodsBatchChangeService;
use Illuminate\Http\JsonResponse;
use Illuminate\Validation\ValidationException;

class GoodsBatchChangeController extends BaseController
{
    /**
     * @return JsonResponse
     * @throws ValidationException
     */
    public function index(): JsonResponse
    {
        $param = $this->getParams();
        return success(['success' => (int) (new GoodsBatchChangeService($this->enterpriseId, $param['change_type']))->run($param['change_file'])]);
    }

    /**
     * @return array
     * @throws ValidationException
     */
    protected function getParams(): array
    {
        $inStr = GoodsBatchChangeService::TYPE_COLLECT_STR;
        return $this->validate(request(), [
            'change_type'   => 'required|integer|in:' . $inStr,
            'change_file'   => 'required|file',
        ], [
            'change_type.*' => '更改类型不合法',
            'change_file.*' => '更改文件不合法',
        ]);
    }
}