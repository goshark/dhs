var vm = new Vue({
    el:"#box",
    data(){
        return {
            ruleForm: {
                username: '',
                password: ''
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
						window.location.href = "login-index"
					}else{
						this.$message({
                            message: sign_res.Msg,
                            type: 'error'
                        });
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
                        this.$message({
                            message: log_res.Msg,
                            type: 'error'
                        });
						return;
					}
                } else {
                    return false;
                }
            });
        }
    }
})