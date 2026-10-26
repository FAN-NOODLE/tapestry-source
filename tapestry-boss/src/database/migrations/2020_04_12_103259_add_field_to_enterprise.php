<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldToEnterprise extends Migration
{
    /**
     * Run the migrations.
     * @author lijian
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise', function (Blueprint $table) {
            // 新版企业字段
            $table->integer('account_limit')->default(0)->after('status')->comment('最大允许创建账户数量');
            $table->integer('account_start_at')->default(0)->after('status')->comment('账户有效期开始日期');
            $table->integer('account_expire_at')->default(0)->after('status')->comment('账户有效期失效日期');
            $table->tinyInteger('version')->default(1)->after('status')->comment('版本 1=标准版，2=门店商城版，3=中心商城版');
            $table->text('available_system')->default('[1]')->comment('企业可用系统,json字符存储[1-CRS后台 2-店长小程序 3-天擎CDP]');
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
            $table->dropColumn('account_limit');
            $table->dropColumn('account_start_at');
            $table->dropColumn('account_expire_at');
            $table->dropColumn('version');
        });
    }
}
