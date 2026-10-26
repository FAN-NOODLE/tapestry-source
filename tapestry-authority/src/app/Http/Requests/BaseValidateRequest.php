<?php


namespace App\Http\Requests;

use Illuminate\Support\Facades\Validator;
use Illuminate\Validation\ValidationException;

class BaseValidateRequest
{
    /**
     * @param array $data
     * @param array $rule
     * @param array $message
     * @throws ValidationException
     */
    protected static function validator($data = [], $rule = [], $message = [])
    {
        $validator = Validator::make($data, $rule, $message);
        if ($validator->fails()) {
            throw new ValidationException($validator);
        }
    }


}