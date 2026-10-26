<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeEnterpriseOpenerpTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise', function (Blueprint $table) {
            // 增加更新字段
            $table->unsignedInteger('erp_status')->default(1)->after('status')->comment('ERP接口调用权限 1：可以 2：不可以');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise', function (Blueprint $table) {
            // 增加更新字段
            $table->dropColumn('erp_status');
        });
    }
}
