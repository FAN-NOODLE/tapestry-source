<?php
/**
 * Created by PhpStorm.
 * User: stupidnb
 * Date: 2018/12/24
 * Time: 11:30
 */

namespace App\Events;


class AddTaskEvent extends Event
{
    /**
     * Create a new event instance.
     *
     * @return void
     */
    public $data;

    public function __construct($data)
    {
        $this->data = $data;
    }

    /**
     * 获取传入数组
     *
     * @return mixed
     */
    public function getData()
    {
        return $this->data;
    }
}
