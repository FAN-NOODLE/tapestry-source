<?php


namespace App\Http\Controllers\V1\Page\PopContent;


use App\Http\Controllers\V1\BaseController;
use CrsSdk\Shop\PagePopContent;

class PopContentController extends BaseController
{
    public function contentDetail()
    {
        $result =  (new PagePopContent())->detail([
            'enterprise_id' => $this->enterpriseId,
            'type' => $this->request->input('type')
        ]);
        return success($result);
    }

    public function contentSave()
    {
        $result =  (new PagePopContent())->save([
            'enterprise_id' => $this->enterpriseId,
            'content' => $this->request->input('content'),
            'type' => $this->request->input('type')
        ]);
        return success($result);
    }
}