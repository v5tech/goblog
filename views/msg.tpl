{{define "msgbox"}}
	{{if .flash.error}}
    	<div class="alert alert-danger alert-dismissible" role="alert">
            <button type="button" class="close" data-dismiss="alert">
                <span aria-hidden="true">&times;</span>
                <span class="sr-only">Close</span>
            </button>
            <strong>{{.flash.error}}</strong>
        </div>
    {{end}}
    {{if .flash.notice}}
        <div class="alert alert-success alert-dismissible" role="alert">
            <button type="button" class="close" data-dismiss="alert">
                <span aria-hidden="true">&times;</span>
                <span class="sr-only">Close</span>
            </button>
            <strong>{{.flash.notice}}</strong>
        </div>
    {{end}}
{{end}}