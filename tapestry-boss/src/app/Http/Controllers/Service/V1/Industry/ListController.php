<?php

namespace App\Http\Controllers\Service\V1\Industry;

use App\Http\Controllers\Service\V1\BaseController;
use App\Models\Industry;

class ListController extends BaseController
{
    /**
     * 行业列表
     *
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function index()
    {
        $this->validate($this->request, [
            'page' => 'required|integer',
            'size' => 'required|integer|max:100',
        ]);
        $input = $this->request->all();

        $industry = new Industry();

        list($total, $list) = $industry->lists($input['page'], $input['size']);

        return success([
            'list'  => $list,
            'total' => $total,
            'page'  => $input['page'],
            'size'  => $input['size']
        ]);
    }
}
