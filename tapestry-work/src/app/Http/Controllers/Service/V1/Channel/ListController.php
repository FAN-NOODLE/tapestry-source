<?php

namespace App\Http\Controllers\Service\V1\Channel;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Channel;

class ListController extends BaseController
{
    /**
     * 渠道列表
     * {domain}/service/v1/channel/list
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'page'          => 'integer',
            'size'          => 'integer',
            'enterprise_id' => 'required|integer',
        ], [
            'page.integer'           => '[参数]page需为integer',
            'size.integer'           => '[参数]size需为integer',
            'enterprise_id.required' => '[参数]enterprise_id必填',
            'enterprise_id.integer'  => '[参数]enterprise_id需为integer',
        ]);

        $input = $this->request->only(['channel', 'page', 'size', 'enterprise_id']);

        $page    = intval($input['page'] ?? Constant::PAGE);
        $size    = intval($input['size'] ?? Constant::SIZE);
        $channel = trim($input['channel'] ?? '');

        $total = Channel::query()
            ->where('enterprise_id', $input['enterprise_id'])
            ->when($channel != '', function ($query) use ($channel) {
                return $query->where(function ($query) use ($channel) {
                    $query->where('channel', 'LIKE', '%' . $channel . '%')
                        ->orWhere('id', 'LIKE', '%' . $channel . '%');
                });
            })->count();

        $data = [];
        if ($total) {
            $data = Channel::query()
                ->where('enterprise_id', $input['enterprise_id'])
                ->when($channel != '', function ($query) use ($channel) {
                    return $query->where(function ($query) use ($channel) {
                        $query->where('channel', 'LIKE', '%' . $channel . '%')
                            ->orWhere('id', 'LIKE', '%' . $channel . '%');
                    });
                })->orderBy('id', 'DESC')
                ->when($size > 0, function ($query) use ($page, $size) {
                    return $query->forpage($page, $size);
                })
                //不分页为下拉列表（不包含自动创建的渠道）
                ->when($size == 0, function ($query) {
                    return $query->where('source_type', Channel::SOURCE_TYPE_DEFAULT);
                })
                //不分页为下拉列表（包含所有渠道）
                ->when($size < 0, function ($query) {
                    return $query;
                })
                ->get(['id', 'channel', 'source_type', 'admin_name', 'created_at'])
                ->toArray();
        }

        return success([
            'data' => $data,
            'meta' => [
                'page'  => $page,
                'size'  => $size,
                'total' => $total
            ]
        ]);
    }
}
