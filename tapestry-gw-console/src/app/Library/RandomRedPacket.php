<?php
namespace App\Library;

class RandomRedPacket
{

    /**
     * 求一个数的平方
     * @param $n
     */
    public function sqr($n){
        return $n*$n;
    }

    /**
     * 生产min和max之间的随机数，但是概率不是平均的，从min到max方向概率逐渐加大。
     * 先平方，然后产生一个平方值范围内的随机数，再开方，这样就产生了一种“膨胀”再“收缩”的效果。
     */
    public function xRandom($bonus_min,$bonus_max){
        $sqr = intval($this->sqr($bonus_max-$bonus_min));
        $rand_num = rand(0, ($sqr-1));
        return intval(sqrt($rand_num));
    }
    /**
     *
     * @param $bonus_count 红包个数
     * @param $bonus_max 每个小红包的最大额
     * @param $bonus_min 每个小红包的最小额
     * @return array 存放生成的每个小红包的值的一维数组
     */
    public function getBonus($bonus_count, $bonus_max, $bonus_min) {
        $result = array();

        $average = ($bonus_max + $bonus_min) / 2;

        for ($i = 0; $i < $bonus_count; $i++) {
            //因为小红包的数量通常是要比大红包的数量要多的，因为这里的概率要调换过来。
            //当随机数>平均值，则产生小红包
            //当随机数<平均值，则产生大红包
            if (rand($bonus_min, $bonus_max) > $average) {
                // 在平均线上减钱
                $temp = $bonus_min + $this->xRandom($bonus_min, $average);
                $result[$i] = $temp;
            } else {
                // 在平均线上加钱
                $temp = $bonus_max - $this->xRandom($average, $bonus_max);
                $result[$i] = $temp;
            }
        }
        return $result;
    }
}
