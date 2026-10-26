<?php

namespace App\Http\Controllers\V1\Account;

use App\Http\Controllers\V1\BaseController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Crypt;

class DecryptTokenController extends BaseController
{
    public function index(Request $request)
    {

        $token = trim($request->post('token', 0));

        $res = [];
        try {
            $res = json_decode(Crypt::decrypt($token), true);
        }catch (\Exception $e){
            notice(100007);
        }

        if(empty($res)){
            notice(100007);
        }

        return success($res);
    }
}
