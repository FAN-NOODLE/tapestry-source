<?php


namespace App\Jobs\Message;


use EasyWeChat\Kernel\Contracts\Arrayable;
use Illuminate\Contracts\Support\Jsonable;
use Illuminate\Support\Facades\Log;
use JsonSerializable;
use function PHPSTORM_META\type;
use ReflectionClass;

abstract class Message implements Jsonable, JsonSerializable
{
    public function toArray()
    {
        $fields = [];
        $ref = new ReflectionClass(static::class);
        foreach ($ref->getProperties() as $key => $property) {
            $value = null;
            if ($property->isPublic()) {
                $value = $this->{$property->getName()};
            } else {
                $getter = 'get' . ucfirst($property->getName());
                if (method_exists($this, $getter)) {
                    $value = $this->{$getter}();
                }
            }
            $fields[$property->getName()] = $value;
        }
        return $fields;
    }

    public function toJson($options = 0)
    {
        return json_encode($this->toArray(), $options);
    }

    public function __toString()
    {
        return $this->toJson();
    }

    public function jsonSerialize()
    {
        return $this->toArray();
    }


}