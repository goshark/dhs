const mysqlittemplate = `
          <div class="squidbox">
              <el-form :label-position="labelPosition"label-width="110px" ref="formLabelAlign" class="demo-ruleForm">
                  <el-form-item label="配置文件" >
                     <el-input type="textarea" v-model="configuration"></el-input>
                  </el-form-item>
                  <el-form-item>
                      <el-button type="primary" @click="submitForm('formLabelAlign')">提交</el-button>
                  </el-form-item>
              </el-form>
          </div>
                `
var mysqlitTemplate = {
    template:mysqlittemplate,
    data:function(){
        return {
            labelPosition: 'left',
        }
    },
	props:["configuration"],
    methods:{
        submitForm(formName) {
           if(this.configuration){
			 	const sendobj ={
		            url:"/squid/userform", 
		            data:{jsoninfo:JSON.stringify({name:"hosts",txt:this.configuration})}
		        }
				var notund = underwaydata(this,"设置从服务器")
				setTimeout(()=>{
						var SendAjaxdata = SendAjax(sendobj,"post")	
						notund.close()
						notifyshow(this,{bool:SendAjaxdata.Success,data:SendAjaxdata.Msg})
	                    this.$refs[formName].resetFields();
						const json = { name:"配置文件返回" }
		                this.$emit("listen",json)
					},500)
			}else{
				this.$message({
		          message: '配置不能为空',
		          type: 'warning'
		        });
			}
           
               
        },
    }
};