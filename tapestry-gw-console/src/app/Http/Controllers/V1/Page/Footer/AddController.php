<?php

namespace App\Http\Controllers\V1\Page\Footer;

use CrsSdk\Shop\Footer;
use Illuminate\Http\JsonResponse;
use App\Http\Controllers\V1\BaseController;

class AddController extends BaseController
{
    /**
     * @return JsonResponse
     */
    public function index()
    {
        $picture = $this->request->input('picture', '');

        (new Footer())->add([
            'picture'       => $picture,
            'enterprise_id' => $this->enterpriseId,
        ]);

        $this->writeOperatorLog('添加页底');

        return success();
    }
}
