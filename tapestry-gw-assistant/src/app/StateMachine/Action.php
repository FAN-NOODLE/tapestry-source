<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/12
 * Time: 18:23
 */

namespace App\StateMachine;


interface Action
{
    public function handle($from, $event, $to, $data, $stateObj, $eventDesc);
}