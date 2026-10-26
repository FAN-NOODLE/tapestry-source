<?php

namespace App\Http\Controllers\Service\V1\Media;

use App\Common\Constant;
use App\Http\Controllers\Service\V1\BaseController;
use App\Models\QYwxMedia;

class ListController extends BaseController
{

    /**
     * 获取企业素材信息
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $input = $this->getParams();

        // 验证
        $page = $this->request->input('page');
        $size = $this->request->input('size');

        $searchText      = isset($input['search_text']) && trim($input['search_text']) ? trim($input['search_text']) : "";
        $type            = $input['type'] ?? 0;
        $page            = $page ?? Constant::PAGE;
        $size            = $size ?? Constant::LIMIT;
        $input['uc_eid'] = $input['uc_eid'] ?? '';

        $sql = QYwxMedia::query()
            ->where([
                'type' => $type
            ])
            ->when('' != $input['uc_eid'], function ($query) use ($input) {
                return $query->where(['uc_eid' => $input['uc_eid']]);
            })
            ->where('status', '>=', 0)
            ->when('' != $searchText, function ($query) use ($searchText) {
                return $query->where('name', 'like', '%' . $searchText . '%');
            });

        $list  = $sql->forpage($page, $size)->get()->toArray();
        $total = $sql->count();

        return success([
            'data'  => $list ?? [],
            'total' => $total ?? 0,
            'page'  => $page,
            'size'  => $size,
        ]);
    }

    /**
     * 参数验证
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    private function getParams()
    {
        $input = $this->request->only([
            'search_text',
            'uc_eid',
            'type',
            'page',
            'size',
        ]);

        $rules = [
            'search_text' => 'string',
            'uc_eid'      => 'string',
            'type'        => 'string',
            'page'        => 'integer',
            'size'        => 'integer',
        ];

        $message = [
            'search_text.*' => 'search_text有误',
            'uc_eid.*'      => 'uc_eid必填',
            'type.*'        => 'type有误',
            'page.*'        => 'page有误',
            'size.*'        => 'size有误',
        ];

        $this->validate($this->request, $rules, $message);

        return $input;
    }

}
