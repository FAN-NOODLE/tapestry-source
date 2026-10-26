<?php

use Illuminate\Support\Facades\Schema;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Database\Migrations\Migration;
use Illuminate\Support\Facades\DB;

class InsertEnterpriseAppTableActivityData extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::table('enterprise_app', function (Blueprint $table) {
            $add = [
                'identity' => 'activity',
                'created_at' => time()
            ];
            DB::table("enterprise_app")->insert($add);
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::table('enterprise_app', function (Blueprint $table) {
            //
        });
    }
}
