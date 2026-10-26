<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/12
 * Time: 16:28
 */

namespace App\StateMachine;

class Transition {
    private $from;
    private $event;
    private $to;
    private $action;

    public function __construct($from, $event, $to, Action $action = null)
    {
        $this->from = $from;
        $this->event = $event;
        $this->to = $to;
        $this->action = $action;
    }

    /**
     * @return mixed
     */
    public function getFrom()
    {
        return $this->from;
    }

    /**
     * @return mixed
     */
    public function getEvent()
    {
        return $this->event;
    }

    /**
     * @return mixed
     */
    public function getTo()
    {
        return $this->to;
    }

    /**
     * @return mixed
     */
    public function getAction()
    {
        return $this->action;
    }
}
