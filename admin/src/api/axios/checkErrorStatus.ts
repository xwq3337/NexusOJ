export function checkErrorStatus(status: number | undefined, message: string | undefined, callback: (errorMessage: string) => any) {
    let errorMessage = message ?? ''
    switch (status) {
        case 200:
            errorMessage = '请求成功';
            break;
        case 201:
            errorMessage = '资源创建成功';
            break;
        case 204:
            errorMessage = '请求成功，但无返回内容';
            break;
        case 301:
            errorMessage = '请求的资源已永久移动到新位置';
            break;
        case 302:
            errorMessage = '请求的资源临时从不同的 URI 响应请求';
            break;
        case 304:
            errorMessage = '资源未修改';
            break;
        case 400:
            errorMessage = '客户端错误，请求格式或参数有误！';
            break;
        case 401:
            errorMessage = '身份认证不通过';
            break;
        case 403:
            errorMessage = '用户得到授权，但是访问是被禁止的!';
            break;
        case 404:
            errorMessage = '未找到目标资源!';
            break;
        case 405:
            errorMessage = '请求方法不被允许';
            break;
        case 408:
            errorMessage = '请求超时';
            break;
        case 409:
            errorMessage = '请求冲突';
            break;
        case 410:
            errorMessage = '请求的资源在服务器上已经不再可用';
            break;
        case 429:
            errorMessage = '请求过多，请稍后再试';
            break;
        case 500:
            errorMessage = '服务器内部错误';
            break;
        case 501:
            errorMessage = '服务器不支持当前请求所需要的某个功能';
            break;
        case 502:
            errorMessage = '无效网关';
            break;
        case 503:
            errorMessage = '服务不可用，服务器暂时过载或维护';
            break;
        case 504:
            errorMessage = '网关超时';
            break;
        default:
            errorMessage = '未知错误';
            break;
    }
    if (errorMessage.length > 0) {
        callback(`checkErrorStatus:${errorMessage}`)
    }
}