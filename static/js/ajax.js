
var SendAjax = function(sendobj,type){
    return $.ajax({
                url:sendobj.url,
                data:sendobj.data,
                type:type,
                async:false,
                dataType: "json"
            }).responseJSON
    }
