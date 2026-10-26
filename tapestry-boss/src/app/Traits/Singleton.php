<?php


namespace App\Traits;


/**
 * Trait Singleton
 * @package App\Traits
 */
trait Singleton
{
    /**
     * @var
     */
    private static $instance;

    /**
     * @param mixed ...$args
     * @return static
     */
    static function getInstance(...$args)
    {
        if(!isset(self::$instance)){
            self::$instance = new static(...$args);
        }
        return self::$instance;
    }
}
