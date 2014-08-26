{{define "nav"}}
<a class="sr-only sr-only-focusable" href="#content">Skip to main content</a>
<header>
	<div class="container">
        <nav class="navbar navbar-default navbar-fixed-top" role="navigation">
            <div class="navbar-header">
                <button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#example-navbar-collapse">
                    <span class="sr-only">Toggle navigation</span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                    <span class="icon-bar"></span>
                </button>
                <a class="navbar-brand fa fa-home" href="http://github.com/astaxie/beego"> Beego</a>
            </div>
            <div class="collapse navbar-collapse" id="example-navbar-collapse">
                <ul class="nav navbar-nav">
                	<li {{if .IsHome}}class="active"{{end}}><a href="/">首页</a></li>
                    <li {{if .IsTopic}}class="active"{{end}}><a href="/topic/add">文章</a></li>
                    <li {{if .IsCategory}}class="active"{{end}}><a href="/category">分类</a></li>
                </ul>
                <ul class="nav navbar-nav navbar-right">
                    {{if .IsLogin}}
                        <li><a href="/user/{{.Username}}" class="fa fa-lock"> {{.Nickname}}</a></li>
                        <li><a href="/logout" class="fa fa-lock"> 退出</a></li>
                    {{else}}
                        <li><a href="/login" class="fa fa-user"> 登录</a></li>
                        <li><a href="/register" class="fa fa-user"> 注册</a></li>
                    {{end}}
                </ul>
            </div>
        </nav>
    </div>
</header>
{{end}}