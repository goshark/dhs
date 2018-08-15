var serveraddMode = { 
    template:`
        <el-dialog
        title="添加"
        :visible.sync="dialogvisible"
        width="30%"
        :before-close="handleCloses">
            <el-form :label-position="labelPosition" :model="formLabelAlign" :rules="rules" ref="formLabelAlign" label-width="100px" class="demo-ruleForm">
                    <el-form-item label="Ip" prop="ip">
                        <el-input v-model="formLabelAlign.ip" placeholder="请输入服务Ip"></el-input>
                    </el-form-item>
                    <el-form-item label="用户" prop="dbuser">
                        <el-input v-model="formLabelAlign.dbuser" placeholder="请输入用户名"></el-input>
                    </el-form-item>
                    <el-form-item label="密码" prop="dbpass">
                        <el-input v-model="formLabelAlign.dbpass" placeholder="请输入密码"></el-input>
                    </el-form-item>
                    <el-form-item label="端口" prop="dbport">
                        <el-input v-model="formLabelAlign.dbport" placeholder="请输入端口号"></el-input>
                    </el-form-item>
                    <el-form-item>
                        <el-button type="primary" @click="submitForm('formLabelAlign')">立即创建</el-button>
                    </el-form-item>
            </el-form>
        </el-dialog>
       `, 
    props:["dialogvisible","operationaldata"],
    data:function(){
		  var validateip = (rule, value, callback) => {
	          var reg=/^(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])\.(\d{1,2}|1\d\d|2[0-4]\d|25[0-5])$/;
	          if(!reg.test(value)){
	          callback(new Error('ip段格式输入有误'));
	          }else{
	          callback()
	          }
          };
        return {
             //server
            labelPosition: 'right',
            formLabelAlign: {
                ip: '',
                dbuser: '',
                dbpass: '',
                dbport: '',
            },
            rules: {
                ip: [
                    { required: true, message: '请输入ip', trigger: 'blur' },{ validator: validateip,trigger: 'blur' },
                ],
                dbuser: [
                    { required: true, message: '请输入账户', trigger: 'blur' }
                ],
                dbpass: [
                    { required: true, message: '请输入用户密码', trigger: 'blur' },
                ],
                dbport: [
                    { required: true, message: '请输入端口', trigger: 'blur' }
                ]
            },
           
        }
    },
    methods:{
         // 关闭模态框
        handleCloses(){
			this.$emit("modelBack",{data:"sadasd"})
            this.$emit("closelisten",{name:"添加关闭"})
			this.$refs["formLabelAlign"].resetFields();
        },
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    var sendobj ={
                        url:"/home/append-list", 
                        data:{jsoninfo:JSON.stringify(this.formLabelAlign)}
                    }
					this.$emit("closelisten",{name:"添加关闭"})
					var notund = underwaydata(this,"添加服务器")
					setTimeout(()=>{
						var returnData = SendAjax(sendobj,"post")
						notund.close()
						notifyshow(this,{bool:returnData.Success,data:returnData.Msg})
						this.$refs[formName].resetFields();
						this.$emit("closelisten",{name:"返回",data:returnData})
					},500)
                } else {
                    console.log('error submit!!');
                    return false;
                }
            });
        },
    }
};


var serverredactMode = {
    template:`
        <el-dialog
        title="编辑修改"
        :visible.sync="redactvisible"
        width="30%"
        :before-close="handleClose">
            <el-form :label-position="labelPosition" :model="operationaldata" ref="operationaldata" label-width="100px" class="demo-ruleForm">
                <el-form-item label="服务Ip" prop="ip">
                    <el-input v-model="operationaldata.ip"></el-input> 
                </el-form-item>
                <el-form-item label="region" prop="name">
                    <el-input v-model="operationaldata.name"></el-input>
                </el-form-item>
                <el-form-item label="type" prop="address">
                    <el-input v-model="operationaldata.address"></el-input>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="submitForm('operationaldata')">立即创建</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>`,
       props:["redactvisible","operationaldata"],
    data:function(){
        return {
             //server
            labelPosition: 'right',
            formLabelAlign: {
                ServerIP: '',
                region: '',
                type: ''
            },
            rules: {
                ServerIP: [
                { required: true, message: '请输入活动名称', trigger: 'blur' },
                { min: 3, max: 5, message: '长度在 3 到 5 个字符', trigger: 'blur' }
                ],
                region: [
                { required: true, message: '请选择活动区域', trigger: 'change' }
                ],
                type: [
                    { required: true, message: '请选择活动区域', trigger: 'change' }
                ],
            },
        }
    },
    created(){
        console.log(this)
    },
    methods:{
         // 关闭模态框
        handleClose(){
            this.$emit("closelisten",{name:"编辑关闭"})
        },
        submitForm(formName) {
            this.$refs[formName].validate((valid) => {
                if (valid) {
                    alert('submit!');
                    const sendobj ={
                        url:"/home/append-list", 
                        data:{jsoninfo:JSON.stringify(this.operationaldata)}
                    }
               		
					var returnData = SendAjax(sendobj,"post")
					notifyshow(this,{bool:returnData.Success,data:returnData.Msg})
	                this.$emit("closelisten",{name:"编辑关闭",data:returnData})
            
                } else {
                    console.log('error submit!!');
                    return false;
                }
            });
        },
    }
};


var serverslaveMode = {
	template:`
        <el-dialog
        title="设置从"
        :visible.sync="slavevisible"
        width="30%"
        :before-close="handleClose">
            <el-form :label-position="labelPosition" :model="operationaldata" ref="operationaldata" label-width="100px" class="demo-ruleForm">
                <el-form-item label="主控IP" prop="ip">
                     <el-select v-model="operationaldata.masterip" placeholder="请选择">
					    <el-option
					      v-for="item in operationaldata.masterIpdata"
					      :key="item.ip"
					      :label="item.ip"
					      :value="item.ip">
					    </el-option>
					  </el-select>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="submitForm('operationaldata')">确认</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>`,
		 data:function(){
		        return {
		             //server
		            labelPosition: 'right',
		        }
		    },
       	props:["slavevisible","operationaldata"],
		methods:{
	         // 关闭模态框
	        handleClose(){
	            this.$emit("closelisten",{name:"设置从关闭"})
	        },
	        submitForm(formName) {
	            this.$refs[formName].validate((valid) => {
	                if (valid) {
	                	delete this.operationaldata['masterIpdata']
			           	this.operationaldata["type"] = "slave"
		                const sendobj ={
		                    url:"/home/action", 
		                    data:{jsoninfo:JSON.stringify(this.operationaldata)}
		                }		
						var notund = underwaydata(this,"设置从服务器")
						this.$emit("closelisten",{name:"设置从关闭"})
						setTimeout(()=>{
							var returnData = SendAjax(sendobj,"post")
							notund.close()
							notifyshow(this,{bool:returnData.Success,data:returnData.Msg})
							this.$emit("closelisten",{name:"返回",data:returnData})
						},500)
	                } else {
	                    console.log('error submit!!');
	                    return false;
	                }
	            });
	        },
    	}
	
	}
	
