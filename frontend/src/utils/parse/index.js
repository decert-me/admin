export function getJwt(token) {
    let strings = token.split("."); //截取token，获取载体
    try {
        var userinfo = JSON.parse(decodeURIComponent(escape(window.atob(strings[1].replace(/-/g, "+").replace(/_/g, "/"))))); //解析，需要吧‘_’,'-'进行转换否则会无法解析
        // 判断token有效期
        const now = parseInt(new Date().getTime() / 1000) + (2 * 60 * 60) ;
        var status = now > userinfo.exp ? false : true;   //false 超时 : true 未超时
    } catch (error) {
        localStorage.clear();
        window.history.go(0);
    }

    return status
}