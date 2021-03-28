## 2.01 武器: postman

一款强大的后端调试工具

#### 1. 基本介绍

```
1. 各项目独立coollection，各组独立fold
2. 支持各类请求，post/get/put/option
3. 支持collection环境变量、全局环境变量

```

#### 2. pre script (login)

```
// 定义发送登录接口请求方法
function sendLoginRequest() {

   //定义请求体
    var data = {
       "name":"admin",
       "password": "1234567",
    }

    // 构造一个 POST raw 格式请求。这里需要改成你们实际登录接口的请求参数。
    const loginRequest = {
        url: 'http://127.0.0.1:8081/api/login',
        method: 'POST',
        header:'Content-Type:application/json',
        // body: JSON.stringify(data)
        body: {
		     mode:'raw',
			 raw:JSON.stringify(data)
        }
    };

    // 发送请求。
    // pm.sendrequest 参考文档: https://www.apifox.cn/help/app/scripts/api-references/pm-reference/#pm-sendrequest
    pm.sendRequest(loginRequest, function (err, res) {
        if (err) {
            console.log(err);
        } else {
            // 读取接口返回的 json 数据。
            // 如果你的 token 信息是存放在 cookie 的，可以使用 pm.cookies.get('token') 方式获取。
            // pm.cookies 参考文档：https://www.apifox.cn/help/app/scripts/api-references/pm-reference/#pm-cookies
            const jsonData = res.json();
            console.log("xxx", jsonData)
            // 将 accessToken 写入环境变量 Authorization

            pm.globals.set("token", jsonData.result.token);
            pm.globals.get("token");
            console.log("xxx 2", pm.globals.get("token"))

            // pm.globals.get("variable_key");
        }
    });
}


// 获取环境变量里的 Authorization
const accessToken = pm.environment.get('token');

// 如 Authorization 没有值,则执行发送登录接口请求
if (!accessToken ) {
    sendLoginRequest();
}

```
