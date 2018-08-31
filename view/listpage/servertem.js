const servertemplate = `
                <div>
                <el-row>
                    <el-button @click="addserver">添加</el-button>
                    <el-dropdown @command="handleCommand">
                        <el-button type="primary" >
                            操作<i class="el-icon-arrow-down el-icon--right"></i>
                        </el-button>
                        <el-dropdown-menu slot="dropdown">
                            <el-dropdown-item command="redact" disabled>编辑</el-dropdown-item>
                            <el-dropdown-item command="main" :disabled="mainbool">设置主</el-dropdown-item>
                            <el-dropdown-item command="follow" :disabled="pagebool">设置从</el-dropdown-item>
                            <el-dropdown-item command="delete" :disabled="pagebool">删除</el-dropdown-item>
                        </el-dropdown-menu>
                    </el-dropdown>
                </el-row>
                <el-table
                ref="multipleTable"
                :data="modelback"
                tooltip-effect="dark"
                style="width: 100%;margin-top:20px"
                @selection-change="handleSelectionChange">
                        <el-table-column
                        type="selection"
                        width="55">
                        </el-table-column>
                        <el-table-column
                        prop="ip"
                        label="服务器Ip"
                        show-overflow-tooltip>
                        </el-table-column>
                        <el-table-column
                        prop="type"
                        label="主从"
                        >
                        <template slot-scope="scope">
							<span v-if="scope.row.type">
								<span v-if="scope.row.type == 'master'">
	                                主
	                            </span>
	                            <span v-else>
	                                从
	                            </span>  
							</span>
							<span v-else>
								无
							</span>     
                        </template>
                        </el-table-column>
						<el-table-column
                        label="备注"
                        >
                        </el-table-column>
                       
                </el-table>
                </div>
                `

var serverTemplate = {
    template:servertemplate,
    data:function(){
        return {
			updatabool:true,
            mainbool:true,
            pagebool:true,
            multipleSelection: []
        }
    },
 	props:["modelback"],
    methods:{
            // 关闭模态框
            handleCommand(command){
                if(command == 'loginout'){
                    remove('name')
                    window.location.href="/user/login.html"
                }else if(command == 'redact'){
                    this.redact()
                }else if(command == 'main'){
                    this.main()
                }else if(command == 'follow'){
                    this.follow()
                }else if(command == 'delete'){
                    this.Delete()
                }
            },
             // 移除
            Delete(){
                this.$confirm('将该服务器从 主从配置中移除, 是否继续?', '提示', {
                  confirmButtonText: '确定',
                  cancelButtonText: '取消',
                  type: 'warning'
                }).then(() => {
                                 var sendobj ={
                    url:"/home/remove-list",
                    data:{jsoninfo:JSON.stringify(this.multipleSelection[0])}
                }
                    var notund = underwaydata(this,"移除服务器")
                    setTimeout(()=>{
                            var returnData = SendAjax(sendobj,"post")
                            notund.close()
                            notifyshow(this,{bool:returnData.Success,data:returnData.Msg})
                            this.$emit("listen",{name:"移除返回"})
                           
                    },500)

                }).catch(() => {
                  this.$message({
                    type: 'info',
                    message: '已取消删除'
                  });
                });
            },
            // 编辑
            redact(){
                const json = {
                    name:"编辑",
                    data:this.multipleSelection[0]
                }
                this.$emit("listen",json)
            },
            // 设置主 
            main(){
                this.multipleSelection[0]["type"] = "master"
                var sendobj ={
                    url:"/home/action", 
                    data:{jsoninfo:JSON.stringify(this.multipleSelection[0])}
                }
				var notund = underwaydata(this,"设置主服务器")
				setTimeout(()=>{
					var returnData = SendAjax(sendobj,"post")
					notund.close()
					notifyshow(this,{bool:returnData.Success,data:returnData.Msg})
                    this.$emit("listen",{name:"设置主返回"})
                    if(!returnData.Success){
                        var sendobjs ={
                            url:"/home/get-db-list", 
                            data:{jsoninfo:JSON.stringify(this.multipleSelection[0])}
                        }
                        var mysqlData = SendAjax(sendobjs,"post")
                        const data = {
                            multipleSelection:this.multipleSelection[0],
                            mysqllist:mysqlData.Msg
                        }
                        const json = {
                            name:"设置主数据库",
                            data:data 
                        }
                        this.$emit("listen",json)
                    }
				},500)
            },
            // 设置从
            follow(){
				var masterIp = []
				for (var i = 0; i < this.modelback.length; i++){
					if(this.modelback[i].type == 'master'){
						masterIp.push(this.modelback[i])
					}
				}
				if(masterIp.length > 0){
					this.multipleSelection[0]['masterIpdata'] = masterIp
					const json = {
	                    name:"设置从",
	                    data:this.multipleSelection[0]
	                }
	                this.$emit("listen",json)
				}else{
					this.$message({
			          message: '请设置一个主控！ 最少有一个主控！！！',
			          type: 'warning'
			        });
				}
            },
            // 添加服务
            addserver(){
                this.name="更改的值"
                const json = {
                    name:"添加服务器",
                }
                this.$emit("listen",json)
            },
        	handleSelectionChange(val) {
	            if(val.length == 1){
					this.updatabool = false
					this.mainbool = false
	                this.pagebool = false
						if(val[0].type == 'master' || val[0].type == "slave"){
							this.mainbool = true
	                    	this.pagebool = true
						}else{
							this.mainbool = false
	                    	this.pagebool = false
                        }
                        
	            }else{
					this.updatabool = true
	                this.mainbool = true
                    this.pagebool = true
                    
	            }
	            this.multipleSelection = val;
        },
        // 查看是否 有主
        mainjudge(){
			var noset = 0
			var setmassav = 0
            for ( var i = 0 ; i < this.modelback.length; i++ ){
                if(this.modelback[i].type){
					if(this.modelback[i].type == "master"){
						setmassav ++
	                }
				}else{
					noset++
				}
            }
			if(noset == this.modelback.length){
				// 设置主
				return '设置主'
			}else if(setmassav == 1){
				// 已经有主这能设置从
				return '设置从'
			}
			
			noset = 0
			setmassav = 0
        }
    }
};





