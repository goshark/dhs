const squidtemplate = `
          <div class="squidbox">
                <el-form :label-position="labelPosition"label-width="110px" ref="formLabelAlign" class="demo-ruleForm">
                    <el-form-item label="HOSTS文件" >
                        <el-input type="textarea" v-model="configuration.hosts"></el-input>
                    </el-form-item>
                    <el-form-item>
                        <el-button class="submitbut" type="primary" @click="submitHosts('formLabelAlign')">提交</el-button>
                    </el-form-item>
                </el-form>
                <el-form :label-position="labelPosition"label-width="110px" ref="formLabelAlign" class="demo-ruleForm">
                    <el-form-item label="squid配置" >
                        <el-input type="textarea" v-model="configuration.config"></el-input>
                    </el-form-item>
                    <el-form-item>
                        <el-button class="submitbut" type="primary" @click="submitConfig('formLabelAlign')">提交</el-button>
                    </el-form-item>
                </el-form>
          </div>
                `
var squidTemplate = {
    template:squidtemplate,
    data:function(){
        return {
            labelPosition: 'left',
        }
    },
	props:["configuration"],
    methods:{
        submitHosts(formName) {
            if(this.configuration){
                  const sendobj ={
                     url:"/squid/userform", 
                     data:{jsoninfo:JSON.stringify({name:"hosts",txt:this.configuration.hosts})}
                 }
                 var notund = underwaydata(this,"修改HOSTS文件")
                 setTimeout(()=>{
                         var SendAjaxdata = SendAjax(sendobj,"post")	
                         notund.close()
                         notifyshow(this,{bool:SendAjaxdata.Success,data:SendAjaxdata.Msg})
                         this.$refs[formName].resetFields();
                         this.$emit("listen",{ name:"配置文件返回" })
                     },500)
             }else{
                 this.$message({
                   message: '配置不能为空',
                   type: 'warning'
                 });
             }
         },
         submitConfig(formName) {
            if(this.configuration){
                  const sendobj ={
                     url:"/squid/userform", 
                     data:{jsoninfo:JSON.stringify({name:"config",txt:this.configuration.config})}
                 }
                 var notund = underwaydata(this,"修改squid配置")
                 setTimeout(()=>{
                         var SendAjaxdata = SendAjax(sendobj,"post")	
                         notund.close()
                         notifyshow(this,{bool:SendAjaxdata.Success,data:SendAjaxdata.Msg})
                         this.$refs[formName].resetFields();
                         this.$emit("listen",{ name:"配置文件返回" })
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