{{define "footer"}}
<footer>
    	<div class="container">
    		<p class="text-center">Beego © 雪山飞鹄</p>
    		<address class="text-center">
			  <strong>Twitter, Inc.</strong><br>
			  795 Folsom Ave, Suite 600<br>
			  San Francisco, CA 94107<br>
			  <abbr title="Phone">P:</abbr> (123) 456-7890
			</address>
    	</div>
    </footer>
    <script src="http://cdn.bootcss.com/jquery/1.11.1/jquery.min.js"></script>
    <script src="http://cdn.bootcss.com/bootstrap/3.2.0/js/bootstrap.min.js"></script>
    <script type="text/javascript">
    function reload(){
        document.getElementById("captcha").src="/captcha?"+Math.random();
    }
    </script>
</body>
</html>
{{end}}