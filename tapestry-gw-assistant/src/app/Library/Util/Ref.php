<?php


namespace App\Library\Util;


class Ref
{

    public static function methods($class,$stop=true)
    {
        $ref = new \ReflectionClass($class);
        $methods = [];
        foreach ($ref->getMethods() as $method) {
            if ($method->isPublic()) {
                $methods[] = $method->getName();
            }
        }
        $stop ? dd($methods) : dump($methods);
    }

    public static function propsClass($class,$stop=true)
    {
        $ref = new \ReflectionClass($class);
        $props = [];
        foreach ($ref->getProperties() as $prop) {
            if ($prop->isPublic()) {
                $props[] = $prop->getName();
            }
        }
        $stop ? dd($props) : dump($props);
    }

    public static function propsObj($class, $stop = true, $objExclude = true)
    {
        $ref = new \ReflectionClass($class);
        $props = [];
        foreach ($ref->getProperties() as $property) {
            $property->setAccessible(true);
            if ($objExclude && is_object($property->getValue($class))) {
                continue;
            }
            $props[$property->getName()] = $property->getValue($class);
        }
        ksort($props);
        Logger::html($props);
        $stop && dd($props);
    }
}