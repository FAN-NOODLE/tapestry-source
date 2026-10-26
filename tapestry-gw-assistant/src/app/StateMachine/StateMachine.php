<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/12
 * Time: 16:28
 */

namespace App\StateMachine;

use App\Exceptions\Error;
use Illuminate\Support\Facades\Log;

abstract class StateMachine
{
    protected $states = [];
    protected $start = [];
    protected $end = [];
    protected $events = [];
    protected $transitions = [];
    protected $eventDesc = '';

    protected $obj;

    protected $currentState = 0;

    protected function createMachine($allowState, $start, $end, $events, $transitions)
    {
        $this->setStates($allowState);
        $this->setStart($start);
        $this->setEnd($end);
        $this->setEvents($events);

        foreach ($transitions as $transition) {
            $this->addTransition($transition);
        }
    }

    private function addTransition(Transition $transition)
    {
        $this->transitions[$transition->getFrom()][$transition->getEvent()] = [
            'to'     => $transition->getTo(),
            'action' => $transition->getAction()
        ];
    }

    private function setStates($states)
    {
        $this->states = $states;
    }

    private function setStart($start)
    {
        $this->start = $start;
    }

    private function setEnd($end)
    {
        $this->end = $end;
    }

    private function setEvents($events)
    {
        $this->events = $events;
    }

    public function getEventDesc($event)
    {
        return $this->events[$event]['desc'] ?? '';
    }

    protected abstract function addLock();

    protected abstract function delLock();

    protected abstract function onEnter($event, $data);

    protected abstract function onExit();

    protected abstract function getStateObject($data, $operator);

    /**
     * 执行方法
     *
     * @param        $data
     * @param string $operator
     *
     * @return array
     */
    public function run($data, $operator='user')
    {
        $event = $data['event'];

        // TODO 分布式锁
        $stateObj = $this->getStateObject($data, $operator);
        try {

            $from     = $stateObj->status;

            //先检查是否为结束状态
            if (in_array($from, $this->end)) {
                throw new Error(codeNum('state_is_end'), codeMsg('state_is_end', [$this->states[$from]['desc']]), 210);
            }
            //检查状态是否被允许

            if (!array_key_exists($from, $this->transitions)) {
                throw new Error(codeNum('state_not_exists'), codeMsg('state_not_exists'), 210);
            }

            $stateAction = $this->transitions[$from];

            if (!array_key_exists($event, $stateAction)) {
                throw new Error(codeNum('event_not_exists'), codeMsg('event_not_exists', [$event]), 210);
            }
            $this->onEnter($event, $data);

            $stateEvent = $stateAction[$event];

            $eventDesc = $this->getEventDesc($event);
            $startActionTime = microtime();
            $res_msg = $stateEvent['action']->handle($from, $event, $stateEvent['to'], $data, $stateObj, $eventDesc);
            $endActionTime = microtime();
            Log::info('记录对象id:'.$stateObj->id.',记录对象状态:'.$from.',事件:'.$event.',处理完成耗时:'.((int)$endActionTime-(int)$startActionTime).'ms');
            // TODO 分布式锁
            $this->delLock();

            return $res_msg;
        } catch (\Exception $exception) {
            // TODO 分布式锁
            $this->delLock();
            Log::info("状态扭转异常",[
                'msg'  => $exception->getMessage(),
                'code' => $exception->getCode(),
                'file' => $exception->getFile(),
                'line' => $exception->getLine(),
            ]);
            throw new Error($exception->getCode(), $exception->getMessage(), 210);
        }
    }
}
