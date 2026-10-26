<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddStartFeeLockTimeToEnterprise extends Migration
{
    /**
     *
     */
    public function up()
    {
        if (Schema::hasTable('enterprise')) {
            Schema::table('enterprise', function (Blueprint $table) {
                $table->integer('start_fee')->default(0)->after('address')->comment('分账起始金额, 大于此金额分账');
                $table->integer('proportion')->default(0)->after('address')->comment('分账比例 精度小数点后两位, 实际百分比为 值/10000');
                $table->integer('lock_time')->default(0)->after('address')->comment('时间分 每单锁定这个时间,到了之后进行分账');
                $table->tinyInteger('separate_switch')->default(1)->after('address')->comment('1未开启 2开启');
            });
        }
    }

    /**
     *
     */
    public function down()
    {
        if (Schema::hasTable('enterprise')) {
            Schema::table('enterprise', function (Blueprint $table) {
                $table->dropColumn('start_fee');
                $table->dropColumn('proportion');
                $table->dropColumn('lock_time');
                $table->dropColumn('separate_switch');
            });
        }
    }
}
