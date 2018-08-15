// 验证用户是否登录
(function(){
    var name = localStorage.getItem("name")
     if(!name){
         window.location.href="/user/login.html"
     }
 })()

