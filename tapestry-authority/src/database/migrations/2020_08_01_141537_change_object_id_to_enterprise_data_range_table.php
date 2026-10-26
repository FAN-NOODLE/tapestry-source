<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeObjectIdToEnterpriseDataRangeTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_data_range', function (Blueprint $table) {
            // 修改字段类型
            $table->string('object_id', 255)->change();
        });

        \Illuminate\Support\Facades\DB::statement("alter table enterprise_data_range modify column type int comment '状态 0-全部门店 1-区域 2-门店  3-商城  4 部门';");
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_data_range', function (Blueprint $table) {
            //
        });
    }
}
