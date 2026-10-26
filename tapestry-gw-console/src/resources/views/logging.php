<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>日志排查</title>
    <style>
        * {
            margin: 0;
            padding: 0;
        }
        .container {
            padding: 50px;
        }
        h1 {
            font-size: 30px;
            line-height: 80px;
        }
        h3 {
            font-size: 20px;
            line-height: 60px;
        }
        h3 span {
            display: inline-block;
            padding-left: 20px;
            padding-right: 20px;
        }
        .delete {
            position: fixed;
            right: 50px;
            top: 74px;
            background-color: #000000;
            border-radius: 50px;
        }
        .delete a {
            padding: 8px 20px;
            color: #fff;
            display: inline-block;
            text-decoration: none;
            font-size: 14px;
        }
        .data pre {
            font-size: 13px;
            line-height: 16px;
        }
        p {
            font-size: 12px;
            padding: 8px 10px;
            border: 1px solid #000;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
<div class="container">
    <h1>日志排查</h1>
    <h3>日志大小：<?=$size?> K <span>|</span> 日志行数：<?=$line?></h3>
    <p><?=$file?></p>
    <div class="delete"><a href="?delete=true&p=<?=date('md')?>">清空</a></div>
    <div class="data">
        <pre><?=$data?></pre>
    </div>
</div>
</body>
</html>
