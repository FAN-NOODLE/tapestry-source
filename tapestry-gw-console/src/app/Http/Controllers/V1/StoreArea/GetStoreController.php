<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Common\Constant;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class GetStoreController extends BaseController
{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {
        return success((new StoreArea())->getStore($this->getParam()));
    }

    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'id'          => 'integer',
            'search_type' => 'integer',
            'search_text' => 'string',
            'filter_ids'  => 'array',
            'page'        => 'integer',
            'size'        => 'integer'
        ], [
            'id.*' => '区域ID不正确'
        ]);
        $page     = $this->request->input('page', Constant::PAGE);
        $pageSize = $isSize = $this->request->input('size', Constant::LIMIT);
        if ($pageSize > 50) {
            $pageSize = Constant::LIMIT;
        }
        //兼容门店弹窗
        $isNotPage = $this->request->input('is_not_page', 0);
        if(!empty($isNotPage)){
            $pageSize = $isSize;
        }
        return [
            'page'        => $page,
            'size'        => $pageSize,
            'is_not_page' => $isNotPage,
            'id'          => $this->request->input('id',0),
            'search_type' => $this->request->input('search_type',0),
            'search_text' => $this->request->input('search_text',''),
            'filter_ids'  => $this->request->input('filter_ids',[]),
            'store_ids'   => $this->request->input('store_ids',[]),
            'enterprise_id' => $this->enterpriseId ?? 0,
        ];
    }
}
