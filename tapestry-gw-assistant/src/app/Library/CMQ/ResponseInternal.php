<?php
/**
 * Created by PhpStorm.
 * Function: xxx
 * User: 龚国新
 * Date: 2019/10/22
 * Time: 21:12
 */

namespace App\Library\CMQ;


class ResponseInternal
{
    public $header;
    public $status;
    public $data;

    public function __construct($status = 0, $header = NULL, $data = "") {
        if ($header == NULL) {
            $header = array();
        }
        $this->status = $status;
        $this->header = $header;
        $this->data = $data;
    }

    public function __toString()
    {
        $info = array("status" => $this->status,
                      "header" => json_encode($this->header),
                      "data" => $this->data);
        return json_encode($info);
    }
}