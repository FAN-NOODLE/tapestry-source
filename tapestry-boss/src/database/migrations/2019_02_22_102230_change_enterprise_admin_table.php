<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class ChangeEnterpriseAdminTable extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_admin', function (Blueprint $table) {
            // 增加更新字段
            $table->unsignedBigInteger('operator_id')->default(0)->after('status')->comment('操作人ID');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_admin', function (Blueprint $table) {
            //
        });
    }
}
