var Vm = new Vue({
    el:"#box",
    data(){
        return {
            // 操作的数据
            operaTionaldata:{},
			modelBackfrom:[],
			modelBackmaster:[],
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
            routerpath:"mainprimary",
            redactVisible:false,
            mastermasterVisible:false,
            dialogVisible:false,
            slaveVisible:false,
            mysqltableVisible:false,
            items: [
                {
                    icon: '/static/img/host-color.png',
                    index: 'server',
                    title: '主从'
                },
                {
                    icon: '/static/img/squid-solo.png',
                    index: 'squid',
                    title: 'squid'
                },
                {
                    icon: '/static/img/host-color.png',
                    index: 'mainprimary',
                    title: '主主'
                }
            ],
        }
    },
    created(){
        if(localStorage.getItem("name")){
            this.username = localStorage.getItem("name")
        }
		this.getfromdata()
		this.getmasterdata()
		this.getQueryhosts()
    },
    methods:{
		// 主 从 务器数据
		getfromdata(){
            // {name:"type",value:"master",op:"="}
            // {name:"type",value:"slave",op:"="}
            // {name:"type",value:"",op:"="}
			const sendobj ={
	            url:"/home/configlist", 
	            data:{jsoninfo:'{"condition":[{"name":"type","value":"master","op":"="},{"name":"type","value":"slave","op":"="},{"name":"type","value":"","op":"="}]}'}
	        }
			var SendAjaxdata = SendAjax(sendobj,"post")
			if(!SendAjaxdata.Success && SendAjaxdata.Msg){
				this.modelBackfrom = SendAjaxdata.Msg
			}
        },
        // 主主务器数据
		getmasterdata(){
            // {name:"type",value:"master-master",op:"="}
            // {name:"type",value:"",op:"="}
			const sendobj ={
	            url:"/home/configlist", 
	            data:{jsoninfo:'{"condition":[{"name":"type","value":"master-master","op":"="},{"name":"type","value":"","op":"="}]}'}
	        }
			var SendAjaxdata = SendAjax(sendobj,"post")
			if(!SendAjaxdata.Success && SendAjaxdata.Msg){
				this.modelBackmaster = SendAjaxdata.Msg
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
            }else if(data.name == '主主关闭'){
                this.mastermasterVisible = false
            }else if(data.name == '编辑关闭'){
                this.redactVisible = false
            }else if(data.name == '设置从关闭'){
                this.slaveVisible = false
            }else if(data.name == '数据库' || data.name == '设置返回忽律数据库'){
                this.mysqltableVisible = false
            }else if(data.name == '返回'){
				this.getfromdata()
		        this.getmasterdata()
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
			}else if(data.name == '设置主主'){
				this.mastermasterVisible = true
				this.operaTionaldata = data.data
			}else if(data.name == '设置主返回' || data.name == '移除返回'){
				this.getfromdata()
		        this.getmasterdata()
			}else if(data.name == '配置文件返回'){
				this.getQueryhosts()
			}else if(data.name == '设置主数据库'){
                this.mysqltableVisible = true
                this.operaTionaldata = data.data
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
        mainprimarytemplate:mainPrimaryTemplate,
        squidtemplate:squidTemplate,
        serveradd:serveraddMode,
        serverredact: serverredactMode,
        serverslave:serverslaveMode,
        mysqltable:mysqltableMode,
        mastermaster:mastermasterMode
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
