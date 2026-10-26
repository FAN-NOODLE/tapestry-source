<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;

class AddFieldsEnterprise extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise', function (Blueprint $table) {
            $table->string('version_tag', 16)->default('')->after('status')->comment('版本标签;正式版,试用版');
            $table->unsignedBigInteger('industry_id')->default(0)->after('separate_switch')->comment('行业ID');
            $table->unsignedTinyInteger('business_source')->default(0)->after('status')->comment('业务来源; 0:-,1:企业微信,10:畅移（直签）,11:畅移（代理）,12:畅移（试用）20:畅鲟（直签）21:畅鲟（代理）22:畅鲟（试用）');
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        //
    }
}
