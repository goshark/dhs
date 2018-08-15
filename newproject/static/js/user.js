var vm = new Vue({
    el:"#box",
    data(){
        return {
            ruleForm: {
                username: 'oyx',
                password: '123456'
            },
            rules: {
                username: [
                    { required: true, message: '请输入用户名', trigger: 'blur' }
                ],
                password: [
                    { required: true, message: '请输入密码', trigger: 'blur' }
                ]
            }
        }
            
    },
    methods: {
        RegisterForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    sign_res = SendAjax({url:"register",data:{jsoninfo:JSON.stringify(this.ruleForm)}},'POST')
					
					if (sign_res.Success == 0){
						alert(sign_res.Msg)
						window.location.href = "login-index"
					}else{
						alert(sign_res.Msg)
					}
                   
                } else {
                    return false;
                }
            });
        },
        loginForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                   	log_res = SendAjax({url:"login",data:{jsoninfo:JSON.stringify(this.ruleForm)}},'POST')
					
					if(log_res.Success == 0){
						window.location.href = "/"
					}else{
						alert(log_res.Msg)
						return;
					}
                } else {
                    return false;
                }
            });
        }
    }
})