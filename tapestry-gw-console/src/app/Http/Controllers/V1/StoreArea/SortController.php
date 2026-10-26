<?php


namespace App\Http\Controllers\V1\StoreArea;


use App\Exceptions\Error;
use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\StoreArea;

class SortController extends BaseController
{
    /**
     * @return array
     * @throws \Illuminate\Validation\ValidationException
     * @throws \Exception
     */
    public function index()
    {

        (new StoreArea())->areaSort($this->getParam());
        return success();
    }

    /**
     * @return mixed
     * @throws \Illuminate\Validation\ValidationException
     */
    public function getParam()
    {
        $this->validate($this->request, [
            'sort' => 'required|json',
        ], [
            'sort.*'=> '排序参数不合法',
        ]);
        $sortArray = json_decode($this->request->sort, true);
        foreach ($sortArray as $item) {
            if (!is_array($item) || !isset($item['id']) || !isset($item['sort'])) {
                throw new Error(1620001, '排序参数不合法');
            }
        }
        return [
        'sort'              => $this->request->sort,
            'enterprise_id' => $this->enterpriseId,
        ];
    }
}
