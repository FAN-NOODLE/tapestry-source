<?php

use Illuminate\Database\Seeder;

class DatabaseSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {

    }

    private function posterSeeder()
    {
        $releaseEnterprise = [13, 120, 102, 73, 86, 100244242505984, 102756213248256];
        $data = [];
        foreach ($releaseEnterprise as $id) {
            array_push($data, [
                "enterprise_id" => $id,
                "poster_release_enterprise" => json_encode(\App\Common\Constant::MARKETING_PATH_NEW),
                "created_at" => time(),
                "updated_at" => time()
            ]);
        }

        \Illuminate\Support\Facades\DB::table("poster_release_enterprise")->insert($data);
    }
}
