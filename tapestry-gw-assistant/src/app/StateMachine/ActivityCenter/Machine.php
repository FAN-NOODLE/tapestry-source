<?php
/**
 * Created by PhpStorm.
 * User: gongguoxin
 * Date: 2019/6/12
 * Time: 16:28
 */

namespace App\StateMachine\ActivityCenter;

use App\Exceptions\Error;
use Closure;

/**
 * Class StateMachine
 *
 * @package Phwoolcon\Fsm
 * @method string do(string $action, mixed &$payload = null)
 * @method string init(string $action)
 * @method string next()
 */
class Machine
{
    protected static $instances = [];
    protected $initState;
    protected $transitions = [];
    protected $currentState;
    protected $previousState;
    protected $history = [];

    /**
     * StateMachine constructor.
     *
     * @param array      $transitions
     * $transitions example:
     *  [
     *      'State1' => [
     *          'action1' => 'State2',
     *          'action2' => 'State3',
     *      ],
     *      'State2' => [
     *          'action3' => 'State3',
     *      ],
     *      'State3' => [
     *          'action4' => 'State1',
     *      ],
     *  ]
     * You may use Closure to do real actions in the target state:
     *  [
     *      'State1' => [
     *          'action1' => 'State1',
     *          'action2' => 'State1',
     *      ],
     *      'State2' => [
     *          'action3' => function ($stateMachine) {
     *              // Do something here
     *              return 'State3';
     *          },
     *      ],
     *      'State3' => [
     *          'action4' => 'State1',
     *      ],
     *  ]
     * @param array|null $history
     */
    public function __construct(array $transitions, array $history = [])
    {
        $this->transitions = $transitions;
        $this->initAction($history);
    }

    public function __call($name, $arguments)
    {
        return call_user_func_array([$this, $name . 'Action'], $arguments);
    }

    /*
    public function addTransition($fromState, $action, $toState)
    {
        $this->transitions[$fromState][$action] = $toState;
    }
    */
    public function canDoAction($action)
    {
        return isset($this->transitions[$this->currentState][$action]);
    }

    public static function create(array $transitions, array $history = [])
    {
        return new static($transitions, $history);
    }

    /**
     * @param string $action
     * @param mixed  $payload
     *
     * @return string Next state
     * @throws Error
     */
    public function doAction($action, &$payload = null)
    {
        if (!$this->canDoAction($action)) {
            throw new Error(
                33333,
                sprintf("Invalid action '%s' for current state '%s'", $action, $this->currentState),
                '210'
            );
        }
        if (($state = $this->transitions[$this->currentState][$action]) instanceof Closure) {
            $state = $state($this, $payload);
        }
        $this->previousState = $this->currentState;
        $this->currentState  = $state;
        $this->history[]     = ['time' => time(), 'action' => $action, 'state' => $state];
        return $state;
    }

    public function getCurrentState()
    {
        return $this->currentState;
    }

    public function getHistory()
    {
        return $this->history;
    }

    public function initAction(array $history = [])
    {
        $fromStates      = array_keys($this->transitions);
        $this->initState = reset($fromStates);
        if ($history) {
            $this->history = $history;
            $currentState  = end($history);
            return $this->currentState = $currentState['state'];
        }
        $this->history[] = ['time' => time(), 'action' => 'init', 'state' => $this->initState];
        return $this->currentState = $this->initState;
    }

    /**
     * @return string Next state
     * @throws Error
     */
//    public function nextAction()
//    {
//        if (!isset($this->transitions[$this->currentState]) ||
//            !is_array($nextTransition = $this->transitions[$this->currentState])
//        ) {
//            throw new Error(
//                200,
//                sprintf('No further actions for current state "%s"', $this->currentState)
//            );
//        }
//        if (count($nextTransition) != 1) {
//            throw new Error(300, 'Unable to call next on a forked state');
//        }
//        return $this->doAction(key($nextTransition));
//    }
    /*
    public function reset()
    {
        $this->currentState = $this->initState;
        $this->previousState = null;
        $this->history = [];
    }
    */

}
