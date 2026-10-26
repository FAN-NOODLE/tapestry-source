<?php


namespace App\Library;


use Exception;

trait Validator
{
    /**
     * @throws Exception
     */
    public function required(array $rules, array $params)
    {
        foreach ($rules as $field => $rule) {
            if (true === $rule[0]) {
                if (!array_key_exists($field, $params)) {
                    $message = "该字段为必填";
                    if (!array_key_exists(1, $rule)) {
                        $message = $rule[1];
                    }
                    throw new Exception("{$field}::{$message}");
                }
            }
        }
    }
}