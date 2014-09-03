{{define "footer"}}
    <footer>
        <div class="container">
            <p class="text-center">Beego © 雪山飞鹄</p>
            <address class="text-center">
                <strong>Twitter, Inc.</strong><br>
			</address>
        </div>
    </footer>
    <script type="text/javascript">
    function reload(){
        document.getElementById("captcha").src="/captcha?"+Math.random();
    }
    </script>
</body>
</html>
{{end}}