<?php

namespace App\Http\Controllers\App\V1\Industry;

use App\Models\Industry;
use App\Http\Controllers\App\V1\MiniProgram\MiniBaseController;

class ListController extends MiniBaseController
{
    /**
     * @return array
     */
    public function index()
    {
        $industry = new Industry();

        $list = $industry->allData();

        return success([
            'data' => (object)$list
        ]);
    }
}