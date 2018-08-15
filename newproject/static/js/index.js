var Vm = new Vue({
    el:"#box",
    data(){
        return {
            // 操作的数据
            operaTionaldata:{},
			// 操作后返回
			modelBack:[],
			// 双击热备 data
            showData:[],
            // 导航展示
            sidebarshow:false,
			configUration:{post:"",hosts:""},
            username: "54166564168419",
            collapse: false,
            fullscreen: false,
            message: 2,
            collapse: false,
            routerpath:"server",
            redactVisible:false,
            dialogVisible:false,
			slaveVisible:false,
            items: [
                {
                    icon: '/static/img/host-color.png',
                    index: 'server',
                    title: '服务器'
                },
                {
                    icon: '/static/img/squid-solo.png',
                    index: 'squid',
                    title: 'squid'
                }
            ],
        }
    },
    created(){
        if(localStorage.getItem("name")){
            this.username = localStorage.getItem("name")
        }
		this.getdata()
		this.getQueryhosts()
    },
    methods:{
		// 双机热备的服务器数据
		getdata(){
			const sendobj ={
	            url:"/home/configlist", 
	            data:{jsoninfo:'{"condition":[]}'}
	        }
			var SendAjaxdata = SendAjax(sendobj,"post")
			if(!SendAjaxdata.Success){
				this.modelBack = SendAjaxdata.Msg
			}
		},
		// 配置文件获取
		getQueryhosts(){
			var sendobj1 ={
	            url:"/squid/query-hosts", 
	            data:{jsoninfo:'{"name":"hosts"}'}
	        }
            var SendAjax1 = SendAjax(sendobj1,"post")	
            if(!SendAjax1.Success){
                this.configUration['hosts'] = SendAjax1.Msg
            }

            var sendobj2 ={
	            url:"/squid/query-hosts", 
	            data:{jsoninfo:'{"name":"config"}'}
	        }
            var SendAjax2 = SendAjax(sendobj2,"post")	
            if(!SendAjax2.Success){
                this.configUration['config'] = SendAjax2.Msg
            }
		},
        showcloselisten(data){
            if(data.name == '添加关闭'){ 
                this.dialogVisible = false
            }else if(data.name == '编辑关闭'){
                this.redactVisible = false
            }else if(data.name == '设置从关闭'){
                this.slaveVisible = false
            }else if(data.name == '返回'){
				this.getdata()
			}
        },
        showlisten(data){
            if(data.name == '添加服务器'){
                this.dialogVisible = true
            }else if(data.name == '编辑'){
                this.redactVisible = true
                this.operaTionaldata = data.data
            }else if(data.name == '设置从'){
				this.slaveVisible = true
				this.operaTionaldata = data.data
			}else if(data.name == '设置主返回'){
				this.getdata()
			}else if(data.name == '配置文件返回'){
				this.getQueryhosts()
			}
        },
        mounted(){
            if(document.body.clientWidth < 1500){
            }
        },
        selectclick(key){
            this.sidebarshow = false
            this.routerpath = key
          },
        handleFullScreen(){
            let element = document.documentElement;
            if (this.fullscreen) {
                if (document.exitFullscreen) {
                    document.exitFullscreen();
                } else if (document.webkitCancelFullScreen) {
                    document.webkitCancelFullScreen();
                } else if (document.mozCancelFullScreen) {
                    document.mozCancelFullScreen();
                } else if (document.msExitFullscreen) {
                    document.msExitFullscreen();
                }
            } else {
                if (element.requestFullscreen) {
                    element.requestFullscreen();
                } else if (element.webkitRequestFullScreen) {
                    element.webkitRequestFullScreen();
                } else if (element.mozRequestFullScreen) {
                    element.mozRequestFullScreen();
                } else if (element.msRequestFullscreen) {
                    // IE11
                    element.msRequestFullscreen();
                }
            }
            this.fullscreen = !this.fullscreen;
        },
        // 关闭模态框
        handleCommand(command){
            if(command == 'loginout'){
                res = SendAjax({url:"/user/logout"})
				if (res.Success == 0){
					window.location.href="/"
				}else{
					alert(res.Msg);
				}
            }
        },
    },
    components:{
        servertemplate:serverTemplate,
        squidtemplate:squidTemplate,
        serveradd:serveraddMode,
        serverredact: serverredactMode,
		serverslave:serverslaveMode
    }
})
// 操作反馈
function notifyshow(_this,data){
		if(!data.bool){
			_this.$notify({
	          title: '成功',
	          message: data.data,
	          type: 'success'
	        });
		}else{
			_this.$notify({
	          title: '失败',
	          message: data.data,
	          type: 'error'
	        });
		}
}
// 正在进行
function underwaydata(_this,name){
	return _this.$notify.info({
        title: '正在进行'+name,
        iconClass:'el-icon-loading',
        duration:0
    });
}

//  移除
function remove(key){
    localStorage.removeItem(key)
}
