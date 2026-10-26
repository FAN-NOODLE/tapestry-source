<?php

namespace App\Http\Controllers\App\V1\AutoRelease;

use App\Http\Controllers\App\V1\BaseController;
use App\Models\EnterpriseGrayscaleReleaseLog;
use App\Models\EnterpriseThirdPartAccount;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;

class ListController extends BaseController
{
    /**
     * 自动发版列表
     *
     * @param Request $request
     * @return array
     * @throws ValidationException
     */
    public function index(Request $request)
    {
        $this->validateParam($request);

        $page = $request->post('page', 1);
        $size = $request->post('limit', 10);
        $branchName = $request->post('release_branch');

        $query = EnterpriseGrayscaleReleaseLog::query()->select([
            'id',
            'enterprise_id',
            'release_branch',
            'status',
            'created_at',
            'cancel_at',
            'can_cancel',
            'remark'
        ]);

        if ($branchName) {
            $query->where('release_branch', 'like', "%{$branchName}%");
        }

        $logs = $query->orderBy('created_at', 'DESC')->paginate($size, ['*'], '', $page);

        $enterpriseIDs = array_column($logs->items(), 'enterprise_id');

        if (empty($enterpriseIDs)) {
            return [
                'data' => [],
                'meta' => [
                    'size'  => $size,
                    'page'  => $page,
                    'total' => 0,
                ]
            ];
        }

        $accounts = EnterpriseThirdPartAccount::query()
            ->whereIn('enterprise_id', $enterpriseIDs)
            ->pluck('name', 'enterprise_id')
            ->toArray();

        $data = [];
        foreach ($logs->items() as $item) {
            $data[] = [
                'id'             => $item['id'],
                'name'           => $accounts[$item['enterprise_id']] ?? '',
                'enterprise_id'  => $item['enterprise_id'],
                'release_branch' => $item['release_branch'],
                'status'         => $item['status'],
                'created_at'     => $item['created_at'],
                'cancel_at'      => $item['cancel_at'],
                'can_cancel'     => $item['can_cancel'],
                'remark'         => $item['remark'],
            ];
        }

        return [
            'data' => $data,
            'meta' => [
                'size'  => $size,
                'page'  => $page,
                'total' => $logs->total(),
            ]
        ];
    }

    /**
     * 参数验证
     *
     * @param Request $request
     * @throws ValidationException
     */
    protected function validateParam(Request $request)
    {
        $rule = [
            'page'  => 'integer',
            'limit' => 'integer',
        ];

        $errMsg = [
            'page.integer'  => '页码必须为正整数',
            'limit.integer' => '每页数量必须为正整数',
        ];

        $this->validate($request, $rule, $errMsg);
    }
}