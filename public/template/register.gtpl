<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <title>Register</title>
    <style type="text/css">
        body{
            margin: 0;
            padding: 0;
            height: 100vh;
            background: #2F323A;
            background-size: cover;
        }
        .form{
            background: #000;
            z-index: 1;
            font-family: "Candara", sans-serif;
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 300px;
            padding: 0 45px 30px 45px;
            text-align: center;
            border-radius: 10px;
            opacity: 0.9;
        }
        .form h2{
            color: #fff;
            font-size: 28px;
            font-weight: 500;
            text-align: center;
            text-transform: uppercase;
        }
        .form .icons i{
            color: #fff;
            font-size: 25px;
            margin: 0 30px 30px 30px;
            transition: 0.8s;
            transition-property: color, transform;
        }
        .form .icons i:hover{
            color: #06C5CF;
            transform: scale(1.3);
        }
        .form input{
            outline: 0;
            background: none;
            font-size: 15px;
            color: #fff;
            text-align: center;
            width: 265px;
            margin-bottom: 30px;
            padding: 15px;
            box-sizing: border-box;
            border: 2.5px solid #2E86DE;
            border-radius: 25px;
            transition: 0.5s;
            transition-property: width;
        }
        .form input:hover{
            width: 300px;
        }
        .form input:focus{
            width: 300px;
        }
        .form button{
            outline: 0;
            background: none;
            color: #fff;
            font-size: 14px;
            text-transform: uppercase;
            width: 150px;
            padding: 15px;
            border: 2.5px solid #10AC84;
            border-radius: 25px;
            cursor: pointer;
            transition: 0.5s;
            transition-property: background, transform;
        }
        .form button:hover, .form button:active, .form button:focus{
            background: #10AC84;
            transform: scale(1.1);
        }
        .form .options{
            color: #bbb;
            font-size: 14px;
            margin: 20px 0 0;
        }
        .form .options a{
            text-decoration: none;
            color: #06C5CF;
        }
    </style>
</head>
<body>
<div class="form">
    <form action="/register" method="post">
        <h2>Entry Task</h2>
        <input type="text" name="username" placeholder="Username" required="required">
        <input type="password" name="password"  placeholder="Password" required="required">
        <button type="submit" name="button">Register</button>
        <p class="options">Already Registered? <a href="/login">Sign in</a></p>
    </form>
</div>
</body>
</html>