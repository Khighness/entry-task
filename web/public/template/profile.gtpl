<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <title>Profile</title>
</head>
<body>
    <form method="post" enctype="multipart/form-data">
        <input type="file" name="profile_picture" />
        <input type="submit" />
    </form>
    <p>姓名：{{.Username}}</p>
</body>
</html>