<?php
/**
 * Created by PhpStorm.
 * User: simple
 * Date: 2018/11/14
 * Time: 10:23 AM
 */

namespace App\Caches;

interface Cache
{
    public function key($params, $prefix);

    public function set($key, $value, $ttl = 3600);

    public function setArray($key, $array, $ttl = 3600);

    public function get($key);

    public function getArray($key);

    public function updateField($key, $field, $value);

    public function decr($key, $num);

    public function incy($key, $num);

    public function decrField($key, $field, $num);

    public function incrField($key, $field, $num);

    public function delete($keys);
}
